package commands

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func testCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	return cmd
}

func TestRunImport_NoSpec_Success(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runImport(tmpDir, srv.URL, "myspec.yaml", "", nil, cmd)
	if err != nil {
		t.Fatalf("runImport() = %v", err)
	}

	if !strings.Contains(buf.String(), "myspec.yaml") {
		t.Errorf("output = %q, want success message with filename", buf.String())
	}

	specPath := filepath.Join(tmpDir, "specs", "myspec.yaml")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Errorf("spec file was not created at %s", specPath)
	}
}

func TestRunImport_NoSpec_MissingArgs(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	err := runImport(tmpDir, "", "", "", nil, cmd)

	if err == nil {
		t.Fatal("runImport() expected error, got nil")
	}
}

func TestRunImport_WithSpec_Success(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	cfgContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	if err := os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runImport(tmpDir, "", "", "", []string{"meteo"}, cmd)
	if err != nil {
		t.Fatalf("runImport() = %v", err)
	}

	if !strings.Contains(buf.String(), "Imported 1 spec files") {
		t.Errorf("output = %q, want success message", buf.String())
	}
}

func TestRunImport_WithSpec_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runImport(tmpDir, "", "", "", []string{"meteo"}, cmd)
	if err == nil {
		t.Fatal("runImport() expected error, got nil")
	}
}

func TestRunImport_WithSpec_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: https://example.com/spec.yaml
`
	if err := os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runImport(tmpDir, "", "", "", []string{"nonexistent"}, cmd)
	if err == nil {
		t.Fatal("runImport() expected error for no matching specs, got nil")
	}
}

func TestRunImport_FromZip(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a swag2mcp backup zip
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	exportDir := t.TempDir()
	exportWs, _ := workspace.New(exportDir)
	if err := exportWs.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	if err := os.WriteFile(exportWs.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	zipPath := filepath.Join(tmpDir, "backup.zip")
	exportSvc, _ := service.New(service.WithWorkspace(exportWs))
	_, exportErr := exportSvc.Export(context.Background(), service.ExportRequest{
		OutputPath: zipPath,
	})
	if exportErr != nil {
		t.Fatalf("Export() = %v", exportErr)
	}

	// Now restore from zip
	restoreDir := t.TempDir()
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runImport(restoreDir, "", "", zipPath, nil, cmd)
	if err != nil {
		t.Fatalf("runImport() = %v", err)
	}

	if !strings.Contains(buf.String(), "Restored successfully") {
		t.Errorf("output = %q, want restore message", buf.String())
	}
}

func TestRunImport_FromZip_DetectByExtension(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a swag2mcp backup zip
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	exportDir := t.TempDir()
	exportWs, _ := workspace.New(exportDir)
	if err := exportWs.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	if err := os.WriteFile(exportWs.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	zipPath := filepath.Join(tmpDir, "backup.zip")
	exportSvc, _ := service.New(service.WithWorkspace(exportWs))
	_, exportErr := exportSvc.Export(context.Background(), service.ExportRequest{
		OutputPath: zipPath,
	})
	if exportErr != nil {
		t.Fatalf("Export() = %v", exportErr)
	}

	// Restore by passing zip as source (detected by .zip extension)
	// Simulate: swag2mcp import /path/to/workspace /path/to/backup.zip
	// args = [restoreDir, zipPath] → parseImportArgs detects zip in args[1]
	restoreDir := t.TempDir()
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	basePath, _, _, zipSource := parseImportArgs([]string{restoreDir, zipPath}, nil, "")
	err := runImport(basePath, "", "", zipSource, nil, cmd)
	if err != nil {
		t.Fatalf("runImport() = %v", err)
	}

	if !strings.Contains(buf.String(), "Restored successfully") {
		t.Errorf("output = %q, want restore message", buf.String())
	}
}
