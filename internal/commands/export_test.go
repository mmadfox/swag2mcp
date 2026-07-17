package commands

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestRunExport_Success(t *testing.T) {
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
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	if err := os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	outputPath := filepath.Join(tmpDir, "backup.zip")
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runExport(tmpDir, outputPath, nil, cmd)
	if err != nil {
		t.Fatalf("runExport() = %v", err)
	}

	if !strings.Contains(buf.String(), "Exported") {
		t.Errorf("output = %q, want success message", buf.String())
	}
	if _, statErr := os.Stat(outputPath); os.IsNotExist(statErr) {
		t.Error("zip file was not created")
	}
}

func TestRunExport_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runExport(tmpDir, "", nil, cmd)
	if err == nil {
		t.Fatal("runExport() expected error, got nil")
	}
}

func TestRunExport_DefaultOutputPath(t *testing.T) {
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
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	if err := os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	outputPath := filepath.Join(tmpDir, "backup.zip")
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runExport(tmpDir, outputPath, nil, cmd)
	if err != nil {
		t.Fatalf("runExport() = %v", err)
	}

	if !strings.Contains(buf.String(), "backup.zip") {
		t.Errorf("output = %q, want .zip in output", buf.String())
	}
	if _, statErr := os.Stat(outputPath); os.IsNotExist(statErr) {
		t.Error("zip file was not created")
	}
}

func TestRunExport_WithSpecFilter(t *testing.T) {
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
	if err := os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	outputPath := filepath.Join(tmpDir, "backup.zip")
	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runExport(tmpDir, outputPath, []string{"petstore"}, cmd)
	if err != nil {
		t.Fatalf("runExport() = %v", err)
	}

	if !strings.Contains(buf.String(), "Exported") {
		t.Errorf("output = %q, want success message", buf.String())
	}
}
