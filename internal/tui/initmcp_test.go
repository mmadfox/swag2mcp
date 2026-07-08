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
