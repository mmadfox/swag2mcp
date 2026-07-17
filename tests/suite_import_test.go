package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestScript_Import_SingleFile(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	stdout, stderr, code := runCommand(t, "import", ws, srv.URL, "petstore.yaml")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "output", stdout+stderr, "petstore.yaml")

	specPath := filepath.Join(ws, "specs", "petstore.yaml")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		t.Errorf("spec file not created at %s", specPath)
	}
}

func TestScript_Import_SingleFile_Duplicate(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	specDir := filepath.Join(ws, "specs")
	if err := os.WriteFile(filepath.Join(specDir, "petstore.yaml"), []byte("existing"), 0644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("openapi: 3.0.0"))
	}))
	t.Cleanup(srv.Close)

	_, _, code := runCommand(t, "import", ws, srv.URL, "petstore.yaml")
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Import_WithSpec(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	specContent := `openapi: 3.0.0
info:
  title: Petstore API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	writeConfig(t, ws, configContent)

	stdout, stderr, code := runCommand(t, "import", ws, "--spec", "petstore")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "output", stdout+stderr, "Imported")

	specPath := filepath.Join(ws, "specs")
	entries, err := os.ReadDir(specPath)
	if err != nil {
		t.Fatalf("read specs dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("no spec files were imported")
	}
}

func TestScript_Import_WithSpec_NoMatch(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: https://example.com/spec.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommand(t, "import", ws, "--spec", "nonexistent")
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Import_WithSpec_NoConfig(t *testing.T) {
	ws := newTestWorkspace(t)

	_, _, code := runCommand(t, "import", ws, "--spec", "petstore")
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Import_MissingArgs(t *testing.T) {
	ws := newTestWorkspace(t)

	_, _, code := runCommand(t, "import", ws)
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Import_FromZip(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	writeConfig(t, ws, configContent)

	// First export the workspace
	exportPath := filepath.Join(ws, "backup.zip")
	stdout, stderr, code := runCommand(t, "export", ws, exportPath)
	assertEqual(t, "export exit code", code, 0)
	assertContains(t, "export output", stdout+stderr, "Exported")

	// Now import from zip into a fresh workspace
	restoreWS := newTestWorkspace(t)
	stdout, stderr, code = runCommand(t, "import", restoreWS, "--from-zip", exportPath)
	assertEqual(t, "import exit code", code, 0)
	assertContains(t, "import output", stdout+stderr, "Restored")

	// Verify specs were restored
	specDir := filepath.Join(restoreWS, "specs")
	entries, err := os.ReadDir(specDir)
	if err != nil {
		t.Fatalf("read specs dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("no spec files were restored from zip")
	}

	// Verify config was restored
	configPath := filepath.Join(restoreWS, "swag2mcp.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config was not restored from zip")
	}
}

func TestScript_Import_FromZip_Invalid(t *testing.T) {
	ws := newTestWorkspace(t)

	invalidZip := filepath.Join(ws, "invalid.zip")
	if err := os.WriteFile(invalidZip, []byte("not a zip"), 0644); err != nil {
		t.Fatalf("write invalid zip: %v", err)
	}

	_, _, code := runCommand(t, "import", ws, "--from-zip", invalidZip)
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Import_FromZip_DetectByExtension(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.petstore.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	writeConfig(t, ws, configContent)

	exportPath := filepath.Join(ws, "backup.zip")
	runCommand(t, "export", ws, exportPath)

	// Import by passing zip path as source (no --from-zip)
	restoreWS := newTestWorkspace(t)
	stdout, stderr, code := runCommand(t, "import", restoreWS, exportPath)
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "import output", stdout+stderr, "Restored")

	specDir := filepath.Join(restoreWS, "specs")
	entries, err := os.ReadDir(specDir)
	if err != nil {
		t.Fatalf("read specs dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("no spec files were restored from zip")
	}
}
