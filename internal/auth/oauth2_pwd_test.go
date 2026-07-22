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

func TestOAuth2PasswordAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("successful token fetch with password grant", func(t *testing.T) {
		t.Parallel()

		var reqCount atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqCount.Add(1)

			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			if r.URL.Path != "/token" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("grant_type") != "password" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("username") != "testuser" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("password") != "testpass" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("client_id") != "test-client" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("client_secret") != "test-secret" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			resp := oauth2TokenResponse{
				AccessToken: "pwd-access-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2PasswordAuthClient{
			Username:     "testuser",
			Password:     "testpass",
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			TokenURL:     srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, int32(1), reqCount.Load(), "expected 1 token request")
		assert.Equal(t, "Bearer pwd-access-token", req.Header.Get(headerAuthorization))
		assert.Equal(t, "Bearer pwd-access-token", info.Headers[headerAuthorization])
	})

	t.Run("caches token and reuses on second Apply", func(t *testing.T) {
		t.Parallel()

		var reqCount atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			reqCount.Add(1)
			resp := oauth2TokenResponse{
				AccessToken: "cached-pwd-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2PasswordAuthClient{
			Username:     "u",
			Password:     "p",
			ClientID:     "c",
			ClientSecret: "s",
			TokenURL:     srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		req2, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		assert.Equal(t, int32(1), reqCount.Load(), "expected 1 token request (cached)")
		assert.Equal(t, "Bearer cached-pwd-token", req2.Header.Get(headerAuthorization))
	})

	t.Run("refetches token after expiration", func(t *testing.T) {
		t.Parallel()

		var reqCount atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			reqCount.Add(1)
			resp := oauth2TokenResponse{
				AccessToken: fmt.Sprintf("pwd-token-%d", reqCount.Load()),
				TokenType:   "Bearer",
				ExpiresIn:   1,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2PasswordAuthClient{
			Username:     "u",
			Password:     "p",
			ClientID:     "c",
			ClientSecret: "s",
			TokenURL:     srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		time.Sleep(1100 * time.Millisecond)

		req2, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		assert.Equal(t, int32(2), reqCount.Load(), "expected 2 token requests (expired)")
		assert.Equal(t, "Bearer pwd-token-2", req2.Header.Get(headerAuthorization))
	})

	t.Run("successful token fetch without client_secret (public client)", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("grant_type") != "password" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("username") != "testuser" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("password") != "testpass" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("client_id") != "public-client" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.Form.Get("client_secret") != "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			resp := oauth2TokenResponse{
				AccessToken: "public-client-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2PasswordAuthClient{
			Username: "testuser",
			Password: "testpass",
			ClientID: "public-client",
			TokenURL: srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "Bearer public-client-token", req.Header.Get(headerAuthorization))
	})

	t.Run("returns error on non-200 response", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2PasswordAuthClient{
			Username:     "u",
			Password:     "bad",
			ClientID:     "c",
			ClientSecret: "s",
			TokenURL:     srv.URL + "/token",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error")
	})

	t.Run("sends scopes when configured", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			if r.Form.Get("scope") != "read write" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := oauth2TokenResponse{AccessToken: "scoped-token", TokenType: "Bearer", ExpiresIn: 3600}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}))
		t.Cleanup(srv.Close)

		client := &OAuth2PasswordAuthClient{
			Username:     "u",
			Password:     "p",
			ClientID:     "c",
			ClientSecret: "s",
			TokenURL:     srv.URL + "/token",
			Scopes:       []string{"read", "write"},
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req, nil), "Apply()")
		assert.Equal(t, "Bearer scoped-token", req.Header.Get(headerAuthorization))
	})
}

func TestOAuth2PasswordAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *OAuth2PasswordAuthClient
		wantErr bool
	}{
		{
			name: "valid",
			client: &OAuth2PasswordAuthClient{
				Username: "u", Password: "p", ClientID: "cid",
				TokenURL: "https://example.com/token",
			},
			wantErr: false,
		},
		{name: "empty", client: &OAuth2PasswordAuthClient{}, wantErr: true},
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

func TestOAuth2PasswordAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	assert.Equal(t, OAuth2Password, client.Type())
}

func TestOAuth2PasswordAuthClient_New_EnvVars(t *testing.T) {
	t.Setenv("TEST_PWD_USER", "env-user")
	t.Setenv("TEST_PWD_PASS", "env-pass")
	t.Setenv("TEST_PWD_CID", "env-cid")
	t.Setenv("TEST_PWD_CS", "env-cs")

	client := &OAuth2PasswordAuthClient{
		Username: "$(TEST_PWD_USER)", Password: "$(TEST_PWD_PASS)",
		ClientID: "$(TEST_PWD_CID)", ClientSecret: "$(TEST_PWD_CS)",
		TokenURL: "https://example.com/token",
	}
	require.NoError(t, client.New(), "New()")
	assert.Equal(t, "env-user", client.Username)
	assert.Equal(t, "env-pass", client.Password)
}

func TestOAuth2PasswordAuthClient_Apply_DefaultExpiry(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := oauth2TokenResponse{AccessToken: "token", TokenType: "Bearer", ExpiresIn: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer token", req.Header.Get(headerAuthorization))
}

func TestOAuth2PasswordAuthClient_Apply_Scopes(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("scope") != "openid profile" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp := oauth2TokenResponse{AccessToken: "scoped-token", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
		Scopes:   []string{"openid", "profile"},
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")
	assert.Equal(t, "Bearer scoped-token", req.Header.Get(headerAuthorization))
}

func TestOAuth2PasswordAuthClient_Apply_EmptyAccessToken(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := oauth2TokenResponse{AccessToken: "", TokenType: "Bearer", ExpiresIn: 3600}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for empty access_token")
}

func TestOAuth2PasswordAuthClient_Apply_Non200(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	client := &OAuth2PasswordAuthClient{
		Username: "u", Password: "p", ClientID: "c",
		TokenURL: srv.URL + "/token",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	err := client.Apply(req, nil)
	require.Error(t, err, "expected error for non-200")
}

func TestOAuth2PasswordAuthClient_SetTokenURL(t *testing.T) {
	t.Parallel()

	client := &OAuth2PasswordAuthClient{}
	client.SetTokenURL("http://localhost:9090/token")
	assert.Equal(t, "http://localhost:9090/token", client.TokenURL)
}
