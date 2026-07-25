package service

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/types"
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

func TestCollectionByID_OrphanSpec(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	orphanIdx, idxErr := index.New()
	if idxErr != nil {
		t.Fatalf("index.New() = %v", idxErr)
	}
	if idxErr = orphanIdx.EnsureIndex(
		&types.Spec{ID: "00000000000000000000000000000000", Domain: "orphan"},
		[]*types.Collection{{ID: "00000000000000000000000000000001", SpecID: "00000000000000000000000000000000"}},
		[]*types.Tag{},
		[]*types.Endpoint{},
	); idxErr != nil {
		t.Fatalf("EnsureIndex() = %v", idxErr)
	}

	svc.index = orphanIdx
	orphanIdx.RemoveSpec("00000000000000000000000000000000")

	_, err := svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: "00000000000000000000000000000001"})
	if err == nil {
		t.Fatal("expected error for orphan spec")
	}
}

func TestCollectionByID_TagsError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, collectionInfo, _, _ := seedTestData(t, svc, t.Name())

	// Remove tags to trigger the else branch (TagsByCollection returns error)
	svc.index.RemoveAllTags()

	resp, err := svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: collectionInfo.ID})
	if err != nil {
		t.Fatalf("CollectionByID() = %v", err)
	}
	if len(resp.Tags) != 0 {
		t.Errorf("Tags = %d, want 0", len(resp.Tags))
	}
}
