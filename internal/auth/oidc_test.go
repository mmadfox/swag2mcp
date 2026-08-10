/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newOIDCTestServer starts an httptest server that serves both the OIDC
// discovery document and the token endpoint.
func newOIDCTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	var tokenCalls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case oidcDiscoveryPath:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(oidcConfig{TokenEndpoint: "http://" + r.Host + "/token"})
		case "/token":
			tokenCalls.Add(1)
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_ = r.ParseForm()
			if r.Form.Get("grant_type") != "client_credentials" ||
				r.Form.Get("client_id") != "test-client" ||
				r.Form.Get("client_secret") != "test-secret" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := oauth2TokenResponse{AccessToken: "oidc-token-123", TokenType: "Bearer", ExpiresIn: 3600}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestOIDCAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &OIDCAuthClient{}
	assert.Equal(t, OpenIDConnect, client.Type())
}

func TestOIDCAuthClient_Apply(t *testing.T) {
	t.Parallel()

	srv := newOIDCTestServer(t)
	client := &OIDCAuthClient{
		Issuer:       srv.URL,
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Equal(t, "Bearer oidc-token-123", req.Header.Get(headerAuthorization))
	assert.Equal(t, "Bearer oidc-token-123", info.Headers[headerAuthorization])
}

func TestOIDCAuthClient_Apply_CachesToken(t *testing.T) {
	t.Parallel()

	srv := newOIDCTestServer(t)
	client := &OIDCAuthClient{
		Issuer:       srv.URL,
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	require.NoError(t, client.New(), "New()")

	req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req1, nil), "first Apply()")

	req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req2, nil), "second Apply()")

	assert.Equal(t, "Bearer oidc-token-123", req2.Header.Get(headerAuthorization))
}

func TestOIDCAuthClient_Apply_DiscoveryError(t *testing.T) {
	t.Parallel()

	// A server that returns a non-2xx for the discovery document.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	client := &OIDCAuthClient{
		Issuer:       srv.URL,
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "openid-connect")
}

func TestOIDCAuthClient_Apply_MissingTokenEndpoint(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)

	client := &OIDCAuthClient{
		Issuer:       srv.URL,
		ClientID:     "test-client",
		ClientSecret: "test-secret",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token_endpoint")
}

func TestOIDCAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *OIDCAuthClient
		wantErr bool
	}{
		{
			name:   "valid",
			client: &OIDCAuthClient{Issuer: "https://issuer.example.com", ClientID: "id", ClientSecret: "secret"},
		},
		{name: "missing issuer", client: &OIDCAuthClient{ClientID: "id", ClientSecret: "secret"}, wantErr: true},
		{
			name:    "missing client id",
			client:  &OIDCAuthClient{Issuer: "https://issuer.example.com", ClientSecret: "secret"},
			wantErr: true,
		},
		{
			name:    "missing secret",
			client:  &OIDCAuthClient{Issuer: "https://issuer.example.com", ClientID: "id"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.client.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestOIDCAuthClient_QueryParamNames(t *testing.T) {
	t.Parallel()

	client := &OIDCAuthClient{}
	assert.Nil(t, client.QueryParamNames())
}
