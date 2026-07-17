package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestExport_Success(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, resp.FileCount)
	assert.Equal(t, outputPath, resp.OutputPath)
	_, statErr := os.Stat(outputPath)
	assert.NoError(t, statErr, "zip file was not created")
}

func TestExport_WithSpecFilter(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
  - domain: store
    llm_title: Store API
    base_url: https://api.store.com
    collections:
      - title: Products
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
		SpecFilter: []string{"petstore"},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, resp.FileCount)
}

func TestExport_NoConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	require.Error(t, err)
}

func TestExport_NoCollections(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections: []
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	require.Error(t, err)
}

func TestExport_DefaultOutputPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(resp.OutputPath, ".zip"))
	_, statErr := os.Stat(outputPath)
	assert.NoError(t, statErr, "zip file was not created")
}

func TestExport_ProducesValidSwag2mcpZip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	outputPath := filepath.Join(tmpDir, "backup.zip")
	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	require.NoError(t, err)
	assert.True(t, workspace.IsSwag2mcpZip(outputPath))
}

func TestExport_LoadConfigError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte("invalid: yaml: ["), 0600))

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	require.Error(t, err)
}

func TestExport_DownloadSpecError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: http://localhost:1/nonexistent
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	require.Error(t, err)
}

func TestExport_DisabledSpec(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    disable: true
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no collections")
}

func TestExport_DisabledCollection(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
        disable: true
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no collections")
}

func TestExport_DuplicateLocation(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
      - title: PetsDup
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, resp.FileCount)
}

func TestExport_CreateZipError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600))

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "nonexistent", "backup.zip"),
	})
	require.Error(t, err)
}
