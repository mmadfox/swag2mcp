package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2ClientCredentialsAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("successful token fetch", func(t *testing.T) {
		t.Parallel()

		var reqCount atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCount.Add(1)
			if r.Method != http.MethodPost || r.URL.Path != "/token" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
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
			resp := oauth2TokenResponse{AccessToken: "access-token-123", TokenType: "Bearer", ExpiresIn: 3600}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2ClientCredentialsAuthClient{
			ClientID: "test-client", ClientSecret: "test-secret", TokenURL: srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, int32(1), reqCount.Load(), "expected 1 token request")
		assert.Equal(t, "Bearer access-token-123", req.Header.Get(headerAuthorization))
		assert.Equal(t, "Bearer access-token-123", info.Headers[headerAuthorization])
	})

	t.Run("caches token and reuses on second Apply", func(t *testing.T) {
		t.Parallel()

		var reqCount atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			reqCount.Add(1)
			resp := oauth2TokenResponse{AccessToken: "cached-token", TokenType: "Bearer", ExpiresIn: 3600}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2ClientCredentialsAuthClient{
			ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		assert.Equal(t, int32(1), reqCount.Load(), "expected 1 token request (cached)")
		assert.Equal(t, "Bearer cached-token", req2.Header.Get(headerAuthorization))
	})

	t.Run("refetches token after expiration", func(t *testing.T) {
		t.Parallel()

		var reqCount atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			reqCount.Add(1)
			resp := oauth2TokenResponse{
				AccessToken: fmt.Sprintf("token-%d", reqCount.Load()),
				TokenType:   "Bearer",
				ExpiresIn:   1,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2ClientCredentialsAuthClient{
			ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		time.Sleep(1100 * time.Millisecond)

		req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		assert.Equal(t, int32(2), reqCount.Load(), "expected 2 token requests (expired)")
		assert.Equal(t, "Bearer token-2", req2.Header.Get(headerAuthorization))
	})

	t.Run("returns error on non-200 response", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2ClientCredentialsAuthClient{
			ClientID: "bad", ClientSecret: "bad", TokenURL: srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error")
	})

	t.Run("returns error on empty access_token", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			resp := oauth2TokenResponse{AccessToken: "", TokenType: "Bearer", ExpiresIn: 3600}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2ClientCredentialsAuthClient{
			ClientID: "c", ClientSecret: "s", TokenURL: srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error for empty access_token")
	})
}

func TestOAuth2ClientCredentialsAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_OAUTH2_CLIENT_ID", "env-client")
	t.Setenv("TEST_OAUTH2_CLIENT_SECRET", "env-secret")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("client_id") != "env-client" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp := oauth2TokenResponse{AccessToken: "env-token", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2ClientCredentialsAuthClient{
		ClientID: "$(TEST_OAUTH2_CLIENT_ID)", ClientSecret: "$(TEST_OAUTH2_CLIENT_SECRET)",
		TokenURL: srv.URL + "/token",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer env-token", req.Header.Get(headerAuthorization))
}
