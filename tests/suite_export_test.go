package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestScript_Export_Success(t *testing.T) {
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

	outputPath := filepath.Join(ws, "backup.zip")
	stdout, stderr, code := runCommand(t, "export", ws, outputPath)
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "output", stdout+stderr, "Exported")

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("zip file not created at %s", outputPath)
	}
}

func TestScript_Export_NoConfig(t *testing.T) {
	ws := newTestWorkspace(t)

	_, _, code := runCommand(t, "export", ws)
	assertNotEqual(t, "exit code", code, 0)
}

func TestScript_Export_DefaultOutputPath(t *testing.T) {
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

	// Run export with CWD set to workspace so default output lands inside ws (temp dir)
	stdout, stderr, code := runCommandInWS(t, ws, "export", ".")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "output", stdout+stderr, ".zip")
}

func TestScript_Export_WithSpecFilter(t *testing.T) {
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
  - domain: store
    llm_title: Store API
    base_url: https://api.store.com
    collections:
      - title: Products
        location: ` + srv.URL + `
`
	writeConfig(t, ws, configContent)

	outputPath := filepath.Join(ws, "backup.zip")
	stdout, stderr, code := runCommand(t, "export", ws, outputPath, "--spec", "petstore")
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "output", stdout+stderr, "Exported")
}

func TestScript_Export_ProducesValidSwag2mcpZip(t *testing.T) {
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

	outputPath := filepath.Join(ws, "backup.zip")
	runCommand(t, "export", ws, outputPath)

	// Verify the zip is a valid swag2mcp backup by importing it
	restoreWS := newTestWorkspace(t)
	stdout, stderr, code := runCommand(t, "import", restoreWS, "--from-zip", outputPath)
	assertEqual(t, "import exit code", code, 0)
	assertContains(t, "import output", stdout+stderr, "Restored")
}
