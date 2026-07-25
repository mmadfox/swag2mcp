package workspace

import (
	"archive/zip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultRoot(t *testing.T) {
	ws, err := New("")
	require.NoError(t, err)
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName)
	assert.Equal(t, want, ws.Root())
}

func TestNew_CustomRoot(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, tmpDir, ws.Root())
}

func TestNew_RelativePath(t *testing.T) {
	ws, err := New("relative/path")
	require.NoError(t, err)
	abs, _ := filepath.Abs("relative/path")
	assert.Equal(t, abs, ws.Root())
}

func TestNewFromBase_Empty(t *testing.T) {
	ws, err := NewFromBase("")
	require.NoError(t, err)
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName)
	assert.Equal(t, want, ws.Root())
}

func TestNewFromBase_Custom(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := NewFromBase(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, tmpDir, ws.Root())
}

func TestSub(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", "mydir")
	assert.Equal(t, want, ws.Sub("mydir"))
}

func TestCacheDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirCache)
	assert.Equal(t, want, ws.CacheDir())
}

func TestSpecsDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirSpecs)
	assert.Equal(t, want, ws.SpecsDir())
}

func TestResponsesDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirResponses)
	assert.Equal(t, want, ws.ResponsesDir())
}

func TestAuthScriptsDir(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirAuthScripts)
	assert.Equal(t, want, ws.AuthScriptsDir())
}

func TestDefaultRoot(t *testing.T) {
	root := DefaultRoot()
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName)
	assert.Equal(t, want, root)
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, DefaultRootName, "swag2mcp.yaml")
	assert.Equal(t, want, path)
}

func TestConfigPathIn(t *testing.T) {
	path := ConfigPathIn("/custom/workspace")
	want := filepath.Join("/custom/workspace", "swag2mcp.yaml")
	assert.Equal(t, want, path)
}

func TestConfigPath(t *testing.T) {
	ws, _ := New("/root")
	path := ws.ConfigPath()
	want := filepath.Join("/root", "swag2mcp.yaml")
	assert.Equal(t, want, path)
}

func TestConfigExists_True(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	cfgPath := ws.ConfigPath()
	require.NoError(t, os.WriteFile(cfgPath, []byte("specs: []"), 0600))
	assert.True(t, ws.ConfigExists())
}

func TestConfigExists_False(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	assert.False(t, ws.ConfigExists())
}

func TestConfigNotExists_True(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	assert.True(t, ws.ConfigNotExists())
}

func TestConfigNotExists_False(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	cfgPath := ws.ConfigPath()
	require.NoError(t, os.WriteFile(cfgPath, []byte("specs: []"), 0600))
	assert.False(t, ws.ConfigNotExists())
}

func TestInit(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())
	for _, dir := range []string{ws.Root(), ws.CacheDir(), ws.SpecsDir(), ws.ResponsesDir(), ws.AuthScriptsDir()} {
		assert.DirExists(t, dir)
	}
}

func TestInit_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())
	require.NoError(t, ws.Init())
}

func TestClean_RemovesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	cacheFile := filepath.Join(ws.CacheDir(), "cached.yaml")
	require.NoError(t, os.WriteFile(cacheFile, []byte("data"), 0600))
	respFile := filepath.Join(ws.ResponsesDir(), "response.json")
	require.NoError(t, os.WriteFile(respFile, []byte("data"), 0600))

	require.NoError(t, ws.Clean())

	assert.NoFileExists(t, cacheFile)
	assert.NoFileExists(t, respFile)
	assert.DirExists(t, ws.CacheDir())
	assert.DirExists(t, ws.ResponsesDir())
}

func TestClean_NoDirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Clean())
}

func TestClean_EmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())
	require.NoError(t, ws.Clean())
}

func TestCleanOldResponses_RemovesOldFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	oldFile := filepath.Join(ws.ResponsesDir(), "old-response.json")
	require.NoError(t, os.WriteFile(oldFile, []byte("old"), 0600))
	oldModTime := time.Now().Add(-72 * time.Hour)
	if err := os.Chtimes(oldFile, oldModTime, oldModTime); err != nil {
		t.Skipf("Chtimes not supported on this filesystem: %v", err)
	}

	freshFile := filepath.Join(ws.ResponsesDir(), "fresh-response.json")
	require.NoError(t, os.WriteFile(freshFile, []byte("fresh"), 0600))

	require.NoError(t, ws.CleanOldResponses(48*time.Hour))

	assert.NoFileExists(t, oldFile)
	assert.FileExists(t, freshFile)
}

func TestCleanOldResponses_NoDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)

	require.NoError(t, ws.CleanOldResponses(48*time.Hour))
}

func TestCleanOldResponses_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	require.NoError(t, ws.CleanOldResponses(48*time.Hour))
}

func TestCleanOldResponses_SkipsSubdirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	subDir := filepath.Join(ws.ResponsesDir(), "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0750))
	subFile := filepath.Join(subDir, "nested.json")
	require.NoError(t, os.WriteFile(subFile, []byte("nested"), 0600))

	oldModTime := time.Now().Add(-72 * time.Hour)
	if err := os.Chtimes(subDir, oldModTime, oldModTime); err != nil {
		t.Skipf("Chtimes not supported on this filesystem: %v", err)
	}

	require.NoError(t, ws.CleanOldResponses(48*time.Hour))

	assert.DirExists(t, subDir)
	assert.FileExists(t, subFile)
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
	require.NoError(t, ws.Init())

	require.NoError(t, ws.EnsureAuthScript("test-domain"))

	scriptPath := ws.AuthScriptPath("test-domain")
	assert.FileExists(t, scriptPath)
}

func TestEnsureAuthScript_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	require.NoError(t, ws.EnsureAuthScript("test-domain"))
	require.NoError(t, ws.EnsureAuthScript("test-domain"))
}

func TestEnsureAuthScript_CreatesDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)

	require.NoError(t, ws.EnsureAuthScript("test-domain"))

	assert.DirExists(t, ws.AuthScriptsDir())
}

func TestRemoveOrphanAuthScripts_RemovesOrphan(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	orphanPath := filepath.Join(ws.AuthScriptsDir(), "orphan.sh")
	require.NoError(t, os.WriteFile(orphanPath, []byte("echo test"), 0600))

	require.NoError(t, ws.RemoveOrphanAuthScripts([]string{"active"}))

	assert.NoFileExists(t, orphanPath)
}

func TestRemoveOrphanAuthScripts_KeepsActive(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	activePath := filepath.Join(ws.AuthScriptsDir(), "active.sh")
	require.NoError(t, os.WriteFile(activePath, []byte("echo test"), 0600))

	require.NoError(t, ws.RemoveOrphanAuthScripts([]string{"active"}))

	assert.FileExists(t, activePath)
}

func TestRemoveOrphanAuthScripts_NoDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)

	require.NoError(t, ws.RemoveOrphanAuthScripts([]string{"active"}))
}

func TestRemoveOrphanAuthScripts_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	require.NoError(t, ws.RemoveOrphanAuthScripts([]string{"active"}))
}

func TestRemoveOrphanAuthScripts_SkipsNonScriptFiles(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	nonScript := filepath.Join(ws.AuthScriptsDir(), "readme.txt")
	require.NoError(t, os.WriteFile(nonScript, []byte("hello"), 0600))

	require.NoError(t, ws.RemoveOrphanAuthScripts([]string{"active"}))

	assert.FileExists(t, nonScript)
}

func TestRemoveOrphanAuthScripts_SkipsSubdirs(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	subDir := filepath.Join(ws.AuthScriptsDir(), "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0750))

	require.NoError(t, ws.RemoveOrphanAuthScripts([]string{"active"}))

	assert.DirExists(t, subDir)
}

// removeContents error: ReadDir on a file returns ENOTDIR.
func TestClean_ReadDirError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	cacheDir := ws.CacheDir()
	require.NoError(t, os.RemoveAll(cacheDir))
	require.NoError(t, os.WriteFile(cacheDir, []byte("not-a-dir"), 0600))

	require.Error(t, ws.Clean())
}

// CleanOldResponses error: ReadDir on a file returns ENOTDIR.
func TestCleanOldResponses_ReadDirError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	respDir := ws.ResponsesDir()
	require.NoError(t, os.RemoveAll(respDir))
	require.NoError(t, os.WriteFile(respDir, []byte("not-a-dir"), 0600))

	require.Error(t, ws.CleanOldResponses(48*time.Hour))
}

// CleanOldResponses error: Remove fails on read-only parent dir.
func TestCleanOldResponses_RemoveError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	oldFile := filepath.Join(ws.ResponsesDir(), "old-response.json")
	require.NoError(t, os.WriteFile(oldFile, []byte("old"), 0600))
	oldModTime := time.Now().Add(-72 * time.Hour)
	if err := os.Chtimes(oldFile, oldModTime, oldModTime); err != nil {
		t.Skipf("Chtimes not supported: %v", err)
	}

	require.NoError(t, os.Chmod(ws.ResponsesDir(), 0500))
	t.Cleanup(func() { os.Chmod(ws.ResponsesDir(), 0750) })

	require.Error(t, ws.CleanOldResponses(48*time.Hour))
}

// EnsureAuthScript error: MkdirAll fails when auth_scripts dir is a file.
func TestEnsureAuthScript_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)

	authDir := ws.AuthScriptsDir()
	require.NoError(t, os.WriteFile(authDir, []byte("not-a-dir"), 0600))

	require.Error(t, ws.EnsureAuthScript("test-domain"))
}

// EnsureAuthScript error: WriteFile fails on read-only dir.
func TestEnsureAuthScript_WriteFileError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	require.NoError(t, os.Chmod(ws.AuthScriptsDir(), 0500))
	t.Cleanup(func() { os.Chmod(ws.AuthScriptsDir(), 0750) })

	require.Error(t, ws.EnsureAuthScript("test-domain"))
}

// RemoveOrphanAuthScripts error: ReadDir on a file returns ENOTDIR.
func TestRemoveOrphanAuthScripts_ReadDirError(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	authDir := ws.AuthScriptsDir()
	require.NoError(t, os.RemoveAll(authDir))
	require.NoError(t, os.WriteFile(authDir, []byte("not-a-dir"), 0600))

	require.Error(t, ws.RemoveOrphanAuthScripts([]string{"active"}))
}

// RemoveOrphanAuthScripts error: Remove fails on read-only parent dir.
func TestRemoveOrphanAuthScripts_RemoveError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	orphanPath := filepath.Join(ws.AuthScriptsDir(), "orphan.sh")
	require.NoError(t, os.WriteFile(orphanPath, []byte("echo test"), 0600))

	require.NoError(t, os.Chmod(ws.AuthScriptsDir(), 0500))
	t.Cleanup(func() { os.Chmod(ws.AuthScriptsDir(), 0750) })

	require.Error(t, ws.RemoveOrphanAuthScripts([]string{"active"}))
}

// Init error: MkdirAll fails when root dir is a file.
func TestInit_MkdirAllError(t *testing.T) {
	tmpDir := t.TempDir()

	rootFile := filepath.Join(tmpDir, "blocker")
	require.NoError(t, os.WriteFile(rootFile, []byte("block"), 0600))
	blocker, _ := New(rootFile)

	require.Error(t, blocker.Init())
}

// DefaultRoot error: [os.UserHomeDir] fails when HOME is unset.
func TestDefaultRoot_NoHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	root := DefaultRoot()
	assert.Equal(t, DefaultRootName, root)
}

// removeContents error: RemoveAll fails when parent dir is read-only.
func TestClean_RemoveAllError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod does not prevent deletion on Windows")
	}

	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	require.NoError(t, ws.Init())

	cacheFile := filepath.Join(ws.CacheDir(), "data.yaml")
	require.NoError(t, os.WriteFile(cacheFile, []byte("data"), 0600))

	require.NoError(t, os.Chmod(ws.CacheDir(), 0000))
	t.Cleanup(func() { os.Chmod(ws.CacheDir(), 0750) })

	require.Error(t, ws.Clean())
}

func TestIsEmpty_NotExists(t *testing.T) {
	ws, _ := New("/nonexistent/path/that/does/not/exist")
	empty, err := ws.IsEmpty()
	require.NoError(t, err)
	assert.True(t, empty)
}

func TestIsEmpty_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	require.NoError(t, err)
	assert.True(t, empty)
}

func TestIsEmpty_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("data"), 0644)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	require.NoError(t, err)
	assert.False(t, empty)
}

func TestIsEmpty_WithSubdir(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	require.NoError(t, err)
	assert.False(t, empty)
}

func TestIsEmpty_OnlyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "swag2mcp.yaml"), []byte("specs: []"), 0644)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	require.NoError(t, err)
	assert.True(t, empty)
}

func TestIsEmpty_ConfigAndOtherFiles(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "swag2mcp.yaml"), []byte("specs: []"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "other.txt"), []byte("data"), 0644)
	ws, _ := New(tmpDir)
	empty, err := ws.IsEmpty()
	require.NoError(t, err)
	assert.False(t, empty)
}

func TestDownloadSpec_FromLocalFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	specPath := filepath.Join(tmpDir, "test.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(specContent), 0600))

	ws, err := New(tmpDir)
	require.NoError(t, err)

	data, err := ws.DownloadSpec(context.Background(), specPath)
	require.NoError(t, err)
	assert.Equal(t, specContent, string(data))
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
	require.NoError(t, err)

	data, err := ws.DownloadSpec(context.Background(), srv.URL)
	require.NoError(t, err)
	assert.Equal(t, specContent, string(data))
}

func TestDownloadSpec_EmptySource(t *testing.T) {
	t.Parallel()

	ws, err := New(t.TempDir())
	require.NoError(t, err)

	_, err = ws.DownloadSpec(context.Background(), "")
	require.Error(t, err)
}

func TestDownloadSpec_FromFileURL(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	specContent := "openapi: 3.0.0\ninfo:\n  title: Test\n"
	specPath := filepath.Join(tmpDir, "test.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(specContent), 0600))

	ws, err := New(tmpDir)
	require.NoError(t, err)

	data, err := ws.DownloadSpec(context.Background(), "file://"+specPath)
	require.NoError(t, err)
	assert.Equal(t, specContent, string(data))
}

func TestSpecPath(t *testing.T) {
	ws, _ := New("/root")
	want := filepath.Join("/root", DirSpecs, "myspec.yaml")
	assert.Equal(t, want, ws.SpecPath("myspec.yaml"))
}

func TestListSpecs_Empty(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)

	names, err := ws.ListSpecs()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestListSpecs_WithFiles(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	specsDir := ws.SpecsDir()
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "a.yaml"), []byte("a"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "b.yaml"), []byte("b"), 0600))

	names, err := ws.ListSpecs()
	require.NoError(t, err)
	assert.Len(t, names, 2)
}

func TestSaveSpec_Success(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	data := []byte("openapi: 3.0.0")
	path, err := ws.SaveSpec("myspec.yaml", data)
	require.NoError(t, err)

	saved, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, string(data), string(saved))
}

func TestSaveSpec_DuplicateError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	data := []byte("openapi: 3.0.0")
	_, err = ws.SaveSpec("myspec.yaml", data)
	require.NoError(t, err)

	_, err = ws.SaveSpec("myspec.yaml", data)
	require.Error(t, err)
}

func TestSaveSpec_EmptyName(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)

	_, err = ws.SaveSpec("", []byte("data"))
	require.Error(t, err)
}

func TestSaveSpec_EmptyData(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)

	_, err = ws.SaveSpec("spec.yaml", nil)
	require.Error(t, err)
}

func TestCreateExportDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)

	exportDir, err := ws.CreateExportDir()
	require.NoError(t, err)
	defer os.RemoveAll(exportDir)

	assert.DirExists(t, exportDir)
}

func TestWriteSpecToExport(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	data := []byte("openapi: 3.0.0")

	require.NoError(t, WriteSpecToExport(exportDir, "myspec.yaml", data))

	specPath := filepath.Join(exportDir, DirSpecs, "myspec.yaml")
	saved, err := os.ReadFile(specPath)
	require.NoError(t, err)
	assert.Equal(t, string(data), string(saved))
}

func TestCopyAuthScriptsToExport(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	scriptContent := "#!/bin/sh\necho test"
	scriptPath := filepath.Join(ws.AuthScriptsDir(), "test.sh")
	require.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0600))

	exportDir := t.TempDir()
	require.NoError(t, ws.CopyAuthScriptsToExport(exportDir))

	saved, err := os.ReadFile(filepath.Join(exportDir, DirAuthScripts, "test.sh"))
	require.NoError(t, err)
	assert.Equal(t, scriptContent, string(saved))
}

func TestCopyAuthScriptsToExport_NoDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)

	exportDir := t.TempDir()
	require.NoError(t, ws.CopyAuthScriptsToExport(exportDir))
}

func TestCopyAuthScriptsToExport_Subdir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	subdir := filepath.Join(ws.AuthScriptsDir(), "subdir")
	require.NoError(t, os.MkdirAll(subdir, 0750))
	require.NoError(t, os.WriteFile(filepath.Join(subdir, "nested.sh"), []byte("#!/bin/sh"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(ws.AuthScriptsDir(), "test.sh"), []byte("#!/bin/sh"), 0600))

	exportDir := t.TempDir()
	require.NoError(t, ws.CopyAuthScriptsToExport(exportDir))

	assert.FileExists(t, filepath.Join(exportDir, DirAuthScripts, "test.sh"))
	assert.NoFileExists(t, filepath.Join(exportDir, DirAuthScripts, "subdir", "nested.sh"))
}

func TestCreateMetaFile(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	require.NoError(t, CreateMetaFile(exportDir, "1.0.0"))

	metaPath := filepath.Join(exportDir, MetaFileName)
	require.FileExists(t, metaPath)
}

func TestCreateZip(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))

	outputPath := filepath.Join(t.TempDir(), "test.zip")
	require.NoError(t, CreateZip(sourceDir, outputPath))

	require.FileExists(t, outputPath)
}

func TestCreateZip_AddsExtension(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	outputPath := filepath.Join(t.TempDir(), "test")
	require.NoError(t, CreateZip(sourceDir, outputPath))

	require.FileExists(t, outputPath+".zip")
}

func TestCreateZip_WithSubdir(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0750))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "subdir", "test.txt"), []byte("hello"), 0600))

	outputPath := filepath.Join(t.TempDir(), "test.zip")
	require.NoError(t, CreateZip(sourceDir, outputPath))
	require.FileExists(t, outputPath)
}

func TestCreateZip_WalkError(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "unreadable.txt"), []byte("secret"), 0000))

	outputPath := filepath.Join(t.TempDir(), "test.zip")
	err := CreateZip(sourceDir, outputPath)
	require.Error(t, err)
}

func TestValidateZip_Valid(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	require.NoError(t, CreateMetaFile(exportDir, "1.0.0"))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(exportDir, zipPath))

	valid, err := ValidateZip(zipPath)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestValidateZip_Invalid(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	valid, err := ValidateZip(zipPath)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestValidateZip_NotAZip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "notazip.zip")
	require.NoError(t, os.WriteFile(path, []byte("not a zip"), 0600))

	_, err := ValidateZip(path)
	require.Error(t, err)
}

func TestExtractZip(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0750))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "subdir", "test.txt"), []byte("hello"), 0600))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	destDir := t.TempDir()
	require.NoError(t, ExtractZip(zipPath, destDir))

	data, err := os.ReadFile(filepath.Join(destDir, "subdir", "test.txt"))
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestExtractZip_ZipSlipPrevented(t *testing.T) {
	t.Parallel()

	destDir := t.TempDir()
	zipPath := filepath.Join(t.TempDir(), "evil.zip")

	zipFile, err := os.Create(filepath.Clean(zipPath))
	require.NoError(t, err)

	zw := zip.NewWriter(zipFile)
	_, err = zw.Create("../evil.txt")
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	require.NoError(t, zipFile.Close())

	err = ExtractZip(zipPath, destDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zip slip")
}

func TestExtractZip_CreateDirError(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0750))
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "subdir", "test.txt"), []byte("hello"), 0600))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	destDir := t.TempDir()
	require.NoError(t, os.Chmod(destDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(destDir, 0755) })

	err := ExtractZip(zipPath, destDir)
	require.Error(t, err)
}

func TestDefaultExportName(t *testing.T) {
	name := DefaultExportName()
	assert.NotEmpty(t, name)
	assert.Equal(t, ".zip", filepath.Ext(name))
}

func TestIsSwag2mcpZip(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	require.NoError(t, CreateMetaFile(exportDir, "1.0.0"))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(exportDir, zipPath))

	assert.True(t, IsSwag2mcpZip(zipPath))
}

func TestIsSwag2mcpZip_NotZip(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "test.txt")
	require.NoError(t, os.WriteFile(path, []byte("hello"), 0600))

	assert.False(t, IsSwag2mcpZip(path))
}

func TestReadMetaFromZip(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	require.NoError(t, CreateMetaFile(exportDir, "2.0.0"))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(exportDir, zipPath))

	meta, err := ReadMetaFromZip(zipPath)
	require.NoError(t, err)
	assert.Equal(t, MetaType, meta.Type)
	assert.Equal(t, "2.0.0", meta.Version)
}

func TestCreateEmptyDirsInExport(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	require.NoError(t, CreateEmptyDirsInExport(exportDir))

	for _, dir := range []string{DirCache, DirResponses} {
		path := filepath.Join(exportDir, dir)
		assert.DirExists(t, path)
	}
}

func TestCopyConfigToExport(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	cfgContent := "specs:\n  - domain: test\n"
	require.NoError(t, os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600))

	exportDir := t.TempDir()
	require.NoError(t, ws.copyConfigToExport(exportDir))

	data, err := os.ReadFile(filepath.Join(exportDir, "swag2mcp.yaml"))
	require.NoError(t, err)
	assert.Equal(t, cfgContent, string(data))
}

func TestCopySpecsToWorkspace(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	specsDir := filepath.Join(exportDir, DirSpecs)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "test.yaml"), []byte("openapi: 3.0.0"), 0600))

	require.NoError(t, ws.CopySpecsToWorkspace(exportDir))

	specPath := filepath.Join(ws.SpecsDir(), "test.yaml")
	assert.FileExists(t, specPath)
}

func TestCopyAuthScriptsToWorkspace(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	authDir := filepath.Join(exportDir, DirAuthScripts)
	require.NoError(t, os.MkdirAll(authDir, 0750))
	require.NoError(t, os.WriteFile(filepath.Join(authDir, "test.sh"), []byte("#!/bin/sh"), 0600))

	require.NoError(t, ws.CopyAuthScriptsToWorkspace(exportDir))

	scriptPath := filepath.Join(ws.AuthScriptsDir(), "test.sh")
	assert.FileExists(t, scriptPath)
}

func TestReadConfigFromExport(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	cfgContent := "specs:\n  - domain: test\n"
	require.NoError(t, os.WriteFile(filepath.Join(exportDir, "swag2mcp.yaml"), []byte(cfgContent), 0600))

	data, err := ReadConfigFromExport(exportDir)
	require.NoError(t, err)
	assert.Equal(t, cfgContent, string(data))
}

func TestWriteSpecToExport_MkdirAllError(t *testing.T) {
	t.Parallel()

	err := WriteSpecToExport("/nonexistent/readonly", "spec.yaml", []byte("data"))
	require.Error(t, err)
}

func TestWriteSpecToExport_WriteFileError(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	specsDir := filepath.Join(exportDir, DirSpecs)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	require.NoError(t, os.Chmod(specsDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(specsDir, 0755) })

	err := WriteSpecToExport(exportDir, "spec.yaml", []byte("data"))
	require.Error(t, err)
}

func TestCopyAuthScriptsToExport_ReadDirError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	require.NoError(t, os.Chmod(ws.AuthScriptsDir(), 0000))
	t.Cleanup(func() { _ = os.Chmod(ws.AuthScriptsDir(), 0755) })

	exportDir := t.TempDir()
	err = ws.CopyAuthScriptsToExport(exportDir)
	require.Error(t, err)
}

func TestCopyAuthScriptsToExport_WriteError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	scriptContent := "#!/bin/sh\necho test"
	scriptPath := filepath.Join(ws.AuthScriptsDir(), "test.sh")
	require.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0600))

	exportDir := t.TempDir()
	authDir := filepath.Join(exportDir, DirAuthScripts)
	require.NoError(t, os.MkdirAll(authDir, 0750))
	require.NoError(t, os.Chmod(authDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(authDir, 0755) })

	err = ws.CopyAuthScriptsToExport(exportDir)
	require.Error(t, err)
}

func TestCopyAuthScriptsToExport_ReadFileError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	require.NoError(t, os.WriteFile(filepath.Join(ws.AuthScriptsDir(), "test.sh"), []byte("#!/bin/sh"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(ws.AuthScriptsDir(), "broken.sh"), []byte("data"), 0000))

	exportDir := t.TempDir()
	err = ws.CopyAuthScriptsToExport(exportDir)
	require.Error(t, err)
}

func TestCreateMetaFile_WriteError(t *testing.T) {
	t.Parallel()

	exportDir := t.TempDir()
	require.NoError(t, os.Chmod(exportDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(exportDir, 0755) })

	err := CreateMetaFile(exportDir, "1.0.0")
	require.Error(t, err)
}

func TestCreateZip_CreateError(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	err := CreateZip(sourceDir, "/nonexistent/backup.zip")
	require.Error(t, err)
}

func TestCreateZip_ReadFileError(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))

	outputPath := filepath.Join(t.TempDir(), "test.zip")
	require.NoError(t, CreateZip(sourceDir, outputPath))
	require.FileExists(t, outputPath)
}

func TestValidateZip_OpenError(t *testing.T) {
	t.Parallel()

	_, err := ValidateZip("/nonexistent/backup.zip")
	require.Error(t, err)
}

func TestValidateZip_InvalidMeta(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, MetaFileName), []byte("{invalid json}"), 0600))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	valid, err := ValidateZip(zipPath)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestExtractZip_OpenError(t *testing.T) {
	t.Parallel()

	err := ExtractZip("/nonexistent/backup.zip", t.TempDir())
	require.Error(t, err)
}

func TestExtractZip_WriteError(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	destDir := t.TempDir()
	require.NoError(t, os.Chmod(destDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(destDir, 0755) })

	err := ExtractZip(zipPath, destDir)
	require.Error(t, err)
}

func TestReadMetaFromZip_OpenError(t *testing.T) {
	t.Parallel()

	_, err := ReadMetaFromZip("/nonexistent/backup.zip")
	require.Error(t, err)
}

func TestReadMetaFromZip_NotFound(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("hello"), 0600))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	_, err := ReadMetaFromZip(zipPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestCreateEmptyDirsInExport_Error(t *testing.T) {
	t.Parallel()

	err := CreateEmptyDirsInExport("/nonexistent/readonly")
	require.Error(t, err)
}

func TestCopyConfigToExport_ReadError(t *testing.T) {
	t.Parallel()

	ws := &Workspace{root: t.TempDir()}
	exportDir := t.TempDir()

	err := ws.copyConfigToExport(exportDir)
	require.Error(t, err)
}

func TestCopyConfigToExport_WriteError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	cfgContent := "specs:\n  - domain: test\n"
	require.NoError(t, os.WriteFile(ws.ConfigPath(), []byte(cfgContent), 0600))

	exportDir := t.TempDir()
	require.NoError(t, os.Chmod(exportDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(exportDir, 0755) })

	err = ws.copyConfigToExport(exportDir)
	require.Error(t, err)
}

func TestReadConfigFromExport_Error(t *testing.T) {
	t.Parallel()

	_, err := ReadConfigFromExport(t.TempDir())
	require.Error(t, err)
}

func TestCopySpecsToWorkspace_NoDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	require.NoError(t, ws.CopySpecsToWorkspace(exportDir))
}

func TestCopySpecsToWorkspace_SaveSpecError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	specsDir := filepath.Join(exportDir, DirSpecs)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	require.NoError(t, os.WriteFile(filepath.Join(specsDir, "test.yaml"), []byte("openapi: 3.0.0"), 0600))

	require.NoError(t, os.WriteFile(filepath.Join(ws.SpecsDir(), "test.yaml"), []byte("existing"), 0600))

	err = ws.CopySpecsToWorkspace(exportDir)
	require.Error(t, err)
}

func TestCopySpecsToWorkspace_ReadFileError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	specsDir := filepath.Join(exportDir, DirSpecs)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "test.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("openapi: 3.0.0"), 0600))
	require.NoError(t, os.Chmod(specFile, 0000))
	t.Cleanup(func() { _ = os.Chmod(specFile, 0600) })

	err = ws.CopySpecsToWorkspace(exportDir)
	require.Error(t, err)
}

func TestCopySpecsToWorkspace_ReadDirError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := filepath.Join(t.TempDir(), "not-a-dir")
	require.NoError(t, os.WriteFile(exportDir, []byte("not a dir"), 0600))

	err = ws.CopySpecsToWorkspace(exportDir)
	require.Error(t, err)
}

func TestCopyAuthScriptsToWorkspace_NoDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	require.NoError(t, ws.CopyAuthScriptsToWorkspace(exportDir))
}

func TestCopyAuthScriptsToWorkspace_ReadDirError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := filepath.Join(t.TempDir(), "not-a-dir")
	require.NoError(t, os.WriteFile(exportDir, []byte("not a dir"), 0600))

	err = ws.CopyAuthScriptsToWorkspace(exportDir)
	require.Error(t, err)
}

func TestCopyAuthScriptsToWorkspace_ReadFileError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	exportDir := t.TempDir()
	authDir := filepath.Join(exportDir, DirAuthScripts)
	require.NoError(t, os.MkdirAll(authDir, 0750))
	scriptFile := filepath.Join(authDir, "test.sh")
	require.NoError(t, os.WriteFile(scriptFile, []byte("#!/bin/sh"), 0600))
	require.NoError(t, os.Chmod(scriptFile, 0000))
	t.Cleanup(func() { _ = os.Chmod(scriptFile, 0600) })

	err = ws.CopyAuthScriptsToWorkspace(exportDir)
	require.Error(t, err)
}

func TestCopyAuthScriptsToWorkspace_WriteError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	require.NoError(t, os.Chmod(ws.AuthScriptsDir(), 0444))
	t.Cleanup(func() { _ = os.Chmod(ws.AuthScriptsDir(), 0755) })

	exportDir := t.TempDir()
	authDir := filepath.Join(exportDir, DirAuthScripts)
	require.NoError(t, os.MkdirAll(authDir, 0750))
	require.NoError(t, os.WriteFile(filepath.Join(authDir, "test.sh"), []byte("#!/bin/sh"), 0600))

	err = ws.CopyAuthScriptsToWorkspace(exportDir)
	require.Error(t, err)
}

func TestCopyAuthScriptsToExport_MkdirAllError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	ws, err := New(tmpDir)
	require.NoError(t, err)
	require.NoError(t, ws.Init())

	require.NoError(t, os.WriteFile(filepath.Join(ws.AuthScriptsDir(), "test.sh"), []byte("#!/bin/sh"), 0600))

	exportDir := t.TempDir()
	require.NoError(t, os.Chmod(exportDir, 0444))
	t.Cleanup(func() { _ = os.Chmod(exportDir, 0755) })

	err = ws.CopyAuthScriptsToExport(exportDir)
	require.Error(t, err)
}

func TestValidateZip_CorruptMetaEntry(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, CreateMetaFile(sourceDir, "1.0.0"))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	f, err := os.OpenFile(filepath.Clean(zipPath), os.O_RDWR, 0600)
	require.NoError(t, err)
	_, err = f.WriteAt([]byte("CORRUPT"), 50)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	_, err = ValidateZip(zipPath)
	require.Error(t, err)
}

func TestReadMetaFromZip_CorruptMeta(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	require.NoError(t, CreateMetaFile(sourceDir, "1.0.0"))

	zipPath := filepath.Join(t.TempDir(), "backup.zip")
	require.NoError(t, CreateZip(sourceDir, zipPath))

	f, err := os.OpenFile(filepath.Clean(zipPath), os.O_RDWR, 0600)
	require.NoError(t, err)
	_, err = f.WriteAt([]byte("CORRUPT"), 50)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	_, err = ReadMetaFromZip(zipPath)
	require.Error(t, err)
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
		assert.Equal(t, tt.expected, result)
	}
}
