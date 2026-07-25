package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRunMCP_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := &cobra.Command{}
	opts := &mcpCmdOpts{}
	err := runMCP(tmpDir, "dev", opts, cmd)
	if err == nil {
		t.Fatal("runMCP() expected error, got nil")
	}
}
