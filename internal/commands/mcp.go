/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/server/mcp"
	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

const (
	transportSSE            = "sse"
	transportStreamableHTTP = "streamable-http"
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
	AuthType       string
	AuthJWKSURL    string
	AuthIssuer     string
	AuthAudience   string
	AuthIntroURL   string
	AuthClientID   string
	AuthClientSec  string
}

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
			return runMCP(basePath, version, &opts, cmd)
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
	cmd.Flags().StringVar(&opts.AuthType, "auth-type", "", "JWT auth type: jwks, oidc, introspection")
	cmd.Flags().StringVar(&opts.AuthJWKSURL, "auth-jwks-url", "", "JWKS URL for JWT auth")
	cmd.Flags().StringVar(&opts.AuthIssuer, "auth-issuer", "", "JWT issuer for token validation")
	cmd.Flags().StringVar(&opts.AuthAudience, "auth-audience", "", "JWT audience for token validation")
	cmd.Flags().StringVar(&opts.AuthIntroURL, "auth-introspection-url", "", "Token introspection URL")
	cmd.Flags().StringVar(&opts.AuthClientID, "auth-client-id", "", "Client ID for introspection auth")
	cmd.Flags().StringVar(&opts.AuthClientSec, "auth-client-secret", "", "Client secret for introspection auth")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runMCP(basePath, version string, opts *mcpCmdOpts, cmd *cobra.Command) error {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return fmt.Errorf("workspace: %w", err)
	}

	configFile := ws.ConfigPath()

	if ws.ConfigNotExists() {
		var err error
		configFile, err = ensureConfigExists(basePath)
		if err != nil {
			return fmt.Errorf("configuration not found at %s: %w", configFile, err)
		}
	}

	var logger *slog.Logger
	if len(opts.Logfile) > 0 {
		f, logErr := os.Create(opts.Logfile)
		if logErr != nil {
			return fmt.Errorf("opening logfile: %w", logErr)
		}
		defer f.Close()
		logger = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{AddSource: true}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
	}
	slog.SetDefault(logger)

	cfg, loadErr := config.Load(configFile)
	if loadErr == nil && cfg.MCP != nil {
		cfg.MCP.Auth.Resolve()
		applyMCPConfig(cmd, cfg, opts)
	}

	var tags []string
	if opts.Tags != "" {
		tags = strings.Split(opts.Tags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	svcOpts := []service.Option{
		service.WithDisableLLMAuth(opts.DisableLLMAuth),
		service.WithVersion(version),
		service.WithLogger(logger),
	}
	if opts.DumpDir != "" {
		svcOpts = append(svcOpts, service.WithDumpDir(opts.DumpDir))
	}

	svc, svcErr := service.New(svcOpts...)
	if svcErr != nil {
		return fmt.Errorf("failed to create service: %w", svcErr)
	}

	if bootErr := svc.Bootstrap(cmd.Context(), service.BootstrapRequest{
		ConfFilePath: configFile,
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

	var authJWT *mcp.JWTConfig
	if opts.AuthType != "" {
		authJWT = &mcp.JWTConfig{
			Type:             opts.AuthType,
			JWKSURL:          opts.AuthJWKSURL,
			Issuer:           opts.AuthIssuer,
			Audience:         opts.AuthAudience,
			IntrospectionURL: opts.AuthIntroURL,
			ClientID:         opts.AuthClientID,
			ClientSecret:     opts.AuthClientSec,
		}
	}

	mcpOpts := mcp.Options{
		Version:   version,
		Logger:    logger,
		Service:   svc,
		Transport: transportType,
		HTTPAddr:  opts.HTTPAddr,
		HTTPPath:  opts.HTTPPath,
		AuthToken: opts.AuthToken,
		AuthJWT:   authJWT,
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	return mcp.Serve(ctx, mcpOpts)
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
	if cfg.MCP.Auth == nil {
		return
	}
	if !cmd.Flags().Changed("auth-token") && cfg.MCP.Auth.Token != "" {
		opts.AuthToken = cfg.MCP.Auth.Token
	}
	if !cmd.Flags().Changed("auth-type") && cfg.MCP.Auth.Type != "" {
		opts.AuthType = cfg.MCP.Auth.Type
	}
	if !cmd.Flags().Changed("auth-jwks-url") && cfg.MCP.Auth.JWKSURL != "" {
		opts.AuthJWKSURL = cfg.MCP.Auth.JWKSURL
	}
	if !cmd.Flags().Changed("auth-issuer") && cfg.MCP.Auth.Issuer != "" {
		opts.AuthIssuer = cfg.MCP.Auth.Issuer
	}
	if !cmd.Flags().Changed("auth-audience") && cfg.MCP.Auth.Audience != "" {
		opts.AuthAudience = cfg.MCP.Auth.Audience
	}
	if !cmd.Flags().Changed("auth-introspection-url") && cfg.MCP.Auth.IntrospectionURL != "" {
		opts.AuthIntroURL = cfg.MCP.Auth.IntrospectionURL
	}
	if !cmd.Flags().Changed("auth-client-id") && cfg.MCP.Auth.ClientID != "" {
		opts.AuthClientID = cfg.MCP.Auth.ClientID
	}
	if !cmd.Flags().Changed("auth-client-secret") && cfg.MCP.Auth.ClientSecret != "" {
		opts.AuthClientSec = cfg.MCP.Auth.ClientSecret
	}
}
