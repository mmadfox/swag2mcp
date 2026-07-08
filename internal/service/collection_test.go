package service

import (
	"context"
	"testing"
)

func TestCollectionsBySpec_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.CollectionsBySpec(context.Background(), CollectionsRequest{SpecID: specInfo.ID})
	if err != nil {
		t.Fatalf("CollectionsBySpec() = %v", err)
	}
	if resp.Spec.ID != specInfo.ID {
		t.Errorf("Spec.ID = %q, want %q", resp.Spec.ID, specInfo.ID)
	}
	if len(resp.Collections) != 1 {
		t.Fatalf("Collections = %d, want 1", len(resp.Collections))
	}
	if resp.Collections[0].Title != "Test Collection" {
		t.Errorf("Title = %q, want %q", resp.Collections[0].Title, "Test Collection")
	}
}

func TestCollectionsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.CollectionsBySpec(context.Background(), CollectionsRequest{SpecID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectionsBySpec_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.CollectionsBySpec(context.Background(), CollectionsRequest{SpecID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectionByID_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, collectionInfo, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: collectionInfo.ID})
	if err != nil {
		t.Fatalf("CollectionByID() = %v", err)
	}
	if resp.Collection.ID != collectionInfo.ID {
		t.Errorf("ID = %q, want %q", resp.Collection.ID, collectionInfo.ID)
	}
	if resp.Collection.Title != "Test Collection" {
		t.Errorf("Title = %q, want %q", resp.Collection.Title, "Test Collection")
	}
	if len(resp.Tags) != 1 {
		t.Errorf("Tags = %d, want 1", len(resp.Tags))
	}
}

func TestCollectionByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectionByID_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}
