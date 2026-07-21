package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeScript(t *testing.T, dir, content string) {
	t.Helper()

	var (
		ext    string
		header string
	)
	if runtime.GOOS == "windows" {
		ext = ".bat"
		header = "@echo off\n"
	} else {
		ext = ".sh"
		header = "#!/bin/sh\n"
	}

	scriptPath := filepath.Join(dir, "auth_scripts", "testdomain"+ext)
	require.NoError(t, os.MkdirAll(filepath.Dir(scriptPath), 0700), "mkdir")
	require.NoError(t, os.WriteFile(scriptPath, []byte(header+content), 0700), "write script")
}

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
		require.NoError(t, client.New(), "New()")

		req, _ := newGetRequest()
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "Bearer script-token-456", req.Header.Get(headerAuthorization))
		assert.Equal(t, "Bearer script-token-456", info.Headers[headerAuthorization])
	})

	t.Run("caches token and reuses on second Apply", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"token": "cached-script-token", "expires_in": 3600}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		require.NoError(t, client.New(), "New()")

		req1, _ := newGetRequest()
		require.NoError(t, client.Apply(req1, nil), "Apply #1")

		req2, _ := newGetRequest()
		require.NoError(t, client.Apply(req2, nil), "Apply #2")

		assert.Equal(t, "Bearer cached-script-token", req2.Header.Get(headerAuthorization))
	})

	t.Run("returns error on invalid JSON output", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo 'not-json'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		require.NoError(t, client.New(), "New()")

		req, _ := newGetRequest()
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error for invalid JSON")
	})

	t.Run("returns error on missing token field", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"foo": "bar"}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		require.NoError(t, client.New(), "New()")

		req, _ := newGetRequest()
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error for missing token")
	})

	t.Run("uses default expires_in when not provided", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `echo '{"token": "no-expiry-token"}'`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		require.NoError(t, client.New(), "New()")

		req, _ := newGetRequest()
		require.NoError(t, client.Apply(req, nil), "Apply()")

		assert.Equal(t, "Bearer no-expiry-token", req.Header.Get(headerAuthorization))
	})

	t.Run("returns error on script execution failure", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeScript(t, dir, `exit 1`)

		client := &ScriptAuthClient{
			Domain:       "testdomain",
			workspaceDir: dir,
		}
		require.NoError(t, client.New(), "New()")

		req, _ := newGetRequest()
		err := client.Apply(req, nil)
		require.Error(t, err, "expected error for script failure")
	})
}

func newGetRequest() (*http.Request, error) {
	return http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
}
