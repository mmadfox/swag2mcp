package index

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestIndex(t *testing.T) *Index {
	t.Helper()
	idx, err := New()
	require.NoError(t, err, "New()")
	return idx
}

func seedTestData(
	t *testing.T,
	idx *Index,
	domain string,
) (*model.Spec, *model.Collection, *model.Tag, *model.Endpoint) {
	t.Helper()

	specID := id.Domain(domain)
	specInfo := &model.Spec{
		ID:      specID,
		Domain:  domain,
		BaseURL: "https://api.example.com",
	}

	collectionID := id.Collection(specID, domain+"/collection")
	collectionInfo := &model.Collection{
		ID:     collectionID,
		SpecID: specID,
		Title:  "Test Collection",
	}

	tagID := id.Tag(specID, collectionID, "test-tag")
	tagInfo := &model.Tag{
		ID:           tagID,
		SpecID:       specID,
		CollectionID: collectionID,
		Name:         "test-tag",
	}

	endpointID := id.Method(specID, collectionID, tagID, "GET", "/test", "testOp")
	endpointInfo := &model.Endpoint{
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

	require.NoError(t, idx.EnsureIndex(
		specInfo,
		[]*model.Collection{collectionInfo},
		[]*model.Tag{tagInfo},
		[]*model.Endpoint{endpointInfo},
	), "EnsureIndex()")
	idx.RefreshSearchReader()

	return specInfo, collectionInfo, tagInfo, endpointInfo
}

func TestNew(t *testing.T) {
	t.Parallel()

	idx, err := New()
	require.NoError(t, err, "New()")
	require.NotNil(t, idx, "New() returned nil")
	assert.Equal(t, 0, idx.Size())
}

func TestEnsureIndex_DuplicateSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	err := idx.EnsureIndex(specInfo, nil, nil, nil)
	require.Error(t, err, "expected error for duplicate spec")
}

func TestEnsureIndex_CollectionMissingSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &model.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &model.Collection{
		ID:     "coll-1",
		SpecID: "spec-2",
		Title:  "Orphan Collection",
	}

	err := idx.EnsureIndex(specInfo, []*model.Collection{collectionInfo}, nil, nil)
	require.Error(t, err, "expected error for collection with missing spec")
}

func TestEnsureIndex_DuplicateCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	err := idx.EnsureIndex(specInfo, []*model.Collection{collectionInfo}, nil, nil)
	require.Error(t, err, "expected error for duplicate collection")
}

func TestEnsureIndex_TagMissingSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &model.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	tagInfo := &model.Tag{
		ID:     "tag-1",
		SpecID: "spec-2",
		Name:   "orphan-tag",
	}

	err := idx.EnsureIndex(specInfo, nil, []*model.Tag{tagInfo}, nil)
	require.Error(t, err, "expected error for tag with missing spec")
}

func TestEnsureIndex_TagMissingCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &model.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &model.Collection{ID: "coll-1", SpecID: "spec-1", Title: "Coll"}
	tagInfo := &model.Tag{
		ID:           "tag-1",
		SpecID:       "spec-1",
		CollectionID: "coll-2",
		Name:         "orphan-tag",
	}

	err := idx.EnsureIndex(specInfo, []*model.Collection{collectionInfo}, []*model.Tag{tagInfo}, nil)
	require.Error(t, err, "expected error for tag with missing collection")
}

func TestEnsureIndex_DuplicateTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	specInfo := &model.Spec{ID: tagInfo.SpecID, Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &model.Collection{ID: tagInfo.CollectionID, SpecID: tagInfo.SpecID, Title: "Coll"}

	err := idx.EnsureIndex(specInfo, []*model.Collection{collectionInfo}, []*model.Tag{tagInfo}, nil)
	require.Error(t, err, "expected error for duplicate tag")
}

func TestEnsureIndex_EndpointMissingSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &model.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &model.Collection{ID: "coll-1", SpecID: "spec-1", Title: "Coll"}
	tagInfo := &model.Tag{ID: "tag-1", SpecID: "spec-1", CollectionID: "coll-1", Name: "tag"}

	endpointInfo := &model.Endpoint{
		ID:           "ep-1",
		SpecID:       "spec-2",
		CollectionID: "coll-1",
		TagID:        "tag-1",
		Name:         "GET",
		Path:         "/test",
	}

	err := idx.EnsureIndex(
		specInfo,
		[]*model.Collection{collectionInfo},
		[]*model.Tag{tagInfo},
		[]*model.Endpoint{endpointInfo},
	)
	require.Error(t, err, "expected error for endpoint with missing spec")
}

func TestEnsureIndex_EndpointMissingCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &model.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	tagInfo := &model.Tag{ID: "tag-1", SpecID: "spec-1", CollectionID: "coll-1", Name: "tag"}

	endpointInfo := &model.Endpoint{
		ID:           "ep-1",
		SpecID:       "spec-1",
		CollectionID: "coll-2",
		TagID:        "tag-1",
		Name:         "GET",
		Path:         "/test",
	}

	err := idx.EnsureIndex(specInfo, nil, []*model.Tag{tagInfo}, []*model.Endpoint{endpointInfo})
	require.Error(t, err, "expected error for endpoint with missing collection")
}

func TestEnsureIndex_EndpointMissingTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo := &model.Spec{ID: "spec-1", Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &model.Collection{ID: "coll-1", SpecID: "spec-1", Title: "Coll"}

	endpointInfo := &model.Endpoint{
		ID:           "ep-1",
		SpecID:       "spec-1",
		CollectionID: "coll-1",
		TagID:        "tag-2",
		Name:         "GET",
		Path:         "/test",
	}

	err := idx.EnsureIndex(specInfo, []*model.Collection{collectionInfo}, nil, []*model.Endpoint{endpointInfo})
	require.Error(t, err, "expected error for endpoint with missing tag")
}

func TestEnsureIndex_DuplicateEndpoint(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, _, endpointInfo := seedTestData(t, idx, t.Name())

	specInfo := &model.Spec{ID: endpointInfo.SpecID, Domain: t.Name(), BaseURL: "https://example.com"}
	collectionInfo := &model.Collection{ID: endpointInfo.CollectionID, SpecID: endpointInfo.SpecID, Title: "Coll"}
	tagInfo := &model.Tag{
		ID:           endpointInfo.TagID,
		SpecID:       endpointInfo.SpecID,
		CollectionID: endpointInfo.CollectionID,
		Name:         "tag",
	}

	err := idx.EnsureIndex(
		specInfo,
		[]*model.Collection{collectionInfo},
		[]*model.Tag{tagInfo},
		[]*model.Endpoint{endpointInfo},
	)
	require.Error(t, err, "expected error for duplicate endpoint")
}

func TestAllSpecs(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	specs := idx.AllSpecs()
	require.Len(t, specs, 1)
	assert.Equal(t, t.Name(), specs[0].Domain)
}

func TestAllSpecs_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specs := idx.AllSpecs()
	assert.Len(t, specs, 0)
}

func TestSpecByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	got, err := idx.SpecByID(specInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, specInfo.ID, got.ID)
}

func TestSpecByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.SpecByID("nonexistent")
	require.Error(t, err)
}

func TestCollectionByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	got, err := idx.CollectionByID(collectionInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, collectionInfo.ID, got.ID)
}

func TestCollectionByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.CollectionByID("nonexistent")
	require.Error(t, err)
}

func TestCollectionsBySpec_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	cols, err := idx.CollectionsBySpec(specInfo.ID)
	require.NoError(t, err)
	require.Len(t, cols, 1)
}

func TestCollectionsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.CollectionsBySpec("nonexistent")
	require.Error(t, err)
}

func TestTagByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	got, err := idx.TagByID(tagInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, tagInfo.ID, got.ID)
}

func TestTagByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.TagByID("nonexistent")
	require.Error(t, err)
}

func TestTagsByCollection_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	tags, err := idx.TagsByCollection(collectionInfo.ID)
	require.NoError(t, err)
	require.Len(t, tags, 1)
}

func TestTagsByCollection_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.TagsByCollection("nonexistent")
	require.Error(t, err)
}

func TestTagsBySpec_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	tags, err := idx.TagsBySpec(specInfo.ID)
	require.NoError(t, err)
	require.Len(t, tags, 1)
}

func TestTagsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.TagsBySpec("nonexistent")
	require.Error(t, err)
}

func TestEndpointsByTag_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	eps, err := idx.EndpointsByTag(tagInfo.ID)
	require.NoError(t, err)
	require.Len(t, eps, 1)
}

func TestEndpointsByTag_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointsByTag("nonexistent")
	require.Error(t, err)
}

func TestEndpointsBySpec_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	eps, err := idx.EndpointsBySpec(specInfo.ID)
	require.NoError(t, err)
	require.Len(t, eps, 1)
}

func TestEndpointsBySpec_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointsBySpec("nonexistent")
	require.Error(t, err)
}

func TestEndpointByCollection_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	eps, err := idx.EndpointByCollection(collectionInfo.ID)
	require.NoError(t, err)
	require.Len(t, eps, 1)
}

func TestEndpointByCollection_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointByCollection("nonexistent")
	require.Error(t, err)
}

func TestEndpointByID_Success(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, _, endpointInfo := seedTestData(t, idx, t.Name())

	got, err := idx.EndpointByID(endpointInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, endpointInfo.ID, got.ID)
}

func TestEndpointByID_NotFound(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.EndpointByID("nonexistent")
	require.Error(t, err)
}

func TestSearch_EmptyQuery(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, err := idx.Search(context.Background(), "", 10)
	require.Error(t, err, "expected error for empty query")
}

func TestSearch_ByMethod(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "method:GET", 10)
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least 1 result")
}

func TestSearch_ByTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "test", 10)
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least 1 result")
}

func TestSearch_ByPath(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "test", 10)
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least 1 result")
}

func TestSearch_BySummary(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "endpoint", 10)
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least 1 result")
}

func TestSearch_NoResults(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "zzzzzznonexistent", 10)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestSearch_Limit(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "test", 1)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(results), 1)
}

func TestSearch_MatchAll(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	results, err := idx.Search(context.Background(), "*", 10)
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least 1 result for match all")
}

func TestSize(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	assert.Equal(t, 0, idx.Size())

	seedTestData(t, idx, t.Name())
	assert.Equal(t, 1, idx.Size())
}

func TestIterateByEndpoints(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	count := 0
	for range idx.IterateByEndpoints() {
		count++
	}
	assert.Equal(t, 1, count)
}

func TestIterateByEndpoints_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	count := 0
	for range idx.IterateByEndpoints() {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestIterateByTags(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	count := 0
	for range idx.IterateByTags() {
		count++
	}
	assert.Equal(t, 1, count)
}

func TestIterateByTags_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	count := 0
	for range idx.IterateByTags() {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestIterateByCollections(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	count := 0
	for range idx.IterateByCollections() {
		count++
	}
	assert.Equal(t, 1, count)
}

func TestIterateByCollections_Empty(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	count := 0
	for range idx.IterateByCollections() {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestClose(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	require.NoError(t, idx.Close())
}

func TestRefreshSearchReader(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	idx.RefreshSearchReader()

	results, err := idx.Search(context.Background(), "*", 10)
	require.NoError(t, err)
	assert.Len(t, results, 0)
}

func TestRemoveSpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	idx.RemoveSpec(specInfo.ID)

	_, err := idx.SpecByID(specInfo.ID)
	require.Error(t, err, "expected error after RemoveSpec")
}

func TestRemoveSpec_NonExistent(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	idx.RemoveSpec("nonexistent")
}

func TestRemoveCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	idx.RemoveCollection(collectionInfo.ID)

	_, err := idx.CollectionByID(collectionInfo.ID)
	require.Error(t, err, "expected error after RemoveCollection")
}

func TestRemoveCollection_NonExistent(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	idx.RemoveCollection("nonexistent")
}

func TestRemoveTag(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, _, tagInfo, _ := seedTestData(t, idx, t.Name())

	idx.RemoveTag(tagInfo.ID)

	_, err := idx.TagByID(tagInfo.ID)
	require.Error(t, err, "expected error after RemoveTag")
}

func TestRemoveTag_NonExistent(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	idx.RemoveTag("nonexistent")
}

func TestRemoveAllTags(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	seedTestData(t, idx, t.Name())

	idx.RemoveAllTags()

	_, err := idx.TagByID("any")
	require.Error(t, err, "expected error after RemoveAllTags")
}

func TestRemoveCollectionsBySpec(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	idx.RemoveCollectionsBySpec(specInfo.ID)

	_, err := idx.CollectionsBySpec(specInfo.ID)
	require.Error(t, err, "expected error after RemoveCollectionsBySpec")
}

func TestRemoveCollectionsBySpec_NonExistent(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	idx.RemoveCollectionsBySpec("nonexistent")
}

func TestAddCollection(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	specInfo, _, _, _ := seedTestData(t, idx, t.Name())

	newColl := &model.Collection{
		ID:     "new-collection-id",
		SpecID: specInfo.ID,
		Title:  "New Collection",
	}
	idx.AddCollection(newColl)

	got, err := idx.CollectionByID(newColl.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Collection", got.Title)
}

func TestAddCollection_Duplicate(t *testing.T) {
	t.Parallel()

	idx := newTestIndex(t)
	_, collectionInfo, _, _ := seedTestData(t, idx, t.Name())

	dup := &model.Collection{
		ID:     collectionInfo.ID,
		SpecID: collectionInfo.SpecID,
		Title:  "Duplicate",
	}
	idx.AddCollection(dup)

	got, err := idx.CollectionByID(collectionInfo.ID)
	require.NoError(t, err)
	assert.Equal(t, "Duplicate", got.Title)
}
