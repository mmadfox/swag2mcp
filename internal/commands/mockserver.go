package commands

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/server/mockserver"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

type mockServerCmdOptions struct {
	TLS     bool
	TLSCert string
	TLSKey  string
}

// NewMockServerCmd creates the mockserver subcommand that starts mock servers for all API specs.
func NewMockServerCmd() *cobra.Command {
	options := mockServerCmdOptions{}

	command := &cobra.Command{
		Use:   "mockserver [path]",
		Short: "Start mock servers for all API specifications",
		Long: `Start mock servers for all API specifications defined in the configuration.

  swag2mcp-mock              — start mocks for ~/.swag2mcp/swag2mcp.yaml
  swag2mcp-mock ./           — start mocks for ./swag2mcp.yaml
  swag2mcp-mock path/to      — start mocks for path/to/swag2mcp.yaml
  swag2mcp-mock --tls        — enable TLS with self-signed certificate

Addresses for mock servers are taken from the base_mock_url field in the
configuration (spec or collection level). Format: "host:port".`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, arguments []string) error {
			basePath := ""
			if len(arguments) > 0 {
				basePath = arguments[0]
			}
			return runMockServer(basePath, &options, command.Context())
		},
	}

	command.Flags().BoolVar(&options.TLS, "tls", false, "Enable TLS with self-signed certificate")
	command.Flags().StringVar(&options.TLSCert, "tls-cert", "", "Path to TLS certificate file")
	command.Flags().StringVar(&options.TLSKey, "tls-key", "", "Path to TLS key file")
	command.SilenceUsage = true
	command.SilenceErrors = true

	return command
}

func runMockServer(basePath string, opts *mockServerCmdOptions, ctx context.Context) error {
	ws, workspaceError := workspace.NewFromBase(basePath)
	if workspaceError != nil {
		return fmt.Errorf("workspace: %w", workspaceError)
	}

	configFile := ws.ConfigPath()

	if ws.ConfigNotExists() {
		return fmt.Errorf("configuration not found at %s", configFile)
	}

	configuration, loadError := config.Load(configFile)
	if loadError != nil {
		return fmt.Errorf("failed to load config: %w", loadError)
	}

	validateOpts := config.ValidateOptions{
		Cache: cache.New(filepath.Dir(configFile)),
	}
	if err := config.ValidateConfig(configuration, validateOpts); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	mockServerOptions := mockserver.Options{
		Config:     configuration,
		ConfigPath: configFile,
		Workspace:  ws,
		TLS:        opts.TLS,
		TLSCert:    opts.TLSCert,
		TLSKey:     opts.TLSKey,
	}

	mockServer := mockserver.New(mockServerOptions)

	return mockServer.Start(ctx)
}
