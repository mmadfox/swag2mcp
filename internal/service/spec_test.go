package service

import (
	"context"
	"testing"
)

func TestSpecs_Empty(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	resp, err := svc.Specs(context.Background())
	if err != nil {
		t.Fatalf("Specs() = %v", err)
	}
	if len(resp.Specs) != 0 {
		t.Errorf("Specs = %d, want 0", len(resp.Specs))
	}
}

func TestSpecs_WithData(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Specs(context.Background())
	if err != nil {
		t.Fatalf("Specs() = %v", err)
	}
	if len(resp.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(resp.Specs))
	}
	if resp.Specs[0].Domain != t.Name() {
		t.Errorf("Domain = %q, want %q", resp.Specs[0].Domain, t.Name())
	}
}

func TestSpecByID_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	resp, err := svc.SpecByID(context.Background(), SpecByIDRequest{ID: specInfo.ID})
	if err != nil {
		t.Fatalf("SpecByID() = %v", err)
	}
	if resp.Spec.ID != specInfo.ID {
		t.Errorf("ID = %q, want %q", resp.Spec.ID, specInfo.ID)
	}
	if resp.Spec.Domain != t.Name() {
		t.Errorf("Domain = %q, want %q", resp.Spec.Domain, t.Name())
	}
	if len(resp.Collections) != 1 {
		t.Errorf("Collections = %d, want 1", len(resp.Collections))
	}
}

func TestSpecByID_NotFound(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.SpecByID(context.Background(), SpecByIDRequest{ID: "00000000000000000000000000000000"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSpecByID_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.SpecByID(context.Background(), SpecByIDRequest{ID: "invalid"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSpecByID_CollectionsError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	specInfo, _, _, _ := seedTestData(t, svc, t.Name())

	// Remove collections to trigger the else branch
	svc.index.RemoveCollectionsBySpec(specInfo.ID)

	resp, err := svc.SpecByID(context.Background(), SpecByIDRequest{ID: specInfo.ID})
	if err != nil {
		t.Fatalf("SpecByID() = %v", err)
	}
	if resp.Spec.ID != specInfo.ID {
		t.Errorf("ID = %q, want %q", resp.Spec.ID, specInfo.ID)
	}
	if len(resp.Collections) != 0 {
		t.Errorf("Collections = %d, want 0", len(resp.Collections))
	}
}
