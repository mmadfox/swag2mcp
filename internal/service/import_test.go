package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestImportService_Import_noSource(t *testing.T) {
	t.Parallel()

	svc := newImportService(NewMockWorkspaceOps(gomock.NewController(t)))
	_, err := svc.Import(context.Background(), ImportRequest{})
	require.Error(t, err)
}

func TestImportService_Import_sourceWithoutName(t *testing.T) {
	t.Parallel()

	svc := newImportService(NewMockWorkspaceOps(gomock.NewController(t)))
	_, err := svc.Import(context.Background(), ImportRequest{Source: "https://example.com/spec.yaml"})
	require.Error(t, err)
}

func TestImportService_importSingle_success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().DownloadSpec(gomock.Any(), "https://example.com/spec.yaml").Return([]byte("data"), nil)
	ws.EXPECT().SaveSpec("test.yaml", []byte("data")).Return("/specs/test.yaml", nil)

	svc := newImportService(ws)
	resp, err := svc.Import(context.Background(), ImportRequest{
		Source: "https://example.com/spec.yaml",
		Name:   "test.yaml",
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	require.Equal(t, "test.yaml", resp.Files[0].Name)
}

func TestImportService_importSingle_downloadError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().DownloadSpec(gomock.Any(), gomock.Any()).Return(nil, errors.New("download failed"))

	svc := newImportService(ws)
	_, err := svc.Import(context.Background(), ImportRequest{
		Source: "https://example.com/spec.yaml",
		Name:   "test.yaml",
	})
	require.Error(t, err)
}

func TestImportService_importSpecs_noConfigPath(t *testing.T) {
	t.Parallel()

	svc := newImportService(NewMockWorkspaceOps(gomock.NewController(t)))
	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter: []string{"api"},
	})
	require.Error(t, err)
}

func TestImportService_importSpecs_noMatch(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)

	svc := newImportService(ws)
	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"nonexistent"},
		ConfFilePath: "testdata/swag2mcp.yaml",
	})
	require.Error(t, err)
}

func TestImportService_importFromZip_invalid(t *testing.T) {
	t.Parallel()

	svc := newImportService(NewMockWorkspaceOps(gomock.NewController(t)))
	_, err := svc.Import(context.Background(), ImportRequest{
		ZipSource: "/nonexistent/backup.zip",
	})
	require.Error(t, err)
}
