package service

import (
	"context"
	"net/http"
	"testing"
)

func TestEndpointsByTag_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, _, tagInfo, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.EndpointsByTag(context.Background(), EndpointsByTagRequest{TagID: tagInfo.ID})
	if err != nil {
		t.Fatalf("EndpointsByTag() = %v", err)
	}
	if len(resp.Endpoints) != 1 {
		t.Fatalf("Endpoints = %d, want 1", len(resp.Endpoints))
	}
	if resp.Endpoints[0].Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", resp.Endpoints[0].Method, http.MethodGet)
	}
	if resp.Endpoints[0].Path != "/test" {
		t.Errorf("Path = %q, want %q", resp.Endpoints[0].Path, "/test")
	}
}

func TestEndpointsByTag_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointsByTag(context.Background(), EndpointsByTagRequest{TagID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsByTag_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointsByTag(context.Background(), EndpointsByTagRequest{TagID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsByCollection_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, collectionInfo, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.EndpointsByCollection(context.Background(), EndpointsByCollectionRequest{CollectionID: collectionInfo.ID})
	if err != nil {
		t.Fatalf("EndpointsByCollection() = %v", err)
	}
	if len(resp.Endpoints) != 1 {
		t.Fatalf("Endpoints = %d, want 1", len(resp.Endpoints))
	}
	if resp.Endpoints[0].Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", resp.Endpoints[0].Method, http.MethodGet)
	}
}

func TestEndpointsByCollection_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointsByCollection(context.Background(), EndpointsByCollectionRequest{CollectionID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsByCollection_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointsByCollection(context.Background(), EndpointsByCollectionRequest{CollectionID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsBySpec_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.EndpointsBySpec(context.Background(), EndpointsBySpecRequest{SpecID: specInfo.ID})
	if err != nil {
		t.Fatalf("EndpointsBySpec() = %v", err)
	}
	if len(resp.Endpoints) != 1 {
		t.Fatalf("Endpoints = %d, want 1", len(resp.Endpoints))
	}
}

func TestEndpointsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointsBySpec(context.Background(), EndpointsBySpecRequest{SpecID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsBySpec_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointsBySpec(context.Background(), EndpointsBySpecRequest{SpecID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointByID_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, _, _, endpointInfo := seedTestData(t, svc, t.Name())

	resp, err := svc.EndpointByID(context.Background(), EndpointByIDRequest{ID: endpointInfo.ID})
	if err != nil {
		t.Fatalf("EndpointByID() = %v", err)
	}
	if resp.Endpoint.ID != endpointInfo.ID {
		t.Errorf("ID = %q, want %q", resp.Endpoint.ID, endpointInfo.ID)
	}
	if resp.Endpoint.Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", resp.Endpoint.Method, http.MethodGet)
	}
	if resp.Endpoint.Path != "/test" {
		t.Errorf("Path = %q, want %q", resp.Endpoint.Path, "/test")
	}
}

func TestEndpointByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointByID(context.Background(), EndpointByIDRequest{ID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointByID_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.EndpointByID(context.Background(), EndpointByIDRequest{ID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}
