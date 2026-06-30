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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// set up logging if we have a logfile
			var logWriter io.Writer
			if len(opts.Logfile) > 0 {
				f, err := os.Create(opts.Logfile)
				if err != nil {
					return fmt.Errorf("opening logfile: %v", err)
				}
				defer f.Close()
				logWriter = f
			}

			// initialize service
			svc, err := service.New()
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}

			// bootstrap service with config
			err = svc.Bootstrap(cmd.Context(), service.BootstrapRequest{
				ConfFilepath: opts.ConfigFile,
			})
			if err != nil {
				return fmt.Errorf("failed to bootstrap service: %w", err)
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
