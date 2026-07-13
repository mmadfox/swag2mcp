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
)

func TestSetWorkspaceDir(t *testing.T) {
	t.Parallel()

	c := New("/old/workspace")
	c.SetWorkspaceDir("/new/workspace")

	expectedCache := filepath.Join("/new/workspace", CacheDirName)
	if c.dir != expectedCache {
		t.Errorf("dir = %q, want %q", c.dir, expectedCache)
	}

	expectedSpecs := filepath.Join("/new/workspace", SpecsDirName)
	if c.specsDir != expectedSpecs {
		t.Errorf("specsDir = %q, want %q", c.specsDir, expectedSpecs)
	}

	if c.workspaceDir != "/new/workspace" {
		t.Errorf("workspaceDir = %q, want %q", c.workspaceDir, "/new/workspace")
	}
}

func TestExpandTilde(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() = %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/config.yaml", filepath.Join(home, "config.yaml")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := expandTilde(tt.input)
		if result != tt.expected {
			t.Errorf("expandTilde(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestRandomTTL(t *testing.T) {
	t.Parallel()

	for range 100 {
		ttl := randomTTL()
		if ttl < MinTTL {
			t.Errorf("TTL %v < MinTTL %v", ttl, MinTTL)
		}
		if ttl > MaxTTL {
			t.Errorf("TTL %v > MaxTTL %v", ttl, MaxTTL)
		}
	}
}

func TestLocationError_Error(t *testing.T) {
	t.Parallel()

	err := &LocationError{
		Location: "/path/to/spec.yaml",
		Type:     "file",
		Err:      errors.New("file not found"),
	}
	errStr := err.Error()
	if errStr == "" {
		t.Fatal("Error() returned empty string")
	}
}

func TestLocationError_Unwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner error")
	err := &LocationError{
		Location: "/path",
		Type:     "file",
		Err:      inner,
	}
	if !errors.Is(err, inner) {
		t.Error("errors.Is() should match inner error")
	}
}

func TestDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	cli := defaultHTTPClient()
	if cli == nil {
		t.Fatal("defaultHTTPClient() returned nil")
	}
	if cli.cli == nil {
		t.Fatal("http.Client is nil")
	}
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
	if err != nil {
		t.Fatalf("Get() = %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("got %q, want %q", string(data), "hello world")
	}
}

func TestHTTPClientGet_Non200(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	_, err := cli.Get(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestHTTPClientGet_EmptyBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	_, err := cli.Get(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestNormalizeLocation_URL(t *testing.T) {
	t.Parallel()

	normalized, stype, err := normalizeLocation("https://example.com/spec.yaml")
	if err != nil {
		t.Fatalf("normalizeLocation() = %v", err)
	}
	if normalized != "https://example.com/spec.yaml" {
		t.Errorf("normalized = %q, want %q", normalized, "https://example.com/spec.yaml")
	}
	if stype != sourceURL {
		t.Errorf("stype = %q, want %q", stype, sourceURL)
	}
}

func TestNormalizeLocation_FileURL(t *testing.T) {
	t.Parallel()

	normalized, stype, err := normalizeLocation("file:///home/user/spec.yaml")
	if err != nil {
		t.Fatalf("normalizeLocation() = %v", err)
	}
	if stype != sourceLocal {
		t.Errorf("stype = %q, want %q", stype, sourceLocal)
	}
	if normalized != "/home/user/spec.yaml" {
		t.Errorf("normalized = %q, want %q", normalized, "/home/user/spec.yaml")
	}
}

func TestNormalizeLocation_LocalPath(t *testing.T) {
	t.Parallel()

	normalized, stype, err := normalizeLocation("/absolute/path/spec.yaml")
	if err != nil {
		t.Fatalf("normalizeLocation() = %v", err)
	}
	if stype != sourceLocal {
		t.Errorf("stype = %q, want %q", stype, sourceLocal)
	}
	if normalized != "/absolute/path/spec.yaml" {
		t.Errorf("normalized = %q, want %q", normalized, "/absolute/path/spec.yaml")
	}
}

func TestNormalizeLocation_InvalidFileURL(t *testing.T) {
	t.Parallel()

	_, _, err := normalizeLocation("ftp://example.com/spec.yaml")
	if err != nil {
		t.Fatalf("normalizeLocation() = %v, want nil (ftp falls through to local path)", err)
	}
}

func TestIsInsideSpecs_Inside(t *testing.T) {
	t.Parallel()

	c := New("/workspace")
	path := filepath.Join("/workspace", SpecsDirName, "subdir", "spec.yaml")
	if !c.isInsideSpecs(path) {
		t.Error("isInsideSpecs() = false, want true")
	}
}

func TestIsInsideSpecs_Outside(t *testing.T) {
	t.Parallel()

	c := New("/workspace")
	if c.isInsideSpecs("/other/path/spec.yaml") {
		t.Error("isInsideSpecs() = true, want false")
	}
}

func TestIsInsideSpecs_EmptySpecsDir(t *testing.T) {
	t.Parallel()

	c := &Cache{specsDir: ""}
	if c.isInsideSpecs("/any/path") {
		t.Error("isInsideSpecs() = true, want false")
	}
}

func TestResolveSpecsPath_Found(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	wsDir := filepath.Join(dir, ".swag2mcp")
	specsDir := filepath.Join(wsDir, SpecsDirName)
	if err := os.MkdirAll(specsDir, 0750); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	specFile := filepath.Join(specsDir, "api.yaml")
	if err := os.WriteFile(specFile, []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	c := New(wsDir)
	got, ok := c.resolveSpecsPath("specs/api.yaml")
	if !ok {
		t.Fatal("resolveSpecsPath() = false, want true")
	}
	if got != specFile {
		t.Errorf("got %q, want %q", got, specFile)
	}
}

func TestResolveSpecsPath_NotFound(t *testing.T) {
	t.Parallel()

	c := New(t.TempDir())
	_, ok := c.resolveSpecsPath("specs/nonexistent.yaml")
	if ok {
		t.Fatal("resolveSpecsPath() = true, want false")
	}
}

func TestResolveSpecsPath_NotSpecsPath(t *testing.T) {
	t.Parallel()

	c := New(t.TempDir())
	_, ok := c.resolveSpecsPath("/absolute/path.yaml")
	if ok {
		t.Fatal("resolveSpecsPath() = true, want false")
	}
}

func TestResolve_localFileOutsideSpecsReCacheOnExpiredMeta(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got1, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	// Manually expire the meta
	hash := cacheKey(specFile)
	metaPath := filepath.Join(c.dir, hash+".meta")
	if mErr := writeMeta(metaPath, fileMeta{
		Source:     specFile,
		SourceType: "local",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}); mErr != nil {
		t.Fatal(mErr)
	}

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

func TestResolve_localFileOutsideSpecsReCacheOnMissingSpecFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	c := New(dir)

	got1, err := c.Resolve(context.Background(), specFile)
	if err != nil {
		t.Fatal(err)
	}

	if rmErr := os.Remove(got1); rmErr != nil {
		t.Fatal(rmErr)
	}

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
	if err != nil {
		t.Fatal(err)
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

	// Exists should fall through to HEAD request
	err = c.Exists(context.Background(), srv.URL)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
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
	if err != nil {
		t.Fatal(err)
	}

	if rmErr := os.Remove(got); rmErr != nil {
		t.Fatal(rmErr)
	}

	if exErr := c.Exists(context.Background(), srv.URL); exErr != nil {
		t.Errorf("expected nil, got %v", exErr)
	}
}

func TestExists_fileURLSpecsRelative(t *testing.T) {
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
	if err := c.Exists(context.Background(), "file://"+specFile); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExists_tildePath(t *testing.T) {
	t.Parallel()
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		t.Fatalf("UserHomeDir() = %v", homeErr)
	}

	tmpFile := filepath.Join(home, ".swag2mcp-test-exists")
	if wrErr := os.WriteFile(tmpFile, []byte("data"), 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}
	defer os.Remove(tmpFile)

	c := New(t.TempDir())
	if exErr := c.Exists(context.Background(), "~/.swag2mcp-test-exists"); exErr != nil {
		t.Errorf("expected nil, got %v", exErr)
	}
}

func TestExists_relativePath(t *testing.T) {
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	if wrErr := os.WriteFile(specFile, []byte("data"), 0644); wrErr != nil {
		t.Fatal(wrErr)
	}

	t.Chdir(dir)

	c := New(dir)
	if exErr := c.Exists(context.Background(), "spec.yaml"); exErr != nil {
		t.Errorf("expected nil, got %v", exErr)
	}
}
