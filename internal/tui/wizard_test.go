package tui

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildConfigYAML(t *testing.T) {
	specs := []SpecInput{
		{
			Domain:      "meteo",
			LLMTitle:    "Open-Meteo Weather API",
			Instruction: "Use the Open-Meteo API to get weather forecasts and related data.",
			BaseURL:     "https://api.open-meteo.com",
			Collections: []CollectionInput{
				{Title: "Forecast", Location: "https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml"},
			},
		},
	}

	data, err := BuildConfigYAML(specs)
	if err != nil {
		t.Fatalf("BuildConfigYAML() error = %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "domain: \"meteo\"") {
		t.Error("missing domain")
	}
	if !strings.Contains(content, "llm_title: \"Open-Meteo Weather API\"") {
		t.Error("missing llm_title")
	}
	if !strings.Contains(content, "base_url: \"https://api.open-meteo.com\"") {
		t.Error("missing base_url")
	}
	if !strings.Contains(content, "location: \"https://raw.githubusercontent.com/mmadfox/swag2mcp/main/specs/meteo/forecast.yml\"") {
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

func TestWriteResult_ConfigPathIsDir(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	specs := []SpecInput{
		{
			Domain:   "test",
			LLMTitle: "Test API",
			BaseURL:  "https://test.example.com/v1",
		},
	}

	if err := WriteResult(tmpDir, wsDir, specs); err != nil {
		t.Fatalf("WriteResult() error = %v", err)
	}

	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Error("config file was not created inside directory")
	}
}

func TestWriteResult_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	blocker := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(blocker, []byte("block"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	cfgPath := filepath.Join(blocker, "swag2mcp.yaml")

	specs := []SpecInput{
		{
			Domain:   "test",
			LLMTitle: "Test API",
			BaseURL:  "https://test.example.com/v1",
		},
	}

	if err := WriteResult(cfgPath, wsDir, specs); err == nil {
		t.Error("WriteResult() expected error, got nil")
	}
}

func TestWriteResult_WriteFileError(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, ".swag2mcp")

	cfgDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(cfgDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	if err := os.Chmod(cfgDir, 0000); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(cfgDir, 0750) })

	cfgPath := filepath.Join(cfgDir, "swag2mcp.yaml")

	specs := []SpecInput{
		{
			Domain:   "test",
			LLMTitle: "Test API",
			BaseURL:  "https://test.example.com/v1",
		},
	}

	if err := WriteResult(cfgPath, wsDir, specs); err == nil {
		t.Error("WriteResult() expected error, got nil")
	}
}

func TestWriteResult_WorkspaceNewError(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")

	specs := []SpecInput{
		{
			Domain:   "test",
			LLMTitle: "Test API",
			BaseURL:  "https://test.example.com/v1",
		},
	}

	if err := WriteResult(cfgPath, "/invalid:\x00path", specs); err == nil {
		t.Error("WriteResult() expected error for invalid workspace path, got nil")
	}
}

func TestWriteResult_WorkspaceInitError(t *testing.T) {
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

	specs := []SpecInput{
		{
			Domain:   "test",
			LLMTitle: "Test API",
			BaseURL:  "https://test.example.com/v1",
		},
	}

	if err := WriteResult(cfgPath, wsDir, specs); err == nil {
		t.Error("WriteResult() expected error for workspace init failure, got nil")
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

func TestHeaderLine_Empty(t *testing.T) {
	if got := headerLine(""); got != "" {
		t.Errorf("headerLine('') = %q, want ''", got)
	}
}

func TestHeaderLine_Ascii(t *testing.T) {
	got := headerLine("Hello")
	want := "─────"
	if got != want {
		t.Errorf("headerLine('Hello') = %q, want %q", got, want)
	}
}

func TestHeaderLine_Unicode(t *testing.T) {
	got := headerLine("Привет")
	want := "──────"
	if got != want {
		t.Errorf("headerLine('Привет') = %q, want %q", got, want)
	}
}
