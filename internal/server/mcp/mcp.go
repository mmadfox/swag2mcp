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
//
//nolint:funlen
func Serve(ctx context.Context, opts Options) error {
	if opts.Service == nil {
		return errors.New("service is required")
	}

	defs, err := opts.Service.MakeToolDefinitions()
	if err != nil {
		return fmt.Errorf("failed to make tool definitions: %w", err)
	}

	var mcpTransport sdkmcp.Transport
	if opts.Logger != nil {
		mcpTransport = &sdkmcp.LoggingTransport{
			Transport: &sdkmcp.StdioTransport{},
			Writer:    opts.Logger,
		}
	} else {
		mcpTransport = &sdkmcp.StdioTransport{}
	}

	srvOpts := &sdkmcp.ServerOptions{
		Instructions: defs.Instruction,
	}

	mcpServer := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    service.Name,
		Version: opts.Version,
	}, srvOpts)

	h := handler{service: opts.Service}
	for _, deftool := range defs.Tools {
		switch deftool.Name {
		case service.EndpointByTag:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleEndpointsByTag)
		case service.EndpointByCollection:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleEndpointsByCollection)
		case service.EndpointBySpec:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleEndpointsBySpec)
		case service.EndpointByID:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleEndpointByID)
		case service.TagByCollection:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleTagsByCollection)

		case service.TagBySpec:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleTagsBySpec)

		case service.TagByID:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleTagByID)

		case service.SpecByID:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleSpecByID)

		case service.SpecList:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleSpecList)

		case service.CollectionBySpec:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleCollectionBySpec)

		case service.CollectionByID:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleCollectionByID)

		case service.Search:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleSearch)

		case service.Inspect:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
				Annotations: &sdkmcp.ToolAnnotations{
					IdempotentHint: true,
					ReadOnlyHint:   true,
				},
			}, h.handleInspect)

		case service.Invoke:
			sdkmcp.AddTool(mcpServer, &sdkmcp.Tool{
				Name:        deftool.Name,
				Description: deftool.Description,
			}, h.handleInvoke)
		}
	}
	return mcpServer.Run(ctx, mcpTransport)
}
