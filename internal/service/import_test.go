/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/yaml.v3"

	"github.com/mmadfox/swag2mcp/internal/config"
)

func writeTestConfig(t *testing.T, dir string, spec config.Spec) string {
	t.Helper()
	cfg := config.Config{Specs: []config.Spec{spec}}
	data, err := yaml.Marshal(cfg)
	require.NoError(t, err)
	path := filepath.Join(dir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(path, data, 0600))
	return path
}

func TestImportService_Import_noSource(t *testing.T) {
	t.Parallel()

	svc := newImportService(NewMockWorkspaceOps(gomock.NewController(t)))
	_, err := svc.Import(context.Background(), ImportRequest{})
	require.Error(t, err)
}

func TestImportService_Import_sourceWithoutName_withFilename(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().DownloadSpec(gomock.Any(), "https://example.com/spec.yaml").Return([]byte("data"), nil)
	ws.EXPECT().SaveSpec("spec.yaml", []byte("data")).Return("/specs/spec.yaml", nil)

	svc := newImportService(ws)
	resp, err := svc.Import(context.Background(), ImportRequest{Source: "https://example.com/spec.yaml"})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	require.False(t, resp.Files[0].Skipped)
}

func TestImportService_importSpecs_errorWhenLocalSpecPathMissing(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	require.NoError(t, os.MkdirAll(specsDir, 0750))

	cfgPath := writeTestConfig(t, tmpDir, config.Spec{
		Domain:   "meteo",
		LLMTitle: "Meteo API",
		BaseURL:  "https://api.meteo.com",
		Collections: []config.Collection{
			{LLMTitle: "Pet Store", Location: filepath.Join(specsDir, "nonexistent.yaml")},
		},
	})

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)

	ws.EXPECT().SpecsDir().Return(specsDir)
	ws.EXPECT().DownloadSpec(gomock.Any(), gomock.Any()).Times(0)
	ws.EXPECT().SaveOrUpdateSpec(gomock.Any(), gomock.Any()).Times(0)

	svc := newImportService(ws)
	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"meteo"},
		ConfFilePath: cfgPath,
	})
	require.Error(t, err)
}

func TestImportService_Import_sourceWithoutName_noFilename(t *testing.T) {
	t.Parallel()

	svc := newImportService(NewMockWorkspaceOps(gomock.NewController(t)))
	_, err := svc.Import(context.Background(), ImportRequest{Source: "https://example.com/"})
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

func TestImportService_importSingle_force(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().DownloadSpec(gomock.Any(), "https://example.com/spec.yaml").Return([]byte("data"), nil)
	ws.EXPECT().SaveOrUpdateSpec("test.yaml", []byte("data")).Return("/specs/test.yaml", nil)

	svc := newImportService(ws)
	resp, err := svc.Import(context.Background(), ImportRequest{
		Source: "https://example.com/spec.yaml",
		Name:   "test.yaml",
		Force:  true,
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

func TestImportService_importSpecs_skipWhenExistsAndUnreachable(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := writeTestConfig(t, tmpDir, config.Spec{
		Domain:   "meteo",
		LLMTitle: "Meteo API",
		BaseURL:  "https://api.meteo.com",
		Collections: []config.Collection{
			{LLMTitle: "Pet Store", Location: "https://example.com/petstore.yaml"},
		},
	})

	specPath := filepath.Join(tmpDir, "specs", "meteo-pet-store.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(specPath), 0750))
	require.NoError(t, os.WriteFile(specPath, []byte("existing"), 0600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)

	ws.EXPECT().SpecsDir().Return(filepath.Join(tmpDir, "specs")).AnyTimes()
	ws.EXPECT().SpecPath("meteo-pet-store.yaml").Return(specPath)
	ws.EXPECT().DownloadSpec(gomock.Any(), gomock.Any()).Return(nil, errors.New("unreachable"))
	ws.EXPECT().SaveOrUpdateSpec(gomock.Any(), gomock.Any()).Times(0)

	svc := newImportService(ws)
	resp, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"meteo"},
		ConfFilePath: cfgPath,
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	require.True(t, resp.Files[0].Skipped)
}

func TestImportService_importSpecs_updateWhenExistsAndReachable(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := writeTestConfig(t, tmpDir, config.Spec{
		Domain:   "meteo",
		LLMTitle: "Meteo API",
		BaseURL:  "https://api.meteo.com",
		Collections: []config.Collection{
			{LLMTitle: "Pet Store", Location: "https://example.com/petstore.yaml"},
		},
	})

	specPath := filepath.Join(tmpDir, "specs", "meteo-pet-store.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(specPath), 0750))
	require.NoError(t, os.WriteFile(specPath, []byte("existing"), 0600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)

	ws.EXPECT().SpecsDir().Return(filepath.Join(tmpDir, "specs")).AnyTimes()
	ws.EXPECT().SpecPath("meteo-pet-store.yaml").Return(specPath)
	ws.EXPECT().DownloadSpec(gomock.Any(), gomock.Any()).Return([]byte("updated"), nil)
	ws.EXPECT().SaveOrUpdateSpec("meteo-pet-store.yaml", []byte("updated")).Return(specPath, nil)

	svc := newImportService(ws)
	resp, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"meteo"},
		ConfFilePath: cfgPath,
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	require.False(t, resp.Files[0].Skipped)
}

func TestImportService_importSpecs_errorWhenNotExistsAndUnreachable(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfgPath := writeTestConfig(t, tmpDir, config.Spec{
		Domain:   "meteo",
		LLMTitle: "Meteo API",
		BaseURL:  "https://api.meteo.com",
		Collections: []config.Collection{
			{LLMTitle: "Pet Store", Location: "https://example.com/petstore.yaml"},
		},
	})

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)

	ws.EXPECT().SpecsDir().Return(filepath.Join(tmpDir, "specs")).AnyTimes()
	ws.EXPECT().SpecPath("meteo-pet-store.yaml").Return(filepath.Join(tmpDir, "nonexistent", "meteo-pet-store.yaml"))
	ws.EXPECT().DownloadSpec(gomock.Any(), gomock.Any()).Return(nil, errors.New("unreachable"))

	svc := newImportService(ws)
	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"meteo"},
		ConfFilePath: cfgPath,
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

func TestIsLocalSpecPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	require.NoError(t, os.MkdirAll(specsDir, 0750))

	tests := []struct {
		name     string
		location string
		want     bool
	}{
		{
			name:     "URL is not local",
			location: "https://example.com/spec.yaml",
			want:     false,
		},
		{
			name:     "relative path inside specs",
			location: filepath.Join("specs", "myspec.yaml"),
			want:     true,
		},
		{
			name:     "absolute path inside specs",
			location: filepath.Join(specsDir, "myspec.yaml"),
			want:     true,
		},
		{
			name:     "path outside specs",
			location: filepath.Join(tmpDir, "other", "spec.yaml"),
			want:     false,
		},
		{
			name:     "path in parent directory",
			location: "../outside.yaml",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isLocalSpecPath(tt.location, specsDir, tmpDir)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestImportService_importSpecs_skipsLocalSpecPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specsDir := filepath.Join(tmpDir, "specs")
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "meteo-pet-store.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("existing"), 0600))

	cfgPath := writeTestConfig(t, tmpDir, config.Spec{
		Domain:   "meteo",
		LLMTitle: "Meteo API",
		BaseURL:  "https://api.meteo.com",
		Collections: []config.Collection{
			{LLMTitle: "Pet Store", Location: specFile},
		},
	})

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)

	ws.EXPECT().SpecsDir().Return(specsDir)
	ws.EXPECT().DownloadSpec(gomock.Any(), gomock.Any()).Times(0)
	ws.EXPECT().SaveOrUpdateSpec(gomock.Any(), gomock.Any()).Times(0)

	svc := newImportService(ws)
	resp, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"meteo"},
		ConfFilePath: cfgPath,
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	require.True(t, resp.Files[0].Skipped)
}
