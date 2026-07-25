package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
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
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(cfg.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(cfg.Specs))
	}
	if cfg.Specs[0].Domain != "test-api" {
		t.Errorf("Domain = %q, want %q", cfg.Specs[0].Domain, "test-api")
	}
}

func TestLoad_FromTildePath(t *testing.T) {
	// Note: no t.Parallel() because t.Setenv is used

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	dir := filepath.Join(tmpHome, ".swag2mcp-test")
	if mkErr := os.MkdirAll(dir, 0750); mkErr != nil {
		t.Fatalf("MkdirAll() = %v", mkErr)
	}

	path := filepath.Join(dir, "config.yaml")
	content := []byte(`specs:
  - domain: test-api
    llm_title: Test API v1
    base_url: https://api.example.com
    collections:
      - llm_title: Main
        location: https://example.com/spec.yaml
`)
	if wrErr := os.WriteFile(path, content, 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}

	tildePath := "~/.swag2mcp-test/config.yaml"
	cfg, err := Load(tildePath)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(cfg.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(cfg.Specs))
	}
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
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("Abs() = %v", err)
	}

	cfg, err := Load(absPath)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(cfg.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(cfg.Specs))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("invalid: [yaml: broken"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
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
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	fileURL := "file://" + path
	cfg, err := Load(fileURL)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(cfg.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(cfg.Specs))
	}
}

func TestLoad_InvalidFileURL(t *testing.T) {
	t.Parallel()

	_, err := Load("ftp://example.com/config.yaml")
	if err == nil {
		t.Fatal("expected error for non-file URL scheme")
	}
}

func TestExpandTilde(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() = %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/config.yaml", filepath.Join(home, "config.yaml")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := expandTilde(tt.input)
		if result != tt.expected {
			t.Errorf("expandTilde(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if len(cfg.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(cfg.Specs))
	}
	if cfg.Specs[0].Domain != "test-api" {
		t.Errorf("Domain = %q, want %q", cfg.Specs[0].Domain, "test-api")
	}
}

func TestLoad_FromHTTPURL_NotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	_, err := Load(srv.URL)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestLoad_FromHTTPURL_InvalidYAML(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid: [yaml: broken"))
	}))
	t.Cleanup(srv.Close)

	_, err := Load(srv.URL)
	if err == nil {
		t.Fatal("expected error for invalid YAML from HTTP")
	}
}

func TestLoad_FromHTTPURL_Unreachable(t *testing.T) {
	t.Parallel()

	_, err := Load("http://127.0.0.1:1/config.yaml")
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestLoad_FromHTTPURL_InvalidURL(t *testing.T) {
	t.Parallel()

	_, err := Load("http://invalid url with spaces")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestLoad_FromFileURL_InvalidScheme(t *testing.T) {
	t.Parallel()

	_, err := Load("ftp://example.com/config.yaml")
	if err == nil {
		t.Fatal("expected error for non-file URL scheme")
	}
}

func TestLoad_FromFileURL_InvalidPath(t *testing.T) {
	t.Parallel()

	_, err := Load("file://%ZZinvalid")
	if err == nil {
		t.Fatal("expected error for invalid percent-encoded path")
	}
}

func TestLoad_FromAbsolutePath_RelativePath(t *testing.T) {
	t.Parallel()

	_, err := loadFromAbsolutePath("relative/path")
	if err == nil {
		t.Fatal("expected error for relative path")
	}
}

func TestExpandTilde_NoMatch(t *testing.T) {
	t.Parallel()

	result := expandTilde("no-tilde")
	if result != "no-tilde" {
		t.Errorf("got %q, want %q", result, "no-tilde")
	}
}

func TestExpandTilde_WindowsBackslash(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() = %v", err)
	}

	result := expandTilde("~\\config.yaml")
	expected := filepath.Join(home, "config.yaml")
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}
