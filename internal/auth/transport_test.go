package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransport_RoundTrip(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &BearerTokenAuthClient{Token: "test-token"}
	require.NoError(t, client.New(), "New()")

	transport := &Transport{
		Base: http.DefaultTransport,
		Auth: client,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err, "RoundTrip()")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Bearer test-token", req.Header.Get(headerAuthorization))
}

func TestTransport_RoundTrip_Error(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "test-token"}
	require.NoError(t, client.New(), "New()")

	transport := &Transport{
		Base: http.DefaultTransport,
		Auth: client,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://nonexistent.example.com", nil)
	_, err := transport.RoundTrip(req)
	require.Error(t, err, "expected error for nonexistent host")
}

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "test-token"}
	require.NoError(t, client.New(), "New()")

	httpClient := newHTTPClient(client)
	require.NotNil(t, httpClient, "newHTTPClient() returned nil")

	transport, ok := httpClient.Transport.(*Transport)
	require.True(t, ok, "Transport type should be *Transport")
	assert.Equal(t, client, transport.Auth)
}
