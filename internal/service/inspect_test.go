package service

import (
	"context"
	"net/http"
	"testing"
)

func TestInspect_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, _, _, endpointInfo := seedTestData(t, svc, t.Name())

	resp, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: endpointInfo.ID})
	if err != nil {
		t.Fatalf("Inspect() = %v", err)
	}
	if resp.ID != endpointInfo.ID {
		t.Errorf("ID = %q, want %q", resp.ID, endpointInfo.ID)
	}
	if resp.Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", resp.Method, "GET")
	}
	if resp.Path != "/test" {
		t.Errorf("Path = %q, want %q", resp.Path, "/test")
	}
	if resp.SpecDomain != t.Name() {
		t.Errorf("SpecDomain = %q, want %q", resp.SpecDomain, t.Name())
	}
	if resp.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q, want %q", resp.BaseURL, "https://api.example.com")
	}
	if resp.FullURL != "https://api.example.com/test" {
		t.Errorf("FullURL = %q, want %q", resp.FullURL, "https://api.example.com/test")
	}
	if resp.Operation == nil {
		t.Fatal("Operation is nil")
	}
	if resp.Operation.Summary != "Test endpoint" {
		t.Errorf("Summary = %q, want %q", resp.Operation.Summary, "Test endpoint")
	}
}

func TestInspect_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInspect_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}
