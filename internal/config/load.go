package config

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.yaml.in/yaml/v3"
)

func Load(confFilepath string) (*Config, error) {
	isURL := strings.HasPrefix(confFilepath, "https://") || strings.HasPrefix(confFilepath, "http://")
	isFileURL := strings.HasPrefix(confFilepath, "file://")
	isPath := strings.HasPrefix(confFilepath, "~")

	switch {
	case isURL:
		return loadFromHTTPURL(confFilepath)
	case isFileURL:
		return loadFromFileURL(confFilepath)
	case isPath:
		return loadPath(confFilepath)
	default:
		if filepath.IsAbs(confFilepath) {
			return loadFromAbsolutePath(confFilepath)
		}
		return loadPath(confFilepath)
	}
}

func loadPath(filepathSpec string) (*Config, error) {
	filepathSpec = expandTilde(filepathSpec)

	data, readErr := os.ReadFile(filepathSpec)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", filepathSpec, readErr)
	}

	var cfg Config
	if parseErr := yaml.Unmarshal(data, &cfg); parseErr != nil {
		return nil, fmt.Errorf("failed to parse config file '%s': %w", filepathSpec, parseErr)
	}

	return &cfg, nil
}

func loadFromHTTPURL(urlStr string) (*Config, error) {
	if _, parseErr := url.Parse(urlStr); parseErr != nil {
		return nil, fmt.Errorf("invalid URL '%s': %w", urlStr, parseErr)
	}
	req, reqErr := http.NewRequestWithContext(context.Background(), http.MethodGet, urlStr, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("failed to create request: %w", reqErr)
	}
	resp, getErr := http.DefaultClient.Do(req)
	if getErr != nil {
		return nil, fmt.Errorf("failed to GET from URL '%s': %w", urlStr, getErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %d for URL '%s'", resp.StatusCode, urlStr)
	}

	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", readErr)
	}

	var cfg Config
	if parseErr := yaml.Unmarshal(data, &cfg); parseErr != nil {
		return nil, fmt.Errorf("failed to parse YAML from URL '%s': %w", urlStr, parseErr)
	}

	return &cfg, nil
}

func loadFromFileURL(rawURL string) (*Config, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid file URL: %w", err)
	}

	if u.Scheme != "file" {
		return nil, fmt.Errorf("expected 'file' scheme, got '%s'", u.Scheme)
	}

	path, err := url.PathUnescape(u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape URL path: %w", err)
	}

	// file:///C:/Users/... → /C:/Users/... on Windows; strip the leading /
	if runtime.GOOS == "windows" && len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	return loadFromAbsolutePath(filepath.FromSlash(path))
}

func loadFromAbsolutePath(path string) (*Config, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("path must be absolute: %s", path)
	}

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read config file at absolute path '%s': %w", path, readErr)
	}

	var cfg Config
	if parseErr := yaml.Unmarshal(data, &cfg); parseErr != nil {
		return nil, fmt.Errorf("failed to parse config file at absolute path '%s': %w", path, parseErr)
	}

	return &cfg, nil
}

// Save writes the config to the given file path as YAML.
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if writeErr := os.WriteFile(path, data, 0600); writeErr != nil {
		return fmt.Errorf("failed to write config to %q: %w", path, writeErr)
	}
	return nil
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
