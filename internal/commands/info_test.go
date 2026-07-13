package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"gopkg.in/yaml.v3"
)

func TestRunInfo_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	var buf strings.Builder
	err := runInfo(tmpDir, &buf)
	if err != nil {
		t.Fatalf("runInfo() = %v", err)
	}
	if buf.Len() == 0 {
		t.Error("runInfo() produced no output")
	}
}

func TestRunInfo_WithConfig(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}
	specFile := filepath.Join(tmpDir, "spec.json")
	if err := os.WriteFile(specFile, []byte(`{"openapi":"3.0.0"}`), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	cfg := config.Config{
		Specs: []config.Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API Service",
				BaseURL:  "https://example.com",
				Collections: []config.Collection{
					{LLMTitle: "Main API", Location: specFile},
				},
			},
		},
	}
	data, _ := yaml.Marshal(cfg)
	if err := os.WriteFile(ws.ConfigPath(), data, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	var buf strings.Builder
	err := runInfo(tmpDir, &buf)
	if err != nil {
		t.Fatalf("runInfo() = %v", err)
	}
	if buf.Len() == 0 {
		t.Error("runInfo() produced no output")
	}
}

func TestRunInfo_InvalidPath(t *testing.T) {
	var buf strings.Builder
	err := runInfo("/nonexistent/path", &buf)
	if err == nil {
		t.Fatal("runInfo() expected error, got nil")
	}
}
