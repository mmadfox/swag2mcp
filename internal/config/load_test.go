package config

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_FromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`specs:
  - domain: test-api
    llm_title: Test API v1
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
`)
	require.NoError(t, os.WriteFile(path, content, 0600))

	cfg, err := Load(path)
	require.NoError(t, err)
	require.Len(t, cfg.Specs, 1)
	assert.Equal(t, "test-api", cfg.Specs[0].Domain)
}

func TestLoad_FromTildePath(t *testing.T) {
	// Note: no t.Parallel() because t.Setenv is used

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	dir := filepath.Join(tmpHome, ".swag2mcp-test")
	require.NoError(t, os.MkdirAll(dir, 0750))

	path := filepath.Join(dir, "config.yaml")
	content := []byte(`specs:
  - domain: test-api
    llm_title: Test API v1
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
`)
	require.NoError(t, os.WriteFile(path, content, 0600))

	tildePath := "~/.swag2mcp-test/config.yaml"
	cfg, err := Load(tildePath)
	require.NoError(t, err)
	require.Len(t, cfg.Specs, 1)
}

func TestLoad_FromAbsolutePath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`specs:
  - domain: test-api
    llm_title: Test API v1
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
`)
	require.NoError(t, os.WriteFile(path, content, 0600))

	absPath, err := filepath.Abs(path)
	require.NoError(t, err)

	cfg, err := Load(absPath)
	require.NoError(t, err)
	require.Len(t, cfg.Specs, 1)
}

func TestLoad_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := Load("/nonexistent/path/config.yaml")
	require.Error(t, err, "expected error for nonexistent file")
}

func TestLoad_InvalidYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte("invalid: [yaml: broken"), 0600))

	_, err := Load(path)
	require.Error(t, err, "expected error for invalid YAML")
}

func TestLoad_FromFileURL(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`specs:
  - domain: test-api
    llm_title: Test API v1
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
`)
	require.NoError(t, os.WriteFile(path, content, 0600))

	fileURL := "file://" + path
	cfg, err := Load(fileURL)
	require.NoError(t, err)
	require.Len(t, cfg.Specs, 1)
}

func TestLoad_InvalidFileURL(t *testing.T) {
	t.Parallel()

	_, err := Load("ftp://example.com/config.yaml")
	require.Error(t, err, "expected error for non-file URL scheme")
}

func TestExpandTilde(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		input    string
		expected string
	}{
		{"~/config.yaml", filepath.Join(home, "config.yaml")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := env.ExpandTilde(tt.input)
		assert.Equal(t, tt.expected, result, "ExpandTilde(%q)", tt.input)
	}
}

func TestLoad_FromHTTPURL_Success(t *testing.T) {
	t.Parallel()

	yamlContent := "specs:\n  - domain: test-api\n    llm_title: Test API v1\n    base_url: https://api.example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(yamlContent))
	}))
	t.Cleanup(srv.Close)

	cfg, err := Load(srv.URL)
	require.NoError(t, err)
	require.Len(t, cfg.Specs, 1)
	assert.Equal(t, "test-api", cfg.Specs[0].Domain)
}

func TestLoad_FromHTTPURL_NotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	_, err := Load(srv.URL)
	require.Error(t, err, "expected error for 404 response")
}

func TestLoad_FromHTTPURL_InvalidYAML(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid: [yaml: broken"))
	}))
	t.Cleanup(srv.Close)

	_, err := Load(srv.URL)
	require.Error(t, err, "expected error for invalid YAML from HTTP")
}

func TestLoad_FromHTTPURL_Unreachable(t *testing.T) {
	t.Parallel()

	_, err := Load("http://127.0.0.1:1/config.yaml")
	require.Error(t, err, "expected error for unreachable URL")
}

func TestLoad_FromHTTPURL_InvalidURL(t *testing.T) {
	t.Parallel()

	_, err := Load("http://invalid url with spaces")
	require.Error(t, err, "expected error for invalid URL")
}

func TestLoad_FromFileURL_InvalidScheme(t *testing.T) {
	t.Parallel()

	_, err := Load("ftp://example.com/config.yaml")
	require.Error(t, err, "expected error for non-file URL scheme")
}

func TestLoad_FromFileURL_InvalidPath(t *testing.T) {
	t.Parallel()

	_, err := Load("file://%ZZinvalid")
	require.Error(t, err, "expected error for invalid percent-encoded path")
}

func TestLoad_FromAbsolutePath_RelativePath(t *testing.T) {
	t.Parallel()

	_, err := loadFromAbsolutePath("relative/path")
	require.Error(t, err, "expected error for relative path")
}

func TestExpandTilde_NoMatch(t *testing.T) {
	t.Parallel()

	result := env.ExpandTilde("no-tilde")
	assert.Equal(t, "no-tilde", result)
}

func TestExpandTilde_WindowsBackslash(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	result := env.ExpandTilde("~\\config.yaml")
	expected := filepath.Join(home, "config.yaml")
	assert.Equal(t, expected, result)
}
