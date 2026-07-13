package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// CacheDirName is the name of the cache subdirectory within the workspace.
	CacheDirName = "cache"
	// SpecsDirName is the name of the specs subdirectory within the workspace.
	SpecsDirName = "specs"
	// MaxTTL is the maximum time-to-live for cached spec files.
	MaxTTL = 48 * time.Hour
	// MinTTL is the minimum time-to-live for cached spec files.
	MinTTL = 1 * time.Hour
)

type sourceType string

const (
	sourceURL   sourceType = "url"
	sourceLocal sourceType = "local"
)

// Cache resolves spec locations to local file paths, caching sources on disk.
type Cache struct {
	dir          string
	specsDir     string
	workspaceDir string
	cli          *httpClient
}

// New creates a cache rooted at workspaceDir/cache.
// The cache directory is created lazily on the first Resolve call.
func New(workspaceDir string) *Cache {
	return &Cache{
		dir:          filepath.Join(workspaceDir, CacheDirName),
		specsDir:     filepath.Join(workspaceDir, SpecsDirName),
		workspaceDir: workspaceDir,
		cli:          defaultHTTPClient(),
	}
}

// SetWorkspaceDir sets the workspace directory and updates all subdirectories.
func (c *Cache) SetWorkspaceDir(workspaceDir string) {
	c.dir = filepath.Join(workspaceDir, CacheDirName)
	c.specsDir = filepath.Join(workspaceDir, SpecsDirName)
	c.workspaceDir = workspaceDir
}

// Resolve takes a location (local path, file:// URL, or http(s):// URL)
// and returns a path to a local file containing the spec data.
//
// Caching rules:
//   - URLs are always cached on disk.
//   - Local paths inside workspaceDir/specs are returned as-is (not cached).
//   - Local paths outside workspaceDir/specs are cached on disk.
func (c *Cache) Resolve(ctx context.Context, location string) (string, error) {
	if location == "" {
		return "", errors.New("empty location")
	}

	normalized, stype, err := normalizeLocation(location)
	if err != nil {
		return "", err
	}

	if stype == sourceLocal {
		if resolved, ok := c.resolveSpecsPath(location); ok {
			return resolved, nil
		}
		if c.isInsideSpecs(normalized) {
			return normalized, nil
		}
	}

	if mkdirErr := os.MkdirAll(c.dir, 0750); mkdirErr != nil {
		return "", fmt.Errorf("create cache dir: %w", mkdirErr)
	}

	hash := cacheKey(normalized)
	specPath := filepath.Join(c.dir, hash+".spec")
	metaPath := filepath.Join(c.dir, hash+".meta")

	if c.hitCache(normalized, stype, metaPath, specPath) {
		return specPath, nil
	}

	data, modTime, err := c.loadSource(ctx, normalized, stype)
	if err != nil {
		return "", err
	}

	if writeErr := os.WriteFile(filepath.Clean(specPath), data, 0600); writeErr != nil {
		return "", fmt.Errorf("write cache file: %w", writeErr)
	}

	ttl := randomTTL()
	meta := fileMeta{
		Source:     normalized,
		SourceType: string(stype),
		CachedAt:   time.Now(),
		ModTime:    modTime,
		TTLSec:     int(ttl.Seconds()),
	}
	if metaErr := writeMeta(metaPath, meta); metaErr != nil {
		return "", fmt.Errorf("write meta file: %w", metaErr)
	}

	return specPath, nil
}

func (c *Cache) isInsideSpecs(path string) bool {
	if c.specsDir == "" {
		return false
	}
	cleanPath := filepath.Clean(path)
	cleanSpecs := filepath.Clean(c.specsDir)
	return strings.HasPrefix(cleanPath, cleanSpecs+string(filepath.Separator)) || cleanPath == cleanSpecs
}

// resolveSpecsPath checks if the location is a relative path starting with
// "specs/" or "./specs/" and resolves it relative to the workspace directory.
// Returns the resolved path only if the file actually exists there.
func (c *Cache) resolveSpecsPath(location string) (string, bool) {
	clean := filepath.Clean(location)
	if !strings.HasPrefix(clean, SpecsDirName+string(filepath.Separator)) && clean != SpecsDirName {
		return "", false
	}
	resolved := filepath.Join(c.workspaceDir, clean)
	if _, err := os.Stat(resolved); err != nil {
		return "", false
	}
	return resolved, true
}

func (c *Cache) hitCache(normalized string, stype sourceType, metaPath, specPath string) bool {
	meta, readErr := readMeta(metaPath)
	if readErr != nil || meta.IsExpired() {
		return false
	}

	switch stype {
	case sourceLocal:
		fi, statErr := os.Stat(normalized)
		if statErr != nil || fi.ModTime().After(meta.ModTime) {
			return false
		}
	case sourceURL:
	}

	if _, specErr := os.Stat(specPath); specErr != nil {
		return false
	}
	return true
}

func (c *Cache) loadSource(ctx context.Context, normalized string, stype sourceType) ([]byte, time.Time, error) {
	switch stype {
	case sourceLocal:
		fi, statErr := os.Stat(normalized)
		if statErr != nil {
			return nil, time.Time{}, fmt.Errorf("stat %s: %w", normalized, statErr)
		}
		data, readErr := os.ReadFile(normalized)
		if readErr != nil {
			return nil, time.Time{}, fmt.Errorf("read %s: %w", normalized, readErr)
		}
		return data, fi.ModTime(), nil

	case sourceURL:
		data, getErr := c.cli.Get(ctx, normalized)
		if getErr != nil {
			return nil, time.Time{}, fmt.Errorf("download %s: %w", normalized, getErr)
		}
		return data, time.Time{}, nil

	default:
		return nil, time.Time{}, fmt.Errorf("unknown source type %q", stype)
	}
}

// Exists checks whether the location is accessible.
// For local paths it checks [os.Stat].
// For file:// URLs it resolves the path and checks [os.Stat].
// For http(s):// URLs it checks the cache first, then does a HEAD request.
// Returns nil if accessible, [LocationError] otherwise.
func (c *Cache) Exists(ctx context.Context, location string) error {
	if location == "" {
		return errors.New("empty location")
	}

	isURL := strings.HasPrefix(location, "https://") || strings.HasPrefix(location, "http://")
	isFileURL := strings.HasPrefix(location, "file://")

	switch {
	case isFileURL:
		path, pathErr := fileURIToPath(location)
		if pathErr != nil {
			return &LocationError{Location: location, Type: "file", Err: pathErr}
		}
		if resolved, ok := c.resolveSpecsPath(path); ok {
			path = resolved
		}
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			return &LocationError{Location: location, Type: "file", Err: fmt.Errorf("file not found at %s", path)}
		}
		return nil

	case isURL:
		return c.existsURL(ctx, location)

	default:
		path := expandTilde(location)
		if resolved, ok := c.resolveSpecsPath(path); ok {
			path = resolved
		}
		if !filepath.IsAbs(path) {
			absPath, absErr := filepath.Abs(path)
			if absErr == nil {
				path = absPath
			}
		}
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			return &LocationError{Location: location, Type: "file", Err: fmt.Errorf("file not found at %s", path)}
		}
		return nil
	}
}

// normalizeLocation converts any location to a canonical absolute path or URL.
func normalizeLocation(location string) (string, sourceType, error) {
	isURL := strings.HasPrefix(location, "https://") || strings.HasPrefix(location, "http://")
	isFileURL := strings.HasPrefix(location, "file://")

	switch {
	case isFileURL:
		path, err := fileURIToPath(location)
		if err != nil {
			return "", "", fmt.Errorf("file URL: %w", err)
		}
		return path, sourceLocal, nil
	case isURL:
		return location, sourceURL, nil
	default:
		path := expandTilde(location)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", "", fmt.Errorf("convert to absolute path: %w", err)
		}
		return absPath, sourceLocal, nil
	}
}

// expandTilde replaces ~/ and ~\ prefix with the user's home directory.
// Works on both Unix and Windows.
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func cacheKey(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:16])
}

func randomTTL() time.Duration {
	n := rand.Int64N(int64(MaxTTL - MinTTL))
	return MinTTL + time.Duration(n)
}

// existsURL checks if a URL is accessible, using cache if available.
func (c *Cache) existsURL(ctx context.Context, url string) error {
	// Check cache first
	if c.dir != "" {
		hash := cacheKey(url)
		specPath := filepath.Join(c.dir, hash+".spec")
		metaPath := filepath.Join(c.dir, hash+".meta")

		meta, readErr := readMeta(metaPath)
		if readErr == nil && !meta.IsExpired() {
			if _, statErr := os.Stat(specPath); statErr == nil {
				return nil
			}
		}
	}

	const fallbackTimeout = 10 * time.Second

	// HEAD request as fallback
	headCtx, cancel := context.WithTimeout(
		ctx,
		fallbackTimeout,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(headCtx, http.MethodHead, url, nil)
	if err != nil {
		return &LocationError{Location: url, Type: "url", Err: fmt.Errorf("create request: %w", err)}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &LocationError{Location: url, Type: "url", Err: fmt.Errorf("unreachable: %w", err)}
	}
	if closeErr := resp.Body.Close(); closeErr != nil {
		slog.Default().DebugContext(headCtx, "closing response body", "error", closeErr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &LocationError{Location: url, Type: "url", Err: fmt.Errorf("unexpected status %d", resp.StatusCode)}
	}

	return nil
}
