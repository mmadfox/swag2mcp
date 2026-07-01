package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

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
			// set up logging if we have a logfile
			var logWriter io.Writer
			if len(opts.Logfile) > 0 {
				f, logErr := os.Create(opts.Logfile)
				if logErr != nil {
					return fmt.Errorf("opening logfile: %w", logErr)
				}
				defer f.Close()
				logWriter = f
			}

			// initialize service
			svc, svcErr := service.New()
			if svcErr != nil {
				return fmt.Errorf("failed to create service: %w", svcErr)
			}

			// bootstrap service with config
			if bootErr := svc.Bootstrap(cmd.Context(), service.BootstrapRequest{
				ConfFilepath: opts.ConfigFile,
			}); bootErr != nil {
				return fmt.Errorf("failed to bootstrap service: %w", bootErr)
			}

			// create mcp server options
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
