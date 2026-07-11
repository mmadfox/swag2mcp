package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScript_Init_CreatesWorkspace(t *testing.T) {
	ws := newTestWorkspace(t)
	_, stderr, code := runCommand(t, "init", ws)
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stderr", stderr, "initialized")

	wsDir := filepath.Join(ws, ".swag2mcp")
	dirs := []string{"cache", "specs", "responses", "auth_scripts"}
	for _, d := range dirs {
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

func TestScript_Init_ForceOverwrite(t *testing.T) {
	ws := newTestWorkspace(t)
	runCommand(t, "init", ws)

	_, stderr, code := runCommand(t, "init", ws)
	assertNotEqual(t, "exit code without -f", code, 0)
	assertContains(t, "stderr", stderr, "already exists")

	_, _, code = runCommand(t, "init", "-f", ws)
	assertEqual(t, "exit code with -f", code, 0)
}

func TestScript_Init_Interactive(t *testing.T) {
	t.Skip("requires TTY")
}

func TestScript_Init_CustomPath(t *testing.T) {
	ws := newTestWorkspace(t)
	customPath := filepath.Join(ws, "custom", "nested", "workspace")
	stdout, stderr, code := runCommand(t, "init", customPath)
	assertEqual(t, "exit code", code, 0)
	assertContains(t, "stdout", stdout+stderr, "initialized")

	configPath := filepath.Join(customPath, ".swag2mcp", "swag2mcp.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("config not created at custom path: %s", configPath)
	}
}
