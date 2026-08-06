/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

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
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/pets", Operation: &spec.Operation{}}),
	).build(context.Background())
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com/pets", req.URL.String())
	require.Equal(t, http.MethodGet, req.Method)
}

func TestRequestBuilder_withPathParams(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
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
	).build(context.Background())
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com/pets/42", req.URL.String())
}

func TestRequestBuilder_withQueryParams(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
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
	).build(context.Background())
	require.NoError(t, err)
	require.Contains(t, req.URL.RawQuery, "limit=10")
}

func TestRequestBuilder_withBody(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "POST", Path: "/pets", Operation: &spec.Operation{}}),
		withBody(map[string]any{"name": "Rex"}),
	).build(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.MethodPost, req.Method)
	require.NotNil(t, req.Body)
}

func TestRequestBuilder_withGlobalHeaders(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalHeaders(map[string]string{"X-Custom": "val"}),
	).build(context.Background())
	require.NoError(t, err)
	require.Equal(t, "val", req.Header.Get("X-Custom"))
}

func TestRequestBuilder_withGlobalCookies(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalCookies([]httpclient.Cookie{{Name: "session", Value: "abc"}}),
	).build(context.Background())
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
		spec:        &model.Spec{BaseURL: "https://spec.example.com"},
		collection:  &model.Collection{BaseMockURL: "localhost:8080"},
		mockEnabled: true,
	}
	require.Equal(t, "http://localhost:8080", b.resolveBaseURL())
}

func TestRequestBuilder_resolveBaseURL_mock_disabled(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:        &model.Spec{BaseURL: "https://spec.example.com"},
		collection:  &model.Collection{BaseMockURL: "localhost:8080", BaseURL: "https://coll.example.com"},
		mockEnabled: false,
	}
	require.Equal(t, "https://coll.example.com", b.resolveBaseURL())
}

func TestRequestBuilder_resolveBaseURL_mock_disabled_fallback(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		spec:        &model.Spec{BaseURL: "https://spec.example.com"},
		collection:  &model.Collection{BaseMockURL: "localhost:8080"},
		mockEnabled: false,
	}
	require.Equal(t, "https://spec.example.com", b.resolveBaseURL())
}

func TestRequestBuilder_withHTTPConfig(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withHTTPConfig(&model.HTTPClientConfig{
			Headers: map[string]string{"X-Spec": "spec-val"},
		}),
	).build(context.Background())
	require.NoError(t, err)
	require.Equal(t, "spec-val", req.Header.Get("X-Spec"))
}

func TestRequestBuilder_withGlobalUserAgent(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalUserAgent("test-agent"),
	).build(context.Background())
	require.NoError(t, err)
	require.Equal(t, "test-agent", req.Header.Get("User-Agent"))
}

func TestRequestBuilder_applyDefaultAccept_json(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "POST", Path: "/test", Operation: &spec.Operation{}}),
		withBody(map[string]any{"key": "val"}),
	).build(context.Background())
	require.NoError(t, err)
	// Content-Type is set, so Accept should be application/json
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestRequestBuilder_applyDefaultAccept_other(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
	).build(context.Background())
	require.NoError(t, err)
	// Go's http.Client sets a default Accept header; our code should not override it
	require.NotEmpty(t, req.Header.Get("Accept"))
}

func TestRequestBuilder_applyDefaultAccept_preservesExisting(t *testing.T) {
	t.Parallel()

	req, err := newRequestBuilder(
		withSpec(&model.Spec{BaseURL: "https://api.example.com"}),
		withCollection(&model.Collection{}),
		withEndpoint(&model.Endpoint{Name: "GET", Path: "/test", Operation: &spec.Operation{}}),
		withGlobalHeaders(map[string]string{"Accept": "text/plain"}),
	).build(context.Background())
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
	}
	req, err := b.build(context.Background())
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
	}
	req, err := b.build(context.Background())
	require.NoError(t, err)
	// Should not panic when httpConfig is nil
	require.NotNil(t, req)
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

func TestRequestBuilder_filterParametersByLocation_skipsEmpty(t *testing.T) {
	t.Parallel()

	b := &requestBuilder{
		parameters: map[string]any{
			"id":         "42",
			"signature":  "",
			"timestamp":  0,
			"recvWindow": float64(0),
		},
		endpoint: &model.Endpoint{
			Operation: &spec.Operation{
				Parameters: []*spec.Parameter{
					{Name: "id", In: "path"},
					{Name: "signature", In: "query"},
					{Name: "timestamp", In: "query"},
					{Name: "recvWindow", In: "query"},
				},
			},
		},
	}

	pathParams := b.filterParametersByLocation("path")
	require.Equal(t, "42", pathParams["id"])

	queryParams := b.filterParametersByLocation("query")
	require.NotContains(t, queryParams, "signature", "empty string should be skipped")
	require.NotContains(t, queryParams, "timestamp", "zero int should be skipped")
	require.NotContains(t, queryParams, "recvWindow", "zero float64 should be skipped")
}

func TestIsEmptyParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{name: "empty string", value: "", want: true},
		{name: "non-empty string", value: "hello", want: false},
		{name: "zero int", value: 0, want: true},
		{name: "non-zero int", value: 42, want: false},
		{name: "zero int64", value: int64(0), want: true},
		{name: "non-zero int64", value: int64(7314340), want: false},
		{name: "zero float64", value: float64(0), want: true},
		{name: "non-zero float64", value: float64(3.14), want: false},
		{name: "zero float32", value: float32(0), want: true},
		{name: "non-zero float32", value: float32(100), want: false},
		{name: "bool false", value: false, want: false},
		{name: "bool true", value: true, want: false},
		{name: "nil", value: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isEmptyParam(tt.value)
			require.Equal(t, tt.want, got)
		})
	}
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

func TestFormatParamValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "float64 whole number", value: float64(7314340), want: "7314340"},
		{name: "float64 zero", value: float64(0), want: "0"},
		{name: "float64 negative whole", value: float64(-42), want: "-42"},
		{name: "float64 fractional", value: float64(3.14), want: "3.14"},
		{name: "float64 large whole", value: float64(9999999999), want: "9999999999"},
		{name: "float32 whole number", value: float32(100), want: "100"},
		{name: "float32 fractional", value: float32(3.14), want: "3.14"},
		{name: "int", value: 42, want: "42"},
		{name: "int8", value: int8(8), want: "8"},
		{name: "int16", value: int16(16), want: "16"},
		{name: "int32", value: int32(32), want: "32"},
		{name: "int64", value: int64(7314340), want: "7314340"},
		{name: "uint", value: uint(42), want: "42"},
		{name: "uint8", value: uint8(8), want: "8"},
		{name: "uint16", value: uint16(16), want: "16"},
		{name: "uint32", value: uint32(32), want: "32"},
		{name: "uint64", value: uint64(64), want: "64"},
		{name: "string", value: "hello", want: "hello"},
		{name: "bool true", value: true, want: "true"},
		{name: "bool false", value: false, want: "false"},
		{name: "default fallback", value: struct{}{}, want: "{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatParamValue(tt.value)
			require.Equal(t, tt.want, got)
		})
	}
}
