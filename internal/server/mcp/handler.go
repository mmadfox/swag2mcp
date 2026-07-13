package mcp

import (
	"context"

	"github.com/mmadfox/swag2mcp/internal/service"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Svc is the service interface that the MCP handler depends on.
// It defines all operations exposed as MCP tools.
type Svc interface {
	Specs(_ context.Context) (service.SpecsResponse, error)
	SpecByID(_ context.Context, req service.SpecByIDRequest) (service.SpecByIDResponse, error)
	CollectionByID(_ context.Context, req service.CollectionByIDRequest) (service.CollectionByIDResponse, error)
	CollectionsBySpec(_ context.Context, req service.CollectionsRequest) (service.CollectionsResponse, error)
	TagsByCollection(_ context.Context, req service.TagsByCollectionRequest) (service.TagsByCollectionResponse, error)
	TagsBySpec(_ context.Context, req service.TagsBySpecRequest) (service.TagsBySpecResponse, error)
	TagByID(_ context.Context, req service.TagByIDRequest) (service.TagByIDResponse, error)
	EndpointsByTag(_ context.Context, req service.EndpointsByTagRequest) (service.EndpointsByTagResponse, error)
	EndpointsByCollection(
		_ context.Context,
		req service.EndpointsByCollectionRequest,
	) (service.EndpointsByCollectionResponse, error)
	EndpointsBySpec(_ context.Context, req service.EndpointsBySpecRequest) (service.EndpointsBySpecResponse, error)
	EndpointByID(_ context.Context, req service.EndpointByIDRequest) (service.EndpointByIDResponse, error)
	Search(ctx context.Context, req service.SearchRequest) (service.SearchResponse, error)
	Inspect(_ context.Context, req service.InspectRequest) (service.InspectResponse, error)
	Invoke(_ context.Context, req service.InvokeRequest) (service.InvokeResponse, error)
	Auth(_ context.Context, req service.AuthRequest) (service.AuthResponse, error)
	Info(_ context.Context) (service.InfoResponse, error)
	MakeToolDefinitions() (service.ToolDefinitions, error)
}

type handler struct {
	service Svc
}

// handleSpecByID handles the spec_by_id tool call.
func (h *handler) handleSpecByID(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.SpecByIDRequest,
) (*sdkmcp.CallToolResult, any, error) {
	spec, err := h.service.SpecByID(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: spec,
	}, nil, nil
}

// handleSpecList handles the spec_list tool call.
func (h *handler) handleSpecList(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	_ any,
) (*sdkmcp.CallToolResult, any, error) {
	specs, err := h.service.Specs(ctx)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: specs,
	}, nil, nil
}

// handleCollectionByID handles the collection_by_id tool call.
func (h *handler) handleCollectionByID(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.CollectionByIDRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.CollectionByID(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleCollectionBySpec handles the collection_by_spec tool call.
func (h *handler) handleCollectionBySpec(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.CollectionsRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.CollectionsBySpec(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleTagsByCollection handles the tag_by_collection tool call.
func (h *handler) handleTagsByCollection(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.TagsByCollectionRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.TagsByCollection(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleTagsBySpec handles the tag_by_spec tool call.
func (h *handler) handleTagsBySpec(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.TagsBySpecRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.TagsBySpec(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleTagByID handles the tag_by_id tool call.
func (h *handler) handleTagByID(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.TagByIDRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.TagByID(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleEndpointByID handles the endpoint_by_id tool call.
func (h *handler) handleEndpointByID(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.EndpointByIDRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.EndpointByID(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleEndpointsByTag handles the endpoint_by_tag tool call.
func (h *handler) handleEndpointsByTag(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.EndpointsByTagRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.EndpointsByTag(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleEndpointsByCollection handles the endpoint_by_collection tool call.
func (h *handler) handleEndpointsByCollection(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.EndpointsByCollectionRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.EndpointsByCollection(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleEndpointsBySpec handles the endpoint_by_spec tool call.
func (h *handler) handleEndpointsBySpec(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.EndpointsBySpecRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.EndpointsBySpec(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleSearch handles the search tool call.
func (h *handler) handleSearch(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.SearchRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.Search(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleInspect handles the inspect tool call.
func (h *handler) handleInspect(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.InspectRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.Inspect(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleInvoke handles the invoke tool call.
func (h *handler) handleInvoke(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.InvokeRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.Invoke(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleAuth handles the auth tool call.
func (h *handler) handleAuth(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.AuthRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.Auth(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

// handleInfo handles the info tool call.
func (h *handler) handleInfo(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	_ any,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.Info(ctx)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}
