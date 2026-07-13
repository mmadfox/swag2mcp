package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestRunUpdate_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := runUpdate(tmpDir)
	if err == nil {
		t.Fatal("runUpdate() expected error for missing config, got nil")
	}
}

func TestRunUpdate_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}
	if err := os.WriteFile(ws.ConfigPath(), []byte("invalid: [yaml"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := runUpdate(tmpDir)
	if err == nil {
		t.Fatal("runUpdate() expected error for invalid config, got nil")
	}
}

func TestCacheSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specDir := filepath.Join(tmpDir, "specs")
	if err := os.MkdirAll(specDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	specFile := filepath.Join(specDir, "test.json")
	if err := os.WriteFile(specFile, []byte(`{"openapi":"3.0.0"}`), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:   "test",
				LLMTitle: "Test",
				BaseURL:  "https://example.com",
				Collections: []config.Collection{
					{LLMTitle: "Main", Location: specFile},
				},
			},
		},
	}

	ca := cache.New(tmpDir)
	total, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		t.Fatalf("cacheSpecs() = %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
}

func TestCacheSpecs_DisabledCollection(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:   "test",
				LLMTitle: "Test",
				BaseURL:  "https://example.com",
				Collections: []config.Collection{
					{LLMTitle: "Disabled", Location: "./nonexistent.json", Disable: true},
				},
			},
		},
	}

	ca := cache.New(tmpDir)
	total, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		t.Fatalf("cacheSpecs() = %v", err)
	}
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
}

func TestCacheSpecs_ScriptAuth(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specFile := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(specFile, []byte(`{"openapi":"3.0.0"}`), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:   "script-api",
				LLMTitle: "Script API",
				BaseURL:  "https://example.com",
				Auth: config.Auth{
					Client: &auth.ScriptAuthClient{Domain: "script-api"},
				},
				Collections: []config.Collection{
					{LLMTitle: "Main", Location: specFile},
				},
			},
		},
	}

	ca := cache.New(tmpDir)
	total, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		t.Fatalf("cacheSpecs() = %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}

	scriptPath := ws.AuthScriptPath("script-api")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Error("auth script was not created for script auth")
	}
}

func TestCleanOrphanAuthScripts(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	orphanPath := filepath.Join(ws.AuthScriptsDir(), "orphan.sh")
	if err := os.WriteFile(orphanPath, []byte("echo test"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cfg := &config.Config{
		Specs: []config.Spec{
			{Domain: "active", LLMTitle: "Active", BaseURL: "https://example.com"},
		},
	}

	if err := cleanOrphanAuthScripts(cfg, ws); err != nil {
		t.Fatalf("cleanOrphanAuthScripts() = %v", err)
	}

	if _, err := os.Stat(orphanPath); !os.IsNotExist(err) {
		t.Error("orphan script was not removed")
	}
}
