package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	content := []byte("specs:\n  - domain: test-api\n    llm_title: Test API v1\n    base_url: https://api.example.com\n    tags: [public]\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	output, err := ListConfig(cfgPath, nil)
	if err != nil {
		t.Fatalf("ListConfig() = %v", err)
	}
	if !strings.Contains(output, "test-api") {
		t.Error("output missing domain")
	}
	if !strings.Contains(output, "Test API v1") {
		t.Error("output missing title")
	}
}

func TestListConfig_WithTags(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	content := []byte("specs:\n  - domain: public-api\n    llm_title: Public API\n    base_url: https://api.example.com\n    tags: [public]\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n  - domain: internal-api\n    llm_title: Internal API\n    base_url: https://internal.example.com\n    tags: [internal]\n    collections:\n      - llm_title: Internal\n        location: https://internal.example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	output, err := ListConfig(cfgPath, []string{"public"})
	if err != nil {
		t.Fatalf("ListConfig() = %v", err)
	}
	if !strings.Contains(output, "public-api") {
		t.Error("output missing public-api")
	}
	if strings.Contains(output, "internal-api") {
		t.Error("output should not contain internal-api")
	}
}

func TestListConfig_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := ListConfig("/nonexistent/config.yaml", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestListConfig_DisabledSpec(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	content := []byte("specs:\n  - domain: disabled-api\n    llm_title: Disabled API\n    base_url: https://api.example.com\n    disable: true\n    collections:\n      - llm_title: Main\n        location: https://example.com/spec.yaml\n")
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	output, err := ListConfig(cfgPath, nil)
	if err != nil {
		t.Fatalf("ListConfig() = %v", err)
	}
	if strings.Contains(output, "disabled-api") {
		t.Error("output should not contain disabled spec")
	}
}
