package workspace

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew_DefaultRoot(t *testing.T) {
	ws, err := New("")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName)
	if ws.Root() != want {
		t.Errorf("Root() = %q, want %q", ws.Root(), want)
	}
}

func TestNew_CustomRoot(t *testing.T) {
	ws, err := New("/tmp/test-workspace")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if ws.Root() != "/tmp/test-workspace" {
		t.Errorf("Root() = %q, want %q", ws.Root(), "/tmp/test-workspace")
	}
}

func TestSub(t *testing.T) {
	ws, _ := New("/root")
	if got := ws.Sub("mydir"); got != "/root/mydir" {
		t.Errorf("Sub() = %q, want %q", got, "/root/mydir")
	}
}

func TestCacheDir(t *testing.T) {
	ws, _ := New("/root")
	if got := ws.CacheDir(); got != "/root/"+DirCache {
		t.Errorf("CacheDir() = %q, want %q", got, "/root/"+DirCache)
	}
}

func TestSpecsDir(t *testing.T) {
	ws, _ := New("/root")
	if got := ws.SpecsDir(); got != "/root/"+DirSpecs {
		t.Errorf("SpecsDir() = %q, want %q", got, "/root/"+DirSpecs)
	}
}

func TestResponsesDir(t *testing.T) {
	ws, _ := New("/root")
	if got := ws.ResponsesDir(); got != "/root/"+DirResponses {
		t.Errorf("ResponsesDir() = %q, want %q", got, "/root/"+DirResponses)
	}
}

func TestAuthScriptsDir(t *testing.T) {
	ws, _ := New("/root")
	if got := ws.AuthScriptsDir(); got != "/root/"+DirAuthScripts {
		t.Errorf("AuthScriptsDir() = %q, want %q", got, "/root/"+DirAuthScripts)
	}
}

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	for _, dir := range []string{ws.Root(), ws.CacheDir(), ws.SpecsDir(), ws.ResponsesDir(), ws.AuthScriptsDir()} {
		if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
			t.Errorf("directory %q was not created", dir)
		}
	}
}

func TestCleanOldResponses_RemovesOldFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	oldFile := filepath.Join(ws.ResponsesDir(), "old-response.json")
	if err := os.WriteFile(oldFile, []byte("old"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	oldModTime := time.Now().Add(-72 * time.Hour)
	if err := os.Chtimes(oldFile, oldModTime, oldModTime); err != nil {
		t.Fatalf("Chtimes() = %v", err)
	}

	freshFile := filepath.Join(ws.ResponsesDir(), "fresh-response.json")
	if err := os.WriteFile(freshFile, []byte("fresh"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.CleanOldResponses(48 * time.Hour); err != nil {
		t.Fatalf("CleanOldResponses() = %v", err)
	}

	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("old file was not removed")
	}
	if _, err := os.Stat(freshFile); os.IsNotExist(err) {
		t.Error("fresh file was removed")
	}
}

func TestCleanOldResponses_NoDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	if err := ws.CleanOldResponses(48 * time.Hour); err != nil {
		t.Fatalf("CleanOldResponses() = %v", err)
	}
}

func TestCleanOldResponses_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	if err := ws.CleanOldResponses(48 * time.Hour); err != nil {
		t.Fatalf("CleanOldResponses() = %v", err)
	}
}
