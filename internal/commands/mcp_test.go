package commands

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
