package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
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

func TestInspect_OrphanSpec(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	orphanIdx, idxErr := index.New()
	if idxErr != nil {
		t.Fatalf("index.New() = %v", idxErr)
	}
	orphanEndpoint := &model.Endpoint{
		ID:           "00000000000000000000000000000001",
		SpecID:       "00000000000000000000000000000000",
		CollectionID: "00000000000000000000000000000002",
		TagID:        "00000000000000000000000000000003",
		Name:         "GET",
		Path:         "/orphan",
		Operation:    &spec.Operation{ID: "orphanOp"},
	}
	if idxErr = orphanIdx.EnsureIndex(
		&model.Spec{ID: "00000000000000000000000000000000", Domain: "orphan"},
		[]*model.Collection{{ID: "00000000000000000000000000000002", SpecID: "00000000000000000000000000000000"}},
		[]*model.Tag{{ID: "00000000000000000000000000000003", SpecID: "00000000000000000000000000000000", CollectionID: "00000000000000000000000000000002"}},
		[]*model.Endpoint{orphanEndpoint},
	); idxErr != nil {
		t.Fatalf("EnsureIndex() = %v", idxErr)
	}

	svc.index = orphanIdx
	orphanIdx.RemoveSpec("00000000000000000000000000000000")

	_, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: orphanEndpoint.ID})
	if err == nil {
		t.Fatal("expected error for orphan spec")
	}
}

func TestInspect_OrphanCollection(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	orphanIdx, idxErr := index.New()
	if idxErr != nil {
		t.Fatalf("index.New() = %v", idxErr)
	}
	orphanEndpoint := &model.Endpoint{
		ID:           "00000000000000000000000000000001",
		SpecID:       "00000000000000000000000000000000",
		CollectionID: "00000000000000000000000000000002",
		TagID:        "00000000000000000000000000000003",
		Name:         "GET",
		Path:         "/orphan",
		Operation:    &spec.Operation{ID: "orphanOp"},
	}
	if idxErr = orphanIdx.EnsureIndex(
		&model.Spec{ID: "00000000000000000000000000000000", Domain: "orphan"},
		[]*model.Collection{{ID: "00000000000000000000000000000002", SpecID: "00000000000000000000000000000000"}},
		[]*model.Tag{{ID: "00000000000000000000000000000003", SpecID: "00000000000000000000000000000000", CollectionID: "00000000000000000000000000000002"}},
		[]*model.Endpoint{orphanEndpoint},
	); idxErr != nil {
		t.Fatalf("EnsureIndex() = %v", idxErr)
	}

	svc.index = orphanIdx
	orphanIdx.RemoveCollection("00000000000000000000000000000002")

	_, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: orphanEndpoint.ID})
	if err == nil {
		t.Fatal("expected error for orphan collection")
	}
}

func TestInspect_CollectionBaseURL(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, collectionInfo, _, endpointInfo := seedTestData(t, svc, t.Name())

	collectionInfo.BaseURL = "https://collection.example.com"

	resp, err := svc.Inspect(context.Background(), InspectRequest{EndpointID: endpointInfo.ID})
	if err != nil {
		t.Fatalf("Inspect() = %v", err)
	}
	if resp.BaseURL != "https://collection.example.com" {
		t.Errorf("BaseURL = %q, want %q", resp.BaseURL, "https://collection.example.com")
	}
	if resp.FullURL != "https://collection.example.com/test" {
		t.Errorf("FullURL = %q, want %q", resp.FullURL, "https://collection.example.com/test")
	}
	// Verify spec BaseURL is different
	if specInfo.BaseURL == "https://collection.example.com" {
		t.Error("spec BaseURL should not be the same as collection BaseURL")
	}
}
