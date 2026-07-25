package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fakeRateLimiter struct{}

func (fakeRateLimiter) Allow(_ string) error { return nil }

func newTestInvokeSvc(t *testing.T, idx IndexReader, ws WorkspaceOps) *invokeService {
	t.Helper()
	ctx := newServiceContext()
	ctx.storeHTTPClient(&http.Client{Transport: http.DefaultTransport})
	ctx.maxResponseSize.Store(config.DefaultMaxResponseSize)
	ctx.storeRateLimiter(fakeRateLimiter{})
	return newInvokeService(ctx, idx, ws, fakeValidator{}, "")
}

func TestInvokeService_Invoke_validationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ctx := newServiceContext()
	ctx.storeRateLimiter(fakeRateLimiter{})
	svc := newInvokeService(ctx, NewMockIndexReader(ctrl), NewMockWorkspaceOps(ctrl), strictValidator{}, "")
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "bad"})
	require.Error(t, err)
}

func TestInvokeService_Invoke_rateLimiterDisabled(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("ep1").Return(&model.Endpoint{ID: "ep1", SpecID: "s1", CollectionID: "c1", Operation: &spec.Operation{}}, nil)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1"}, nil)
	idx.EXPECT().CollectionByID("c1").Return(&model.Collection{ID: "c1"}, nil)

	ctx := newServiceContext()
	ctx.disableRateLimiter.Store(true)
	svc := newInvokeService(ctx, idx, NewMockWorkspaceOps(ctrl), fakeValidator{}, "")
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "ep1"})
	require.Error(t, err) // buildRequest will fail due to missing base URL, but rate limiter was skipped
}

type rateLimitReached struct{}

func (rateLimitReached) Allow(_ string) error {
	return errors.New(`rate limit exceeded for endpoint "ep1": try again in 8 seconds`)
}

func TestInvokeService_Invoke_rateLimitError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ctx := newServiceContext()
	ctx.storeRateLimiter(rateLimitReached{})
	svc := newInvokeService(ctx, NewMockIndexReader(ctrl), NewMockWorkspaceOps(ctrl), fakeValidator{}, "")
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "ep1"})
	require.Error(t, err)

	var llmErr *LLMError
	require.ErrorAs(t, err, &llmErr)
	require.Equal(t, "rate_limit", llmErr.Code)
	require.Contains(t, llmErr.Message, "try again in 8 seconds")
	require.Contains(t, llmErr.Hint, "Wait for the cooldown period")
}

func TestInvokeService_Invoke_paramValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("ep1").Return(&model.Endpoint{
		ID: "ep1", SpecID: "s1", CollectionID: "c1",
		Operation: &spec.Operation{
			Parameters: []*spec.Parameter{{Name: "id", In: "path", Required: true}},
		},
	}, nil)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1"}, nil)
	idx.EXPECT().CollectionByID("c1").Return(&model.Collection{ID: "c1"}, nil)

	svc := newTestInvokeSvc(t, idx, NewMockWorkspaceOps(ctrl))
	_, err := svc.Invoke(context.Background(), InvokeRequest{
		EndpointID: "ep1",
		Parameters: map[string]any{"unknown": "val"},
	})
	require.Error(t, err)
}

func TestInvokeService_Invoke_bodyValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("ep1").Return(&model.Endpoint{
		ID: "ep1", SpecID: "s1", CollectionID: "c1",
		Operation: &spec.Operation{
			RequestBody: &spec.RequestBody{
				Required: true,
				Content: map[string]*spec.MediaType{
					"application/json": {
						Schema: &spec.Schema{
							Type:       "object",
							Properties: map[string]*spec.Schema{"name": {Type: "string"}},
						},
					},
				},
			},
		},
	}, nil)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1"}, nil)
	idx.EXPECT().CollectionByID("c1").Return(&model.Collection{ID: "c1"}, nil)

	svc := newTestInvokeSvc(t, idx, NewMockWorkspaceOps(ctrl))
	_, err := svc.Invoke(context.Background(), InvokeRequest{
		EndpointID:  "ep1",
		RequestBody: map[string]any{"unknown": "val"},
	})
	require.Error(t, err)
}

func TestInvokeService_Invoke_endpointNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("missing").Return(nil, errNotFound("endpoint", "missing"))

	svc := newTestInvokeSvc(t, idx, NewMockWorkspaceOps(ctrl))
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "missing"})
	require.Error(t, err)
}

func TestInvokeService_Invoke_specNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("ep1").Return(&model.Endpoint{ID: "ep1", SpecID: "s1", Operation: &spec.Operation{}}, nil)
	idx.EXPECT().SpecByID("s1").Return(nil, errNotFound("spec", "s1"))

	svc := newTestInvokeSvc(t, idx, NewMockWorkspaceOps(ctrl))
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "ep1"})
	require.Error(t, err)
}

func TestInvokeService_Invoke_collectionNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("ep1").Return(&model.Endpoint{ID: "ep1", SpecID: "s1", CollectionID: "c1", Operation: &spec.Operation{}}, nil)
	idx.EXPECT().SpecByID("s1").Return(&model.Spec{ID: "s1"}, nil)
	idx.EXPECT().CollectionByID("c1").Return(nil, errNotFound("collection", "c1"))

	svc := newTestInvokeSvc(t, idx, NewMockWorkspaceOps(ctrl))
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "ep1"})
	require.Error(t, err)
}

func TestInvokeService_Invoke_nilOperation(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().EndpointByID("ep1").Return(&model.Endpoint{ID: "ep1", SpecID: "s1", CollectionID: "c1"}, nil)

	svc := newTestInvokeSvc(t, idx, NewMockWorkspaceOps(ctrl))
	_, err := svc.Invoke(context.Background(), InvokeRequest{EndpointID: "ep1"})
	require.Error(t, err)
}

func TestInvokeService_buildRequest(t *testing.T) {
	t.Parallel()

	svc := newTestInvokeSvc(t, NewMockIndexReader(gomock.NewController(t)), NewMockWorkspaceOps(gomock.NewController(t)))
	req, err := svc.buildRequest(
		context.Background(),
		&model.Spec{BaseURL: "https://api.example.com"},
		&model.Collection{},
		&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}},
		InvokeRequest{},
	)
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, "https://api.example.com/test", req.URL.String())
}

func TestInvokeService_executeRequest_success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	svc := newTestInvokeSvc(t, NewMockIndexReader(gomock.NewController(t)), NewMockWorkspaceOps(gomock.NewController(t)))
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	resp, err := svc.executeRequest(context.Background(), req, &model.Spec{}, &model.Endpoint{Name: "GET", Path: "/"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestInvokeService_executeRequest_withAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	svc := newTestInvokeSvc(t, NewMockIndexReader(gomock.NewController(t)), NewMockWorkspaceOps(gomock.NewController(t)))
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)

	resp, err := svc.executeRequest(context.Background(), req, &model.Spec{
		Auth: &noopAuth{},
	}, &model.Endpoint{Name: "GET", Path: "/"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestInvokeService_saveLargeResponse(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(filepath.Join(tmpDir, "responses")).AnyTimes()

	svc := newTestInvokeSvc(t, NewMockIndexReader(ctrl), ws)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	body := []byte(`{"data": "large response content"}`)
	result, err := svc.saveLargeResponse(resp, body, "test-domain", &model.Endpoint{Name: "GET", Path: "/test"}, 10)
	require.NoError(t, err)
	require.NotNil(t, result.FileRef)
	require.FileExists(t, result.FileRef.Path)
}

func TestInvokeService_dumpRequest(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestInvokeSvc(t, NewMockIndexReader(gomock.NewController(t)), NewMockWorkspaceOps(gomock.NewController(t)))
	svc.dumpDir = tmpDir

	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	require.NoError(t, err)

	svc.dumpRequest(req, "test-domain")
	// Verify a dump file was created
	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
}

func TestInvokeService_dumpRequest_emptyDir(t *testing.T) {
	t.Parallel()

	svc := newTestInvokeSvc(t, NewMockIndexReader(gomock.NewController(t)), NewMockWorkspaceOps(gomock.NewController(t)))
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	require.NoError(t, err)

	// Should not panic when dumpDir is empty
	svc.dumpRequest(req, "test-domain")
}

type noopAuth struct{}

func (noopAuth) New() error { return nil }

func (noopAuth) Type() auth.Type { return auth.NoAuth }

func (noopAuth) Apply(_ *http.Request, _ *auth.Info) error { return nil }

func (noopAuth) Validate() error { return nil }
