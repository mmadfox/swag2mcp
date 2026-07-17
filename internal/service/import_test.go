package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	resp, err := svc.Import(context.Background(), ImportRequest{
		Source: srv.URL,
		Name:   "myspec.yaml",
	})
	if err != nil {
		t.Fatalf("Import() = %v", err)
	}

	if len(resp.Files) != 1 {
		t.Fatalf("Files = %d, want 1", len(resp.Files))
	}
	if resp.Files[0].Name != "myspec.yaml" {
		t.Errorf("Name = %q, want %q", resp.Files[0].Name, "myspec.yaml")
	}

	saved, err := os.ReadFile(resp.Files[0].SavedPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(saved) != specContent {
		t.Errorf("content = %q, want %q", string(saved), specContent)
	}
}

func TestImport_SingleFile_FromLocalPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	specPath := filepath.Join(tmpDir, "source.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	resp, err := svc.Import(context.Background(), ImportRequest{
		Source: specPath,
		Name:   "myspec.yaml",
	})
	if err != nil {
		t.Fatalf("Import() = %v", err)
	}

	if len(resp.Files) != 1 {
		t.Fatalf("Files = %d, want 1", len(resp.Files))
	}

	saved, err := os.ReadFile(resp.Files[0].SavedPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(saved) != specContent {
		t.Errorf("content = %q, want %q", string(saved), specContent)
	}
}

func TestImport_SingleFile_DuplicateError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specsDir := svc.ws.SpecsDir()
	if err := os.WriteFile(filepath.Join(specsDir, "myspec.yaml"), []byte("existing"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := svc.Import(context.Background(), ImportRequest{
		Source: "https://example.com/spec.yaml",
		Name:   "myspec.yaml",
	})
	if err == nil {
		t.Fatal("Import() expected error for duplicate name, got nil")
	}
}

func TestImport_ValidationError_EmptySource(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	_, err := svc.Import(context.Background(), ImportRequest{
		Source: "",
		Name:   "",
	})
	if err == nil {
		t.Fatal("Import() expected error, got nil")
	}
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
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

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
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	resp, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	if err != nil {
		t.Fatalf("Import() = %v", err)
	}

	if len(resp.Files) != 1 {
		t.Fatalf("Files = %d, want 1", len(resp.Files))
	}
	if resp.Files[0].Name != "petstore-pets.yaml" {
		t.Errorf("Name = %q, want %q", resp.Files[0].Name, "petstore-pets.yaml")
	}
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
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: cfgPath,
	})
	if err != nil {
		t.Fatalf("Import() = %v", err)
	}

	cfg, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}

	if !strings.Contains(string(cfg), "specs/petstore-pets.yaml") {
		t.Errorf("config does not contain updated location: %s", string(cfg))
	}
}

func TestImport_WithSpecFilter_NoMatch(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: https://example.com/spec.yaml
`
	cfgPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"nonexistent"},
		ConfFilePath: cfgPath,
	})
	if err == nil {
		t.Fatal("Import() expected error for no matching specs, got nil")
	}
}

func TestImport_WithSpecFilter_EmptyConfFilePath(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	_, err := svc.Import(context.Background(), ImportRequest{
		SpecFilter:   []string{"petstore"},
		ConfFilePath: "",
	})
	if err == nil {
		t.Fatal("Import() expected error for empty confFilepath, got nil")
	}
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

	// First create a workspace with a spec and export it
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)
	if err := svc.ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	if err := os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	zipPath := filepath.Join(tmpDir, "backup.zip")
	_, exportErr := svc.Export(context.Background(), ExportRequest{
		OutputPath: zipPath,
	})
	if exportErr != nil {
		t.Fatalf("Export() = %v", exportErr)
	}

	// Now import from the zip into a fresh workspace
	restoreDir := t.TempDir()
	restoreSvc := newTestService(t)
	restoreSvc.ws, _ = workspace.New(restoreDir)

	resp, importErr := restoreSvc.Import(context.Background(), ImportRequest{
		ZipSource: zipPath,
	})
	if importErr != nil {
		t.Fatalf("Import() = %v", importErr)
	}

	if len(resp.Files) == 0 {
		t.Fatal("no files imported from zip")
	}

	specPath := filepath.Join(restoreDir, "specs", resp.Files[0].Name)
	if _, statErr := os.Stat(specPath); os.IsNotExist(statErr) {
		t.Errorf("spec file was not restored at %s", specPath)
	}

	if _, statErr := os.Stat(restoreSvc.ws.ConfigPath()); os.IsNotExist(statErr) {
		t.Error("config was not restored")
	}
}

func TestImport_FromZip_InvalidZip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	invalidZip := filepath.Join(tmpDir, "not-swag2mcp.zip")
	if err := os.WriteFile(invalidZip, []byte("not a zip"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := svc.Import(context.Background(), ImportRequest{
		ZipSource: invalidZip,
	})
	if err == nil {
		t.Fatal("Import() expected error for invalid zip, got nil")
	}
}

func TestImport_FromZip_NotSwag2mcpZip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	// Create a valid zip but without swag2mcp.meta
	regularZip := filepath.Join(tmpDir, "regular.zip")
	sourceDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	if err := workspace.CreateZip(sourceDir, regularZip); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	_, err := svc.Import(context.Background(), ImportRequest{
		ZipSource: regularZip,
	})
	if err == nil {
		t.Fatal("Import() expected error for non-swag2mcp zip, got nil")
	}
}
