package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildConfigYAML(t *testing.T) {
	specs := []SpecInput{
		{
			Domain:      "petstore",
			LLMTitle:    "Petstore API",
			Instruction: "Use this API to manage pets.",
			BaseURL:     "https://petstore.swagger.io/v2",
			Collections: []CollectionInput{
				{Title: "Petstore Swagger", Location: "https://petstore.swagger.io/v2/swagger.json"},
			},
		},
	}

	data, err := BuildConfigYAML(specs)
	if err != nil {
		t.Fatalf("BuildConfigYAML() error = %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "domain: \"petstore\"") {
		t.Error("missing domain")
	}
	if !strings.Contains(content, "llm_title: \"Petstore API\"") {
		t.Error("missing llm_title")
	}
	if !strings.Contains(content, "base_url: \"https://petstore.swagger.io/v2\"") {
		t.Error("missing base_url")
	}
	if !strings.Contains(content, "location: \"https://petstore.swagger.io/v2/swagger.json\"") {
		t.Error("missing collection location")
	}
}

func TestBuildConfigYAML_WithAuth(t *testing.T) {
	specs := []SpecInput{
		{
			Domain:   "secured-api",
			LLMTitle: "Secured API",
			BaseURL:  "https://api.example.com/v1",
			AuthType: "bearer",
			AuthConfig: map[string]string{
				"token": "my-secret-token",
			},
			Collections: []CollectionInput{
				{Title: "Main", Location: "./specs/main.json"},
			},
		},
	}

	data, err := BuildConfigYAML(specs)
	if err != nil {
		t.Fatalf("BuildConfigYAML() error = %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "auth:") {
		t.Error("missing auth block")
	}
	if !strings.Contains(content, "type: \"bearer\"") {
		t.Error("missing auth type")
	}
	if !strings.Contains(content, "token: \"my-secret-token\"") {
		t.Error("missing auth config")
	}
}

func TestBuildConfigYAML_NoAuth(t *testing.T) {
	specs := []SpecInput{
		{
			Domain:   "no-auth-api",
			LLMTitle: "No Auth API",
			BaseURL:  "https://api.example.com/v1",
			AuthType: "none",
		},
	}

	data, err := BuildConfigYAML(specs)
	if err != nil {
		t.Fatalf("BuildConfigYAML() error = %v", err)
	}

	content := string(data)
	// Check that no uncommented auth block appears for "none" type.
	// The template may contain commented example blocks with "auth:".
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.Contains(trimmed, "auth:") {
			t.Errorf("unexpected auth block for none type: %s", line)
		}
	}
}

func TestBuildConfigYAML_EmptySpecs(t *testing.T) {
	data, err := BuildConfigYAML(nil)
	if err != nil {
		t.Fatalf("BuildConfigYAML() error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "specs:") {
		t.Error("missing specs key")
	}
}

func TestBuildConfigYAML_NoWorkspaceDir(t *testing.T) {
	data, err := BuildConfigYAML(nil)
	if err != nil {
		t.Fatalf("BuildConfigYAML() error = %v", err)
	}
	content := string(data)
	if strings.Contains(content, "workspace_dir") {
		t.Errorf("unexpected workspace_dir in output:\n%s", content)
	}
}

func TestWriteResult(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	specs := []SpecInput{
		{
			Domain:   "test",
			LLMTitle: "Test API",
			BaseURL:  "https://test.example.com/v1",
			Collections: []CollectionInput{
				{Title: "Test Collection", Location: "./specs/test.json"},
			},
		},
	}

	if err := WriteResult(cfgPath, wsDir, specs); err != nil {
		t.Fatalf("WriteResult() error = %v", err)
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	for _, dir := range []string{wsDir, filepath.Join(wsDir, "cache"), filepath.Join(wsDir, "specs"), filepath.Join(wsDir, "responses")} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("workspace subdirectory %q was not created", dir)
		}
	}
}

func TestAuthMethodsList(t *testing.T) {
	list := authMethodsList()
	for _, m := range availableAuthMethods {
		if !strings.Contains(list, m.Type) {
			t.Errorf("auth method %s not found in list", m.Type)
		}
	}
}

func TestAuthFieldsFor(t *testing.T) {
	fields := authFieldsFor("basic")
	if len(fields) != 2 {
		t.Errorf("basic auth should have 2 fields, got %d", len(fields))
	}
	fields = authFieldsFor("none")
	if len(fields) != 0 {
		t.Errorf("none auth should have 0 fields, got %d", len(fields))
	}
	fields = authFieldsFor("unknown")
	if len(fields) != 0 {
		t.Errorf("unknown auth should have 0 fields, got %d", len(fields))
	}
}
