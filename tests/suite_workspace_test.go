package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScript_Workspace_DirectoryStructure(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	expectedDirs := []string{"cache", "specs", "responses", "auth_scripts"}
	wsDir := filepath.Join(ws, ".swag2mcp")
	for _, d := range expectedDirs {
		info, err := os.Stat(filepath.Join(wsDir, d))
		if err != nil {
			t.Errorf("missing directory %s: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", d)
		}
	}

	configPath := filepath.Join(wsDir, "swag2mcp.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("missing swag2mcp.yaml")
	}
}

func TestScript_Workspace_CleanRemovesCacheAndResponses(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	wsDir := filepath.Join(ws, ".swag2mcp")
	cacheFile := filepath.Join(wsDir, "cache", "test.cache")
	_ = os.MkdirAll(filepath.Dir(cacheFile), 0755)
	_ = os.WriteFile(cacheFile, []byte("data"), 0644)
	respFile := filepath.Join(wsDir, "responses", "test.json")
	_ = os.MkdirAll(filepath.Dir(respFile), 0755)
	_ = os.WriteFile(respFile, []byte("{}"), 0644)
	specFile := filepath.Join(wsDir, "specs", "test.yaml")
	_ = os.MkdirAll(filepath.Dir(specFile), 0755)
	_ = os.WriteFile(specFile, []byte("spec: test"), 0644)

	runCommandInWS(t, ws, "clean", ".")

	if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
		t.Errorf("cache file was not removed")
	}
	if _, err := os.Stat(respFile); !os.IsNotExist(err) {
		t.Errorf("response file was not removed")
	}
	if _, err := os.Stat(specFile); os.IsNotExist(err) {
		t.Errorf("spec file was removed but should be preserved")
	}
}

func TestScript_Workspace_UpdateReCachesSpecs(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	stdout, stderr, code := runCommandInWS(t, ws, "update", ".")
	if code != 0 {
		t.Fatalf("update failed (exit %d):\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
}

func TestScript_Workspace_UpdateWithInvalidConfig(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	configContent := `specs:
  - domain: "INVALID DOMAIN"
    llm_title: Bad
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./nonexistent.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "update", ".")
	if code == 0 {
		t.Errorf("expected update to fail with invalid config")
	}
}

func TestScript_Workspace_OldResponsesCleaned(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	wsDir := filepath.Join(ws, ".swag2mcp")
	oldResp := filepath.Join(wsDir, "responses", "old.json")
	_ = os.MkdirAll(filepath.Dir(oldResp), 0755)
	_ = os.WriteFile(oldResp, []byte("old"), 0644)

	past := time.Now().Add(-49 * time.Hour)
	_ = os.Chtimes(oldResp, past, past)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	runCommandInWS(t, ws, "mcp", ".", "--auth-token", "test")

	if _, err := os.Stat(oldResp); !os.IsNotExist(err) {
		t.Errorf("old response file (>48h) was not cleaned on mcp start")
	}
}

func TestScript_Workspace_RecentResponsesPreserved(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	wsDir := filepath.Join(ws, ".swag2mcp")
	recentResp := filepath.Join(wsDir, "responses", "recent.json")
	_ = os.MkdirAll(filepath.Dir(recentResp), 0755)
	_ = os.WriteFile(recentResp, []byte("recent"), 0644)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	runCommandInWS(t, ws, "mcp", ".", "--auth-token", "test")

	if _, err := os.Stat(recentResp); os.IsNotExist(err) {
		t.Errorf("recent response file (<48h) was incorrectly cleaned")
	}
}
