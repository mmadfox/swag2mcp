package mcp

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mmadfox/swag2mcp/internal/service"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Options holds the options for the MCP server.
type Options struct {
	Version string
	Logger  io.Writer
	Service svc
}

// Serve starts the MCP server.
func Serve(ctx context.Context, opts Options) error {
	if opts.Service == nil {
		return errors.New("service is required")
	}

	defs, err := opts.Service.MakeToolDefinitions()
	if err != nil {
		return fmt.Errorf("failed to make tool definitions: %w", err)
	}

	mcpServer := newServer(defs, opts)
	h := handler{service: opts.Service}
	registerTools(mcpServer, defs.Tools, h)

	return mcpServer.Run(ctx, newTransport(opts))
}

func newTransport(opts Options) sdkmcp.Transport {
	if opts.Logger != nil {
		return &sdkmcp.LoggingTransport{
			Transport: &sdkmcp.StdioTransport{},
			Writer:    opts.Logger,
		}
	}
	return &sdkmcp.StdioTransport{}
}

func newServer(defs service.ToolDefinitions, opts Options) *sdkmcp.Server {
	return sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    service.Name,
		Version: opts.Version,
	}, &sdkmcp.ServerOptions{
		Instructions: defs.Instruction,
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
