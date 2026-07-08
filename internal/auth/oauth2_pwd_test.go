package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

//nolint:gocognit
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if reqCount.Load() != 1 {
			t.Errorf("expected 1 token request, got %d", reqCount.Load())
		}
		if v := req.Header.Get("Authorization"); v != "Bearer pwd-access-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer pwd-access-token")
		}
		if v := info.Headers["Authorization"]; v != "Bearer pwd-access-token" {
			t.Errorf("info.Headers[Authorization] = %q, want %q", v, "Bearer pwd-access-token")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		req2, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		if reqCount.Load() != 1 {
			t.Errorf("expected 1 token request (cached), got %d", reqCount.Load())
		}
		if v := req2.Header.Get("Authorization"); v != "Bearer cached-pwd-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer cached-pwd-token")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		req2, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		if reqCount.Load() != 2 {
			t.Errorf("expected 2 token requests (expired), got %d", reqCount.Load())
		}
		if v := req2.Header.Get("Authorization"); v != "Bearer pwd-token-2" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer pwd-token-2")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.Header.Get("Authorization"); v != "Bearer public-client-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer public-client-token")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req, nil); err != nil {
			t.Fatalf("Apply() = %v", err)
		}
		if v := req.Header.Get("Authorization"); v != "Bearer scoped-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer scoped-token")
		}
	})
}
