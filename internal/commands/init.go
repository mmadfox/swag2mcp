package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/initmcp"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newInitCmd() *cobra.Command {
	opts := struct {
		ConfigPath   string
		WorkspaceDir string
		Force        bool
	}{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize workspace and configuration",
		Long: `Initialize workspace and configuration.

Without flags, starts an interactive wizard that guides you through
setting up the workspace and adding API specifications.

With --config-path and --workspace-dir, creates the workspace and
writes the example configuration file non-interactively.

Use --force to overwrite an existing configuration.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			hasFlags := opts.ConfigPath != "" || opts.WorkspaceDir != ""

			if hasFlags {
				if opts.ConfigPath == "" {
					opts.ConfigPath = workspace.DefaultConfigPath()
				}
				if opts.WorkspaceDir == "" {
					opts.WorkspaceDir = workspace.DefaultRoot()
				}
				if !opts.Force {
					if _, err := os.Stat(opts.ConfigPath); err == nil {
						return fmt.Errorf("configuration already exists at %s\n  Use --force to overwrite", opts.ConfigPath)
					}
				}
				if err := initmcp.Setup(opts.ConfigPath, opts.WorkspaceDir); err != nil {
					return fmt.Errorf("init: %w", err)
				}
				cmd.Printf("Configuration written to %s\n", opts.ConfigPath)
				cmd.Printf("Workspace initialized at %s\n", opts.WorkspaceDir)
				return nil
			}

			if !opts.Force {
				configPath := workspace.DefaultConfigPath()
				if _, err := os.Stat(configPath); err == nil {
					return fmt.Errorf("configuration already exists at %s\n  Use --force to overwrite", configPath)
				}
			}

			configPath, workspaceDir, _, err := initmcp.RunTUI()
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
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing configuration")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
