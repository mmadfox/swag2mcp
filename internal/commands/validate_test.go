package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"gopkg.in/yaml.v3"
)

func TestRunValidate_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	var buf strings.Builder
	err := runValidate(tmpDir, "", &buf)
	if err == nil {
		t.Fatal("runValidate() expected error, got nil")
	}
}

func TestRunValidate_ValidConfig(t *testing.T) {
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
	err := runValidate(tmpDir, "", &buf)
	if err != nil {
		t.Fatalf("runValidate() = %v", err)
	}
	if !strings.Contains(buf.String(), "Configuration is valid") {
		t.Errorf("output = %q, want valid message", buf.String())
	}
}

func TestRunValidate_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}
	if err := os.WriteFile(ws.ConfigPath(), []byte("invalid: [yaml"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	var buf strings.Builder
	err := runValidate(tmpDir, "", &buf)
	if err == nil {
		t.Fatal("runValidate() expected error, got nil")
	}
}
