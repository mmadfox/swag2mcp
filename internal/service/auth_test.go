package service

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

func TestAuth_Disabled(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, WithDisableLLMAuth(true))
	seedTestData(t, svc, t.Name())

	resp, err := svc.Auth(context.Background(), AuthRequest{SpecID: "00000000000000000000000000000000"})
	if err != nil {
		t.Fatalf("Auth() = %v", err)
	}
	if resp.Token != "" {
		t.Errorf("Token = %q, want empty", resp.Token)
	}
}

func TestAuth_NoAuthOnSpec(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.Auth(context.Background(), AuthRequest{SpecID: specInfo.ID})
	if err != nil {
		t.Fatalf("Auth() = %v", err)
	}
	if resp.Token != "" {
		t.Errorf("Token = %q, want empty", resp.Token)
	}
}

func TestAuth_WithBearer(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	authenticator := &auth.BearerTokenAuthClient{Token: "test-token"}
	if err := authenticator.New(); err != nil {
		t.Fatalf("authenticator.New() = %v", err)
	}
	specInfo, _, _, _ := seedTestDataWithAuth(t, svc, t.Name(), authenticator)

	resp, err := svc.Auth(context.Background(), AuthRequest{SpecID: specInfo.ID})
	if err != nil {
		t.Fatalf("Auth() = %v", err)
	}
	if resp.Token == "" {
		t.Fatal("Token is empty, expected a bearer token")
	}
}

func TestAuth_SpecNotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.Auth(context.Background(), AuthRequest{SpecID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}
