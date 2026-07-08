package service

import (
	"context"
	"testing"
)

func TestTagsByCollection_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, collectionInfo, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.TagsByCollection(context.Background(), TagsByCollectionRequest{CollectionID: collectionInfo.ID})
	if err != nil {
		t.Fatalf("TagsByCollection() = %v", err)
	}
	if len(resp.Tags) != 1 {
		t.Fatalf("Tags = %d, want 1", len(resp.Tags))
	}
	if resp.Tags[0].Title != "test-tag" {
		t.Errorf("Title = %q, want %q", resp.Tags[0].Title, "test-tag")
	}
}

func TestTagsByCollection_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.TagsByCollection(context.Background(), TagsByCollectionRequest{CollectionID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagsByCollection_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.TagsByCollection(context.Background(), TagsByCollectionRequest{CollectionID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagByID_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, _, tagInfo, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.TagByID(context.Background(), TagByIDRequest{ID: tagInfo.ID})
	if err != nil {
		t.Fatalf("TagByID() = %v", err)
	}
	if resp.Tag.ID != tagInfo.ID {
		t.Errorf("ID = %q, want %q", resp.Tag.ID, tagInfo.ID)
	}
	if resp.Tag.Title != "test-tag" {
		t.Errorf("Title = %q, want %q", resp.Tag.Title, "test-tag")
	}
}

func TestTagByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.TagByID(context.Background(), TagByIDRequest{ID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagByID_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.TagByID(context.Background(), TagByIDRequest{ID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagsBySpec_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.TagsBySpec(context.Background(), TagsBySpecRequest{SpecID: specInfo.ID})
	if err != nil {
		t.Fatalf("TagsBySpec() = %v", err)
	}
	if len(resp.Tags) != 1 {
		t.Fatalf("Tags = %d, want 1", len(resp.Tags))
	}
}

func TestTagsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.TagsBySpec(context.Background(), TagsBySpecRequest{SpecID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagsBySpec_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.TagsBySpec(context.Background(), TagsBySpecRequest{SpecID: "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}
