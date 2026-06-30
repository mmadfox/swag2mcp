package cache

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	CacheDirName = "cache"
	MaxTTL       = 48 * time.Hour
	MinTTL       = 1 * time.Hour
)

// Cache resolves spec locations to local file paths, caching remote URLs on disk.
type Cache struct {
	dir string
	cli *httpClient
}

// New creates a cache rooted at workspaceDir/cache.
// The cache directory is created lazily on the first Resolve call.
func New(workspaceDir string) *Cache {
	return &Cache{
		dir: filepath.Join(workspaceDir, CacheDirName),
		cli: defaultHTTPClient(),
	}
}

// SetWorkspaceDir sets the root directory for the cache.
func (c *Cache) SetWorkspaceDir(dir string) {
	c.dir = filepath.Join(dir, CacheDirName)
}

// Resolve takes a location (local path, file:// URL, or http(s):// URL)
// and returns a path to a local file containing the spec data.
//
// For local paths the input is returned unchanged.
// For remote URLs the file is downloaded (or served from cache if TTL is still valid).
func (c *Cache) Resolve(location string) (string, error) {
	if location == "" {
		return "", fmt.Errorf("empty location")
	}

	isURL := strings.HasPrefix(location, "https://") || strings.HasPrefix(location, "http://")
	isFileURL := strings.HasPrefix(location, "file://")

	switch {
	case isFileURL:
		return resolveFileURL(location)
	case isURL:
		return c.resolveURL(location)
	default:
		// Local path — expand ~ and return as-is
		return resolveLocalPath(location)
	}
}

func resolveFileURL(rawURL string) (string, error) {
	path, err := fileURIToPath(rawURL)
	if err != nil {
		return "", fmt.Errorf("file URL: %w", err)
	}
	return path, nil
}

func resolveLocalPath(location string) (string, error) {
	location = expandTilde(location)
	absPath, err := filepath.Abs(location)
	if err != nil {
		return "", fmt.Errorf("convert to absolute path: %w", err)
	}
	return absPath, nil
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

func (c *Cache) resolveURL(url string) (string, error) {
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}

	hash := cacheKey(url)
	specPath := filepath.Join(c.dir, hash+".spec")
	metaPath := filepath.Join(c.dir, hash+".meta")

	// Check existing cache
	meta, err := readMeta(metaPath)
	if err == nil && !meta.IsExpired() {
		_, err := os.Stat(specPath)
		if err == nil {
			return specPath, nil
		}
	}

	// Download
	data, err := c.cli.Get(url)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", url, err)
	}

	if err := os.WriteFile(specPath, data, 0644); err != nil {
		return "", fmt.Errorf("write cache file: %w", err)
	}

	ttl := randomTTL()
	meta = fileMeta{
		URL:      url,
		CachedAt: time.Now(),
		TTLSec:   int(ttl.Seconds()),
	}
	if err := writeMeta(metaPath, meta); err != nil {
		return "", fmt.Errorf("write meta file: %w", err)
	}

	return specPath, nil
}

func cacheKey(rawURL string) string {
	h := sha256.Sum256([]byte(rawURL))
	return fmt.Sprintf("%x", h[:16])
}

func randomTTL() time.Duration {
	n := rand.Int63n(int64(MaxTTL - MinTTL))
	return MinTTL + time.Duration(n)
}
