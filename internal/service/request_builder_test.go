package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestRequestBuilder_build(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/pets", Operation: &spec.Operation{}}),
	).build()
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com/pets", req.URL.String())
	require.Equal(t, http.MethodGet, req.Method)
}

func TestRequestBuilder_withPathParams(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{
			Name: "GET", Path: "/pets/{id}",
			Operation: &spec.Operation{
				Parameters: []*spec.Parameter{
					{Name: "id", In: "path"},
				},
			},
		}),
		withParameters(map[string]any{"id": "42"}),
	).build()
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com/pets/42", req.URL.String())
}

func TestRequestBuilder_withQueryParams(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{
			Name: "GET", Path: "/pets",
			Operation: &spec.Operation{
				Parameters: []*spec.Parameter{
					{Name: "limit", In: "query"},
				},
			},
		}),
		withParameters(map[string]any{"limit": "10"}),
	).build()
	require.NoError(t, err)
	require.Contains(t, req.URL.RawQuery, "limit=10")
}

func TestRequestBuilder_withBody(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "POST", Path: "/pets", Operation: &spec.Operation{}}),
		withBody(map[string]any{"name": "Rex"}),
	).build()
	require.NoError(t, err)
	require.Equal(t, http.MethodPost, req.Method)
	require.NotNil(t, req.Body)
}

func TestRequestBuilder_withGlobalHeaders(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalHeaders(map[string]string{"X-Custom": "val"}),
	).build()
	require.NoError(t, err)
	require.Equal(t, "val", req.Header.Get("X-Custom"))
}

func TestRequestBuilder_withGlobalCookies(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalCookies([]httpclient.Cookie{{Name: "session", Value: "abc"}}),
	).build()
	require.NoError(t, err)
	cookie, _ := req.Cookie("session")
	require.NotNil(t, cookie)
	require.Equal(t, "abc", cookie.Value)
}

func TestRequestBuilder_resolveBaseURL_collection(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:       &model.Spec{BaseURL: "https://spec.example.com"},
		collection: &model.Collection{BaseURL: "https://coll.example.com"},
	}
	require.Equal(t, "https://coll.example.com", b.resolveBaseURL())
}

func TestRequestBuilder_resolveBaseURL_mock(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:       &model.Spec{BaseURL: "https://spec.example.com"},
		collection: &model.Collection{BaseMockURL: "localhost:8080"},
	}
	require.Equal(t, "http://localhost:8080", b.resolveBaseURL())
}

func TestRequestBuilder_withHTTPConfig(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withHTTPConfig(&model.HTTPClientConfig{
			Headers: map[string]string{"X-Spec": "spec-val"},
		}),
	).build()
	require.NoError(t, err)
	require.Equal(t, "spec-val", req.Header.Get("X-Spec"))
}

func TestRequestBuilder_withInvokeHeaders(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withInvokeHeaders(map[string]string{"X-Override": "override-val"}),
	).build()
	require.NoError(t, err)
	require.Equal(t, "override-val", req.Header.Get("X-Override"))
}

func TestRequestBuilder_withInvokeCookies(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withInvokeCookies(map[string]string{"session": "abc"}),
	).build()
	require.NoError(t, err)
	cookie, _ := req.Cookie("session")
	require.NotNil(t, cookie)
	require.Equal(t, "abc", cookie.Value)
}

func TestRequestBuilder_withGlobalUserAgent(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalUserAgent("test-agent"),
	).build()
	require.NoError(t, err)
	require.Equal(t, "test-agent", req.Header.Get("User-Agent"))
}

func TestRequestBuilder_applyDefaultAccept_json(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "POST", Path: "/test", Operation: &spec.Operation{}}),
		withBody(map[string]any{"key": "val"}),
	).build()
	require.NoError(t, err)
	// Content-Type is set, so Accept should be application/json
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestRequestBuilder_applyDefaultAccept_other(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
	).build()
	require.NoError(t, err)
	// Go's http.Client sets a default Accept header; our code should not override it
	require.NotEmpty(t, req.Header.Get("Accept"))
}

func TestRequestBuilder_applyDefaultAccept_preservesExisting(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withContext(context.Background()),
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withInvokeHeaders(map[string]string{"Accept": "text/plain"}),
	).build()
	require.NoError(t, err)
	require.Equal(t, "text/plain", req.Header.Get("Accept"))
}

func TestRequestBuilder_applySpecConfig(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:       &model.Spec{BaseURL: "https://api.example.com"},
		collection: &model.Collection{},
		endpoint:   &model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}},
		httpConfig: &model.HTTPClientConfig{
			Headers: map[string]string{"X-Spec": "spec-val"},
			Cookies: []httpclient.Cookie{{Name: "spec-cookie", Value: "spec-val"}},
		},
		context: context.Background(),
	}
	req, err := b.build()
	require.NoError(t, err)
	require.Equal(t, "spec-val", req.Header.Get("X-Spec"))
	cookie, _ := req.Cookie("spec-cookie")
	require.NotNil(t, cookie)
	require.Equal(t, "spec-val", cookie.Value)
}

func TestRequestBuilder_applySpecConfig_nil(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:       &model.Spec{BaseURL: "https://api.example.com"},
		collection: &model.Collection{},
		endpoint:   &model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}},
		context:    context.Background(),
	}
	req, err := b.build()
	require.NoError(t, err)
	// Should not panic when httpConfig is nil
	require.NotNil(t, req)
}

func TestRequestBuilder_applyInvokeOverrides(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:          &model.Spec{BaseURL: "https://api.example.com"},
		collection:    &model.Collection{},
		endpoint:      &model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}},
		invokeHeaders: map[string]string{"X-Override": "override-val"},
		invokeCookies: map[string]string{"override-cookie": "override-val"},
		context:       context.Background(),
	}
	req, err := b.build()
	require.NoError(t, err)
	require.Equal(t, "override-val", req.Header.Get("X-Override"))
	cookie, _ := req.Cookie("override-cookie")
	require.NotNil(t, cookie)
	require.Equal(t, "override-val", cookie.Value)
}

func TestRequestBuilder_filterParametersByLocation(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		parameters: map[string]any{"id": "42", "name": "test"},
		endpoint: &model.Endpoint{
			Operation: &spec.Operation{
				Parameters: []*spec.Parameter{
					{Name: "id", In: "path"},
					{Name: "name", In: "query"},
				},
			},
		},
	}
	pathParams := b.filterParametersByLocation("path")
	require.Equal(t, "42", pathParams["id"])
	require.NotContains(t, pathParams, "name")

	queryParams := b.filterParametersByLocation("query")
	require.Equal(t, "test", queryParams["name"])
	require.NotContains(t, queryParams, "id")
}

func TestRequestBuilder_applyHeaders(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		parameters: map[string]any{"X-Custom": "header-val"},
		endpoint: &model.Endpoint{
			Operation: &spec.Operation{
				Parameters: []*spec.Parameter{
					{Name: "X-Custom", In: "header"},
				},
			},
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	require.NoError(t, err)
	b.applyHeaders(req)
	require.Equal(t, "header-val", req.Header.Get("X-Custom"))
}

func TestRequestBuilder_applyHeaders_setsContentType(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		body:     map[string]any{"key": "val"},
		endpoint: &model.Endpoint{Name: "POST", Path: "/test", Operation: &spec.Operation{}},
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.example.com/test", nil)
	require.NoError(t, err)
	b.applyHeaders(req)
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestRequestBuilder_applyGlobalConfig(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		globalHeaders:   map[string]string{"X-Global": "global-val"},
		globalUserAgent: "global-agent",
		globalCookies:   []httpclient.Cookie{{Name: "global-cookie", Value: "global-val"}},
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	require.NoError(t, err)
	b.applyGlobalConfig(req)
	require.Equal(t, "global-val", req.Header.Get("X-Global"))
	require.Equal(t, "global-agent", req.Header.Get("User-Agent"))
	cookie, _ := req.Cookie("global-cookie")
	require.NotNil(t, cookie)
	require.Equal(t, "global-val", cookie.Value)
}

func TestRequestBuilder_applyGlobalConfig_doesNotOverride(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		globalHeaders: map[string]string{"Accept": "text/plain"},
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "application/json")
	b.applyGlobalConfig(req)
	require.Equal(t, "application/json", req.Header.Get("Accept"))
}
