// Package mockserver provides mock HTTP servers for API specifications and
// authentication methods defined in the swag2mcp configuration.
//
// The MockServer starts two kinds of servers:
//   - Auth mock servers: two global servers (OAuth2 on port 9090, Digest on
//     port 9091) that simulate the real auth flow. Other auth types (Basic,
//     Bearer, API Key, Script) do not need a mock server — the MCP server
//     applies authentication automatically via applyMockAuthURLs.
//   - API mock servers: one per collection, each on a separate port. They parse
//     the OpenAPI/Swagger spec and respond to requests with randomly generated
//     data that conforms to the response schema.
//
// Addresses for mock servers are taken from the base_mock_url field in the
// configuration (spec or collection level). The format is "host:port".
package mockserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// Options holds configuration for the MockServer.
type Options struct {
	// Config is the parsed swag2mcp configuration containing specs and collections.
	Config *config.Config

	// ConfigPath is the path to the configuration file.
	ConfigPath string

	// Workspace is the swag2mcp workspace directory.
	Workspace *workspace.Workspace

	// TLS enables TLS for all mock servers. If TLSCert and TLSKey are empty,
	// a self-signed certificate is generated.
	TLS bool

	// TLSCert is the path to a TLS certificate file. If empty and TLS is true,
	// a self-signed certificate is used.
	TLSCert string

	// TLSKey is the path to a TLS key file. If empty and TLS is true,
	// a self-signed certificate is used.
	TLSKey string
}

// MockServer manages all mock servers (auth and API) for the configured specs.
type MockServer struct {
	options     Options
	authServers []*authMockServer
	apiServers  []*apiMockServer
	logger      *slog.Logger
	tlsConfig   *tls.Config
	mu          sync.Mutex
}

// New creates a new MockServer with the given options.
func New(options Options) *MockServer {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	var tlsConfig *tls.Config
	if options.TLS {
		var err error
		tlsConfig, err = newTLSConfig(options.TLSCert, options.TLSKey)
		if err != nil {
			logger.Error("failed to create TLS config, falling back to HTTP", "error", err)
		}
	}

	return &MockServer{
		options:   options,
		logger:    logger,
		tlsConfig: tlsConfig,
	}
}

// Start launches all mock servers and blocks until a shutdown signal is received.
// It parses the configuration, creates auth and API mock servers for each spec,
// prints a summary table, and waits for SIGINT or SIGTERM.
// Addresses are taken from base_mock_url fields in the configuration.
func (m *MockServer) Start(ctx context.Context) error {
	if !m.options.Config.MockEnabled {
		return fmt.Errorf(
			"mock server mode is disabled in configuration\n\n"+
				"Configuration file: %s\n\n"+
				"To use the mock server, set 'mock_enabled: true' in your swag2mcp.yaml "+
				"and configure 'base_mock_url' for each collection\n\n"+
				"Example:\n"+
				"  mock_enabled: true\n"+
				"  specs:\n"+
				"    - domain: petstore\n"+
				"      base_url: https://petstore.swagger.io/v2\n"+
				"      collections:\n"+
				"        - location: specs/petstore.json\n"+
				"          base_mock_url: localhost:8080",
			m.options.ConfigPath,
		)
	}

	m.startAuthServers()
	m.startAPIServers()

	if len(m.authServers) == 0 && len(m.apiServers) == 0 {
		return errors.New("no mock servers to start — check your configuration")
	}

	startContext, startCancel := context.WithCancel(ctx)
	defer startCancel()

	m.startAll(startContext)

	m.printSummary()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChannel:
		m.logger.InfoContext(ctx, "shutting down...")
	case <-ctx.Done():
	}

	m.shutdownAll()

	return nil
}

const (
	defaultOAuth2Port = 9090
	defaultDigestPort = 9091
	defaultHMACPort   = 9092
)

func (m *MockServer) startAuthServers() {
	oauth2Port := defaultOAuth2Port
	digestPort := defaultDigestPort
	hmacPort := defaultHMACPort
	if m.options.Config.MockAuth != nil {
		if m.options.Config.MockAuth.OAuth2Port > 0 {
			oauth2Port = m.options.Config.MockAuth.OAuth2Port
		}
		if m.options.Config.MockAuth.DigestPort > 0 {
			digestPort = m.options.Config.MockAuth.DigestPort
		}
		if m.options.Config.MockAuth.HMACPort > 0 {
			hmacPort = m.options.Config.MockAuth.HMACPort
		}
	}

	addr := fmt.Sprintf("127.0.0.1:%d", oauth2Port)
	server := newAuthMockServer(authServerOAuth2, addr, m.tlsConfig, m.logger)
	m.mu.Lock()
	m.authServers = append(m.authServers, server)
	m.mu.Unlock()

	addr = fmt.Sprintf("127.0.0.1:%d", digestPort)
	server = newAuthMockServer(authServerDigest, addr, m.tlsConfig, m.logger)
	m.mu.Lock()
	m.authServers = append(m.authServers, server)
	m.mu.Unlock()

	addr = fmt.Sprintf("127.0.0.1:%d", hmacPort)
	server = newAuthMockServer(authServerHMAC, addr, m.tlsConfig, m.logger)
	m.mu.Lock()
	m.authServers = append(m.authServers, server)
	m.mu.Unlock()
}

func (m *MockServer) startAPIServers() {
	for specIndex := range m.options.Config.Specs {
		specConfig := &m.options.Config.Specs[specIndex]
		if specConfig.Disable {
			continue
		}

		for collectionIndex := range specConfig.Collections {
			collectionConfig := &specConfig.Collections[collectionIndex]
			if collectionConfig.Disable {
				continue
			}

			mockAddr := collectionConfig.BaseMockURL

			apiServer := newAPIMockServer(
				specConfig,
				collectionConfig,
				mockAddr,
				m.tlsConfig,
				m.logger,
				m.options.Workspace,
			)
			if apiServer == nil {
				continue
			}

			m.mu.Lock()
			m.apiServers = append(m.apiServers, apiServer)
			m.mu.Unlock()
		}
	}
}

func (m *MockServer) startAll(ctx context.Context) {
	for _, authServer := range m.authServers {
		authServer.start(ctx)
	}

	for _, apiServer := range m.apiServers {
		apiServer.start(ctx)
	}
}

func (m *MockServer) shutdownAll() {
	for _, authServer := range m.authServers {
		authServer.shutdown()
	}
	for _, apiServer := range m.apiServers {
		apiServer.shutdown()
	}
}

func (m *MockServer) printSummary() {
	var output strings.Builder

	output.WriteString("\n")
	output.WriteString("swag2mcp-mock\n")
	output.WriteString("━━━━━━━━━━━━━━\n\n")

	if len(m.authServers) > 0 {
		output.WriteString("Auth Mocks:\n")
		for _, authServer := range m.authServers {
			protocol := "http"
			if m.options.TLS {
				protocol = "https"
			}
			line := fmt.Sprintf("  %s → %s://%s\n",
				authServer.serverType,
				protocol,
				authServer.addr,
			)
			output.WriteString(line)
		}
		output.WriteString("\n")
	}

	if len(m.apiServers) > 0 {
		output.WriteString("API Mocks:\n")
		for _, apiServer := range m.apiServers {
			protocol := "http"
			if m.options.TLS {
				protocol = "https"
			}
			line := fmt.Sprintf("  %s / %s → %s://%s\n",
				apiServer.specDomain,
				apiServer.collectionTitle,
				protocol,
				apiServer.addr,
			)
			output.WriteString(line)
		}
	}

	output.WriteString("\n")

	if _, err := os.Stdout.WriteString(output.String()); err != nil {
		m.logger.Warn("failed to write mock server output", "error", err)
	}
}

// extractHostPort extracts the "host:port" portion from an address that may
// include a path suffix (e.g. "127.0.0.1:9000/v1/smev" → "127.0.0.1:9000").
func extractHostPort(addr string) string {
	hostPort, _, _ := strings.Cut(addr, "/")
	return hostPort
}
