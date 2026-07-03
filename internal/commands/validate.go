package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newValidateCmd() *cobra.Command {
	opts := struct {
		ConfigPath string
	}{}

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration file",
		Long: `Validate the configuration file and report any issues.

  swag2mcp validate -c ./swag2mcp.yaml`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			configPath := opts.ConfigPath
			if configPath == "" {
				configPath = workspace.DefaultConfigPath()
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			if err := cfg.Validate(config.NewFilter(nil)); err != nil {
				return fmt.Errorf("❌ Configuration is invalid:\n  %w\n  File: %s", err, configPath)
			}

			cmd.Printf("✅ Configuration is valid.\n")
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to configuration file")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
