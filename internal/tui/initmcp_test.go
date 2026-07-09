package tui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	if err := Setup(cfgPath, wsDir); err != nil {
		t.Fatalf("Setup() = %v", err)
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("config file not created")
	}
	if _, err := os.Stat(wsDir); os.IsNotExist(err) {
		t.Error("workspace dir not created")
	}
}

func TestSetup_InvalidWorkspace(t *testing.T) {
	t.Parallel()

	err := Setup("/nonexistent/config.yaml", "/invalid:\x00path")
	if err == nil {
		t.Fatal("expected error for invalid workspace path")
	}
}

func TestWriteConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	data := []byte("specs:\n  - domain: test\n    llm_title: Test API\n    base_url: https://example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")

	if err := WriteConfig(cfgPath, data); err != nil {
		t.Fatalf("WriteConfig() = %v", err)
	}

	read, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(read) != string(data) {
		t.Errorf("written data differs")
	}
}

func TestWriteConfig_CreatesDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "subdir", "swag2mcp.yaml")
	data := []byte("specs: []")

	if err := WriteConfig(cfgPath, data); err != nil {
		t.Fatalf("WriteConfig() = %v", err)
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("config file not created in new subdirectory")
	}
}

func TestExampleConfig(t *testing.T) {
	t.Parallel()

	data := ExampleConfig()
	if len(data) == 0 {
		t.Fatal("ExampleConfig() returned empty data")
	}
}

func TestSetup_WriteConfigError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(cfgDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	if err := os.Chmod(cfgDir, 0000); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(cfgDir, 0750) })

	cfgPath := filepath.Join(cfgDir, "swag2mcp.yaml")
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	if err := Setup(cfgPath, wsDir); err == nil {
		t.Error("Setup() expected error for read-only config dir, got nil")
	}
}

func TestSetup_WorkspaceInitError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	if err := os.MkdirAll(wsDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	blocker := filepath.Join(wsDir, "cache")
	if err := os.WriteFile(blocker, []byte("block"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := Setup(cfgPath, wsDir); err == nil {
		t.Error("Setup() expected error for workspace init failure, got nil")
	}
}

func TestWriteConfig_MkdirAllError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	blocker := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(blocker, []byte("block"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	cfgPath := filepath.Join(blocker, "swag2mcp.yaml")

	if err := WriteConfig(cfgPath, []byte("specs: []")); err == nil {
		t.Error("WriteConfig() expected error for blocked dir, got nil")
	}
}
