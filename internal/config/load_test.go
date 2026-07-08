package config

import (
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
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() = %v", err)
	}

	dir := filepath.Join(home, ".swag2mcp-test")
	if mkErr := os.MkdirAll(dir, 0750); mkErr != nil {
		t.Fatalf("MkdirAll() = %v", mkErr)
	}
	defer os.RemoveAll(dir)

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
