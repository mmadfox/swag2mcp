/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
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

func TestAddSpecFromYAML_DuplicateDomain(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: dup\n    llm_title: Original API\n    base_url: https://original.example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("domain: dup\nllm_title: Duplicate API\nbase_url: https://dup.example.com\ncollections:\n  - llm_title: Coll\n    location: https://dup.example.com/spec.yaml\n")

	err := AddSpecFromYAML(cfgPath, yamlData)
	if err == nil {
		t.Fatal("expected error for duplicate domain")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestAddCollectionFromYAML_DuplicateLocation(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test-api\n    llm_title: Test API\n    base_url: https://example.com\n    collections:\n      - llm_title: Existing\n        location: https://example.com/existing.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	yamlData := []byte("spec_domain: test-api\nllm_title: Duplicate Collection\nlocation: https://example.com/existing.yaml\n")

	err := AddCollectionFromYAML(cfgPath, yamlData)
	if err == nil {
		t.Fatal("expected error for duplicate location")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
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

func TestAddSpecTUI_DuplicateDomain(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: dup\n    llm_title: Original API\n    base_url: https://original.example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	// AtomicWriteConfig callback should detect duplicate domain
	err := AtomicWriteConfig(cfgPath, func(cfg *config.Config) error {
		domainMap := make(map[string]struct{}, len(cfg.Specs))
		for _, sp := range cfg.Specs {
			domainMap[sp.Domain] = struct{}{}
		}
		if _, ok := domainMap["dup"]; ok {
			return fmt.Errorf("spec with domain %q already exists", "dup")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected error for duplicate domain")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestAddCollectionTUI_DuplicateLocation(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: test-api\n    llm_title: Test API\n    base_url: https://example.com\n    collections:\n      - llm_title: Existing\n        location: https://example.com/existing.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	// AtomicWriteConfig callback should detect duplicate location via map
	err := AtomicWriteConfig(cfgPath, func(cfg *config.Config) error {
		specMap := make(map[string]int, len(cfg.Specs))
		for i, sp := range cfg.Specs {
			specMap[sp.Domain] = i
		}
		idx, ok := specMap["test-api"]
		if !ok {
			return fmt.Errorf("spec with domain %q not found", "test-api")
		}
		locMap := make(map[string]struct{}, len(cfg.Specs[idx].Collections))
		for _, c := range cfg.Specs[idx].Collections {
			locMap[c.Location] = struct{}{}
		}
		if _, ok := locMap["https://example.com/existing.yaml"]; ok {
			return fmt.Errorf("collection with location %q already exists in spec %q", "https://example.com/existing.yaml", "test-api")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected error for duplicate location")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestAddCollectionTUI_SpecNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	initialData := []byte("specs:\n  - domain: existing\n    llm_title: Existing API\n    base_url: https://example.com\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, initialData, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	err := AtomicWriteConfig(cfgPath, func(cfg *config.Config) error {
		specMap := make(map[string]int, len(cfg.Specs))
		for i, sp := range cfg.Specs {
			specMap[sp.Domain] = i
		}
		if _, ok := specMap["nonexistent"]; !ok {
			return fmt.Errorf("spec with domain %q not found", "nonexistent")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected error for nonexistent spec")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}
