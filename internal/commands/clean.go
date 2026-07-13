package commands

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean [path]",
		Short: "Remove temporary data (cache and responses)",
		Long: `Remove cached remote specs and invocation responses.

  swag2mcp clean              — clean ~/.swag2mcp/{cache,responses}
  swag2mcp clean ./           — clean ./.swag2mcp/{cache,responses}
  swag2mcp clean path/to      — clean path/to/.swag2mcp/{cache,responses}`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			return runClean(basePath, cmd.OutOrStdout())
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runClean(basePath string, w io.Writer) error {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return fmt.Errorf("workspace: %w", err)
	}

	if err := ws.Clean(); err != nil {
		return fmt.Errorf("clean: %w", err)
	}

	if ws.ConfigExists() {
		cfg, loadErr := config.Load(ws.ConfigPath())
		if loadErr == nil {
			var activeDomains []string
			for spec := range cfg.Iterate(nil) {
				activeDomains = append(activeDomains, spec.Domain)
			}
			if oErr := ws.RemoveOrphanAuthScripts(activeDomains); oErr != nil {
				return fmt.Errorf("remove orphan auth scripts: %w", oErr)
			}
		}
	}

	fmt.Fprintln(w, "✅ Removed contents")
	return nil
}
