package cache

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetWorkspaceDir(t *testing.T) {
	t.Parallel()

	c := New("/old/workspace")
	c.SetWorkspaceDir("/new/workspace")

	expectedCache := filepath.Join("/new/workspace", CacheDirName)
	assert.Equal(t, expectedCache, c.dir)

	expectedSpecs := filepath.Join("/new/workspace", SpecsDirName)
	assert.Equal(t, expectedSpecs, c.specsDir)

	assert.Equal(t, "/new/workspace", c.workspaceDir)
}

func TestExpandTilde(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		input    string
		expected string
	}{
		{"~/config.yaml", filepath.Join(home, "config.yaml")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := env.ExpandTilde(tt.input)
		assert.Equal(t, tt.expected, result, "ExpandTilde(%q)", tt.input)
	}
}

func TestRandomTTL(t *testing.T) {
	t.Parallel()

	for range 100 {
		ttl := randomTTL()
		assert.GreaterOrEqual(t, ttl, MinTTL, "TTL < MinTTL")
		assert.LessOrEqual(t, ttl, MaxTTL, "TTL > MaxTTL")
	}
}

func TestLocationError_Error(t *testing.T) {
	t.Parallel()

	err := &LocationError{
		Location: "/path/to/spec.yaml",
		Type:     "file",
		Err:      errors.New("file not found"),
	}
	assert.NotEmpty(t, err.Error())
}

func TestLocationError_Unwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner error")
	err := &LocationError{
		Location: "/path",
		Type:     "file",
		Err:      inner,
	}
	assert.True(t, errors.Is(err, inner), "errors.Is() should match inner error")
}

func TestDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	cli := defaultHTTPClient()
	require.NotNil(t, cli, "defaultHTTPClient() returned nil")
	require.NotNil(t, cli.cli, "http.Client is nil")
}

func TestHTTPClientGet_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	data, err := cli.Get(context.Background(), srv.URL)
	require.NoError(t, err, "Get()")
	assert.Equal(t, "hello world", string(data))
}

func TestHTTPClientGet_Non200(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	_, err := cli.Get(context.Background(), srv.URL)
	require.Error(t, err, "expected error for 404")
}

func TestHTTPClientGet_EmptyBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	_, err := cli.Get(context.Background(), srv.URL)
	require.Error(t, err, "expected error for empty body")
}

func TestNormalizeLocation_URL(t *testing.T) {
	t.Parallel()

	normalized, stype, err := normalizeLocation("https://example.com/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/spec.yaml", normalized)
	assert.Equal(t, sourceURL, stype)
}

func TestNormalizeLocation_FileURL(t *testing.T) {
	t.Parallel()

	normalized, stype, err := normalizeLocation("file:///home/user/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceLocal, stype)
	assert.Equal(t, "/home/user/spec.yaml", normalized)
}

func TestNormalizeLocation_LocalPath(t *testing.T) {
	t.Parallel()

	normalized, stype, err := normalizeLocation("/absolute/path/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceLocal, stype)
	assert.Equal(t, "/absolute/path/spec.yaml", normalized)
}

func TestNormalizeLocation_InvalidFileURL(t *testing.T) {
	t.Parallel()

	_, _, err := normalizeLocation("ftp://example.com/spec.yaml")
	require.NoError(t, err, "ftp falls through to local path")
}

func TestIsInsideSpecs_Inside(t *testing.T) {
	t.Parallel()

	c := New("/workspace")
	path := filepath.Join("/workspace", SpecsDirName, "subdir", "spec.yaml")
	assert.True(t, c.isInsideSpecs(path))
}

func TestIsInsideSpecs_Outside(t *testing.T) {
	t.Parallel()

	c := New("/workspace")
	assert.False(t, c.isInsideSpecs("/other/path/spec.yaml"))
}

func TestIsInsideSpecs_EmptySpecsDir(t *testing.T) {
	t.Parallel()

	c := &Cache{specsDir: ""}
	assert.False(t, c.isInsideSpecs("/any/path"))
}

func TestResolveSpecsPath_Found(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "api.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0600))

	c := New(wsDir)
	got, ok := c.resolveSpecsPath("specs/api.yaml")
	require.True(t, ok, "resolveSpecsPath() = false, want true")
	assert.Equal(t, specFile, got)
}

func TestResolveSpecsPath_NotFound(t *testing.T) {
	t.Parallel()

	c := New(t.TempDir())
	_, ok := c.resolveSpecsPath("specs/nonexistent.yaml")
	assert.False(t, ok, "resolveSpecsPath() = true, want false")
}

func TestResolveSpecsPath_NotSpecsPath(t *testing.T) {
	t.Parallel()

	c := New(t.TempDir())
	_, ok := c.resolveSpecsPath("/absolute/path.yaml")
	assert.False(t, ok, "resolveSpecsPath() = true, want false")
}

func TestResolve_localFileOutsideSpecsReCacheOnExpiredMeta(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Manually expire the meta
	hash := cacheKey(specFile)
	metaPath := filepath.Join(c.dir, hash+".meta")
	require.NoError(t, writeMeta(metaPath, fileMeta{
		Source:     specFile,
		SourceType: "local",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}))

	got2, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	assert.Equal(t, got1, got2, "expected same cache path")

	data, err := os.ReadFile(got2)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestResolve_localFileOutsideSpecsReCacheOnMissingSpecFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	require.NoError(t, os.Remove(got1))

	got2, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	assert.Equal(t, got1, got2, "expected same cache path")

	data, err := os.ReadFile(got2)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestExists_urlCachedExpired(t *testing.T) {
	t.Parallel()
	var callCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer srv.Close()

	c := New(t.TempDir())

	// First call — download via Resolve
	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	// Manually expire the meta
	hash := cacheKey(srv.URL)
	metaPath := filepath.Join(c.dir, hash+".meta")
	require.NoError(t, writeMeta(metaPath, fileMeta{
		Source:     srv.URL,
		SourceType: "url",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}))

	// Exists should fall through to HEAD request
	require.NoError(t, c.Exists(context.Background(), srv.URL))
}

func TestExists_urlCachedMissingSpecFile(t *testing.T) {
	t.Parallel()
	var callCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer srv.Close()

	c := New(t.TempDir())

	// First call — download via Resolve
	got, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	require.NoError(t, os.Remove(got))

	require.NoError(t, c.Exists(context.Background(), srv.URL))
}

func TestExists_fileURLSpecsRelative(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "layers.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(wsDir)
	require.NoError(t, c.Exists(context.Background(), "file://"+specFile))
}

func TestExists_tildePath(t *testing.T) {
	t.Parallel()
	home, homeErr := os.UserHomeDir()
	require.NoError(t, homeErr)

	tmpFile := filepath.Join(home, ".swag2mcp-test-exists")
	require.NoError(t, os.WriteFile(tmpFile, []byte("data"), 0600))
	defer os.Remove(tmpFile)

	c := New(t.TempDir())
	require.NoError(t, c.Exists(context.Background(), "~/.swag2mcp-test-exists"))
}

func TestExists_relativePath(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	t.Chdir(dir)

	c := New(dir)
	require.NoError(t, c.Exists(context.Background(), "spec.yaml"))
}

func TestClassifyLocation_httpURL(t *testing.T) {
	t.Parallel()
	stype, path, err := classifyLocation("http://example.com/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceURL, stype)
	assert.Equal(t, "http://example.com/spec.yaml", path)
}

func TestClassifyLocation_httpsURL(t *testing.T) {
	t.Parallel()
	stype, path, err := classifyLocation("https://example.com/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceURL, stype)
	assert.Equal(t, "https://example.com/spec.yaml", path)
}

func TestClassifyLocation_fileURL(t *testing.T) {
	t.Parallel()
	stype, path, err := classifyLocation("file:///home/user/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceLocal, stype)
	assert.Equal(t, "/home/user/spec.yaml", path)
}

func TestClassifyLocation_fileURLBadScheme(t *testing.T) {
	t.Parallel()
	_, _, err := classifyLocation("file://invalid|path")
	require.Error(t, err)
	var locErr *LocationError
	require.True(t, errors.As(err, &locErr))
}

func TestClassifyLocation_localPath(t *testing.T) {
	t.Parallel()
	stype, path, err := classifyLocation("/absolute/path/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceLocal, stype)
	assert.Equal(t, "/absolute/path/spec.yaml", path)
}

func TestClassifyLocation_tildePath(t *testing.T) {
	t.Parallel()
	stype, path, err := classifyLocation("~/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceLocal, stype)
	assert.Contains(t, path, "spec.yaml")
}

func TestExistsFile_found(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(dir)
	require.NoError(t, c.existsFile(specFile))
}

func TestExistsFile_notFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.existsFile("/nonexistent/path/spec.yaml")
	require.Error(t, err)
	var locErr *LocationError
	require.True(t, errors.As(err, &locErr))
	assert.Equal(t, "file", locErr.Type)
}

func TestExistsFile_specsRelative(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "api.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0600))

	c := New(wsDir)
	require.NoError(t, c.existsFile("specs/api.yaml"))
}

func TestClassifyLocation_ftpAsLocal(t *testing.T) {
	t.Parallel()
	stype, path, err := classifyLocation("ftp://example.com/spec.yaml")
	require.NoError(t, err)
	assert.Equal(t, sourceLocal, stype)
	assert.Equal(t, "ftp://example.com/spec.yaml", path)
}

func TestLoadSource_unknownType(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, _, err := c.loadSource(context.Background(), "/tmp/test.yaml", sourceType("unknown"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown source type")
}
