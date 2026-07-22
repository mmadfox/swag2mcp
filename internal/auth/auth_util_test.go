package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

func TestSetAuthHeader_EmptyValue(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	setAuthHeader(req, nil, headerAuthorization, "")
	assert.Empty(t, req.Header.Get(headerAuthorization), "header should not be set for empty value")
}

func TestSetAuthHeader_WithInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	var info Info
	setAuthHeader(req, &info, headerAuthorization, "Bearer token")
	assert.Equal(t, "Bearer token", req.Header.Get(headerAuthorization))
	assert.Equal(t, "Bearer token", info.Headers[headerAuthorization])
}

func TestSetAuthHeader_NilInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	setAuthHeader(req, nil, "X-Key", "value")
	assert.Equal(t, "value", req.Header.Get("X-Key"))
}

func TestSetAuthQuery_EmptyValue(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	setAuthQuery(req, nil, "key", "")
	assert.Empty(t, req.URL.Query().Get("key"), "query param should not be set for empty value")
}

func TestSetAuthQuery_WithInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	setAuthQuery(req, &info, "api_key", "secret")
	assert.Equal(t, "secret", req.URL.Query().Get("api_key"))
	assert.Equal(t, "secret", info.QueryParams["api_key"])
}

func TestSetAuthQuery_NilInfo(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	setAuthQuery(req, nil, "key", "val")
	assert.Equal(t, "val", req.URL.Query().Get("key"))
}

func TestDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	cli, err := httpclient.NewDefault()
	require.NoError(t, err, "NewDefault()")
	require.NotNil(t, cli, "NewDefault() returned nil")
	assert.NotZero(t, cli.Timeout, "Timeout should be set")
}

func TestInfo_HeadersNil(t *testing.T) {
	t.Parallel()

	var info Info
	assert.Nil(t, info.Headers, "Headers should be nil initially")
	assert.Nil(t, info.QueryParams, "QueryParams should be nil initially")
}
