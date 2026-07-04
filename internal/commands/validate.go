package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newValidateCmd() *cobra.Command {
	opts := struct {
		ConfigPath string
		Tags       string
	}{}

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration file",
		Long: `Validate the configuration file and report any issues.

  swag2mcp validate -c ./swag2mcp.yaml
  swag2mcp validate -c ./swag2mcp.yaml --tags=public,internal`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			configPath := opts.ConfigPath
			if configPath == "" {
				configPath = workspace.DefaultConfigPath()
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			var tags []string
			if opts.Tags != "" {
				tags = strings.Split(opts.Tags, ",")
				for i := range tags {
					tags[i] = strings.TrimSpace(tags[i])
				}
			}

			validateOpts := config.ValidateOptions{
				Cache: cache.New(cfg.WorkspaceDir),
				Tags:  tags,
			}

			if err := config.ValidateConfig(cfg, validateOpts); err != nil {
				cmd.Printf("❌ %s\n", err)
				os.Exit(1)
			}

			cmd.Printf("✅ Configuration is valid.\n")
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&opts.Tags, "tags", "t", "", "Filter specs by tags (comma-separated)")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
