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
)

func TestResolve_localPathOutsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Should return a path inside the cache directory
	cacheDir := filepath.Join(dir, CacheDirName)
	if !stringsHasPrefix(got, cacheDir) {
		t.Errorf("expected path in cache dir %q, got %q", cacheDir, got)
	}

	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("got content %q, want %q", string(data), "hello")
	}
}

func TestResolve_localPathInsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specsDir := filepath.Join(dir, SpecsDirName)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatal(err)
	}
	specFile := filepath.Join(specsDir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Should return the original path, not cached
	if got != specFile {
		t.Errorf("expected original path %q, got %q", specFile, got)
	}
}

func TestResolve_localPathSpecsRelative(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatal(err)
	}
	specFile := filepath.Join(specsDir, "layers.yaml")
	if err := os.WriteFile(specFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(wsDir)

	got, err := c.Resolve(context.Background(), "specs/layers.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Should resolve to workspaceDir/specs/layers.yaml and return as-is (not cached)
	if got != specFile {
		t.Errorf("expected %q, got %q", specFile, got)
	}
}

func TestResolve_localPathSpecsRelativeNested(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName, "subdir")
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatal(err)
	}
	specFile := filepath.Join(specsDir, "api.yaml")
	if err := os.WriteFile(specFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(wsDir)

	got, err := c.Resolve(context.Background(), "specs/subdir/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if got != specFile {
		t.Errorf("expected %q, got %q", specFile, got)
	}
}

func TestResolve_fileURLOutsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve(context.Background(), "file://"+specFile)
	if err != nil {
		t.Fatal(err)
	}

	cacheDir := filepath.Join(dir, CacheDirName)
	if !stringsHasPrefix(got, cacheDir) {
		t.Errorf("expected path in cache dir %q, got %q", cacheDir, got)
	}

	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("got content %q, want %q", string(data), "hello")
	}
}

func TestResolve_fileURLInsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specsDir := filepath.Join(dir, SpecsDirName)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatal(err)
	}
	specFile := filepath.Join(specsDir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve(context.Background(), "file://"+specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Should return the original path, not cached
	if got != specFile {
		t.Errorf("expected original path %q, got %q", specFile, got)
	}
}

func TestResolve_emptyLocation(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, err := c.Resolve(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty location")
	}
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
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(got)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "openapi: 3.0.0" {
		t.Errorf("got content %q, want %q", string(data), "openapi: 3.0.0")
	}

	// Meta file should exist
	metaPath := got[:len(got)-len(".spec")] + ".meta"
	if _, statErr := os.Stat(metaPath); statErr != nil {
		t.Errorf("meta file not found: %v", statErr)
	}

	meta, readErr := readMeta(metaPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if meta.Source != srv.URL {
		t.Errorf("got source %q, want %q", meta.Source, srv.URL)
	}
	if meta.SourceType != "url" {
		t.Errorf("got source_type %q, want %q", meta.SourceType, "url")
	}
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
	if err != nil {
		t.Fatal(err)
	}

	// Second call — should serve from cache
	got2, err := c.Resolve(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 server call, got %d", callCount)
	}

	data, err := os.ReadFile(got2)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "openapi: 3.0.0" {
		t.Errorf("got %q, want %q", string(data), "openapi: 3.0.0")
	}
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
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 1 {
		t.Fatalf("expected 1 call, got %d", callCount)
	}

	// Manually expire the meta
	hash := cacheKey(srv.URL)
	metaPath := filepath.Join(c.dir, hash+".meta")
	if mErr := writeMeta(metaPath, fileMeta{
		Source:     srv.URL,
		SourceType: "url",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}); mErr != nil {
		t.Fatal(mErr)
	}

	// Second call — should re-download
	_, err = c.Resolve(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 server calls, got %d", callCount)
	}
}

func TestResolve_downloadError(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, err := c.Resolve(context.Background(), "https://nonexistent.example.com/spec")
	if err == nil {
		t.Fatal("expected error for non-existent URL")
	}
}

func TestResolve_non200Status(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(t.TempDir())
	_, err := c.Resolve(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 404")
	}
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
	if err != nil {
		t.Fatal(err)
	}

	cacheDir := filepath.Join(baseDir, CacheDirName)
	if _, statErr := os.Stat(cacheDir); statErr != nil {
		t.Errorf("cache dir not created: %v", statErr)
	}
}

func TestResolve_localFileOutsideSpecsCaching(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Meta file should exist
	metaPath := got[:len(got)-len(".spec")] + ".meta"
	meta, readErr := readMeta(metaPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if meta.SourceType != "local" {
		t.Errorf("got source_type %q, want %q", meta.SourceType, "local")
	}
	if meta.ModTime.IsZero() {
		t.Error("expected non-zero mod_time for local file")
	}
}

func TestResolve_localFileOutsideSpecsServeFromCache(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	// First call — cache
	got1, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Second call — should serve from cache (file unchanged)
	got2, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	if got1 != got2 {
		t.Errorf("expected same cache path, got %q and %q", got1, got2)
	}

	data, err := os.ReadFile(got2)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", string(data), "hello")
	}
}

func TestResolve_localFileOutsideSpecsReCacheOnModtimeChange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	// First call — cache
	got1, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Modify the file
	if err = os.WriteFile(specFile, []byte("world"), 0644); err != nil {
		t.Fatal(err)
	}

	// Second call — should re-cache because modtime changed
	got2, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Path should be the same (same hash key), but content updated
	if got1 != got2 {
		t.Errorf("expected same cache path, got %q and %q", got1, got2)
	}

	data, err := os.ReadFile(got2)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "world" {
		t.Errorf("got %q, want %q", string(data), "world")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	c := New("/tmp/swag2mcp")
	expected := filepath.Join("/tmp/swag2mcp", CacheDirName)
	if c.dir != expected {
		t.Errorf("got dir %q, want %q", c.dir, expected)
	}
	expectedSpecs := filepath.Join("/tmp/swag2mcp", SpecsDirName)
	if c.specsDir != expectedSpecs {
		t.Errorf("got specsDir %q, want %q", c.specsDir, expectedSpecs)
	}
}

func TestCacheKey(t *testing.T) {
	t.Parallel()
	k1 := cacheKey("https://example.com/spec.yaml")
	k2 := cacheKey("https://example.com/spec.yaml")
	k3 := cacheKey("https://example.com/other.yaml")

	if k1 != k2 {
		t.Error("same URL should produce same key")
	}
	if k1 == k3 {
		t.Error("different URLs should produce different keys")
	}
	if len(k1) != 32 {
		t.Errorf("expected 32-char hex key, got %d", len(k1))
	}
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
	if err := writeMeta(metaPath, m); err != nil {
		t.Fatal(err)
	}

	got, err := readMeta(metaPath)
	if err != nil {
		t.Fatal(err)
	}

	if got.Source != m.Source {
		t.Errorf("got Source %q, want %q", got.Source, m.Source)
	}
	if got.SourceType != m.SourceType {
		t.Errorf("got SourceType %q, want %q", got.SourceType, m.SourceType)
	}
	if got.TTLSec != m.TTLSec {
		t.Errorf("got TTL %d, want %d", got.TTLSec, m.TTLSec)
	}
}

func TestMeta_expired(t *testing.T) {
	t.Parallel()
	m := fileMeta{
		CachedAt: time.Now().Add(-2 * time.Hour),
		TTLSec:   3600, // 1 hour
	}
	if !m.IsExpired() {
		t.Error("expected expired meta")
	}

	m2 := fileMeta{
		CachedAt: time.Now().Add(-30 * time.Minute),
		TTLSec:   3600, // 1 hour
	}
	if m2.IsExpired() {
		t.Error("expected non-expired meta")
	}
}

func TestMeta_readNotFound(t *testing.T) {
	t.Parallel()
	_, err := readMeta("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for non-existent meta")
	}
}

func TestFileURIToPath(t *testing.T) {
	t.Parallel()
	t.Run("unix path", func(t *testing.T) {
		t.Parallel()
		got, err := fileURIToPath("file:///home/user/spec.yaml")
		if err != nil {
			t.Fatal(err)
		}
		want := "/home/user/spec.yaml"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("windows path", func(t *testing.T) {
		t.Parallel()
		if runtime.GOOS != "windows" {
			t.Skip("skipping windows test on non-windows")
		}
		got, err := fileURIToPath("file:///C:/Users/user/spec.yaml")
		if err != nil {
			t.Fatal(err)
		}
		want := `C:\Users\user\spec.yaml`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("bad scheme", func(t *testing.T) {
		t.Parallel()
		_, err := fileURIToPath("https://example.com/spec")
		if err == nil {
			t.Fatal("expected error for non-file scheme")
		}
	})
}

func TestFileURIToPath_badScheme(t *testing.T) {
	t.Parallel()
	_, err := fileURIToPath("https://example.com/spec")
	if err == nil {
		t.Fatal("expected error for non-file scheme")
	}
}

func TestExists_emptyLocation(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty location")
	}
}

func TestExists_localFileExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)
	if err := c.Exists(context.Background(), specFile); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_localFileSpecsRelative(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatal(err)
	}
	specFile := filepath.Join(specsDir, "layers.yaml")
	if err := os.WriteFile(specFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(wsDir)
	if err := c.Exists(context.Background(), "specs/layers.yaml"); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_localFileNotFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "/nonexistent/path/spec.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
	var locErr *LocationError
	if !errors.As(err, &locErr) {
		t.Fatal("expected LocationError")
	}
	if locErr.Type != "file" {
		t.Errorf("expected type 'file', got %q", locErr.Type)
	}
}

func TestExists_fileURLExists(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)
	err := c.Exists(context.Background(), "file://"+specFile)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_fileURLNotFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "file:///nonexistent/path/spec.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
	var locErr *LocationError
	if !errors.As(err, &locErr) {
		t.Fatal("expected LocationError")
	}
	if locErr.Type != "file" {
		t.Errorf("expected type 'file', got %q", locErr.Type)
	}
}

func TestExists_fileURLBadScheme(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "https://example.com/spec")
	if err == nil {
		t.Fatal("expected error")
	}
	// HTTPS URLs go through existsURL, not file URL path
	var locErr *LocationError
	if !errors.As(err, &locErr) {
		t.Fatal("expected LocationError")
	}
}

func TestExists_urlOK(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(t.TempDir())
	if err := c.Exists(context.Background(), srv.URL); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_urlNotFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(t.TempDir())
	err := c.Exists(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error")
	}
	var locErr *LocationError
	if !errors.As(err, &locErr) {
		t.Fatal("expected LocationError")
	}
	if locErr.Type != "url" {
		t.Errorf("expected type 'url', got %q", locErr.Type)
	}
}

func TestExists_urlUnreachable(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists(context.Background(), "https://nonexistent.example.com/spec")
	if err == nil {
		t.Fatal("expected error")
	}
	var locErr *LocationError
	if !errors.As(err, &locErr) {
		t.Fatal("expected LocationError")
	}
	if locErr.Type != "url" {
		t.Errorf("expected type 'url', got %q", locErr.Type)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	// Second call — should use cache
	if existsErr := c.Exists(context.Background(), srv.URL); existsErr != nil {
		t.Errorf("expected nil, got %v", existsErr)
	}

	// Resolve made 1 call, Exists should not make another
	if callCount != 1 {
		t.Errorf("expected 1 server call, got %d", callCount)
	}
}

func TestExists_LocationErrorFields(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())

	err := c.Exists(context.Background(), "/nonexistent/path")
	var locErr *LocationError
	if !errors.As(err, &locErr) {
		t.Fatal("expected LocationError")
	}
	if locErr.Location != "/nonexistent/path" {
		t.Errorf("expected Location %q, got %q", "/nonexistent/path", locErr.Location)
	}
	if locErr.Type != "file" {
		t.Errorf("expected Type 'file', got %q", locErr.Type)
	}
	if locErr.Err == nil {
		t.Fatal("expected non-nil Err")
	}
}

// stringsHasPrefix is a helper to avoid importing strings in tests
// for a single call.
func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
