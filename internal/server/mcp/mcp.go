package mcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// TransportType defines the MCP server transport.
type TransportType int

const (
	TransportStdio TransportType = iota
	TransportSSE
	TransportStreamableHTTP
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 10 * time.Second
	tokenExpiry       = 300 * time.Hour
)

// TokenVerifier verifies bearer tokens for HTTP transport auth.
type TokenVerifier func(ctx context.Context, token string, req *http.Request) (*auth.TokenInfo, error)

// Options holds the options for the MCP server.
type Options struct {
	Version string
	Logger  *slog.Logger
	Service svc

	Transport TransportType

	HTTPAddr string
	HTTPPath string

	AuthToken    string
	AuthVerifier TokenVerifier
}

func (o Options) httpAddr() string {
	if o.HTTPAddr != "" {
		return o.HTTPAddr
	}
	return ":8080"
}

func (o Options) httpPath() string {
	if o.HTTPPath != "" {
		return o.HTTPPath
	}
	return "/mcp"
}

// Serve starts the MCP server with the configured transport.
func Serve(ctx context.Context, opts Options) error {
	if opts.Service == nil {
		return errors.New("service is required")
	}

	if opts.Logger == nil {
		opts.Logger = slog.New(slog.DiscardHandler)
	}

	defs, err := opts.Service.MakeToolDefinitions()
	if err != nil {
		return fmt.Errorf("failed to make tool definitions: %w", err)
	}

	switch opts.Transport {
	case TransportStdio:
		return serveStdio(ctx, defs, opts)
	case TransportSSE, TransportStreamableHTTP:
		return serveHTTP(ctx, defs, opts)
	default:
		return fmt.Errorf("unsupported transport: %v", opts.Transport)
	}
}

func serveStdio(ctx context.Context, defs service.ToolDefinitions, opts Options) error {
	mcpServer := newServer(defs, opts)
	h := handler{service: opts.Service}
	registerTools(mcpServer, defs.Tools, h)

	transport := newStdioTransport(opts)

	opts.Logger.InfoContext(ctx, "MCP server started",
		"transport", "stdio",
	)

	return mcpServer.Run(ctx, transport)
}

func newStdioTransport(opts Options) sdkmcp.Transport {
	t := &sdkmcp.StdioTransport{}
	if opts.Logger != nil {
		return &sdkmcp.LoggingTransport{
			Transport: t,
			Writer:    newSlogWriter(opts.Logger),
		}
	}
	return t
}

func serveHTTP(ctx context.Context, defs service.ToolDefinitions, opts Options) error {
	getServer := func(_ *http.Request) *sdkmcp.Server {
		srv := newServer(defs, opts)
		h := handler{service: opts.Service}
		registerTools(srv, defs.Tools, h)
		return srv
	}

	var handler http.Handler
	switch opts.Transport {
	case TransportSSE:
		handler = sdkmcp.NewSSEHandler(getServer, &sdkmcp.SSEOptions{})
	case TransportStreamableHTTP:
		handler = sdkmcp.NewStreamableHTTPHandler(getServer, &sdkmcp.StreamableHTTPOptions{
			Logger: opts.Logger,
		})
	case TransportStdio:
		return errors.New("stdio transport is not supported for HTTP server")
	}

	handler = applyAuthMiddleware(handler, opts)

	mux := http.NewServeMux()
	mux.Handle(opts.httpPath(), handler)

	srv := &http.Server{
		Addr:              opts.httpAddr(),
		Handler:           withLogging(mux, opts.Logger),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	transportName := "sse"
	if opts.Transport == TransportStreamableHTTP {
		transportName = "streamable-http"
	}

	opts.Logger.InfoContext(ctx, "MCP server started",
		"transport", transportName,
		"addr", opts.httpAddr(),
		"path", opts.httpPath(),
		"auth", opts.AuthToken != "" || opts.AuthVerifier != nil,
	)

	return srv.ListenAndServe()
}

func applyAuthMiddleware(next http.Handler, opts Options) http.Handler {
	if opts.AuthVerifier != nil {
		return auth.RequireBearerToken(
			func(ctx context.Context, token string, req *http.Request) (*auth.TokenInfo, error) {
				return opts.AuthVerifier(ctx, token, req)
			}, nil,
		)(next)
	}
	if opts.AuthToken != "" {
		verifier := func(_ context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
			if token == opts.AuthToken {
				return &auth.TokenInfo{
					UserID:     "swag2mcp-user",
					Expiration: time.Now().Add(tokenExpiry),
				}, nil
			}
			return nil, auth.ErrInvalidToken
		}
		return auth.RequireBearerToken(verifier, nil)(next)
	}
	return next
}

func withLogging(next http.Handler, logger *slog.Logger) http.Handler {
	if logger == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "MCP request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
	})
}

func newSlogWriter(logger *slog.Logger) io.Writer {
	return &slogWriter{logger: logger}
}

type slogWriter struct {
	logger *slog.Logger
}

func (w *slogWriter) Write(p []byte) (int, error) {
	w.logger.Info(string(p))
	return len(p), nil
}

func newTransport(opts Options) sdkmcp.Transport {
	return newStdioTransport(opts)
}

func newServer(defs service.ToolDefinitions, opts Options) *sdkmcp.Server {
	return sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    service.Name,
		Version: opts.Version,
	}, &sdkmcp.ServerOptions{
		Instructions: defs.Instruction,
		Logger:       opts.Logger,
	})
}

func registerTools(mcpServer *sdkmcp.Server, tools []service.Tool, h handler) {
	type reg struct {
		add      func(t *sdkmcp.Tool)
		readOnly bool
	}

	toolRegistrations := map[string]reg{
		service.EndpointByTag: {
			addTool[service.EndpointsByTagRequest](mcpServer, h.handleEndpointsByTag),
			true,
		},
		service.EndpointByCollection: {
			addTool[service.EndpointsByCollectionRequest](mcpServer, h.handleEndpointsByCollection),
			true,
		},
		service.EndpointBySpec: {
			addTool[service.EndpointsBySpecRequest](mcpServer, h.handleEndpointsBySpec),
			true,
		},
		service.EndpointByID: {
			addTool[service.EndpointByIDRequest](mcpServer, h.handleEndpointByID),
			true,
		},
		service.TagByCollection: {
			addTool[service.TagsByCollectionRequest](mcpServer, h.handleTagsByCollection),
			true,
		},
		service.TagBySpec: {
			addTool[service.TagsBySpecRequest](mcpServer, h.handleTagsBySpec),
			true,
		},
		service.TagByID: {
			addTool[service.TagByIDRequest](mcpServer, h.handleTagByID),
			true,
		},
		service.SpecByID: {
			addTool[service.SpecByIDRequest](mcpServer, h.handleSpecByID),
			true,
		},
		service.SpecList: {
			addTool[any](mcpServer, h.handleSpecList),
			true,
		},
		service.CollectionBySpec: {
			addTool[service.CollectionsRequest](mcpServer, h.handleCollectionBySpec),
			true,
		},
		service.CollectionByID: {
			addTool[service.CollectionByIDRequest](mcpServer, h.handleCollectionByID),
			true,
		},
		service.Search: {
			addTool[service.SearchRequest](mcpServer, h.handleSearch),
			true,
		},
		service.Inspect: {
			addTool[service.InspectRequest](mcpServer, h.handleInspect),
			true,
		},
		service.Invoke: {
			addTool[service.InvokeRequest](mcpServer, h.handleInvoke),
			false,
		},
		service.Auth: {
			addTool[service.AuthRequest](mcpServer, h.handleAuth),
			false,
		},
		service.Info: {
			addTool[any](mcpServer, h.handleInfo),
			true,
		},
	}

	for _, deftool := range tools {
		r, ok := toolRegistrations[deftool.Name]
		if !ok {
			continue
		}

		tool := &sdkmcp.Tool{
			Name:        deftool.Name,
			Description: deftool.Description,
		}
		if r.readOnly {
			tool.Annotations = &sdkmcp.ToolAnnotations{
				IdempotentHint: true,
				ReadOnlyHint:   true,
			}
		}

		r.add(tool)
	}
}

func addTool[In any](
	s *sdkmcp.Server,
	fn func(context.Context, *sdkmcp.CallToolRequest, In) (*sdkmcp.CallToolResult, any, error),
) func(t *sdkmcp.Tool) {
	return func(tool *sdkmcp.Tool) {
		sdkmcp.AddTool[In, any](s, tool, fn)
	}
}
