package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
)

func TestAtomicWriteConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test\n    llm_title: Test API v1\n    base_url: https://example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	err := AtomicWriteConfig(cfgPath, func(cfg *config.Config) error {
		cfg.Specs[0].Domain = "updated"
		return nil
	})
	if err != nil {
		t.Fatalf("AtomicWriteConfig() = %v", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() = %v", err)
	}
	if cfg.Specs[0].Domain != "updated" {
		t.Errorf("Domain = %q, want %q", cfg.Specs[0].Domain, "updated")
	}
}

func TestAtomicWriteConfig_LoadError(t *testing.T) {
	t.Parallel()

	err := AtomicWriteConfig("/nonexistent/config.yaml", func(_ *config.Config) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error for nonexistent config")
	}
}

func TestAtomicWriteConfig_ValidationError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test\n    llm_title: Test API v1\n    base_url: https://example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	err := AtomicWriteConfig(cfgPath, func(cfg *config.Config) error {
		cfg.Specs = nil
		return nil
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
