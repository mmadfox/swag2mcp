package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/initpkg"
)

func newInitCmd() *cobra.Command {
	opts := struct {
		ConfigPath   string
		WorkspaceDir string
	}{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize workspace and configuration",
		Long: `Initialize workspace and configuration.

Without flags, starts an interactive wizard that guides you through
setting up the workspace and adding API specifications.

With --config-path and --workspace-dir, creates the workspace and
writes the example configuration file non-interactively.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			hasFlags := opts.ConfigPath != "" || opts.WorkspaceDir != ""

			if hasFlags {
				if opts.ConfigPath == "" {
					opts.ConfigPath = "swag2mcp.yaml"
				}
				if opts.WorkspaceDir == "" {
					opts.WorkspaceDir = ".swag2mcp"
				}
				if err := initpkg.Setup(opts.ConfigPath, opts.WorkspaceDir); err != nil {
					return fmt.Errorf("init: %w", err)
				}
				cmd.Printf("Configuration written to %s\n", opts.ConfigPath)
				cmd.Printf("Workspace initialized at %s\n", opts.WorkspaceDir)
				return nil
			}

			configPath, workspaceDir, _, err := initpkg.RunTUI()
			if err != nil {
				return fmt.Errorf("init wizard: %w", err)
			}

			cmd.Printf("\n  ✅ Configuration written to: %s\n", configPath)
			cmd.Printf("  ✅ Workspace initialized at: %s\n", workspaceDir)
			cmd.Println("  Run `swag2mcp mcp` to start the server.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config-path", "c", "", "Path to write the configuration file (non-interactive)")
	cmd.Flags().StringVarP(&opts.WorkspaceDir, "workspace-dir", "w", "", "Workspace directory path (non-interactive)")

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
