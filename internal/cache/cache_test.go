package cache

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	c := New("/tmp/swag2mcp")
	expected := filepath.Join("/tmp/swag2mcp", CacheDirName)
	assert.Equal(t, expected, c.dir)
	expectedSpecs := filepath.Join("/tmp/swag2mcp", SpecsDirName)
	assert.Equal(t, expectedSpecs, c.specsDir)
}

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

func TestSetHTTPClient(t *testing.T) {
	t.Parallel()

	c := New(t.TempDir())
	cli := &http.Client{Timeout: 5}
	c.SetHTTPClient(cli)
	assert.Same(t, cli, c.cli.client)
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

func TestRandomTTL(t *testing.T) {
	t.Parallel()

	for range 100 {
		ttl := randomTTL()
		assert.GreaterOrEqual(t, ttl, MinTTL, "TTL < MinTTL")
		assert.LessOrEqual(t, ttl, MaxTTL, "TTL > MaxTTL")
	}
}

func TestResolve_localPathOutsideSpecs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

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

	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

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

	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	hash := cacheKey(srv.URL)
	metaPath := filepath.Join(c.dir, hash+".meta")
	require.NoError(t, writeMeta(metaPath, fileMeta{
		Source:     srv.URL,
		SourceType: "url",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}))

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

	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

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

	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(specFile, []byte("world"), 0644))

	got2, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

	assert.Equal(t, got1, got2, "expected same cache path")

	data, err := os.ReadFile(got2)
	require.NoError(t, err)
	assert.Equal(t, "world", string(data))
}

func TestResolve_localFileOutsideSpecsReCacheOnExpiredMeta(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("hello"), 0644))

	c := New(dir)

	got1, err := c.Resolve(context.Background(), specFile)
	require.NoError(t, err)

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

func TestNormalizeLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		location    string
		wantStype   sourceType
		wantPath    string
		wantErr     bool
		errContains string
	}{
		{
			name:      "URL",
			location:  "https://example.com/spec.yaml",
			wantStype: sourceURL,
			wantPath:  "https://example.com/spec.yaml",
		},
		{
			name:      "file URL",
			location:  "file:///home/user/spec.yaml",
			wantStype: sourceLocal,
			wantPath:  "/home/user/spec.yaml",
		},
		{
			name:      "local path",
			location:  "/absolute/path/spec.yaml",
			wantStype: sourceLocal,
			wantPath:  "/absolute/path/spec.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			normalized, stype, err := normalizeLocation(tt.location)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantStype, stype)
			assert.Equal(t, tt.wantPath, normalized)
		})
	}
}

func TestClassifyLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		location    string
		wantStype   sourceType
		wantPath    string
		wantErr     bool
		errContains string
	}{
		{
			name:      "http URL",
			location:  "http://example.com/spec.yaml",
			wantStype: sourceURL,
			wantPath:  "http://example.com/spec.yaml",
		},
		{
			name:      "https URL",
			location:  "https://example.com/spec.yaml",
			wantStype: sourceURL,
			wantPath:  "https://example.com/spec.yaml",
		},
		{
			name:      "file URL",
			location:  "file:///home/user/spec.yaml",
			wantStype: sourceLocal,
			wantPath:  "/home/user/spec.yaml",
		},
		{
			name:        "file URL bad scheme",
			location:    "file://invalid|path",
			wantErr:     true,
			errContains: "file",
		},
		{
			name:      "local path",
			location:  "/absolute/path/spec.yaml",
			wantStype: sourceLocal,
			wantPath:  "/absolute/path/spec.yaml",
		},
		{
			name:      "tilde path",
			location:  "~/spec.yaml",
			wantStype: sourceLocal,
			wantPath:  "spec.yaml",
		},
		{
			name:      "ftp as local",
			location:  "ftp://example.com/spec.yaml",
			wantStype: sourceLocal,
			wantPath:  "ftp://example.com/spec.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stype, path, err := classifyLocation(tt.location)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantStype, stype)
			assert.Contains(t, path, tt.wantPath)
		})
	}
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

func TestLoadSource_unknownType(t *testing.T) {
	t.Parallel()
	c := New(t.TempDir())
	_, _, err := c.loadSource(context.Background(), "/tmp/test.yaml", sourceType("unknown"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown source type")
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

	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	require.NoError(t, c.Exists(context.Background(), srv.URL))

	assert.Equal(t, 1, callCount, "expected 1 server call")
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

	_, err := c.Resolve(context.Background(), srv.URL)
	require.NoError(t, err)

	hash := cacheKey(srv.URL)
	metaPath := filepath.Join(c.dir, hash+".meta")
	require.NoError(t, writeMeta(metaPath, fileMeta{
		Source:     srv.URL,
		SourceType: "url",
		CachedAt:   time.Now().Add(-2 * MaxTTL),
		TTLSec:     60,
	}))

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

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
