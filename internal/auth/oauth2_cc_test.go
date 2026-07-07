package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

//nolint:gocognit
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if reqCount.Load() != 1 {
			t.Errorf("expected 1 token request, got %d", reqCount.Load())
		}
		if v := req.Header.Get("Authorization"); v != "Bearer access-token-123" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer access-token-123")
		}
		if v := info.Headers["Authorization"]; v != "Bearer access-token-123" {
			t.Errorf("info.Headers[Authorization] = %q, want %q", v, "Bearer access-token-123")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		if reqCount.Load() != 1 {
			t.Errorf("expected 1 token request (cached), got %d", reqCount.Load())
		}
		if v := req2.Header.Get("Authorization"); v != "Bearer cached-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer cached-token")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		time.Sleep(1100 * time.Millisecond)

		req2, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		if reqCount.Load() != 2 {
			t.Errorf("expected 2 token requests (expired), got %d", reqCount.Load())
		}
		if v := req2.Header.Get("Authorization"); v != "Bearer token-2" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer token-2")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error, got nil")
		}
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
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error for empty access_token, got nil")
		}
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
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}
	if v := req.Header.Get("Authorization"); v != "Bearer env-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer env-token")
	}
}
