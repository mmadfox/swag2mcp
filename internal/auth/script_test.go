package auth

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func writeScript(t *testing.T, dir, content string) {
	t.Helper()

	var ext string
	var header string
	if runtime.GOOS == "windows" {
		ext = ".bat"
		header = "@echo off\n"
	} else {
		ext = ".sh"
		header = "#!/bin/sh\n"
	}

	scriptPath := filepath.Join(dir, "auth_scripts", "testdomain"+ext)
	if err := os.MkdirAll(filepath.Dir(scriptPath), 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(scriptPath, []byte(header+content), 0700); err != nil {
		t.Fatalf("write script: %v", err)
	}
}

//nolint:gocognit
func TestScriptAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("successful script execution", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"token": "script-token-456", "expires_in": 3600}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := newGetRequest()
		var info Info
		if err := client.Apply(req, &info); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.Header.Get("Authorization"); v != "Bearer script-token-456" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer script-token-456")
		}
		if v := info.Headers["Authorization"]; v != "Bearer script-token-456" {
			t.Errorf("info.Headers[Authorization] = %q, want %q", v, "Bearer script-token-456")
		}
	})

	t.Run("caches token and reuses on second Apply", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"token": "cached-script-token", "expires_in": 3600}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req1, _ := newGetRequest()
		if err := client.Apply(req1, nil); err != nil {
			t.Fatalf("Apply #1 = %v", err)
		}

		req2, _ := newGetRequest()
		if err := client.Apply(req2, nil); err != nil {
			t.Fatalf("Apply #2 = %v", err)
		}

		if v := req2.Header.Get("Authorization"); v != "Bearer cached-script-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer cached-script-token")
		}
	})

	t.Run("returns error on invalid JSON output", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo 'not-json'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := newGetRequest()
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})

	t.Run("returns error on missing token field", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"foo": "bar"}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := newGetRequest()
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error for missing token, got nil")
		}
	})

	t.Run("uses default expires_in when not provided", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"token": "no-expiry-token"}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := newGetRequest()
		if err := client.Apply(req, nil); err != nil {
			t.Fatalf("Apply() = %v", err)
		}

		if v := req.Header.Get("Authorization"); v != "Bearer no-expiry-token" {
			t.Errorf("Authorization = %q, want %q", v, "Bearer no-expiry-token")
		}
	})

	t.Run("returns error on script execution failure", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `exit 1`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		if err := client.New(); err != nil {
			t.Fatalf("New() = %v", err)
		}

		req, _ := newGetRequest()
		if err := client.Apply(req, nil); err == nil {
			t.Fatal("expected error for script failure, got nil")
		}
	})
}

func newGetRequest() (*http.Request, error) {
	return http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
}
