package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/server/mcp"
	"github.com/mmadfox/swag2mcp/internal/service"
)

func newMCPCmd(version string) *cobra.Command {
	opts := struct {
		ConfigFile string
		Logfile    string
	}{}

	cmd := &cobra.Command{
		Use:   "mcp [mcp-flags]",
		Short: "Start the swag2mcp server in headless mode",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var logWriter io.Writer
			if len(opts.Logfile) > 0 {
				f, logErr := os.Create(opts.Logfile)
				if logErr != nil {
					return fmt.Errorf("opening logfile: %w", logErr)
				}
				defer f.Close()
				logWriter = f
			}

			if opts.ConfigFile != "" {
				cfg, loadErr := config.Load(opts.ConfigFile)
				if loadErr == nil {
					validateOpts := config.ValidateOptions{
						Cache: cache.New(cfg.WorkspaceDir),
					}
					if err := config.ValidateConfig(cfg, validateOpts); err != nil {
						fmt.Fprintf(os.Stderr, "⚠️  Configuration warnings:\n%s\n", err)
					}
				}
			}

			svc, svcErr := service.New()
			if svcErr != nil {
				return fmt.Errorf("failed to create service: %w", svcErr)
			}

			if bootErr := svc.Bootstrap(cmd.Context(), service.BootstrapRequest{
				ConfFilepath: opts.ConfigFile,
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

	cmd.Flags().StringVarP(&opts.ConfigFile, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&opts.Logfile, "logfile", "f", "", "Filename to log to; if unset, logs to stderr")

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
