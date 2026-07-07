package auth

import (
	"context"
	"net/http"
	"testing"
)

func TestNoAuthClient_Apply(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	req.Header.Set("X-Custom", "should-stay")

	var info Info
	if err := client.Apply(req, &info); err != nil {
		t.Fatalf("Apply() = %v", err)
	}

	if v := req.Header.Get("Authorization"); v != "" {
		t.Errorf("Authorization = %q, want empty", v)
	}
	if v := req.Header.Get("X-Custom"); v != "should-stay" {
		t.Errorf("X-Custom = %q, want %q", v, "should-stay")
	}
	if info.Headers != nil {
		t.Errorf("info.Headers = %v, want nil", info.Headers)
	}
	if info.QueryParams != nil {
		t.Errorf("info.QueryParams = %v, want nil", info.QueryParams)
	}
}
