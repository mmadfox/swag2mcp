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

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	if err != nil {
		t.Fatalf("Export() = %v", err)
	}

	if resp.FileCount != 1 {
		t.Errorf("FileCount = %d, want 1", resp.FileCount)
	}
	if resp.OutputPath != outputPath {
		t.Errorf("OutputPath = %q, want %q", resp.OutputPath, outputPath)
	}
	if _, statErr := os.Stat(outputPath); os.IsNotExist(statErr) {
		t.Error("zip file was not created")
	}
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
	if err := os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
		SpecFilter: []string{"petstore"},
	})
	if err != nil {
		t.Fatalf("Export() = %v", err)
	}

	if resp.FileCount != 1 {
		t.Errorf("FileCount = %d, want 1", resp.FileCount)
	}
}

func TestExport_NoConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t)
	svc.ws, _ = workspace.New(tmpDir)

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	if err == nil {
		t.Fatal("Export() expected error, got nil")
	}
}

func TestExport_NoCollections(t *testing.T) {
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
    collections: []
`
	if err := os.WriteFile(svc.ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: filepath.Join(tmpDir, "backup.zip"),
	})
	if err == nil {
		t.Fatal("Export() expected error, got nil")
	}
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

	outputPath := filepath.Join(tmpDir, "backup.zip")
	resp, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	if err != nil {
		t.Fatalf("Export() = %v", err)
	}

	if !strings.HasSuffix(resp.OutputPath, ".zip") {
		t.Errorf("OutputPath = %q, want .zip suffix", resp.OutputPath)
	}
	if _, statErr := os.Stat(outputPath); os.IsNotExist(statErr) {
		t.Error("zip file was not created")
	}
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

	outputPath := filepath.Join(tmpDir, "backup.zip")
	_, err := svc.Export(context.Background(), ExportRequest{
		OutputPath: outputPath,
	})
	if err != nil {
		t.Fatalf("Export() = %v", err)
	}

	if !workspace.IsSwag2mcpZip(outputPath) {
		t.Error("exported zip is not a valid swag2mcp backup")
	}
}
