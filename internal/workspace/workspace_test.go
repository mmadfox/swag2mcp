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

func TestNew_RelativePath(t *testing.T) {
	ws, err := New("relative/path")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	abs, _ := filepath.Abs("relative/path")
	if ws.Root() != abs {
		t.Errorf("Root() = %q, want %q", ws.Root(), abs)
	}
}

func TestNewFromBase_Empty(t *testing.T) {
	ws, err := NewFromBase("")
	if err != nil {
		t.Fatalf("NewFromBase() = %v", err)
	}
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName)
	if ws.Root() != want {
		t.Errorf("Root() = %q, want %q", ws.Root(), want)
	}
}

func TestNewFromBase_Custom(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := NewFromBase(tmpDir)
	if err != nil {
		t.Fatalf("NewFromBase() = %v", err)
	}
	want := filepath.Join(tmpDir, DefaultRootName)
	if ws.Root() != want {
		t.Errorf("Root() = %q, want %q", ws.Root(), want)
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

func TestDefaultRoot(t *testing.T) {
	root := DefaultRoot()
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName)
	if root != want {
		t.Errorf("DefaultRoot() = %q, want %q", root, want)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName, "swag2mcp.yaml")
	if path != want {
		t.Errorf("DefaultConfigPath() = %q, want %q", path, want)
	}
}

func TestConfigPathIn(t *testing.T) {
	path := ConfigPathIn("/custom/workspace")
	want := "/custom/workspace/swag2mcp.yaml"
	if path != want {
		t.Errorf("ConfigPathIn() = %q, want %q", path, want)
	}
}

func TestConfigPath(t *testing.T) {
	ws, _ := New("/root")
	path := ws.ConfigPath()
	want := "/root/swag2mcp.yaml"
	if path != want {
		t.Errorf("ConfigPath() = %q, want %q", path, want)
	}
}

func TestConfigExists_True(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	cfgPath := ws.ConfigPath()
	if err := os.WriteFile(cfgPath, []byte("specs: []"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	if !ws.ConfigExists() {
		t.Error("ConfigExists() = false, want true")
	}
}

func TestConfigExists_False(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if ws.ConfigExists() {
		t.Error("ConfigExists() = true, want false")
	}
}

func TestConfigNotExists_True(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if !ws.ConfigNotExists() {
		t.Error("ConfigNotExists() = false, want true")
	}
}

func TestConfigNotExists_False(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	cfgPath := ws.ConfigPath()
	if err := os.WriteFile(cfgPath, []byte("specs: []"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	if ws.ConfigNotExists() {
		t.Error("ConfigNotExists() = true, want false")
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

func TestInit_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("first Init() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("second Init() = %v", err)
	}
}

func TestClean_RemovesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cacheFile := filepath.Join(ws.CacheDir(), "cached.yaml")
	if err := os.WriteFile(cacheFile, []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	respFile := filepath.Join(ws.ResponsesDir(), "response.json")
	if err := os.WriteFile(respFile, []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.Clean(); err != nil {
		t.Fatalf("Clean() = %v", err)
	}

	if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
		t.Error("cache file was not removed")
	}
	if _, err := os.Stat(respFile); !os.IsNotExist(err) {
		t.Error("response file was not removed")
	}
	if _, err := os.Stat(ws.CacheDir()); os.IsNotExist(err) {
		t.Error("cache dir was removed")
	}
	if _, err := os.Stat(ws.ResponsesDir()); os.IsNotExist(err) {
		t.Error("responses dir was removed")
	}
}

func TestClean_NoDirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Clean(); err != nil {
		t.Fatalf("Clean() = %v", err)
	}
}

func TestClean_EmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}
	if err := ws.Clean(); err != nil {
		t.Fatalf("Clean() = %v", err)
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

func TestCleanOldResponses_SkipsSubdirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	subDir := filepath.Join(ws.ResponsesDir(), "subdir")
	if err := os.MkdirAll(subDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	subFile := filepath.Join(subDir, "nested.json")
	if err := os.WriteFile(subFile, []byte("nested"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.CleanOldResponses(0); err != nil {
		t.Fatalf("CleanOldResponses() = %v", err)
	}

	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("subdirectory was removed")
	}
	if _, err := os.Stat(subFile); os.IsNotExist(err) {
		t.Error("nested file was removed")
	}
}

func TestAuthScriptPath_Unix(t *testing.T) {
	ws, _ := New("/root")
	path := ws.AuthScriptPath("my-api")
	want := filepath.Join(ws.AuthScriptsDir(), "my-api.sh")
	if path != want {
		t.Errorf("AuthScriptPath() = %q, want %q", path, want)
	}
}

func TestEnsureAuthScript_Creates(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	if err := ws.EnsureAuthScript("test-domain"); err != nil {
		t.Fatalf("EnsureAuthScript() = %v", err)
	}

	scriptPath := ws.AuthScriptPath("test-domain")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Error("auth script was not created")
	}
}

func TestEnsureAuthScript_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	if err := ws.EnsureAuthScript("test-domain"); err != nil {
		t.Fatalf("first EnsureAuthScript() = %v", err)
	}
	if err := ws.EnsureAuthScript("test-domain"); err != nil {
		t.Fatalf("second EnsureAuthScript() = %v", err)
	}
}

func TestEnsureAuthScript_CreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)

	if err := ws.EnsureAuthScript("test-domain"); err != nil {
		t.Fatalf("EnsureAuthScript() = %v", err)
	}

	if _, err := os.Stat(ws.AuthScriptsDir()); os.IsNotExist(err) {
		t.Error("auth scripts dir was not created")
	}
}

func TestRemoveOrphanAuthScripts_RemovesOrphan(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	orphanPath := filepath.Join(ws.AuthScriptsDir(), "orphan.sh")
	if err := os.WriteFile(orphanPath, []byte("echo test"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err != nil {
		t.Fatalf("RemoveOrphanAuthScripts() = %v", err)
	}

	if _, err := os.Stat(orphanPath); !os.IsNotExist(err) {
		t.Error("orphan script was not removed")
	}
}

func TestRemoveOrphanAuthScripts_KeepsActive(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	activePath := filepath.Join(ws.AuthScriptsDir(), "active.sh")
	if err := os.WriteFile(activePath, []byte("echo test"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err != nil {
		t.Fatalf("RemoveOrphanAuthScripts() = %v", err)
	}

	if _, err := os.Stat(activePath); os.IsNotExist(err) {
		t.Error("active script was removed")
	}
}

func TestRemoveOrphanAuthScripts_NoDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err != nil {
		t.Fatalf("RemoveOrphanAuthScripts() = %v", err)
	}
}

func TestRemoveOrphanAuthScripts_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err != nil {
		t.Fatalf("RemoveOrphanAuthScripts() = %v", err)
	}
}

func TestRemoveOrphanAuthScripts_SkipsNonScriptFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	nonScript := filepath.Join(ws.AuthScriptsDir(), "readme.txt")
	if err := os.WriteFile(nonScript, []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err != nil {
		t.Fatalf("RemoveOrphanAuthScripts() = %v", err)
	}

	if _, err := os.Stat(nonScript); os.IsNotExist(err) {
		t.Error("non-script file was removed")
	}
}

func TestRemoveOrphanAuthScripts_SkipsSubdirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	subDir := filepath.Join(ws.AuthScriptsDir(), "subdir")
	if err := os.MkdirAll(subDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err != nil {
		t.Fatalf("RemoveOrphanAuthScripts() = %v", err)
	}

	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("subdirectory was removed")
	}
}
