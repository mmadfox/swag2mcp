package index

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

func newTestIndex(t *testing.T) *Index {
	t.Helper()
	idx, err := New()
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	return idx
}

func seedTestData(
	t *testing.T,
	idx *Index,
	domain string,
) (*types.Spec, *types.Collection, *types.Tag, *types.Endpoint) {
	t.Helper()

	specID := id.Domain(domain)
	specInfo := &types.Spec{
		ID:      specID,
		Domain:  domain,
		BaseURL: "https://api.example.com",
	}

	collectionID := id.Collection(specID, domain+"/collection")
	collectionInfo := &types.Collection{
		ID:     collectionID,
		SpecID: specID,
		Title:  "Test Collection",
	}

	tagID := id.Tag(specID, collectionID, "test-tag")
	tagInfo := &types.Tag{
		ID:           tagID,
		SpecID:       specID,
		CollectionID: collectionID,
		Name:         "test-tag",
	}

	endpointID := id.Method(specID, collectionID, tagID, "GET", "/test", "testOp")
	endpointInfo := &types.Endpoint{
		ID:           endpointID,
		SpecID:       specID,
		CollectionID: collectionID,
		TagID:        tagID,
		Tag:          "test-tag",
		Name:         "GET",
		Path:         "/test",
		Operation: &spec.Operation{
			ID:          "testOp",
			Summary:     "Test endpoint",
			Description: "A test endpoint",
		},
	}

	if err := idx.EnsureIndex(
		specInfo,
		[]*types.Collection{collectionInfo},
		[]*types.Tag{tagInfo},
		[]*types.Endpoint{endpointInfo},
	); err != nil {
		t.Fatalf("EnsureIndex() = %v", err)
	}
	idx.RefreshSearchReader()

	return specInfo, collectionInfo, tagInfo, endpointInfo
}

func TestNew(t *testing.T) {
	t.Parallel()

	idx, err := New()
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if idx == nil {
		t.Fatal("New() returned nil")
	}
	if idx.Size() != 0 {
		t.Errorf("Size() = %d, want 0", idx.Size())
	}
}

func TestEnsureIndex_DuplicateSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	err := idx.EnsureIndex(specInfo, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for duplicate spec")
	}
}

func TestEnsureIndex_CollectionMissingSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &types.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &types.Collection{
		ID:     "coll-1",
		SpecID: "spec-2",
		Title:  "Orphan Collection",
	}

	err := idx.EnsureIndex(specInfo, []*types.Collection{collectionInfo}, nil, nil)
	if err == nil {
		t.Fatal("expected error for collection with missing spec")
	}
}

func TestEnsureIndex_DuplicateCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	err := idx.EnsureIndex(specInfo, []*types.Collection{collectionInfo}, nil, nil)
	if err == nil {
		t.Fatal("expected error for duplicate collection")
	}
}

func TestEnsureIndex_TagMissingSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &types.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	tagInfo := &types.Tag{
		ID:     "tag-1",
		SpecID: "spec-2",
		Name:   "orphan-tag",
	}

	err := idx.EnsureIndex(specInfo, nil, []*types.Tag{tagInfo}, nil)
	if err == nil {
		t.Fatal("expected error for tag with missing spec")
	}
}

func TestEnsureIndex_TagMissingCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &types.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &types.Collection{ID: "coll-1", SpecID: "spec-1", Title: "Coll"}
	tagInfo := &types.Tag{
		ID:           "tag-1",
		SpecID:       "spec-1",
		CollectionID: "coll-2",
		Name:         "orphan-tag",
	}

	err := idx.EnsureIndex(specInfo, []*types.Collection{collectionInfo}, []*types.Tag{tagInfo}, nil)
	if err == nil {
		t.Fatal("expected error for tag with missing collection")
	}
}

func TestEnsureIndex_DuplicateTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	specInfo := &types.Spec{ID: tagInfo.SpecID, Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &types.Collection{ID: tagInfo.CollectionID, SpecID: tagInfo.SpecID, Title: "Coll"}

	err := idx.EnsureIndex(specInfo, []*types.Collection{collectionInfo}, []*types.Tag{tagInfo}, nil)
	if err == nil {
		t.Fatal("expected error for duplicate tag")
	}
}

func TestEnsureIndex_EndpointMissingSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &types.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &types.Collection{ID: "coll-1", SpecID: "spec-1", Title: "Coll"}
	tagInfo := &types.Tag{ID: "tag-1", SpecID: "spec-1", CollectionID: "coll-1", Name: "tag"}

	endpointInfo := &types.Endpoint{
		ID:           "ep-1",
		SpecID:       "spec-2",
		CollectionID: "coll-1",
		TagID:        "tag-1",
		Name:         "GET",
		Path:         "/test",
	}

	err := idx.EnsureIndex(
		specInfo,
		[]*types.Collection{collectionInfo},
		[]*types.Tag{tagInfo},
		[]*types.Endpoint{endpointInfo},
	)
	if err == nil {
		t.Fatal("expected error for endpoint with missing spec")
	}
}

func TestEnsureIndex_EndpointMissingCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &types.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	tagInfo := &types.Tag{ID: "tag-1", SpecID: "spec-1", CollectionID: "coll-1", Name: "tag"}

	endpointInfo := &types.Endpoint{
		ID:           "ep-1",
		SpecID:       "spec-1",
		CollectionID: "coll-2",
		TagID:        "tag-1",
		Name:         "GET",
		Path:         "/test",
	}

	err := idx.EnsureIndex(specInfo, nil, []*types.Tag{tagInfo}, []*types.Endpoint{endpointInfo})
	if err == nil {
		t.Fatal("expected error for endpoint with missing collection")
	}
}

func TestEnsureIndex_EndpointMissingTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &types.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &types.Collection{ID: "coll-1", SpecID: "spec-1", Title: "Coll"}

	endpointInfo := &types.Endpoint{
		ID:           "ep-1",
		SpecID:       "spec-1",
		CollectionID: "coll-1",
		TagID:        "tag-2",
		Name:         "GET",
		Path:         "/test",
	}

	err := idx.EnsureIndex(specInfo, []*types.Collection{collectionInfo}, nil, []*types.Endpoint{endpointInfo})
	if err == nil {
		t.Fatal("expected error for endpoint with missing tag")
	}
}

func TestEnsureIndex_DuplicateEndpoint(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, _, endpointInfo := seedTestData(t, idx, t.Name())

	specInfo := &types.Spec{ID: endpointInfo.SpecID, Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &types.Collection{ID: endpointInfo.CollectionID, SpecID: endpointInfo.SpecID, Title: "Coll"}
	tagInfo := &types.Tag{
		ID:           endpointInfo.TagID,
		SpecID:       endpointInfo.SpecID,
		CollectionID: endpointInfo.CollectionID,
		Name:         "tag",
	}

	err := idx.EnsureIndex(
		specInfo,
		[]*types.Collection{collectionInfo},
		[]*types.Tag{tagInfo},
		[]*types.Endpoint{endpointInfo},
	)
	if err == nil {
		t.Fatal("expected error for duplicate endpoint")
	}
}

func TestAllSpecs(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	specs := idx.AllSpecs()
	if len(specs) != 1 {
		t.Fatalf("AllSpecs() = %d, want 1", len(specs))
	}
	if specs[0].Domain != t.Name() {
		t.Errorf("Domain = %q, want %q", specs[0].Domain, t.Name())
	}
}

func TestAllSpecs_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specs := idx.AllSpecs()
	if len(specs) != 0 {
		t.Errorf("AllSpecs() = %d, want 0", len(specs))
	}
}

func TestSpecByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	got, err := idx.SpecByID(specInfo.ID)
	if err != nil {
		t.Fatalf("SpecByID() = %v", err)
	}
	if got.ID != specInfo.ID {
		t.Errorf("ID = %q, want %q", got.ID, specInfo.ID)
	}
}

func TestSpecByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.SpecByID("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectionByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	got, err := idx.CollectionByID(collectionInfo.ID)
	if err != nil {
		t.Fatalf("CollectionByID() = %v", err)
	}
	if got.ID != collectionInfo.ID {
		t.Errorf("ID = %q, want %q", got.ID, collectionInfo.ID)
	}
}

func TestCollectionByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.CollectionByID("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCollectionsBySpec_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	cols, err := idx.CollectionsBySpec(specInfo.ID)
	if err != nil {
		t.Fatalf("CollectionsBySpec() = %v", err)
	}
	if len(cols) != 1 {
		t.Fatalf("Collections = %d, want 1", len(cols))
	}
}

func TestCollectionsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.CollectionsBySpec("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	got, err := idx.TagByID(tagInfo.ID)
	if err != nil {
		t.Fatalf("TagByID() = %v", err)
	}
	if got.ID != tagInfo.ID {
		t.Errorf("ID = %q, want %q", got.ID, tagInfo.ID)
	}
}

func TestTagByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.TagByID("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagsByCollection_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	tags, err := idx.TagsByCollection(collectionInfo.ID)
	if err != nil {
		t.Fatalf("TagsByCollection() = %v", err)
	}
	if len(tags) != 1 {
		t.Fatalf("Tags = %d, want 1", len(tags))
	}
}

func TestTagsByCollection_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.TagsByCollection("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTagsBySpec_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	tags, err := idx.TagsBySpec(specInfo.ID)
	if err != nil {
		t.Fatalf("TagsBySpec() = %v", err)
	}
	if len(tags) != 1 {
		t.Fatalf("Tags = %d, want 1", len(tags))
	}
}

func TestTagsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.TagsBySpec("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsByTag_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	eps, err := idx.EndpointsByTag(tagInfo.ID)
	if err != nil {
		t.Fatalf("EndpointsByTag() = %v", err)
	}
	if len(eps) != 1 {
		t.Fatalf("Endpoints = %d, want 1", len(eps))
	}
}

func TestEndpointsByTag_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointsByTag("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointsBySpec_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	eps, err := idx.EndpointsBySpec(specInfo.ID)
	if err != nil {
		t.Fatalf("EndpointsBySpec() = %v", err)
	}
	if len(eps) != 1 {
		t.Fatalf("Endpoints = %d, want 1", len(eps))
	}
}

func TestEndpointsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointsBySpec("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointByCollection_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	eps, err := idx.EndpointByCollection(collectionInfo.ID)
	if err != nil {
		t.Fatalf("EndpointByCollection() = %v", err)
	}
	if len(eps) != 1 {
		t.Fatalf("Endpoints = %d, want 1", len(eps))
	}
}

func TestEndpointByCollection_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointByCollection("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEndpointByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, _, endpointInfo := seedTestData(t, idx, t.Name())

	got, err := idx.EndpointByID(endpointInfo.ID)
	if err != nil {
		t.Fatalf("EndpointByID() = %v", err)
	}
	if got.ID != endpointInfo.ID {
		t.Errorf("ID = %q, want %q", got.ID, endpointInfo.ID)
	}
}

func TestEndpointByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointByID("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.Search(context.Background(), "", 10)
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestSearch_ByMethod(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "method:GET", 10)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_ByTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "test", 10)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_ByPath(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "test", 10)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_BySummary(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "endpoint", 10)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}
}

func TestSearch_NoResults(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "zzzzzznonexistent", 10)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Results = %d, want 0", len(results))
	}
}

func TestSearch_Limit(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "test", 1)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) > 1 {
		t.Errorf("Results = %d, want <= 1", len(results))
	}
}

func TestSearch_MatchAll(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "*", 10)
	if err != nil {
		t.Fatalf("Search() = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 result for match all")
	}
}

func TestSize(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	if idx.Size() != 0 {
		t.Errorf("Size() = %d, want 0", idx.Size())
	}

	seedTestData(t, idx, t.Name())
	if idx.Size() != 1 {
		t.Errorf("Size() = %d, want 1", idx.Size())
	}
}

func TestIterateByEndpoints(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	count := 0
	for range idx.IterateByEndpoints() {
		count++
	}
	if count != 1 {
		t.Errorf("IterateByEndpoints count = %d, want 1", count)
	}
}

func TestIterateByEndpoints_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	count := 0
	for range idx.IterateByEndpoints() {
		count++
	}
	if count != 0 {
		t.Errorf("IterateByEndpoints count = %d, want 0", count)
	}
}

func TestIterateByTags(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	count := 0
	for range idx.IterateByTags() {
		count++
	}
	if count != 1 {
		t.Errorf("IterateByTags count = %d, want 1", count)
	}
}

func TestIterateByTags_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	count := 0
	for range idx.IterateByTags() {
		count++
	}
	if count != 0 {
		t.Errorf("IterateByTags count = %d, want 0", count)
	}
}

func TestIterateByCollections(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	count := 0
	for range idx.IterateByCollections() {
		count++
	}
	if count != 1 {
		t.Errorf("IterateByCollections count = %d, want 1", count)
	}
}

func TestIterateByCollections_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	count := 0
	for range idx.IterateByCollections() {
		count++
	}
	if count != 0 {
		t.Errorf("IterateByCollections count = %d, want 0", count)
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	if err := idx.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}
}

func TestRefreshSearchReader(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	idx.RefreshSearchReader()

	results, err := idx.Search(context.Background(), "*", 10)
	if err != nil {
		t.Fatalf("Search() after RefreshSearchReader = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Results = %d, want 0", len(results))
	}
}
