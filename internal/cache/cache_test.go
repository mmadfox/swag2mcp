package cache

import (
	"context"
	"errors"
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

func TestResolve_localPathOutsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Should return a path inside the cache directory
	cacheDir := filepath.Join(dir, CacheDirName)
	assert.True(t, stringsHasPrefix(got, cacheDir), "expected path in cache dir %q, got %q", cacheDir, got)

	data, err := os.ReadFile(got)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestResolve_localPathInsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specsDir := filepath.Join(dir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Should return the original path, not cached
	assert.Equal(t, specFile, got)
}

func TestResolve_localPathSpecsRelative(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "layers.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(wsDir)

	got, err := c.Resolve(context.Background(), "specs/layers.yaml")
	require.NoError(t, err)

	// Should resolve to workspaceDir/specs/layers.yaml and return as-is (not cached)
	assert.Equal(t, specFile, got)
}

func TestResolve_localPathSpecsRelativeNested(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName, "subdir")
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "api.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(wsDir)

	got, err := c.Resolve(context.Background(), "specs/subdir/api.yaml")
	require.NoError(t, err)

	assert.Equal(t, specFile, got)
}

func TestResolve_fileURLOutsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got, err := c.Resolve(context.Background(), "file://"+specFile)
	require.NoError(t, err)

	cacheDir := filepath.Join(dir, CacheDirName)
	assert.True(t, stringsHasPrefix(got, cacheDir), "expected path in cache dir %q, got %q", cacheDir, got)

	data, err := os.ReadFile(got)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestResolve_fileURLInsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specsDir := filepath.Join(dir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got, err := c.Resolve(context.Background(), "file://"+specFile)
	require.NoError(t, err)

	// Should return the original path, not cached
	assert.Equal(t, specFile, got)
}

func TestResolve_emptyLocation(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, err := c.Resolve(context.Background(), "")
	require.Error(t, err, "expected error for empty location")
}

func TestResolve_downloadAndCache(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("openapi: 3.0.0"))
	}))
	defer srv.Close()

	c := New(t.TempDir())

	got, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	data, err := os.ReadFile(got)
	require.NoError(t, err)
	assert.Equal(t, "openapi: 3.0.0", string(data))

	// Meta file should exist
	metaPath := got[:len(got)-len(".spec")] + ".meta"
	_, statErr := os.Stat(metaPath)
	require.NoError(t, statErr, "meta file not found")

	meta, readErr := readMeta(metaPath)
	require.NoError(t, readErr)
	assert.Equal(t, srv.URL, meta.Source)
	assert.Equal(t, "url", meta.SourceType)
}

func TestResolve_serveFromCache(t *testing.T) {
	t.Parallel()
	var callCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("openapi: 3.0.0"))
	}))
	defer srv.Close()

	c := New(t.TempDir())

	// First call — download
	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	// Second call — should serve from cache
	got2, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	assert.Equal(t, 1, callCount, "expected 1 server call")

	data, err := os.ReadFile(got2)
	require.NoError(t, err)
	assert.Equal(t, "openapi: 3.0.0", string(data))
}

func TestResolve_reDownloadAfterExpiry(t *testing.T) {
	t.Parallel()
	var callCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("openapi: 3.0.0"))
	}))
	defer srv.Close()

	c := New(t.TempDir())

	// First call — download
	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Manually expire the meta
	hash := cacheKey(srv.URL)
	metaPath := filepath.Join(c.dir, hash+".meta")
	require.NoError(t, writeMeta(metaPath, fileMeta{
		Source:     srv.URL,
		SourceType: "url",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}))

	// Second call — should re-download
	_, err = c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)
	assert.Equal(t, 2, callCount, "expected 2 server calls")
}

func TestResolve_downloadError(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, err := c.Resolve(context.Background(), "https://nonexistent.example.com/spec")
	require.Error(t, err, "expected error for non-existent URL")
}

func TestResolve_non200Status(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(t.TempDir())
	_, err := c.Resolve(context.Background(), srv.URL)
	require.Error(t, err, "expected error for 404")
}

func TestResolve_cacheDirCreated(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer srv.Close()

	baseDir := t.TempDir()
	c := New(baseDir)

	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	cacheDir := filepath.Join(baseDir, CacheDirName)
	_, statErr := os.Stat(cacheDir)
	require.NoError(t, statErr, "cache dir not created")
}

func TestResolve_localFileOutsideSpecsCaching(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Meta file should exist
	metaPath := got[:len(got)-len(".spec")] + ".meta"
	meta, readErr := readMeta(metaPath)
	require.NoError(t, readErr)
	assert.Equal(t, "local", meta.SourceType)
	assert.False(t, meta.ModTime.IsZero(), "expected non-zero mod_time for local file")
}

func TestResolve_localFileOutsideSpecsServeFromCache(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	// First call — cache
	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Second call — should serve from cache (file unchanged)
	got2, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	assert.Equal(t, got1, got2, "expected same cache path")

	data, err := os.ReadFile(got2)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestResolve_localFileOutsideSpecsReCacheOnModtimeChange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	// First call — cache
	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Modify the file
	require.NoError(t, os.WriteFile(specFile, []byte("world"), 0644))

	// Second call — should re-cache because modtime changed
	got2, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	// Path should be the same (same hash key), but content updated
	assert.Equal(t, got1, got2, "expected same cache path")

	data, err := os.ReadFile(got2)
	require.NoError(t, err)
	assert.Equal(t, "world", string(data))
}

func TestNew(t *testing.T) {
	t.Parallel()
	c := New("/tmp/swag2mcp")
	expected := filepath.Join("/tmp/swag2mcp", CacheDirName)
	assert.Equal(t, expected, c.dir)
	expectedSpecs := filepath.Join("/tmp/swag2mcp", SpecsDirName)
	assert.Equal(t, expectedSpecs, c.specsDir)
}

func TestCacheKey(t *testing.T) {
	t.Parallel()
	k1 := cacheKey("https://example.com/spec.yaml")
	k2 := cacheKey("https://example.com/spec.yaml")
	k3 := cacheKey("https://example.com/other.yaml")

	assert.Equal(t, k1, k2, "same URL should produce same key")
	assert.NotEqual(t, k1, k3, "different URLs should produce different keys")
	assert.Len(t, k1, 32, "expected 32-char hex key")
}

func TestMeta(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	metaPath := filepath.Join(dir, "test.meta")

	m := fileMeta{
		Source:     "https://example.com/spec",
		SourceType: "url",
		CachedAt:   time.Now(),
		TTLSec:     3600,
	}
	require.NoError(t, writeMeta(metaPath, m))

	got, err := readMeta(metaPath)
	require.NoError(t, err)

	assert.Equal(t, m.Source, got.Source)
	assert.Equal(t, m.SourceType, got.SourceType)
	assert.Equal(t, m.TTLSec, got.TTLSec)
}

func TestMeta_expired(t *testing.T) {
	t.Parallel()
	m := fileMeta{
		CachedAt: time.Now().Add(-2 * time.Hour),
		TTLSec:   3600, // 1 hour
	}
	assert.True(t, m.IsExpired(), "expected expired meta")

	m2 := fileMeta{
		CachedAt: time.Now().Add(-30 * time.Minute),
		TTLSec:   3600, // 1 hour
	}
	assert.False(t, m2.IsExpired(), "expected non-expired meta")
}

func TestMeta_readNotFound(t *testing.T) {
	t.Parallel()
	_, err := readMeta("/nonexistent/path")
	require.Error(t, err, "expected error for non-existent meta")
}

func TestFileURIToPath(t *testing.T) {
	t.Parallel()
	t.Run("unix path", func(t *testing.T) {
		t.Parallel()
		got, err := fileURIToPath("file:///home/user/spec.yaml")
		require.NoError(t, err)
		assert.Equal(t, "/home/user/spec.yaml", got)
	})
	t.Run("windows path", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS != "windows" {
			t.Skip("skipping windows test on non-windows")
		}
		got, err := fileURIToPath("file:///C:/Users/user/spec.yaml")
		require.NoError(t, err)
		assert.Equal(t, `C:\Users\user\spec.yaml`, got)
	})
	t.Run("bad scheme", func(t *testing.T) {
		t.Parallel()
		_, err := fileURIToPath("https://example.com/spec")
		require.Error(t, err, "expected error for non-file scheme")
	})
}

func TestFileURIToPath_badScheme(t *testing.T) {
	t.Parallel()
	_, err := fileURIToPath("https://example.com/spec")
	require.Error(t, err, "expected error for non-file scheme")
}

func TestExists_emptyLocation(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "")
	require.Error(t, err, "expected error for empty location")
}

func TestExists_localFileExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(dir)
	require.NoError(t, c.Exists(context.Background(), specFile))
}

func TestExists_localFileSpecsRelative(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	require.NoError(t, os.MkdirAll(specsDir, 0750))
	specFile := filepath.Join(specsDir, "layers.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(wsDir)
	require.NoError(t, c.Exists(context.Background(), "specs/layers.yaml"))
}

func TestExists_localFileNotFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "/nonexistent/path/spec.yaml")
	require.Error(t, err)

	var locErr *LocationError
	require.True(t, errors.As(err, &locErr), "expected LocationError")
	assert.Equal(t, "file", locErr.Type)
}

func TestExists_fileURLExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("data"), 0644))

	c := New(dir)
	require.NoError(t, c.Exists(context.Background(), "file://"+specFile))
}

func TestExists_fileURLNotFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "file:///nonexistent/path/spec.yaml")
	require.Error(t, err)

	var locErr *LocationError
	require.True(t, errors.As(err, &locErr), "expected LocationError")
	assert.Equal(t, "file", locErr.Type)
}

func TestExists_fileURLBadScheme(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "https://example.com/spec")
	require.Error(t, err)

	// HTTPS URLs go through existsURL, not file URL path
	var locErr *LocationError
	require.True(t, errors.As(err, &locErr), "expected LocationError")
}

func TestExists_urlOK(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(t.TempDir())
	require.NoError(t, c.Exists(context.Background(), srv.URL))
}

func TestExists_urlNotFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(t.TempDir())
	err := c.Exists(context.Background(), srv.URL)
	require.Error(t, err)

	var locErr *LocationError
	require.True(t, errors.As(err, &locErr), "expected LocationError")
	assert.Equal(t, "url", locErr.Type)
}

func TestExists_urlUnreachable(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "https://nonexistent.example.com/spec")
	require.Error(t, err)

	var locErr *LocationError
	require.True(t, errors.As(err, &locErr), "expected LocationError")
	assert.Equal(t, "url", locErr.Type)
}

func TestExists_urlCached(t *testing.T) {
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

	// Second call — should use cache
	require.NoError(t, c.Exists(context.Background(), srv.URL))

	// Resolve made 1 call, Exists should not make another
	assert.Equal(t, 1, callCount, "expected 1 server call")
}

func TestExists_LocationErrorFields(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())

	err := c.Exists(context.Background(), "/nonexistent/path")
	var locErr *LocationError
	require.True(t, errors.As(err, &locErr), "expected LocationError")
	assert.Equal(t, "/nonexistent/path", locErr.Location)
	assert.Equal(t, "file", locErr.Type)
	require.NotNil(t, locErr.Err)
}

// stringsHasPrefix is a helper to avoid importing strings in tests
// for a single call.
func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
