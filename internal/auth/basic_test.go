package auth

import (
	"context"
	"net/http"
	"testing"
)

func TestBasicAuthClient_Apply(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{
		Username: "alice",
		Password: "secret123",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("expected BasicAuth to be set")
	}
	if user != "alice" {
		t.Errorf("username = %q, want %q", user, "alice")
	}
	if pass != "secret123" {
		t.Errorf("password = %q, want %q", pass, "secret123")
	}

	if v := info.Headers["Authorization"]; v == "" {
		t.Error("info.Headers[Authorization] is empty")
	}
}

func TestBasicAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_BASIC_USER", "bob")
	t.Setenv("TEST_BASIC_PASS", "bobpass")

	client := &BasicAuthClient{
		Username: "$(TEST_BASIC_USER)",
		Password: "$(TEST_BASIC_PASS)",
	}
	if err := client.New(); err != nil {
		t.Fatalf("New() = %v", err)
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	if err := client.Apply(req, nil); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Fatal("expected BasicAuth to be set")
	}
	if user != "bob" {
		t.Errorf("username = %q, want %q", user, "bob")
	}
	if pass != "bobpass" {
		t.Errorf("password = %q, want %q", pass, "bobpass")
	}
}
