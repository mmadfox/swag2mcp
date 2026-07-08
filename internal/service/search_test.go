package service

import (
	"context"
	"net/http"
	"testing"
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
