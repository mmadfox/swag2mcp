package commands

import (
	"context"
	"fmt"
	"log/slog"
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

type mcpCmdOpts struct {
	Logfile        string
	Tags           string
	DisableLLMAuth bool
	DumpDir        string
	Transport      string
	HTTPAddr       string
	HTTPPath       string
	AuthToken      string
}

const (
	transportSSE            = "sse"
	transportStreamableHTTP = "streamable-http"
)

func newMCPCmd(version string) *cobra.Command {
	opts := mcpCmdOpts{}

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

			var logger *slog.Logger
			if len(opts.Logfile) > 0 {
				f, logErr := os.Create(opts.Logfile)
				if logErr != nil {
					return fmt.Errorf("opening logfile: %w", logErr)
				}
				defer f.Close()
				logger = slog.New(slog.NewTextHandler(f, nil))
			} else {
				logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
			}

			cfg, loadErr := config.Load(configFile)
			if loadErr == nil {
				if cfg.MCP != nil {
					cfg.MCP.Auth.Resolve()
				}
				validateOpts := config.ValidateOptions{
					Cache: cache.New(filepath.Dir(configFile)),
				}
				if err := config.ValidateConfig(cfg, validateOpts); err != nil {
					fmt.Fprintf(os.Stderr, "⚠️  Configuration warnings:\n%s\n", err)
				}
			}

			// Apply MCP config from YAML as fallback when CLI flags are not set
			applyMCPConfig(cmd, cfg, &opts)

			var tags []string
			if opts.Tags != "" {
				tags = strings.Split(opts.Tags, ",")
				for i := range tags {
					tags[i] = strings.TrimSpace(tags[i])
				}
			}

			svcOpts := []service.NewOption{
				service.WithDisableLLMAuth(opts.DisableLLMAuth),
			}
			if opts.DumpDir != "" {
				svcOpts = append(svcOpts, service.WithDumpDir(opts.DumpDir))
			}

			svc, svcErr := service.New(svcOpts...)
			if svcErr != nil {
				return fmt.Errorf("failed to create service: %w", svcErr)
			}

			if bootErr := svc.Bootstrap(cmd.Context(), service.BootstrapRequest{
				ConfFilepath: configFile,
				Tags:         tags,
			}); bootErr != nil {
				return fmt.Errorf("failed to bootstrap service: %w", bootErr)
			}

			if cleanErr := ws.CleanOldResponses(workspace.DefaultResponseMaxAge); cleanErr != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Failed to clean old responses: %s\n", cleanErr)
			}

			transportType := mcp.TransportStdio
			switch opts.Transport {
			case transportSSE:
				transportType = mcp.TransportSSE
			case transportStreamableHTTP:
				transportType = mcp.TransportStreamableHTTP
			}

			mcpOpts := mcp.Options{
				Version:   version,
				Logger:    logger,
				Service:   svc,
				Transport: transportType,
				HTTPAddr:  opts.HTTPAddr,
				HTTPPath:  opts.HTTPPath,
				AuthToken: opts.AuthToken,
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			return mcp.Serve(ctx, mcpOpts)
		},
	}

	cmd.Flags().StringVarP(&opts.Logfile, "logfile", "f", "", "Filename to log to; if unset, logs to stderr")
	cmd.Flags().StringVarP(&opts.Tags, "tags", "t", "", "Filter specs by tags (comma-separated)")
	cmd.Flags().BoolVar(&opts.DisableLLMAuth, "disable-llm-auth", true, "Disable LLM auth token retrieval")
	cmd.Flags().StringVar(&opts.DumpDir, "dump-dir", "", "Directory to dump HTTP requests for debugging")
	cmd.Flags().StringVar(&opts.Transport, "transport", "stdio", "MCP transport: stdio, sse, streamable-http")
	cmd.Flags().StringVar(&opts.HTTPAddr, "http-addr", ":8080", "HTTP server address (for sse/streamable-http)")
	cmd.Flags().StringVar(&opts.HTTPPath, "http-path", "/mcp", "HTTP path for MCP handler")
	cmd.Flags().StringVar(&opts.AuthToken, "auth-token", "", "Bearer token for HTTP transport auth")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

// applyMCPConfig applies MCP settings from YAML config as fallback
// when the corresponding CLI flags were not explicitly set.
func applyMCPConfig(cmd *cobra.Command, cfg *config.Config, opts *mcpCmdOpts) {
	if cfg == nil || cfg.MCP == nil {
		return
	}
	if !cmd.Flags().Changed("transport") && cfg.MCP.Transport != "" {
		opts.Transport = cfg.MCP.Transport
	}
	if !cmd.Flags().Changed("http-addr") && cfg.MCP.Addr != "" {
		opts.HTTPAddr = cfg.MCP.Addr
	}
	if !cmd.Flags().Changed("http-path") && cfg.MCP.Path != "" {
		opts.HTTPPath = cfg.MCP.Path
	}
	if !cmd.Flags().Changed("auth-token") && cfg.MCP.Auth != nil && cfg.MCP.Auth.Token != "" {
		opts.AuthToken = cfg.MCP.Auth.Token
	}
}
