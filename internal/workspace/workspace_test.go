package workspace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
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
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if ws.Root() != tmpDir {
		t.Errorf("Root() = %q, want %q", ws.Root(), tmpDir)
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
	if ws.Root() != tmpDir {
		t.Errorf("Root() = %q, want %q", ws.Root(), tmpDir)
	}
}

func TestSub(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", "mydir")
	if got := ws.Sub("mydir"); got != want {
		t.Errorf("Sub() = %q, want %q", got, want)
	}
}

func TestCacheDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirCache)
	if got := ws.CacheDir(); got != want {
		t.Errorf("CacheDir() = %q, want %q", got, want)
	}
}

func TestSpecsDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirSpecs)
	if got := ws.SpecsDir(); got != want {
		t.Errorf("SpecsDir() = %q, want %q", got, want)
	}
}

func TestResponsesDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirResponses)
	if got := ws.ResponsesDir(); got != want {
		t.Errorf("ResponsesDir() = %q, want %q", got, want)
	}
}

func TestAuthScriptsDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirAuthScripts)
	if got := ws.AuthScriptsDir(); got != want {
		t.Errorf("AuthScriptsDir() = %q, want %q", got, want)
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
	want := filepath.Join("/custom/workspace", "swag2mcp.yaml")
	if path != want {
		t.Errorf("ConfigPathIn() = %q, want %q", path, want)
	}
}

func TestConfigPath(t *testing.T) {
	ws, _ := New("/root")
	path := ws.ConfigPath()
	want := filepath.Join("/root", "swag2mcp.yaml")
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
		t.Skipf("Chtimes not supported on this filesystem: %v", err)
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

	oldModTime := time.Now().Add(-72 * time.Hour)
	if err := os.Chtimes(subDir, oldModTime, oldModTime); err != nil {
		t.Skipf("Chtimes not supported on this filesystem: %v", err)
	}

	if err := ws.CleanOldResponses(48 * time.Hour); err != nil {
		t.Fatalf("CleanOldResponses() = %v", err)
	}

	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("subdirectory was removed")
	}
	if _, err := os.Stat(subFile); os.IsNotExist(err) {
		t.Error("nested file was removed")
	}
}

func TestAuthScriptPath(t *testing.T) {
	ws, _ := New("/root")
	path := ws.AuthScriptPath("my-api")
	ext := filepath.Ext(path)
	if ext != ".sh" && ext != ".bat" {
		t.Errorf("AuthScriptPath() extension = %q, want .sh or .bat", ext)
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

// removeContents error: ReadDir on a file returns ENOTDIR.
func TestClean_ReadDirError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cacheDir := ws.CacheDir()
	if err := os.RemoveAll(cacheDir); err != nil {
		t.Fatalf("RemoveAll() = %v", err)
	}
	if err := os.WriteFile(cacheDir, []byte("not-a-dir"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.Clean(); err == nil {
		t.Error("Clean() expected error, got nil")
	}
}

// CleanOldResponses error: ReadDir on a file returns ENOTDIR.
func TestCleanOldResponses_ReadDirError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	respDir := ws.ResponsesDir()
	if err := os.RemoveAll(respDir); err != nil {
		t.Fatalf("RemoveAll() = %v", err)
	}
	if err := os.WriteFile(respDir, []byte("not-a-dir"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.CleanOldResponses(48 * time.Hour); err == nil {
		t.Error("CleanOldResponses() expected error, got nil")
	}
}

// CleanOldResponses error: Remove fails on read-only parent dir.
func TestCleanOldResponses_RemoveError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	oldFile := filepath.Join(ws.ResponsesDir(), "old-response.json")
	if err := os.WriteFile(oldFile, []byte("old"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	oldModTime := time.Now().Add(-72 * time.Hour)
	if err := os.Chtimes(oldFile, oldModTime, oldModTime); err != nil {
		t.Skipf("Chtimes not supported: %v", err)
	}

	if err := os.Chmod(ws.ResponsesDir(), 0500); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(ws.ResponsesDir(), 0750) })

	if err := ws.CleanOldResponses(48 * time.Hour); err == nil {
		t.Error("CleanOldResponses() expected error, got nil")
	}
}

// EnsureAuthScript error: MkdirAll fails when auth_scripts dir is a file.
func TestEnsureAuthScript_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)

	authDir := ws.AuthScriptsDir()
	if err := os.WriteFile(authDir, []byte("not-a-dir"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.EnsureAuthScript("test-domain"); err == nil {
		t.Error("EnsureAuthScript() expected error, got nil")
	}
}

// EnsureAuthScript error: WriteFile fails on read-only dir.
func TestEnsureAuthScript_WriteFileError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	if err := os.Chmod(ws.AuthScriptsDir(), 0500); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(ws.AuthScriptsDir(), 0750) })

	if err := ws.EnsureAuthScript("test-domain"); err == nil {
		t.Error("EnsureAuthScript() expected error, got nil")
	}
}

// RemoveOrphanAuthScripts error: ReadDir on a file returns ENOTDIR.
func TestRemoveOrphanAuthScripts_ReadDirError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	authDir := ws.AuthScriptsDir()
	if err := os.RemoveAll(authDir); err != nil {
		t.Fatalf("RemoveAll() = %v", err)
	}
	if err := os.WriteFile(authDir, []byte("not-a-dir"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err == nil {
		t.Error("RemoveOrphanAuthScripts() expected error, got nil")
	}
}

// RemoveOrphanAuthScripts error: Remove fails on read-only parent dir.
func TestRemoveOrphanAuthScripts_RemoveError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	orphanPath := filepath.Join(ws.AuthScriptsDir(), "orphan.sh")
	if err := os.WriteFile(orphanPath, []byte("echo test"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := os.Chmod(ws.AuthScriptsDir(), 0500); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(ws.AuthScriptsDir(), 0750) })

	if err := ws.RemoveOrphanAuthScripts([]string{"active"}); err == nil {
		t.Error("RemoveOrphanAuthScripts() expected error, got nil")
	}
}

// Init error: MkdirAll fails when root dir is a file.
func TestInit_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()

	rootFile := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(rootFile, []byte("block"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	blocker, _ := New(rootFile)

	if err := blocker.Init(); err == nil {
		t.Error("Init() expected error, got nil")
	}
}

// DefaultRoot error: [os.UserHomeDir] fails when HOME is unset.
func TestDefaultRoot_NoHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	root := DefaultRoot()
	if root != DefaultRootName {
		t.Errorf("DefaultRoot() = %q, want %q", root, DefaultRootName)
	}
}

// removeContents error: RemoveAll fails when parent dir is read-only.
func TestClean_RemoveAllError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cacheFile := filepath.Join(ws.CacheDir(), "data.yaml")
	if err := os.WriteFile(cacheFile, []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := os.Chmod(ws.CacheDir(), 0000); err != nil {
		t.Fatalf("Chmod() = %v", err)
	}
	t.Cleanup(func() { os.Chmod(ws.CacheDir(), 0750) })

	if err := ws.Clean(); err == nil {
		t.Error("Clean() expected error, got nil")
	}
}

func TestIsEmpty_NotExists(t *testing.T) {
	ws, _ := New("/nonexistent/path/that/does/not/exist")
	empty, err := ws.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() = %v", err)
	}
	if !empty {
		t.Error("IsEmpty() = false, want true for nonexistent dir")
	}
}

func TestIsEmpty_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() = %v", err)
	}
	if !empty {
		t.Error("IsEmpty() = false, want true for empty dir")
	}
}

func TestIsEmpty_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("data"), 0644)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() = %v", err)
	}
	if empty {
		t.Error("IsEmpty() = true, want false for dir with files")
	}
}

func TestIsEmpty_WithSubdir(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() = %v", err)
	}
	if empty {
		t.Error("IsEmpty() = true, want false for dir with subdir")
	}
}

func TestIsEmpty_OnlyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "swag2mcp.yaml"), []byte("specs: []"), 0644)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() = %v", err)
	}
	if !empty {
		t.Error("IsEmpty() = false, want true when only swag2mcp.yaml exists")
	}
}

func TestIsEmpty_ConfigAndOtherFiles(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "swag2mcp.yaml"), []byte("specs: []"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "other.txt"), []byte("data"), 0644)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	if err != nil {
		t.Fatalf("IsEmpty() = %v", err)
	}
	if empty {
		t.Error("IsEmpty() = true, want false when other files exist alongside config")
	}
}

func TestDownloadSpec_FromLocalFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	specPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	data, err := ws.DownloadSpec(context.Background(), specPath)
	if err != nil {
		t.Fatalf("DownloadSpec() = %v", err)
	}
	if string(data) != specContent {
		t.Errorf("DownloadSpec() = %q, want %q", string(data), specContent)
	}
}

func TestDownloadSpec_FromURL(t *testing.T) {
	t.Parallel()

	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	t.Cleanup(srv.Close)

	ws, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	data, err := ws.DownloadSpec(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("DownloadSpec() = %v", err)
	}
	if string(data) != specContent {
		t.Errorf("DownloadSpec() = %q, want %q", string(data), specContent)
	}
}

func TestDownloadSpec_EmptySource(t *testing.T) {
	t.Parallel()

	ws, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	_, err = ws.DownloadSpec(context.Background(), "")
	if err == nil {
		t.Fatal("DownloadSpec() expected error, got nil")
	}
}

func TestDownloadSpec_FromFileURL(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	specPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(specPath, []byte(specContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	data, err := ws.DownloadSpec(context.Background(), "file://"+specPath)
	if err != nil {
		t.Fatalf("DownloadSpec() = %v", err)
	}
	if string(data) != specContent {
		t.Errorf("DownloadSpec() = %q, want %q", string(data), specContent)
	}
}

func TestSpecPath(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirSpecs, "myspec.yaml")
	if got := ws.SpecPath("myspec.yaml"); got != want {
		t.Errorf("SpecPath() = %q, want %q", got, want)
	}
}

func TestListSpecs_Empty(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	names, err := ws.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs() = %v", err)
	}
	if len(names) != 0 {
		t.Errorf("ListSpecs() = %v, want empty", names)
	}
}

func TestListSpecs_WithFiles(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	specsDir := ws.SpecsDir()
	if err := os.WriteFile(filepath.Join(specsDir, "a.yaml"), []byte("a"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}
	if err := os.WriteFile(filepath.Join(specsDir, "b.yaml"), []byte("b"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	names, err := ws.ListSpecs()
	if err != nil {
		t.Fatalf("ListSpecs() = %v", err)
	}
	if len(names) != 2 {
		t.Errorf("ListSpecs() = %v, want 2 files", names)
	}
}

func TestSaveSpec_Success(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	data := []byte("openapi: 3.0.0")
	path, err := ws.SaveSpec("myspec.yaml", data)
	if err != nil {
		t.Fatalf("SaveSpec() = %v", err)
	}

	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(saved) != string(data) {
		t.Errorf("saved content = %q, want %q", string(saved), string(data))
	}
}

func TestSaveSpec_DuplicateError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	data := []byte("openapi: 3.0.0")
	if _, err := ws.SaveSpec("myspec.yaml", data); err != nil {
		t.Fatalf("first SaveSpec() = %v", err)
	}

	_, err = ws.SaveSpec("myspec.yaml", data)
	if err == nil {
		t.Fatal("second SaveSpec() expected error, got nil")
	}
}

func TestSaveSpec_EmptyName(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	_, err = ws.SaveSpec("", []byte("data"))
	if err == nil {
		t.Fatal("SaveSpec() expected error, got nil")
	}
}

func TestSaveSpec_EmptyData(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	_, err = ws.SaveSpec("spec.yaml", nil)
	if err == nil {
		t.Fatal("SaveSpec() expected error, got nil")
	}
}

func TestCreateExportDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	exportDir, err := ws.CreateExportDir()
	if err != nil {
		t.Fatalf("CreateExportDir() = %v", err)
	}
	defer os.RemoveAll(exportDir)

	if _, statErr := os.Stat(exportDir); os.IsNotExist(statErr) {
		t.Error("export dir was not created")
	}
}

func TestWriteSpecToExport(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	data := []byte("openapi: 3.0.0")

	if err := WriteSpecToExport(exportDir, "myspec.yaml", data); err != nil {
		t.Fatalf("WriteSpecToExport() = %v", err)
	}

	specPath := filepath.Join(exportDir, DirSpecs, "myspec.yaml")
	saved, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(saved) != string(data) {
		t.Errorf("content = %q, want %q", string(saved), string(data))
	}
}

func TestCopyAuthScriptsToExport(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	scriptContent := "#!/bin/sh\necho test"
	scriptPath := filepath.Join(ws.AuthScriptsDir(), "test.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	exportDir := t.TempDir()
	if err := ws.CopyAuthScriptsToExport(exportDir); err != nil {
		t.Fatalf("CopyAuthScriptsToExport() = %v", err)
	}

	saved, err := os.ReadFile(filepath.Join(exportDir, DirAuthScripts, "test.sh"))
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(saved) != scriptContent {
		t.Errorf("content = %q, want %q", string(saved), scriptContent)
	}
}

func TestCopyAuthScriptsToExport_NoDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	exportDir := t.TempDir()
	if err := ws.CopyAuthScriptsToExport(exportDir); err != nil {
		t.Fatalf("CopyAuthScriptsToExport() = %v", err)
	}
}

func TestCreateMetaFile(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	if err := CreateMetaFile(exportDir, "1.0.0"); err != nil {
		t.Fatalf("CreateMetaFile() = %v", err)
	}

	metaPath := filepath.Join(exportDir, MetaFileName)
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Fatal("meta file was not created")
	}
}

func TestCreateZip(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	outputPath := filepath.Join(t.TempDir(), "test.zip")
	if err := CreateZip(sourceDir, outputPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("zip file was not created")
	}
}

func TestCreateZip_AddsExtension(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	outputPath := filepath.Join(t.TempDir(), "test")
	if err := CreateZip(sourceDir, outputPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	if _, err := os.Stat(outputPath + ".zip"); os.IsNotExist(err) {
		t.Fatal("zip file without extension was not created")
	}
}

func TestValidateZip_Valid(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	if err := CreateMetaFile(exportDir, "1.0.0"); err != nil {
		t.Fatalf("CreateMetaFile() = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	if err := CreateZip(exportDir, zipPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	valid, err := ValidateZip(zipPath)
	if err != nil {
		t.Fatalf("ValidateZip() = %v", err)
	}
	if !valid {
		t.Error("ValidateZip() = false, want true")
	}
}

func TestValidateZip_Invalid(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	if err := CreateZip(sourceDir, zipPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	valid, err := ValidateZip(zipPath)
	if err != nil {
		t.Fatalf("ValidateZip() = %v", err)
	}
	if valid {
		t.Error("ValidateZip() = true, want false")
	}
}

func TestValidateZip_NotAZip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "notazip.zip")
	if err := os.WriteFile(path, []byte("not a zip"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	_, err := ValidateZip(path)
	if err == nil {
		t.Fatal("ValidateZip() expected error, got nil")
	}
}

func TestExtractZip(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "subdir", "test.txt"), []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	if err := CreateZip(sourceDir, zipPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	destDir := t.TempDir()
	if err := ExtractZip(zipPath, destDir); err != nil {
		t.Fatalf("ExtractZip() = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(destDir, "subdir", "test.txt"))
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("content = %q, want %q", string(data), "hello")
	}
}

func TestDefaultExportName(t *testing.T) {
	name := DefaultExportName()
	if len(name) == 0 {
		t.Error("DefaultExportName() returned empty string")
	}
	if filepath.Ext(name) != ".zip" {
		t.Errorf("extension = %q, want .zip", filepath.Ext(name))
	}
}

func TestIsSwag2mcpZip(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	if err := CreateMetaFile(exportDir, "1.0.0"); err != nil {
		t.Fatalf("CreateMetaFile() = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	if err := CreateZip(exportDir, zipPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	if !IsSwag2mcpZip(zipPath) {
		t.Error("IsSwag2mcpZip() = false, want true")
	}
}

func TestIsSwag2mcpZip_NotZip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if IsSwag2mcpZip(path) {
		t.Error("IsSwag2mcpZip() = true, want false")
	}
}

func TestReadMetaFromZip(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	if err := CreateMetaFile(exportDir, "2.0.0"); err != nil {
		t.Fatalf("CreateMetaFile() = %v", err)
	}

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	if err := CreateZip(exportDir, zipPath); err != nil {
		t.Fatalf("CreateZip() = %v", err)
	}

	meta, err := ReadMetaFromZip(zipPath)
	if err != nil {
		t.Fatalf("ReadMetaFromZip() = %v", err)
	}
	if meta.Type != MetaType {
		t.Errorf("Type = %q, want %q", meta.Type, MetaType)
	}
	if meta.Version != "2.0.0" {
		t.Errorf("Version = %q, want %q", meta.Version, "2.0.0")
	}
}

func TestCreateEmptyDirsInExport(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	if err := CreateEmptyDirsInExport(exportDir); err != nil {
		t.Fatalf("CreateEmptyDirsInExport() = %v", err)
	}

	for _, dir := range []string{DirCache, DirResponses} {
		path := filepath.Join(exportDir, dir)
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			t.Errorf("dir %s was not created", dir)
		}
	}
}

func TestCopyConfigToExport(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	cfgContent := "specs:\n  - domain: test\n"
	if err := os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	exportDir := t.TempDir()
	if err := ws.CopyConfigToExport(exportDir); err != nil {
		t.Fatalf("CopyConfigToExport() = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(exportDir, "swag2mcp.yaml"))
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(data) != cfgContent {
		t.Errorf("content = %q, want %q", string(data), cfgContent)
	}
}

func TestCopySpecsToWorkspace(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	exportDir := t.TempDir()
	specsDir := filepath.Join(exportDir, DirSpecs)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	if err := os.WriteFile(filepath.Join(specsDir, "test.yaml"), []byte("openapi: 3.0.0"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.CopySpecsToWorkspace(exportDir); err != nil {
		t.Fatalf("CopySpecsToWorkspace() = %v", err)
	}

	specPath := filepath.Join(ws.SpecsDir(), "test.yaml")
	if _, statErr := os.Stat(specPath); os.IsNotExist(statErr) {
		t.Error("spec was not copied to workspace")
	}
}

func TestCopyAuthScriptsToWorkspace(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	exportDir := t.TempDir()
	authDir := filepath.Join(exportDir, DirAuthScripts)
	if err := os.MkdirAll(authDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	if err := os.WriteFile(filepath.Join(authDir, "test.sh"), []byte("#!/bin/sh"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := ws.CopyAuthScriptsToWorkspace(exportDir); err != nil {
		t.Fatalf("CopyAuthScriptsToWorkspace() = %v", err)
	}

	scriptPath := filepath.Join(ws.AuthScriptsDir(), "test.sh")
	if _, statErr := os.Stat(scriptPath); os.IsNotExist(statErr) {
		t.Error("auth script was not copied to workspace")
	}
}

func TestReadConfigFromExport(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	cfgContent := "specs:\n  - domain: test\n"
	if err := os.WriteFile(filepath.Join(exportDir, "swag2mcp.yaml"), []byte(cfgContent), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	data, err := ReadConfigFromExport(exportDir)
	if err != nil {
		t.Fatalf("ReadConfigFromExport() = %v", err)
	}
	if string(data) != cfgContent {
		t.Errorf("content = %q, want %q", string(data), cfgContent)
	}
}

func TestEnsureZipExt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"backup.zip", "backup.zip"},
		{"backup", "backup.zip"},
		{"backup.tar.gz", "backup.tar.gz.zip"},
	}

	for _, tt := range tests {
		result := ensureZipExt(tt.input)
		if result != tt.expected {
			t.Errorf("ensureZipExt(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
