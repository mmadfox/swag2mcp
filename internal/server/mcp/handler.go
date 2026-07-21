package mcp

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"

	"github.com/mmadfox/swag2mcp/internal/service"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// CatalogReader provides read-only access to the API catalog (specs, collections, tags, endpoints).
type CatalogReader interface {
	Specs(_ context.Context) (service.SpecsResponse, error)
	SpecByID(_ context.Context, req service.SpecByIDRequest) (service.SpecByIDResponse, error)
	CollectionsBySpec(_ context.Context, req service.CollectionsRequest) (service.CollectionsResponse, error)
	CollectionByID(_ context.Context, req service.CollectionByIDRequest) (service.CollectionByIDResponse, error)
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
}

// EndpointExplorer provides search and inspection capabilities.
type EndpointExplorer interface {
	Search(ctx context.Context, req service.SearchRequest) (service.SearchResponse, error)
	Inspect(_ context.Context, req service.InspectRequest) (service.InspectResponse, error)
}

// EndpointExecutor provides API invocation and authentication.
type EndpointExecutor interface {
	Invoke(_ context.Context, req service.InvokeRequest) (service.InvokeResponse, error)
	Auth(_ context.Context, req service.AuthRequest) (service.AuthResponse, error)
}

// SystemInfo provides runtime information and tool definitions.
type SystemInfo interface {
	Info(_ context.Context) (service.InfoResponse, error)
	MakeToolDefinitions() (service.ToolDefinitions, error)
}

// ResponseManager provides access to saved response files.
type ResponseManager interface {
	ResponseOutline(_ context.Context, req service.ResponseOutlineRequest) (service.ResponseOutlineResponse, error)
	ResponseCompress(_ context.Context, req service.ResponseCompressRequest) (service.ResponseCompressResponse, error)
	ResponseSlice(_ context.Context, req service.ResponseSliceRequest) (service.ResponseSliceResponse, error)
}

// Svc is the combined service interface consumed by MCP tool handlers.
type Svc interface {
	CatalogReader
	EndpointExplorer
	EndpointExecutor
	SystemInfo
	ResponseManager
}

type handler struct {
	service Svc
}

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

func (h *handler) handleResponseOutline(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.ResponseOutlineRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.ResponseOutline(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

func (h *handler) handleResponseCompress(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.ResponseCompressRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.ResponseCompress(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}

func (h *handler) handleResponseSlice(
	ctx context.Context,
	_ *sdkmcp.CallToolRequest,
	req service.ResponseSliceRequest,
) (*sdkmcp.CallToolResult, any, error) {
	resp, err := h.service.ResponseSlice(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return &sdkmcp.CallToolResult{
		StructuredContent: resp,
	}, nil, nil
}
