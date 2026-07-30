/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestRunMCP_NoConfigCreatesIt(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := &cobra.Command{}
	opts := &mcpCmdOpts{}
	err := runMCP(tmpDir, "dev", opts, cmd)
	if err == nil {
		t.Fatal("runMCP() expected error after autofill, got nil")
	}
	// Config was created but bootstrap fails because example config has remote URLs.
	// The key is it no longer fails with "configuration not found".
	ws, _ := workspace.New(tmpDir)
	if !ws.ConfigExists() {
		t.Fatal("config should have been created")
	}
}
