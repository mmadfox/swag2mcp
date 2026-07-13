package commands

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/tui"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [path]",
		Short: "Interactive API explorer",
		Long: `Interactive API explorer for searching, browsing, and invoking endpoints.

  swag2mcp run              — run ~/.swag2mcp/swag2mcp.yaml
  swag2mcp run ./           — run ./swag2mcp.yaml
  swag2mcp run path/to      — run path/to/swag2mcp.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			return runRun(basePath, cmd.Context())
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runRun(basePath string, ctx context.Context) error {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return err
	}

	configFile := ws.ConfigPath()

	if ws.ConfigNotExists() {
		configFile, err = ensureConfigExists(basePath)
		if err != nil {
			return err
		}
	}

	svc, svcErr := service.New()
	if svcErr != nil {
		return svcErr
	}

	if err := svc.Bootstrap(ctx, service.BootstrapRequest{
		ConfFilepath: configFile,
	}); err != nil {
		return err
	}

	return tui.RunExplorer(svc, svc.Workspace())
}
