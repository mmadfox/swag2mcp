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

	"github.com/mmadfox/swag2mcp/internal/env"
	"go.yaml.in/yaml/v3"
)

// Load reads a YAML configuration from the given file path, HTTP URL, or file URL.
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

// loadPath reads a config file from a local path, expanding ~ if present.
func loadPath(filepathSpec string) (*Config, error) {
	filepathSpec = env.ExpandTilde(filepathSpec)

	data, err := os.ReadFile(filepathSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", filepathSpec, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file '%s': %w", filepathSpec, err)
	}

	return &cfg, nil
}

// loadFromHTTPURL fetches and parses a config from an HTTP(S) URL.
func loadFromHTTPURL(urlStr string) (*Config, error) {
	if _, err := url.Parse(urlStr); err != nil {
		return nil, fmt.Errorf("invalid URL '%s': %w", urlStr, err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to GET from URL '%s': %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status %d for URL '%s'", resp.StatusCode, urlStr)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from URL '%s': %w", urlStr, err)
	}

	return &cfg, nil
}

// loadFromFileURL parses a file:// URL and loads the config from the local path.
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

// loadFromAbsolutePath reads a config file from an absolute local path.
func loadFromAbsolutePath(path string) (*Config, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("path must be absolute: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at absolute path '%s': %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file at absolute path '%s': %w", path, err)
	}

	return &cfg, nil
}

// Save writes the config to the given file path as YAML.
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config to %q: %w", path, err)
	}
	return nil
}
