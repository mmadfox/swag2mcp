package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmadfox/swag2mcp/internal/model"
)

func TestSearch_ByMethod(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "method:GET", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Endpoints)
	assert.Equal(t, http.MethodGet, resp.Endpoints[0].Method)
}

func TestSearch_ByTag(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Endpoints)
}

func TestSearch_ByPath(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Endpoints)
}

func TestSearch_BySummary(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "summary:\"Test endpoint\"", Limit: 10})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Endpoints)
}

func TestSearch_NoResults(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "method:POST", Limit: 10})
	require.NoError(t, err)
	assert.Empty(t, resp.Endpoints)
}

func TestSearch_ValidationError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, err := svc.Search(context.Background(), SearchRequest{Query: "", Limit: 0})
	require.Error(t, err)
}

func TestSearch_Limit(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 1})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(resp.Endpoints), 1)
}

func TestSearch_IndexError(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.Search(ctx, SearchRequest{Query: "test", Limit: 10})
	require.Error(t, err)
}

func TestSearch_OrphanEndpoints(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	svc.index.RemoveAllTags()

	_, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	require.Error(t, err)
}

func TestMapEndpointsToSearchItems_Success(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	_, _, _, endpointInfo := seedTestData(t, svc, t.Name())

	items, err := svc.mapEndpointsToSearchItems([]*model.Endpoint{endpointInfo})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, endpointInfo.ID, items[0].ID)
	assert.NotEmpty(t, items[0].SpecDomain)
	assert.NotEmpty(t, items[0].CollectionTitle)
	assert.NotEmpty(t, items[0].TagName)
}

func TestMapEndpointsToSearchItems_Empty(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	items, err := svc.mapEndpointsToSearchItems(nil)
	require.NoError(t, err)
	assert.Empty(t, items)
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
	require.Error(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
}
