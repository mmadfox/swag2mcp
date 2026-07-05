package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/tui"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newInitCmd() *cobra.Command {
	opts := struct {
		Interactive bool
		Force       bool
	}{}

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize workspace and configuration",
		Long: `Initialize workspace and configuration.

  swag2mcp init              — create ~/.swag2mcp/swag2mcp.yaml
  swag2mcp init ./           — create ./.swag2mcp/swag2mcp.yaml
  swag2mcp init path/to      — create path/to/.swag2mcp/swag2mcp.yaml
  swag2mcp init -i           — interactive wizard
  swag2mcp init -f           — force overwrite existing configuration`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := "."
			if len(args) > 0 {
				basePath = args[0]
			}

			absBase, err := filepath.Abs(basePath)
			if err != nil {
				return fmt.Errorf("resolve path: %w", err)
			}

			workspaceDir := filepath.Join(absBase, workspace.DefaultRootName)
			configPath := workspace.ConfigPathIn(workspaceDir)

			if opts.Interactive {
				if !opts.Force {
					if _, err := os.Stat(configPath); err == nil {
						return fmt.Errorf("configuration already exists at %s\n  Use --force to overwrite", configPath)
					}
				}

				cfgPath, wsDir, _, err := tui.RunTUI()
				if err != nil {
					return fmt.Errorf("init wizard: %w", err)
				}

				cmd.Printf("\n✅ Configuration written to: %s\n", cfgPath)
				cmd.Printf("✅ Workspace initialized at: %s\n", wsDir)
				cmd.Println("Run `swag2mcp mcp` to start the server.")
				return nil
			}

			if !opts.Force {
				if _, err := os.Stat(configPath); err == nil {
					return fmt.Errorf("configuration already exists at %s\n  Use --force to overwrite", configPath)
				}
			}

			if err := tui.Setup(configPath, workspaceDir); err != nil {
				return fmt.Errorf("init: %w", err)
			}

			cmd.Printf("✅ Configuration written to %s\n", configPath)
			cmd.Printf("✅ Workspace initialized at %s\n", workspaceDir)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Run interactive wizard")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Overwrite existing configuration")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
