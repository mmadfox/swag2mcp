package cache

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestResolve_localPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve(specFile)
	if err != nil {
		t.Fatal(err)
	}
	if got != specFile {
		t.Errorf("got %q, want %q", got, specFile)
	}
}

func TestResolve_fileURL(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got, err := c.Resolve("file://" + specFile)
	if err != nil {
		t.Fatal(err)
	}
	if got != specFile {
		t.Errorf("got %q, want %q", got, specFile)
	}
}

func TestResolve_emptyLocation(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, err := c.Resolve("")
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

	got, err := c.Resolve(srv.URL)
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
	_, err := c.Resolve(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Second call — should serve from cache
	got2, err := c.Resolve(srv.URL)
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
	_, err := c.Resolve(srv.URL)
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
		URL:      srv.URL,
		CachedAt: time.Now().Add(-2 * MaxTTL),
		TTLSec:   60,
	}); mErr != nil {
		t.Fatal(mErr)
	}

	// Second call — should re-download
	_, err = c.Resolve(srv.URL)
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
	_, err := c.Resolve("https://nonexistent.example.com/spec")
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
	_, err := c.Resolve(srv.URL)
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

	_, err := c.Resolve(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	cacheDir := filepath.Join(baseDir, CacheDirName)
	if _, statErr := os.Stat(cacheDir); statErr != nil {
		t.Errorf("cache dir not created: %v", statErr)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	c := New("/tmp/swag2mcp")
	expected := filepath.Join("/tmp/swag2mcp", CacheDirName)
	if c.dir != expected {
		t.Errorf("got dir %q, want %q", c.dir, expected)
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
		URL:      "https://example.com/spec",
		CachedAt: time.Now(),
		TTLSec:   3600,
	}
	if err := writeMeta(metaPath, m); err != nil {
		t.Fatal(err)
	}

	got, err := readMeta(metaPath)
	if err != nil {
		t.Fatal(err)
	}

	if got.URL != m.URL {
		t.Errorf("got URL %q, want %q", got.URL, m.URL)
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
	err := c.Exists("")
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
	if err := c.Exists(specFile); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_localFileNotFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists("/nonexistent/path/spec.yaml")
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
	err := c.Exists("file://" + specFile)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_fileURLNotFound(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	err := c.Exists("file:///nonexistent/path/spec.yaml")
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
	err := c.Exists("https://example.com/spec")
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
	if err := c.Exists(srv.URL); err != nil {
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
	err := c.Exists(srv.URL)
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
	err := c.Exists("https://nonexistent.example.com/spec")
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
	_, err := c.Resolve(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Second call — should use cache
	if existsErr := c.Exists(srv.URL); existsErr != nil {
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

	err := c.Exists("/nonexistent/path")
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
