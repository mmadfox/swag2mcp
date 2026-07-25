package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"github.com/spf13/cobra"
)

func TestResolveBasePath_Empty(t *testing.T) {
	if got := resolveBasePath(nil); got != "" {
		t.Errorf("resolveBasePath(nil) = %q, want ''", got)
	}
}

func TestResolveBasePath_EmptySlice(t *testing.T) {
	if got := resolveBasePath([]string{}); got != "" {
		t.Errorf("resolveBasePath([]) = %q, want ''", got)
	}
}

func TestResolveBasePath_WithArg(t *testing.T) {
	if got := resolveBasePath([]string{"/tmp"}); got != "/tmp" {
		t.Errorf("resolveBasePath([/tmp]) = %q, want /tmp", got)
	}
}

func TestReadYAMLInput_Stdin(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() = %v", err)
	}
	defer r.Close()

	oldStdin := os.Stdin
	os.Stdin = r                           //nolint:reassign // test helper temporarily replaces os.Stdin to mock user input
	defer func() { os.Stdin = oldStdin }() //nolint:reassign // restore original stdin

	w.WriteString("domain: test\n")
	w.Close()

	data, err := readYAMLInput("-")
	if err != nil {
		t.Fatalf("readYAMLInput('-') = %v", err)
	}
	if string(data) != "domain: test\n" {
		t.Errorf("got %q, want %q", string(data), "domain: test\n")
	}
}

func TestReadYAMLInput_Inline(t *testing.T) {
	data, err := readYAMLInput("domain: test")
	if err != nil {
		t.Fatalf("readYAMLInput() = %v", err)
	}
	if string(data) != "domain: test" {
		t.Errorf("got %q, want %q", string(data), "domain: test")
	}
}

func TestEnsureConfigExists_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.NewFromBase(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}
	cfgPath := ws.ConfigPath()
	if err := os.WriteFile(cfgPath, []byte("specs: []"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	got, err := ensureConfigExists(tmpDir)
	if err != nil {
		t.Fatalf("ensureConfigExists() = %v", err)
	}
	if got == "" {
		t.Error("ensureConfigExists() returned empty path")
	}
	if _, statErr := os.Stat(got); os.IsNotExist(statErr) {
		t.Error("config file does not exist at returned path")
	}
}

func TestEnsureConfigExists_NotExists(t *testing.T) {
	tmpDir := t.TempDir()

	got, err := ensureConfigExists(tmpDir)
	if err != nil {
		t.Fatalf("ensureConfigExists() = %v", err)
	}
	if !strings.HasSuffix(got, "swag2mcp.yaml") {
		t.Errorf("got %q, want path ending with swag2mcp.yaml", got)
	}
	if _, statErr := os.Stat(got); os.IsNotExist(statErr) {
		t.Error("config file was not created")
	}
}

func TestEnsureConfigExists_WorkspaceError(t *testing.T) {
	_, err := ensureConfigExists("/invalid:\x00path")
	if err == nil {
		t.Fatal("ensureConfigExists() expected error, got nil")
	}
}

func TestEnsureAuthScripts_NoScriptAuth(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API",
				BaseURL:  "https://example.com",
				Auth: config.Auth{
					Client: &auth.BearerTokenAuthClient{Token: "token"},
				},
			},
		},
	}

	ensureAuthScripts(cfg, ws)

	scriptPath := ws.AuthScriptPath("test-api")
	if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
		t.Error("auth script should not exist for non-script auth")
	}
}

func TestEnsureAuthScripts_ScriptAuth(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
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
			},
		},
	}

	ensureAuthScripts(cfg, ws)

	scriptPath := ws.AuthScriptPath("script-api")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Error("auth script should exist for script auth")
	}
}

func TestApplyMCPConfig_NilConfig(t *testing.T) {
	cmd := &cobra.Command{}
	opts := &mcpCmdOpts{Transport: "stdio", HTTPAddr: ":8080", HTTPPath: "/mcp"}

	applyMCPConfig(cmd, nil, opts)

	if opts.Transport != "stdio" {
		t.Errorf("Transport = %q, want stdio", opts.Transport)
	}
}

func TestApplyMCPConfig_NilMCP(t *testing.T) {
	cmd := &cobra.Command{}
	cfg := &config.Config{}
	opts := &mcpCmdOpts{Transport: "stdio"}

	applyMCPConfig(cmd, cfg, opts)

	if opts.Transport != "stdio" {
		t.Errorf("Transport = %q, want stdio", opts.Transport)
	}
}

func TestApplyMCPConfig_AppliesFromConfig(t *testing.T) {
	cmd := &cobra.Command{}
	cfg := &config.Config{
		MCP: &config.MCPConfig{
			Transport: transportSSE,
			Addr:      ":9090",
			Path:      "/api/mcp",
			Auth:      &config.MCPAuthConfig{Token: "secret"},
		},
	}
	opts := &mcpCmdOpts{}

	applyMCPConfig(cmd, cfg, opts)

	if opts.Transport != transportSSE {
		t.Errorf("Transport = %q, want %q", opts.Transport, transportSSE)
	}
	if opts.HTTPAddr != ":9090" {
		t.Errorf("HTTPAddr = %q, want :9090", opts.HTTPAddr)
	}
	if opts.HTTPPath != "/api/mcp" {
		t.Errorf("HTTPPath = %q, want /api/mcp", opts.HTTPPath)
	}
	if opts.AuthToken != "secret" {
		t.Errorf("AuthToken = %q, want secret", opts.AuthToken)
	}
}

func TestApplyMCPConfig_CLIFlagsOverride(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("transport", "stdio", "")
	cmd.Flags().String("http-addr", ":8080", "")
	cmd.Flags().String("http-path", "/mcp", "")
	cmd.Flags().String("auth-token", "", "")
	cmd.Flags().Set("transport", transportStreamableHTTP)

	cfg := &config.Config{
		MCP: &config.MCPConfig{
			Transport: transportSSE,
			Addr:      ":9090",
			Path:      "/api/mcp",
			Auth:      &config.MCPAuthConfig{Token: "secret"},
		},
	}
	opts := &mcpCmdOpts{Transport: transportStreamableHTTP}

	applyMCPConfig(cmd, cfg, opts)

	if opts.Transport != transportStreamableHTTP {
		t.Errorf("Transport = %q, want %q (CLI flag should win)", opts.Transport, transportStreamableHTTP)
	}
	if opts.HTTPAddr != ":9090" {
		t.Errorf("HTTPAddr = %q, want :9090 (config fallback)", opts.HTTPAddr)
	}
}
