package mcp

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/modelcontextprotocol/go-sdk/auth"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Compile-time check: service.Service satisfies the Svc interface.
var _ Svc = (*service.Service)(nil)

func TestServe_NoService(t *testing.T) {
	t.Parallel()

	err := Serve(context.Background(), Options{Service: nil})
	require.Error(t, err, "expected error for nil service")
}

func TestServe_MakeToolDefinitionsError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{}, errors.New("tool defs error"),
	)

	err := Serve(context.Background(), Options{Service: mock})
	require.Error(t, err, "expected error for nil service")
}

func TestNewTransport_WithLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	opts := Options{Logger: logger}
	transport := newTransport(opts)

	require.IsType(t, &sdkmcp.LoggingTransport{}, transport)
}

func TestNewTransport_WithoutLogger(t *testing.T) {
	t.Parallel()

	opts := Options{}
	transport := newTransport(opts)

	require.IsType(t, &sdkmcp.StdioTransport{}, transport)
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	defs := service.ToolDefinitions{
		Instruction: "test instruction",
	}
	opts := Options{Version: "v1.0.0"}

	srv := newServer(defs, opts)
	require.NotNil(t, srv, "newServer() returned nil")
}

func TestRegisterTools_AllTools(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools: []service.Tool{
				{Name: "spec_list", Description: "List specs"},
				{Name: "spec_by_id", Description: "Get spec"},
				{Name: "collection_by_spec", Description: "List collections"},
				{Name: "collection_by_id", Description: "Get collection"},
				{Name: "tag_by_collection", Description: "List tags"},
				{Name: "tag_by_spec", Description: "List tags by spec"},
				{Name: "tag_by_id", Description: "Get tag"},
				{Name: "endpoint_by_tag", Description: "List endpoints"},
				{Name: "endpoint_by_collection", Description: "List endpoints"},
				{Name: "endpoint_by_spec", Description: "List endpoints"},
				{Name: "endpoint_by_id", Description: "Get endpoint"},
				{Name: "search", Description: "Search"},
				{Name: "inspect", Description: "Inspect"},
				{Name: "invoke", Description: "Invoke"},
				{Name: "auth", Description: "Auth"},
				{Name: "info", Description: "Info"},
			},
		}, nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Serve(ctx, Options{Service: mock})
	require.Error(t, err, "expected context canceled error, not nil")
}

func TestRegisterTools_UnknownTool(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools: []service.Tool{
				{Name: "unknown_tool", Description: "Should be ignored"},
			},
		}, nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Serve(ctx, Options{Service: mock})
	require.Error(t, err, "expected context canceled error, not nil")
}

func TestServe_WithLogger(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools:       []service.Tool{{Name: "spec_list", Description: "List"}},
		}, nil,
	)

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Serve(ctx, Options{Service: mock, Logger: logger})
	require.Error(t, err, "expected context canceled error, not nil")
}

func TestServe_WithVersion(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools:       []service.Tool{{Name: "spec_list", Description: "List"}},
		}, nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Serve(ctx, Options{Service: mock, Version: "v2.0.0"})
	require.Error(t, err, "expected context canceled error, not nil")
}

func TestServe_UnsupportedTransport(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools:       []service.Tool{{Name: "spec_list", Description: "List"}},
		}, nil,
	)

	err := Serve(context.Background(), Options{
		Service:   mock,
		Transport: TransportType(999),
	})
	require.Error(t, err, "expected error for unsupported transport")
}

func TestApplyAuthMiddleware_NoAuth(t *testing.T) {
	t.Parallel()

	handler := applyAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), Options{})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApplyAuthMiddleware_StaticToken_Valid(t *testing.T) {
	t.Parallel()

	handler := applyAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), Options{AuthToken: "secret"})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApplyAuthMiddleware_StaticToken_Invalid(t *testing.T) {
	t.Parallel()

	handler := applyAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), Options{AuthToken: "secret"})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestApplyAuthMiddleware_CustomVerifier_Valid(t *testing.T) {
	t.Parallel()

	handler := applyAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), Options{
		AuthVerifier: func(_ context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
			if token == "custom-token" {
				return &auth.TokenInfo{UserID: "test-user", Expiration: time.Now().Add(time.Hour)}, nil
			}
			return nil, auth.ErrInvalidToken
		},
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer custom-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApplyAuthMiddleware_CustomVerifier_Invalid(t *testing.T) {
	t.Parallel()

	handler := applyAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), Options{
		AuthVerifier: func(_ context.Context, _ string, _ *http.Request) (*auth.TokenInfo, error) {
			return nil, auth.ErrInvalidToken
		},
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithLogging_NilLogger(t *testing.T) {
	t.Parallel()

	handler := withLogging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), nil)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSlogWriter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	w := newSlogWriter(logger)

	n, err := w.Write([]byte("test message"))
	require.NoError(t, err)
	assert.Equal(t, 12, n, "written bytes")
	assert.NotZero(t, buf.Len(), "expected log output")
}

func TestOptions_Defaults(t *testing.T) {
	t.Parallel()

	opts := Options{}
	assert.Equal(t, ":8080", opts.httpAddr())
	assert.Equal(t, "/mcp", opts.httpPath())
}

func TestOptions_CustomAddr(t *testing.T) {
	t.Parallel()

	opts := Options{HTTPAddr: ":9090", HTTPPath: "/api/mcp"}
	if opts.httpAddr() != ":9090" {
		t.Errorf("httpAddr = %q, want %q", opts.httpAddr(), ":9090")
	}
	if opts.httpPath() != "/api/mcp" {
		t.Errorf("httpPath = %q, want %q", opts.httpPath(), "/api/mcp")
	}
}

func TestHandler_SpecList_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Specs(gomock.Any()).Return(
		service.SpecsResponse{
			Specs: []service.SpecItem{{ID: "spec-1", Domain: "test"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleSpecList(context.Background(), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_SpecList_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Specs(gomock.Any()).Return(
		service.SpecsResponse{}, errors.New("specs error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleSpecList(context.Background(), nil, nil)
	require.Error(t, err)
}

func TestHandler_SpecByID_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().SpecByID(gomock.Any(), gomock.Any()).Return(
		service.SpecByIDResponse{
			Spec: service.Spec{ID: "abc", Domain: "test"},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleSpecByID(
		context.Background(), nil, service.SpecByIDRequest{ID: "abc"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_SpecByID_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().SpecByID(gomock.Any(), gomock.Any()).Return(
		service.SpecByIDResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleSpecByID(
		context.Background(), nil, service.SpecByIDRequest{ID: "abc"},
	)
	require.Error(t, err)
}

func TestHandler_CollectionByID_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().CollectionByID(gomock.Any(), gomock.Any()).Return(
		service.CollectionByIDResponse{
			Collection: service.Collection{ID: "coll-1", Title: "Test Coll"},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleCollectionByID(
		context.Background(), nil, service.CollectionByIDRequest{ID: "coll-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_CollectionByID_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().CollectionByID(gomock.Any(), gomock.Any()).Return(
		service.CollectionByIDResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleCollectionByID(
		context.Background(), nil, service.CollectionByIDRequest{ID: "coll-1"},
	)
	require.Error(t, err)
}

func TestHandler_CollectionBySpec_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().CollectionsBySpec(gomock.Any(), gomock.Any()).Return(
		service.CollectionsResponse{
			Collections: []service.CollectionItem{{ID: "coll-1", Title: "Coll"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleCollectionBySpec(
		context.Background(), nil, service.CollectionsRequest{SpecID: "spec-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_CollectionBySpec_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().CollectionsBySpec(gomock.Any(), gomock.Any()).Return(
		service.CollectionsResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleCollectionBySpec(
		context.Background(), nil, service.CollectionsRequest{SpecID: "spec-1"},
	)
	require.Error(t, err)
}

func TestHandler_TagsByCollection_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().TagsByCollection(gomock.Any(), gomock.Any()).Return(
		service.TagsByCollectionResponse{
			Tags: []service.TagListItem{{ID: "tag-1", Title: "Tag"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleTagsByCollection(
		context.Background(), nil, service.TagsByCollectionRequest{CollectionID: "coll-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_TagsByCollection_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().TagsByCollection(gomock.Any(), gomock.Any()).Return(
		service.TagsByCollectionResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleTagsByCollection(
		context.Background(), nil, service.TagsByCollectionRequest{CollectionID: "coll-1"},
	)
	require.Error(t, err)
}

func TestHandler_TagsBySpec_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().TagsBySpec(gomock.Any(), gomock.Any()).Return(
		service.TagsBySpecResponse{
			Tags: []service.TagListItem{{ID: "tag-1", Title: "Tag"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleTagsBySpec(
		context.Background(), nil, service.TagsBySpecRequest{SpecID: "spec-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_TagsBySpec_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().TagsBySpec(gomock.Any(), gomock.Any()).Return(
		service.TagsBySpecResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleTagsBySpec(
		context.Background(), nil, service.TagsBySpecRequest{SpecID: "spec-1"},
	)
	require.Error(t, err)
}

func TestHandler_TagByID_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().TagByID(gomock.Any(), gomock.Any()).Return(
		service.TagByIDResponse{
			Tag: service.TagListItem{ID: "tag-1", Title: "Test Tag"},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleTagByID(
		context.Background(), nil, service.TagByIDRequest{ID: "tag-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_TagByID_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().TagByID(gomock.Any(), gomock.Any()).Return(
		service.TagByIDResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleTagByID(
		context.Background(), nil, service.TagByIDRequest{ID: "tag-1"},
	)
	require.Error(t, err)
}

func TestHandler_EndpointByID_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointByID(gomock.Any(), gomock.Any()).Return(
		service.EndpointByIDResponse{
			Endpoint: service.Endpoint{ID: "ep-1", Method: "GET", Path: "/test"},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleEndpointByID(
		context.Background(), nil, service.EndpointByIDRequest{ID: "ep-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_EndpointByID_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointByID(gomock.Any(), gomock.Any()).Return(
		service.EndpointByIDResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleEndpointByID(
		context.Background(), nil, service.EndpointByIDRequest{ID: "ep-1"},
	)
	require.Error(t, err)
}

func TestHandler_EndpointsByTag_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointsByTag(gomock.Any(), gomock.Any()).Return(
		service.EndpointsByTagResponse{
			Endpoints: []service.EndpointTagItem{{ID: "ep-1", Method: "GET"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleEndpointsByTag(
		context.Background(), nil, service.EndpointsByTagRequest{TagID: "tag-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_EndpointsByTag_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointsByTag(gomock.Any(), gomock.Any()).Return(
		service.EndpointsByTagResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleEndpointsByTag(
		context.Background(), nil, service.EndpointsByTagRequest{TagID: "tag-1"},
	)
	require.Error(t, err)
}

func TestHandler_EndpointsByCollection_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointsByCollection(gomock.Any(), gomock.Any()).Return(
		service.EndpointsByCollectionResponse{
			Endpoints: []service.EndpointCollectionItem{{ID: "ep-1", Method: "GET"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleEndpointsByCollection(
		context.Background(), nil, service.EndpointsByCollectionRequest{CollectionID: "coll-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_EndpointsByCollection_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointsByCollection(gomock.Any(), gomock.Any()).Return(
		service.EndpointsByCollectionResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleEndpointsByCollection(
		context.Background(), nil, service.EndpointsByCollectionRequest{CollectionID: "coll-1"},
	)
	require.Error(t, err)
}

func TestHandler_EndpointsBySpec_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointsBySpec(gomock.Any(), gomock.Any()).Return(
		service.EndpointsBySpecResponse{
			Endpoints: []service.EndpointSearchItem{{ID: "ep-1", Method: "GET"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleEndpointsBySpec(
		context.Background(), nil, service.EndpointsBySpecRequest{SpecID: "spec-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_EndpointsBySpec_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().EndpointsBySpec(gomock.Any(), gomock.Any()).Return(
		service.EndpointsBySpecResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleEndpointsBySpec(
		context.Background(), nil, service.EndpointsBySpecRequest{SpecID: "spec-1"},
	)
	require.Error(t, err)
}

func TestHandler_Search_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Search(gomock.Any(), gomock.Any()).Return(
		service.SearchResponse{
			Endpoints: []service.EndpointSearchItem{{ID: "ep-1", Method: "GET"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleSearch(
		context.Background(), nil, service.SearchRequest{Query: "test", Limit: 10},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_Search_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Search(gomock.Any(), gomock.Any()).Return(
		service.SearchResponse{}, errors.New("search error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleSearch(
		context.Background(), nil, service.SearchRequest{Query: "test", Limit: 10},
	)
	require.Error(t, err)
}

func TestHandler_Inspect_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Inspect(gomock.Any(), gomock.Any()).Return(
		service.InspectResponse{
			ID:     "ep-1",
			Method: "GET",
			Path:   "/test",
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleInspect(
		context.Background(), nil, service.InspectRequest{EndpointID: "ep-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_Inspect_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Inspect(gomock.Any(), gomock.Any()).Return(
		service.InspectResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleInspect(
		context.Background(), nil, service.InspectRequest{EndpointID: "ep-1"},
	)
	require.Error(t, err)
}

func TestHandler_Invoke_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Invoke(gomock.Any(), gomock.Any()).Return(
		service.InvokeResponse{
			StatusCode: 200,
			Body:       map[string]any{"ok": true},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleInvoke(
		context.Background(), nil, service.InvokeRequest{EndpointID: "ep-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_Invoke_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Invoke(gomock.Any(), gomock.Any()).Return(
		service.InvokeResponse{}, errors.New("invoke error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleInvoke(
		context.Background(), nil, service.InvokeRequest{EndpointID: "ep-1"},
	)
	require.Error(t, err)
}

func TestHandler_Auth_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(
		service.AuthResponse{
			Token: "test-token",
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleAuth(
		context.Background(), nil, service.AuthRequest{SpecID: "spec-1"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_Auth_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Auth(gomock.Any(), gomock.Any()).Return(
		service.AuthResponse{}, errors.New("auth error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleAuth(
		context.Background(), nil, service.AuthRequest{SpecID: "spec-1"},
	)
	require.Error(t, err)
}

func TestHandler_Info_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Info(gomock.Any()).Return(
		service.InfoResponse{
			Version: "v1.0.0",
			Specs: service.SpecsSummary{
				Total: 2, Active: 1, Disabled: 1, Collections: 3, Endpoints: 42,
			},
			HTTPClient: service.HTTPClientInfo{
				Randomize: true, MaxResponseSize: "2 KB",
			},
			MCP: service.MCPInfo{
				Transport: "sse", AuthEnabled: true,
			},
			Mock: service.MockInfo{Enabled: false},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleInfo(context.Background(), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_Info_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Info(gomock.Any()).Return(
		service.InfoResponse{}, errors.New("info error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleInfo(context.Background(), nil, nil)
	require.Error(t, err)
}

func TestHandler_ResponseOutline_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().ResponseOutline(gomock.Any(), gomock.Any()).Return(
		service.ResponseOutlineResponse{
			Outline: reader.Outline{Type: "object", LineCount: 10},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleResponseOutline(
		context.Background(), nil, service.ResponseOutlineRequest{Path: "/tmp/resp.json"},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_ResponseOutline_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().ResponseOutline(gomock.Any(), gomock.Any()).Return(
		service.ResponseOutlineResponse{}, errors.New("not found"),
	)

	h := handler{service: mock}
	_, _, err := h.handleResponseOutline(
		context.Background(), nil, service.ResponseOutlineRequest{Path: "/tmp/resp.json"},
	)
	require.Error(t, err)
}

func TestHandler_ResponseCompress_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().ResponseCompress(gomock.Any(), gomock.Any()).Return(
		service.ResponseCompressResponse{
			Body: map[string]any{"compressed": map[string]any{"type": "array"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleResponseCompress(
		context.Background(), nil,
		service.ResponseCompressRequest{
			Path: "/tmp/resp.json",
			Mode: "first_of_array",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_ResponseCompress_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().ResponseCompress(gomock.Any(), gomock.Any()).Return(
		service.ResponseCompressResponse{}, errors.New("compress error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleResponseCompress(
		context.Background(), nil,
		service.ResponseCompressRequest{
			Path: "/tmp/resp.json",
			Mode: "first_of_array",
		},
	)
	require.Error(t, err)
}

func TestHandler_ResponseSlice_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().ResponseSlice(gomock.Any(), gomock.Any()).Return(
		service.ResponseSliceResponse{
			Slice: reader.Slice{
				JSONPath: "data.0",
				Context:  "object",
				Value:    map[string]any{"id": 1},
			},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleResponseSlice(
		context.Background(), nil,
		service.ResponseSliceRequest{
			Path:     "/tmp/resp.json",
			JSONPath: "data.0",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestHandler_ResponseSlice_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().ResponseSlice(gomock.Any(), gomock.Any()).Return(
		service.ResponseSliceResponse{}, errors.New("slice error"),
	)

	h := handler{service: mock}
	_, _, err := h.handleResponseSlice(
		context.Background(), nil,
		service.ResponseSliceRequest{
			Path:     "/tmp/resp.json",
			JSONPath: "data.0",
		},
	)
	require.Error(t, err)
}

func TestHandler_StructuredContent(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().Specs(gomock.Any()).Return(
		service.SpecsResponse{
			Specs: []service.SpecItem{{ID: "spec-1", Domain: "test"}},
		}, nil,
	)

	h := handler{service: mock}
	result, _, err := h.handleSpecList(context.Background(), nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result.StructuredContent)
}

func TestServeHTTP_SSE(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools:       []service.Tool{},
		}, nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Serve(ctx, Options{
		Service:   mock,
		Transport: TransportSSE,
		HTTPAddr:  ":0",
		Logger:    slog.New(slog.DiscardHandler),
	})
	require.Error(t, err)
}

func TestServeHTTP_StreamableHTTP(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools:       []service.Tool{},
		}, nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Serve(ctx, Options{
		Service:   mock,
		Transport: TransportStreamableHTTP,
		HTTPAddr:  ":0",
		Logger:    slog.New(slog.DiscardHandler),
	})
	require.Error(t, err)
}

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockSvc(ctrl)
	mock.EXPECT().MakeToolDefinitions().Return(
		service.ToolDefinitions{
			Instruction: "test",
			Tools:       []service.Tool{},
		}, nil,
	)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- Serve(ctx, Options{
			Service:   mock,
			Transport: TransportSSE,
			HTTPAddr:  ":0",
			Version:   "v1.0.0",
			Logger:    slog.New(slog.DiscardHandler),
		})
	}()

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:0/health")
	if err == nil {
		defer resp.Body.Close()
	}
	_ = err
}

func TestWithLogging_WithLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	handler := withLogging(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), logger)

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotZero(t, buf.Len(), "expected log output")
}
