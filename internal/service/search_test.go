package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
)

func TestSearch_ByMethod(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "method:GET", Limit: 10})
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(resp.Endpoints) == 0 {
		t.Fatal("expected at least 1 result")
	}
	if resp.Endpoints[0].Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", resp.Endpoints[0].Method, "GET")
	}
}

func TestSearch_ByTag(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(resp.Endpoints) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_ByPath(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(resp.Endpoints) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_BySummary(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "summary:\"Test endpoint\"", Limit: 10})
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(resp.Endpoints) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_NoResults(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "method:POST", Limit: 10})
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(resp.Endpoints) != 0 {
		t.Errorf("Endpoints = %d, want 0", len(resp.Endpoints))
	}
}

func TestSearch_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.Search(context.Background(), SearchRequest{Query: "", Limit: 0})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSearch_Limit(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 1})
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(resp.Endpoints) > 1 {
		t.Errorf("Endpoints = %d, want <= 1", len(resp.Endpoints))
	}
}

func TestMapEndpointsToSearchItems_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, _, _, endpointInfo := seedTestData(t, svc, t.Name())

	items, err := svc.mapEndpointsToSearchItems([]*model.Endpoint{endpointInfo})
	if err != nil {
		t.Fatalf("mapEndpointsToSearchItems() = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	if items[0].ID != endpointInfo.ID {
		t.Errorf("ID = %q, want %q", items[0].ID, endpointInfo.ID)
	}
	if items[0].SpecDomain == "" {
		t.Error("SpecDomain is empty")
	}
	if items[0].CollectionTitle == "" {
		t.Error("CollectionTitle is empty")
	}
	if items[0].TagName == "" {
		t.Error("TagName is empty")
	}
}

func TestMapEndpointsToSearchItems_Empty(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	items, err := svc.mapEndpointsToSearchItems(nil)
	if err != nil {
		t.Fatalf("mapEndpointsToSearchItems() = %v", err)
	}
	if len(items) != 0 {
		t.Errorf("items = %d, want 0", len(items))
	}
}

func TestMapEndpointsToSearchItems_OrphanSpec(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	orphanEndpoint := &model.Endpoint{
		ID:           "orphan",
		SpecID:       "00000000000000000000000000000000",
		CollectionID: "00000000000000000000000000000000",
		TagID:        "00000000000000000000000000000000",
	}

	_, err := svc.mapEndpointsToSearchItems([]*model.Endpoint{orphanEndpoint})
	if err == nil {
		t.Fatal("expected error for orphan spec")
	}
}

func TestMapEndpointsToSearchItems_OrphanCollection(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	orphanEndpoint := &model.Endpoint{
		ID:           "orphan",
		SpecID:       specInfo.ID,
		CollectionID: "00000000000000000000000000000000",
		TagID:        "00000000000000000000000000000000",
	}

	_, err := svc.mapEndpointsToSearchItems([]*model.Endpoint{orphanEndpoint})
	if err == nil {
		t.Fatal("expected error for orphan collection")
	}
}

func TestMapEndpointsToSearchItems_OrphanTag(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, collectionInfo, _, _ := seedTestData(t, svc, t.Name())

	orphanEndpoint := &model.Endpoint{
		ID:           "orphan",
		SpecID:       specInfo.ID,
		CollectionID: collectionInfo.ID,
		TagID:        "00000000000000000000000000000000",
	}

	_, err := svc.mapEndpointsToSearchItems([]*model.Endpoint{orphanEndpoint})
	if err == nil {
		t.Fatal("expected error for orphan tag")
	}
}

func TestSearch_OrphanEndpoints(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	// Remove all tags to make Search return orphan endpoints
	svc.index.RemoveAllTags()

	_, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	if err == nil {
		t.Fatal("expected error for orphan endpoints through Search")
	}
}
