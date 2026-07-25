package auth

import (
	"context"
	"net/http"
	"testing"
)

func TestBearerTokenAuthClient_Apply(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "my-bearer-token"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if v := req.Header.Get("Authorization"); v != "Bearer my-bearer-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer my-bearer-token")
	}
	if v := info.Headers["Authorization"]; v != "Bearer my-bearer-token" {
		t.Errorf("info.Headers[Authorization] = %q, want %q", v, "Bearer my-bearer-token")
	}
}

func TestBearerTokenAuthClient_Apply_EmptyToken(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: ""}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if v := req.Header.Get("Authorization"); v != "" {
		t.Errorf("Authorization = %q, want empty", v)
	}
	if info.Headers != nil {
		t.Errorf("info.Headers = %v, want nil", info.Headers)
	}
}

func TestBearerTokenAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_BEARER_TOKEN", "env-bearer-token")

	client := &BearerTokenAuthClient{Token: "$(TEST_BEARER_TOKEN)"}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if v := req.Header.Get("Authorization"); v != "Bearer env-bearer-token" {
		t.Errorf("Authorization = %q, want %q", v, "Bearer env-bearer-token")
	}
}
