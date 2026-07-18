package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestImport_SingleFile_FromURL(t *testing.T) {
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

	resp, err := svc.Import(context.Background(), ImportRequest{
		Source: srv.URL,
		Name:   "myspec.yaml",
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	assert.Equal(t, "myspec.yaml", resp.Files[0].Name)

	saved, err := os.ReadFile(resp.Files[0].SavedPath)
	require.NoError(t, err)
	assert.Equal(t, specContent, string(saved))
}

func TestImport_SingleFile_FromLocalPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	specPath := filepath.Join(tmpDir, "source.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(specContent), 0600))

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	resp, err := svc.Import(context.Background(), ImportRequest{
		Source: specPath,
		Name:   "myspec.yaml",
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)

	saved, err := os.ReadFile(resp.Files[0].SavedPath)
	require.NoError(t, err)
	assert.Equal(t, specContent, string(saved))
}

func TestImport_SingleFile_DuplicateError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, svc.ws.Init())

	specsDir := svc.ws.SpecsDir()
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "myspec.yaml"), []byte("existing"), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		Source: "https://example.com/spec.yaml",
		Name:   "myspec.yaml",
	})
	require.Error(t, err)
}

func TestImport_ValidationError_EmptySource(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	_, err := svc.Import(context.Background(), ImportRequest{
		Source: "",
		Name:   "",
	})
	require.Error(t, err)
}

func TestImport_SourceWithoutName(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	_, err := svc.Import(context.Background(), ImportRequest{
		Source: "https://example.com/spec.yaml",
		Name:   "",
	})
	require.Error(t, err)
}

func TestImport_WithSpecFilter(t *testing.T) {
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
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	resp, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	require.NoError(t, err)
	require.Len(t, resp.Files, 1)
	assert.Equal(t, "petstore-pets.yaml", resp.Files[0].Name)
}

func TestImport_WithSpecFilter_UpdatesConfig(t *testing.T) {
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
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	require.NoError(t, err)

	cfg, err := os.ReadFile(cfgPath)
	require.NoError(t, err)
	assert.Contains(t, string(cfg), "specs/petstore-pets.yaml")
}

func TestImport_WithSpecFilter_NoMatch(t *testing.T) {
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
        location: https://example.com/spec.yaml
`
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"nonexistent"},
		ConfFilePath: cfgPath,
	})
	require.Error(t, err)
}

func TestImport_WithSpecFilter_EmptyConfFilePath(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: "",
	})
	require.Error(t, err)
}

func TestImport_FromZip(t *testing.T) {
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

	zipPath := filepath.Join(tmpDir, "backup.zip")
	_, exportErr := svc.Export(context.Background(), ExportRequest{
		OutputPath: zipPath,
	})
	require.NoError(t, exportErr)

	restoreDir := t.TempDir()
	restoreSvc := newTestService(t)
	restoreSvc.ws, _ = workspace.New(restoreDir)

	resp, importErr := restoreSvc.Import(context.Background(), ImportRequest{
		ZipSource: zipPath,
	})
	require.NoError(t, importErr)
	require.NotEmpty(t, resp.Files)

	specPath := filepath.Join(restoreDir, "specs", resp.Files[0].Name)
	_, statErr := os.Stat(specPath)
	require.NoError(t, statErr, "spec file was not restored at %s", specPath)

	_, statErr = os.Stat(restoreSvc.ws.ConfigPath())
	assert.NoError(t, statErr, "config was not restored")
}

func TestImport_FromZip_InvalidZip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	invalidZip := filepath.Join(tmpDir, "not-swag2mcp.zip")
	require.NoError(t, os.WriteFile(invalidZip, []byte("not a zip"), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		ZipSource: invalidZip,
	})
	require.Error(t, err)
}

func TestImport_FromZip_NotSwag2mcpZip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	regularZip := filepath.Join(tmpDir, "regular.zip")
	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))
	require.NoError(t, workspace.CreateZip(sourceDir, regularZip))

	_, err := svc.Import(context.Background(), ImportRequest{
		ZipSource: regularZip,
	})
	require.Error(t, err)
}

func TestImport_FromZip_ExtractError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	// Create a valid swag2mcp zip first
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	exportSvc := newTestService(t)
	exportSvc.ws, _ = workspace.New(tmpDir)
	require.NoError(t, exportSvc.ws.Init())

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	require.NoError(t, os.WriteFile(exportSvc.ws.ConfigPath(), []byte(cfgContent), 0600))

	zipPath := filepath.Join(tmpDir, "backup.zip")
	_, exportErr := exportSvc.Export(context.Background(), ExportRequest{
		OutputPath: zipPath,
	})
	require.NoError(t, exportErr)

	// Truncate the zip to make it corrupt
	require.NoError(t, os.Truncate(zipPath, 10))

	_, err := svc.Import(context.Background(), ImportRequest{
		ZipSource: zipPath,
	})
	require.Error(t, err)
}

func TestImport_WithSpecFilter_DisabledSpec(t *testing.T) {
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
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "No matching specs")
}

func TestImport_WithSpecFilter_DisabledCollection(t *testing.T) {
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
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "No matching specs")
}

func TestImport_WithSpecFilter_DownloadError(t *testing.T) {
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
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	require.Error(t, err)
}

func TestImport_WithSpecFilter_SaveSpecError(t *testing.T) {
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

	// Pre-create the file so SaveSpec fails with duplicate
	specsDir := svc.ws.SpecsDir()
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "petstore-pets.yaml"), []byte("existing"), 0600))

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(cfgContent), 0600))

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	require.Error(t, err)
}

func TestSpecFileNameBase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		want     string
	}{
		{name: "URL with path", location: "https://example.com/path/to/spec.yaml", want: "spec.yaml"},
		{name: "URL root path", location: "https://example.com/", want: defaultSpecName},
		{name: "URL no path", location: "https://example.com", want: defaultSpecName},
		{name: "local path no ext", location: "/home/user/spec", want: "spec"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := specFileNameBase(tt.location)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSpecFileName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		domain   string
		title    string
		location string
		want     string
	}{
		{name: "title matches domain", domain: "pokeapi", title: "PokéAPI", location: "https://example.com/pokeapi.yaml", want: "pokeapi.yaml"},
		{name: "title differs from domain", domain: "petstore", title: "Pet Operations", location: "https://example.com/pet.yaml", want: "petstore-pet-operations.yaml"},
		{name: "empty title falls back to location base", domain: "weather", title: "", location: "https://example.com/forecast.json", want: "weather-forecast.json"},
		{name: "title differs by case only", domain: "pokeapi", title: "pokeapi", location: "https://example.com/pokeapi.yaml", want: "pokeapi.yaml"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := specFileName(tt.domain, tt.title, tt.location)
			assert.Equal(t, tt.want, got)
		})
	}
}
