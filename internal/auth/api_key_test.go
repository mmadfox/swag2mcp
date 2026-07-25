package auth

import (
	"context"
	"net/http"
	"testing"
)

//nolint:gocognit
func TestAPIKeyAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("sets header when In is header", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "X-Api-Key",
			Value: "my-api-key-value",
			In:    "header",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.Header.Get("X-Api-Key"); v != "my-api-key-value" {
			t.Errorf("X-Api-Key = %q, want %q", v, "my-api-key-value")
		}
		if v := info.Headers["X-Api-Key"]; v != "my-api-key-value" {
			t.Errorf("info.Headers[X-Api-Key] = %q, want %q", v, "my-api-key-value")
		}
	})

	t.Run("sets query param when In is query", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "api_key",
			Value: "query-key-value",
			In:    "query",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.URL.Query().Get("api_key"); v != "query-key-value" {
			t.Errorf("query api_key = %q, want %q", v, "query-key-value")
		}
		if v := info.QueryParams["api_key"]; v != "query-key-value" {
			t.Errorf("info.QueryParams[api_key] = %q, want %q", v, "query-key-value")
		}
	})

	t.Run("defaults to header when In is unknown", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "X-Auth",
			Value: "fallback-value",
			In:    "unknown",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		if err := client.Apply(req, nil); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.Header.Get("X-Auth"); v != "fallback-value" {
			t.Errorf("X-Auth = %q, want %q", v, "fallback-value")
		}
	})

	t.Run("does not set empty value", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "X-Key",
			Value: "",
			In:    "header",
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.Header.Get("X-Key"); v != "" {
			t.Errorf("X-Key = %q, want empty", v)
		}
		if info.Headers != nil {
			t.Errorf("info.Headers = %v, want nil", info.Headers)
		}
	})
}
