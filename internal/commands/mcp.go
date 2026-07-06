package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/server/mcp"
	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newMCPCmd(version string) *cobra.Command {
	opts := struct {
		Logfile        string
		Tags           string
		DisableLLMAuth bool
	}{}

	cmd := &cobra.Command{
		Use:   "mcp [path]",
		Short: "Start the swag2mcp server in headless mode",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}

			ws, err := workspace.NewFromBase(basePath)
			if err != nil {
				return fmt.Errorf("workspace: %w", err)
			}

			configFile := ws.ConfigPath()

			if ws.ConfigNotExists() {
				return fmt.Errorf("configuration not found at %s", configFile)
			}

			var logWriter io.Writer
			if len(opts.Logfile) > 0 {
				f, logErr := os.Create(opts.Logfile)
				if logErr != nil {
					return fmt.Errorf("opening logfile: %w", logErr)
				}
				defer f.Close()
				logWriter = f
			}

			cfg, loadErr := config.Load(configFile)
			if loadErr == nil {
				validateOpts := config.ValidateOptions{
					Cache: cache.New(filepath.Dir(configFile)),
				}
				if err := config.ValidateConfig(cfg, validateOpts); err != nil {
					fmt.Fprintf(os.Stderr, "⚠️  Configuration warnings:\n%s\n", err)
				}
			}

			var tags []string
			if opts.Tags != "" {
				tags = strings.Split(opts.Tags, ",")
				for i := range tags {
					tags[i] = strings.TrimSpace(tags[i])
				}
			}

			svc, svcErr := service.New(
				service.WithDisableLLMAuth(opts.DisableLLMAuth),
			)
			if svcErr != nil {
				return fmt.Errorf("failed to create service: %w", svcErr)
			}

			if bootErr := svc.Bootstrap(cmd.Context(), service.BootstrapRequest{
				ConfFilepath: configFile,
				Tags:         tags,
			}); bootErr != nil {
				return fmt.Errorf("failed to bootstrap service: %w", bootErr)
			}

			mcpOpts := mcp.Options{
				Version: version,
				Logger:  logWriter,
				Service: svc,
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			return mcp.Serve(ctx, mcpOpts)
		},
	}

	cmd.Flags().StringVarP(&opts.Logfile, "logfile", "f", "", "Filename to log to; if unset, logs to stderr")
	cmd.Flags().StringVarP(&opts.Tags, "tags", "t", "", "Filter specs by tags (comma-separated)")
	cmd.Flags().BoolVar(&opts.DisableLLMAuth, "disable-llm-auth", false, "Disable LLM auth token retrieval")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
