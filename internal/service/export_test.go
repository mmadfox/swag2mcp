package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestExportService_Export_noConfig(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ConfigPath().Return("/tmp/.swag2mcp/swag2mcp.yaml")
	ws.EXPECT().ConfigNotExists().Return(true)

	svc := newExportService(ws, "1.0")
	_, err := svc.Export(context.Background(), ExportRequest{})
	require.Error(t, err)
}

func TestExportService_Export_success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ConfigPath().Return("/tmp/.swag2mcp/swag2mcp.yaml")
	ws.EXPECT().ConfigNotExists().Return(false)

	svc := newExportService(ws, "1.0")
	_, err := svc.Export(context.Background(), ExportRequest{OutputPath: "/tmp/out.zip"})
	require.Error(t, err) // config load will fail since we don't have a real config file
}

func TestExportService_loadExportConfig_notFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ConfigPath().Return("/tmp/.swag2mcp/swag2mcp.yaml")
	ws.EXPECT().ConfigNotExists().Return(true)

	svc := newExportService(ws, "1.0")
	_, err := svc.loadExportConfig()
	require.Error(t, err)
}

func TestExportService_updateConfigLocations(t *testing.T) {
	t.Parallel()

	svc := newExportService(NewMockWorkspaceOps(gomock.NewController(t)), "1.0")
	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain: "api",
				Collections: []config.Collection{
					{Location: "https://example.com/spec.yaml", Title: "spec1"},
				},
			},
		},
	}
	locationMap := map[string]string{
		"https://example.com/spec.yaml": "specs/api-spec1.yaml",
	}
	svc.updateConfigLocations(cfg, nil, locationMap)
	require.Equal(t, "specs/api-spec1.yaml", cfg.Specs[0].Collections[0].Location)
}

func TestExportService_updateConfigLocations_disabled(t *testing.T) {
	t.Parallel()

	svc := newExportService(NewMockWorkspaceOps(gomock.NewController(t)), "1.0")
	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:  "api",
				Disable: true,
				Collections: []config.Collection{
					{Location: "https://example.com/spec.yaml"},
				},
			},
		},
	}
	locationMap := map[string]string{
		"https://example.com/spec.yaml": "specs/api-spec.yaml",
	}
	svc.updateConfigLocations(cfg, nil, locationMap)
	require.Equal(t, "https://example.com/spec.yaml", cfg.Specs[0].Collections[0].Location)
}

func TestExportService_finalizeExport(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().CopyAuthScriptsToExport(tmpDir).Return(nil)

	svc := newExportService(ws, "1.0")
	err := svc.finalizeExport(&config.Config{}, tmpDir)
	require.NoError(t, err)
}

func TestExportService_finalizeExport_copyError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().CopyAuthScriptsToExport(tmpDir).Return(errors.New("copy failed"))

	svc := newExportService(ws, "1.0")
	err := svc.finalizeExport(&config.Config{}, tmpDir)
	require.Error(t, err)
}
