package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddSpecFromYAML(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: existing\n    llm_title: Existing API\n    base_url: https://existing.example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("domain: new-api\nllm_title: New API\nbase_url: https://new.example.com\ncollections:\n  - llm_title: New Coll\n    location: https://new.example.com/spec.yaml\n")

	if err := AddSpecFromYAML(cfgPath, yamlData); err != nil {
		t.Fatalf("AddSpecFromYAML() = %v", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "new-api") {
		t.Error("missing new spec domain")
	}
	if !strings.Contains(content, "existing") {
		t.Error("missing existing spec")
	}
}

func TestAddSpecFromYAML_MissingDomain(t *testing.T) {
	t.Parallel()

	err := AddSpecFromYAML("/nonexistent/config.yaml", []byte("llm_title: Test\nbase_url: https://example.com\n"))
	if err == nil {
		t.Fatal("expected error for missing domain")
	}
}

func TestAddSpecFromYAML_MissingTitle(t *testing.T) {
	t.Parallel()

	err := AddSpecFromYAML("/nonexistent/config.yaml", []byte("domain: test\nbase_url: https://example.com\n"))
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestAddSpecFromYAML_MissingBaseURL(t *testing.T) {
	t.Parallel()

	err := AddSpecFromYAML("/nonexistent/config.yaml", []byte("domain: test\nllm_title: Test\n"))
	if err == nil {
		t.Fatal("expected error for missing base_url")
	}
}

func TestAddSpecFromYAML_InvalidYAML(t *testing.T) {
	t.Parallel()

	err := AddSpecFromYAML("/nonexistent/config.yaml", []byte("invalid: [yaml"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestAddCollectionFromYAML(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test-api\n    llm_title: Test API\n    base_url: https://example.com\n    collections:\n      - llm_title: Existing\n        location: https://example.com/existing.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("spec_domain: test-api\nllm_title: New Collection\nlocation: https://example.com/new.yaml\n")

	if err := AddCollectionFromYAML(cfgPath, yamlData); err != nil {
		t.Fatalf("AddCollectionFromYAML() = %v", err)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "New Collection") {
		t.Error("missing new collection")
	}
}

func TestAddCollectionFromYAML_MissingSpecDomain(t *testing.T) {
	t.Parallel()

	err := AddCollectionFromYAML("/nonexistent/config.yaml", []byte("llm_title: Test\nlocation: https://example.com\n"))
	if err == nil {
		t.Fatal("expected error for missing spec_domain")
	}
}

func TestAddCollectionFromYAML_MissingTitle(t *testing.T) {
	t.Parallel()

	err := AddCollectionFromYAML("/nonexistent/config.yaml", []byte("spec_domain: test\nlocation: https://example.com\n"))
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestAddCollectionFromYAML_MissingLocation(t *testing.T) {
	t.Parallel()

	err := AddCollectionFromYAML("/nonexistent/config.yaml", []byte("spec_domain: test\nllm_title: Test\n"))
	if err == nil {
		t.Fatal("expected error for missing location")
	}
}

func TestAddCollectionFromYAML_SpecNotFound(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: existing\n    llm_title: Existing API\n    base_url: https://example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("spec_domain: nonexistent\nllm_title: New\nlocation: https://example.com/new.yaml\n")

	err := AddCollectionFromYAML(cfgPath, yamlData)
	if err == nil {
		t.Fatal("expected error for nonexistent spec domain")
	}
}

func TestResolveConfigPath_Default(t *testing.T) {
	t.Parallel()

	path := resolveConfigPath("")
	if path == "" {
		t.Fatal("resolveConfigPath('') returned empty")
	}
}

func TestResolveConfigPath_Directory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	path := resolveConfigPath(tmpDir)
	if path != filepath.Join(tmpDir, "swag2mcp.yaml") {
		t.Errorf("got %q, want %q", path, filepath.Join(tmpDir, "swag2mcp.yaml"))
	}
}

func TestAddSpecFromYAML_AtomicWriteError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test\n    llm_title: Test API v1\n    base_url: https://example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("domain: new-api\nllm_title: New API\nbase_url: https://new.example.com\ncollections:\n  - llm_title: New Coll\n    location: https://new.example.com/spec.yaml\n")

	if err := os.Chmod(tmpDir, 0000); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(tmpDir, 0750) })

	if err := AddSpecFromYAML(cfgPath, yamlData); err == nil {
		t.Error("AddSpecFromYAML() expected error for read-only dir, got nil")
	}
}

func TestAddCollectionFromYAML_AtomicWriteError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test-api\n    llm_title: Test API\n    base_url: https://example.com\n    collections:\n      - llm_title: Existing\n        location: https://example.com/existing.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("spec_domain: test-api\nllm_title: New Collection\nlocation: https://example.com/new.yaml\n")

	if err := os.Chmod(tmpDir, 0000); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(tmpDir, 0750) })

	if err := AddCollectionFromYAML(cfgPath, yamlData); err == nil {
		t.Error("AddCollectionFromYAML() expected error for read-only dir, got nil")
	}
}
