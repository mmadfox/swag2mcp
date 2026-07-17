package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [path]",
		Short: "Show detailed configuration and runtime information",
		Long: `Show a comprehensive summary of the swag2mcp runtime: version,
configuration, active specs, HTTP client settings, MCP transport,
auth methods, and mock mode status.

  swag2mcp info              — show info for ~/.swag2mcp/swag2mcp.yaml
  swag2mcp info ./           — show info for ./swag2mcp.yaml
  swag2mcp info path/to      — show info for path/to/swag2mcp.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			return runInfo(basePath, cmd.OutOrStdout(), cmd.Context())
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runInfo(basePath string, w io.Writer, ctx context.Context) error {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return err
	}

	configPath := ws.ConfigPath()

	if ws.ConfigNotExists() {
		configPath, err = ensureConfigExists(basePath)
		if err != nil {
			return err
		}
	}

	svc, svcErr := service.New(
		service.WithVersion(Version),
		service.WithIndexNoFullText(),
	)
	if svcErr != nil {
		return fmt.Errorf("failed to create service: %w", svcErr)
	}

	if err := svc.Bootstrap(ctx, service.BootstrapRequest{
		ConfFilePath: configPath,
	}); err != nil {
		return fmt.Errorf("failed to bootstrap: %w", err)
	}

	info, infoErr := svc.Info(ctx)
	if infoErr != nil {
		return fmt.Errorf("failed to get info: %w", infoErr)
	}

	out, marshalErr := json.MarshalIndent(info, "", "  ")
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal info: %w", marshalErr)
	}

	_, err = io.WriteString(w, string(out)+"\n")
	return err
}
