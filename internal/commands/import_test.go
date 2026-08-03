/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

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
	"github.com/stretchr/testify/require"

	"github.com/mmadfox/swag2mcp/internal/service"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func testCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	return cmd
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{name: "https URL", s: "https://example.com/spec.yaml", want: true},
		{name: "http URL", s: "http://example.com/spec.yaml", want: true},
		{name: "local path", s: "./local-spec.yaml", want: false},
		{name: "absolute path", s: "/home/user/spec.yaml", want: false},
		{name: "filename only", s: "myspec.yaml", want: false},
		{name: "empty string", s: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isURL(tt.s)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseSingleOrZipArgs_PathSourceURL(t *testing.T) {
	args := []string{"/path/to/ws", "https://example.com/spec.yaml"}
	got := parseSingleOrZipArgs(args)
	require.Equal(t, importModeSingle, got.mode)
	require.Equal(t, "/path/to/ws", got.basePath)
	require.Equal(t, "https://example.com/spec.yaml", got.source)
	require.Empty(t, got.name)
}

func TestParseSingleOrZipArgs_SourceName(t *testing.T) {
	args := []string{"./local-spec.yaml", "myspec"}
	got := parseSingleOrZipArgs(args)
	require.Equal(t, importModeSingle, got.mode)
	require.Empty(t, got.basePath)
	require.Equal(t, "./local-spec.yaml", got.source)
	require.Equal(t, "myspec", got.name)
}

func TestParseSingleOrZipArgs_SourceOnly(t *testing.T) {
	args := []string{"https://example.com/spec.yaml"}
	got := parseSingleOrZipArgs(args)
	require.Equal(t, importModeSingle, got.mode)
	require.Empty(t, got.basePath)
	require.Equal(t, "https://example.com/spec.yaml", got.source)
	require.Empty(t, got.name)
}

func TestParseSingleOrZipArgs_FullPathSourceName(t *testing.T) {
	args := []string{"/path/to/ws", "./local-spec.yaml", "myspec"}
	got := parseSingleOrZipArgs(args)
	require.Equal(t, importModeSingle, got.mode)
	require.Equal(t, "/path/to/ws", got.basePath)
	require.Equal(t, "./local-spec.yaml", got.source)
	require.Equal(t, "myspec", got.name)
}

func TestParseSingleOrZipArgs_PathZip(t *testing.T) {
	args := []string{"/path/to/ws", "/path/to/backup.zip"}
	got := parseSingleOrZipArgs(args)
	require.Equal(t, importModeZip, got.mode)
	require.Equal(t, "/path/to/ws", got.basePath)
	require.Equal(t, "/path/to/backup.zip", got.zipFile)
}

func TestParseSingleOrZipArgs_Empty(t *testing.T) {
	got := parseSingleOrZipArgs(nil)
	require.Equal(t, importModeSingle, got.mode)
	require.Empty(t, got.basePath)
	require.Empty(t, got.source)
	require.Empty(t, got.name)
}

func TestRunImport_NoSpec_ForceOverwrites(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specDir := filepath.Join(tmpDir, "specs")
	existingPath := filepath.Join(specDir, "myspec.yaml")
	if err := os.WriteFile(existingPath, []byte("old"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	specContent := "openapi: 3.0.0\ninfo:\n  title: Updated\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	parsed := importArgs{
		mode:     importModeSingle,
		basePath: tmpDir,
		source:   srv.URL + "/spec.yaml",
		name:     "myspec.yaml",
	}
	err := runImport(parsed, nil, true, cmd)
	if err != nil {
		t.Fatalf("runImport() with force = %v", err)
	}

	if !strings.Contains(buf.String(), "myspec.yaml") {
		t.Errorf("output = %q, want success message with filename", buf.String())
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(data) != specContent {
		t.Errorf("file content = %q, want %q", string(data), specContent)
	}
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

	parsed := importArgs{
		mode:     importModeSingle,
		basePath: tmpDir,
		source:   srv.URL + "/spec.yaml",
		name:     "myspec.yaml",
	}
	err := runImport(parsed, nil, false, cmd)
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
	parsed := importArgs{mode: importModeSingle, basePath: tmpDir}
	err := runImport(parsed, nil, false, cmd)

	if err == nil {
		t.Fatal("runImport() expected error, got nil")
	}
}

func TestRunImport_NoSpec_NameDerivedFromURL(t *testing.T) {
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

	parsed := importArgs{
		mode:     importModeSingle,
		basePath: tmpDir,
		source:   srv.URL + "/specs/petstore.yaml",
	}
	err := runImport(parsed, nil, false, cmd)
	if err != nil {
		t.Fatalf("runImport() = %v", err)
	}

	if !strings.Contains(buf.String(), "petstore-specs.yaml") {
		t.Errorf("output = %q, want success message with derived filename", buf.String())
	}

	specPath := filepath.Join(tmpDir, "specs", "petstore-specs.yaml")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Errorf("spec file was not created at %s", specPath)
	}
}

func TestRunImport_NoSpec_NameNotDerivedFromURL_NoFilename(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	parsed := importArgs{
		mode:     importModeSingle,
		basePath: tmpDir,
		source:   "https://example.com/",
	}
	err := runImport(parsed, nil, false, cmd)
	if err == nil {
		t.Fatal("runImport() expected error for URL without filename, got nil")
	}
}

func TestRunImport_NoSpec_InvalidExtension(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cmd := testCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	parsed := importArgs{
		mode:     importModeSingle,
		basePath: tmpDir,
		source:   "https://example.com/spec.html",
		name:     "spec.yaml",
	}
	err := runImport(parsed, nil, false, cmd)
	if err == nil {
		t.Fatal("runImport() expected error for invalid extension, got nil")
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

	parsed := importArgs{mode: importModeBulk, basePath: tmpDir}
	err := runImport(parsed, []string{"meteo"}, false, cmd)
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

	parsed := importArgs{mode: importModeBulk, basePath: tmpDir}
	err := runImport(parsed, []string{"meteo"}, false, cmd)
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

	parsed := importArgs{mode: importModeBulk, basePath: tmpDir}
	err := runImport(parsed, []string{"nonexistent"}, false, cmd)
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

	parsed := importArgs{mode: importModeZip, basePath: restoreDir, zipFile: zipPath}
	err := runImport(parsed, nil, false, cmd)
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

	parsed := parseImportArgs([]string{restoreDir, zipPath}, "", false)
	err := runImport(parsed, nil, false, cmd)
	if err != nil {
		t.Fatalf("runImport() = %v", err)
	}

	if !strings.Contains(buf.String(), "Restored successfully") {
		t.Errorf("output = %q, want restore message", buf.String())
	}
}
