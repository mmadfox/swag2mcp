package workspace

import (
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
