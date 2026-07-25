package commands

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newValidateCmd() *cobra.Command {
	opts := struct {
		Tags string
	}{}

	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate the configuration file",
		Long: `Validate the configuration file and report any issues.

  swag2mcp validate              — validate ~/.swag2mcp/swag2mcp.yaml
  swag2mcp validate ./           — validate ./swag2mcp.yaml
  swag2mcp validate path/to      — validate path/to/swag2mcp.yaml
  swag2mcp validate --tags=public,internal`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			return runValidate(basePath, opts.Tags, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&opts.Tags, "tags", "t", "", "Filter specs by tags (comma-separated)")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runValidate(basePath, tagsFilter string, w io.Writer) error {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return fmt.Errorf("workspace: %w", err)
	}

	configPath := ws.ConfigPath()

	if ws.ConfigNotExists() {
		return fmt.Errorf("configuration not found at %s", configPath)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var tags []string
	if tagsFilter != "" {
		tags = strings.Split(tagsFilter, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	validateOpts := config.ValidateOptions{
		Cache: cache.New(filepath.Dir(configPath)),
		Tags:  tags,
	}

	if err := config.ValidateConfig(cfg, validateOpts); err != nil {
		fmt.Fprintf(w, "❌ %s\n", err)
		return err
	}

	fmt.Fprintf(w, "✅ Configuration is valid.\n")
	return nil
}
