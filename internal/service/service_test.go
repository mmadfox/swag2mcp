package service

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/workspace"
	"github.com/stretchr/testify/require"
)

func TestService_New(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)
	require.NotNil(t, svc)
	require.NotNil(t, svc.Workspace())
}

func TestService_NewWithOptions(t *testing.T) {
	t.Parallel()

	svc, err := New(
		WithVersion("test-version"),
		WithDisableLLMAuth(true),
		WithDumpDir("/tmp/dumps"),
		WithIndexNoFullText(),
	)
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestService_Specs(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	resp, err := svc.Specs(context.Background())
	require.NoError(t, err)
	require.Empty(t, resp.Specs)
}

func TestService_Search(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	resp, err := svc.Search(context.Background(), SearchRequest{Query: "test", Limit: 10})
	require.NoError(t, err)
	require.Empty(t, resp.Endpoints)
}

func TestService_Info(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	resp, err := svc.Info(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, resp.Workspace)
}

func TestService_Import_noSource(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.Import(context.Background(), ImportRequest{})
	require.Error(t, err)
}

func TestService_ResponseOutline_invalid(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.ResponseOutline(context.Background(), ResponseOutlineRequest{})
	require.Error(t, err)
}

func TestService_ResponseCompress_invalid(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.ResponseCompress(context.Background(), ResponseCompressRequest{})
	require.Error(t, err)
}

func TestService_ResponseSlice_invalid(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.ResponseSlice(context.Background(), ResponseSliceRequest{})
	require.Error(t, err)
}

func TestService_SpecByID_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.SpecByID(context.Background(), SpecByIDRequest{ID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_CollectionsBySpec_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.CollectionsBySpec(context.Background(), CollectionsRequest{SpecID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_CollectionByID_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.CollectionByID(context.Background(), CollectionByIDRequest{ID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_TagsByCollection_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.TagsByCollection(context.Background(), TagsByCollectionRequest{CollectionID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_TagByID_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.TagByID(context.Background(), TagByIDRequest{ID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_TagsBySpec_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.TagsBySpec(context.Background(), TagsBySpecRequest{SpecID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_EndpointsByTag_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.EndpointsByTag(context.Background(), EndpointsByTagRequest{TagID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_EndpointsByCollection_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.EndpointsByCollection(context.Background(), EndpointsByCollectionRequest{CollectionID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_EndpointsBySpec_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.EndpointsBySpec(context.Background(), EndpointsBySpecRequest{SpecID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_EndpointByID_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.EndpointByID(context.Background(), EndpointByIDRequest{ID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_Inspect_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.Inspect(context.Background(), InspectRequest{EndpointID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_Auth_notFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.Auth(context.Background(), AuthRequest{SpecID: "00000000000000000000000000000000"})
	require.Error(t, err)
}

func TestService_MakeToolDefinitions(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	defs, err := svc.MakeToolDefinitions()
	require.NoError(t, err)
	require.NotEmpty(t, defs.Tools)
}

func TestService_Invoke_validationError(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	_, err = svc.Invoke(context.Background(), InvokeRequest{})
	require.Error(t, err)
}

func TestService_Export_noConfig(t *testing.T) {
	t.Parallel()

	ws, wsErr := workspace.New(t.TempDir())
	require.NoError(t, wsErr)
	svc, err := New(WithWorkspace(ws))
	require.NoError(t, err)

	_, err = svc.Export(context.Background(), ExportRequest{})
	require.Error(t, err)
}
