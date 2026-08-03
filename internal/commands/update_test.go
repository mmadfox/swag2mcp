/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRunUpdate_NoConfigCreatesIt(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := runUpdate(tmpDir)
	if err != nil {
		t.Fatalf("runUpdate() = %v", err)
	}
	ws, _ := workspace.New(tmpDir)
	if !ws.ConfigExists() {
		t.Fatal("config should have been created")
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

func TestRunUpdate_Output(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specFile := filepath.Join(tmpDir, "spec.json")
	if err := os.WriteFile(specFile, []byte(`{"openapi":"3.0.0"}`), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:   "nspd",
				LLMTitle: "NSPD Services",
				BaseURL:  "https://example.com",
				Collections: []config.Collection{
					{LLMTitle: "Actors", Location: specFile},
					{LLMTitle: "Core", Location: specFile},
				},
			},
			{
				Domain:   "meteo",
				LLMTitle: "Meteo API",
				BaseURL:  "https://meteo.example.com",
				Collections: []config.Collection{
					{LLMTitle: "Forecast", Location: specFile},
				},
			},
		},
	}
	data, _ := yaml.Marshal(cfg)
	if err := os.WriteFile(ws.ConfigPath(), data, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	result, err := runUpdate(tmpDir)
	if err != nil {
		t.Fatalf("runUpdate() = %v", err)
	}
	require.Equal(t, 2, result.total)
	require.Len(t, result.specs, 2)
	require.Equal(t, "nspd", result.specs[0].domain)
	require.Equal(t, 2, result.specs[0].collections)
	require.Equal(t, "meteo", result.specs[1].domain)
	require.Equal(t, 1, result.specs[1].collections)
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
	remote, local, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		t.Fatalf("cacheSpecs() = %v", err)
	}
	if remote+local != 1 {
		t.Errorf("total = %d, want 1", remote+local)
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
	remote, local, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		t.Fatalf("cacheSpecs() = %v", err)
	}
	if remote+local != 0 {
		t.Errorf("total = %d, want 0", remote+local)
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
	remote, local, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		t.Fatalf("cacheSpecs() = %v", err)
	}
	if remote+local != 1 {
		t.Errorf("total = %d, want 1", remote+local)
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
