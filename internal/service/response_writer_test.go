package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
)

func TestResolveMaxResponseSize_nil(t *testing.T) {
	t.Parallel()

	require.Equal(t, defaultMaxResponseSize, resolveMaxResponseSize(nil))
}

func TestResolveMaxResponseSize_zero(t *testing.T) {
	t.Parallel()

	zero := 0
	require.Equal(t, defaultMaxResponseSize, resolveMaxResponseSize(&zero))
}

func TestResolveMaxResponseSize_overflow(t *testing.T) {
	t.Parallel()

	big := 20_000_000
	require.Equal(t, maxAllowedResponseSize, resolveMaxResponseSize(&big))
}

func TestResolveMaxResponseSize_valid(t *testing.T) {
	t.Parallel()

	val := 2048
	require.Equal(t, 2048, resolveMaxResponseSize(&val))
}

func TestFormatSize_bytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, "500 B", formatSize(500))
}

func TestFormatSize_kb(t *testing.T) {
	t.Parallel()

	require.Equal(t, "2.0 KB", formatSize(2048))
}

func TestFormatSize_mb(t *testing.T) {
	t.Parallel()

	require.Equal(t, "1.0 MB", formatSize(1048576))
}

func TestRandomSuffix_length(t *testing.T) {
	t.Parallel()

	suf := randomSuffix(6)
	require.Len(t, suf, 6)
}

func TestRandomSuffix_fallback(t *testing.T) {
	t.Parallel()

	suf := randomSuffix(0)
	require.Len(t, suf, 0)
}

func TestMergeHTTPClientConfigs_bothNil(t *testing.T) {
	t.Parallel()

	require.Nil(t, mergeHTTPClientConfigs(nil, nil))
}

func TestMergeHTTPClientConfigs_specOnly(t *testing.T) {
	t.Parallel()

	sp := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Api": "key"},
	}
	result := mergeHTTPClientConfigs(sp, nil)
	require.Equal(t, "key", result.Headers["X-Api"])
}

func TestMergeHTTPClientConfigs_collectionOverrides(t *testing.T) {
	t.Parallel()

	sp := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Api": "spec"},
	}
	coll := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Api": "coll"},
	}
	result := mergeHTTPClientConfigs(sp, coll)
	require.Equal(t, "coll", result.Headers["X-Api"])
}

func TestNewInvokeResponse_json(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	body := []byte(`{"key": "value"}`)
	result := newInvokeResponse(resp, body)
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.NotNil(t, result.Body)
}

func TestNewInvokeResponse_string(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Header:     http.Header{},
	}
	body := []byte("error occurred")
	result := newInvokeResponse(resp, body)
	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Equal(t, "error occurred", result.Body)
}

func TestNewInvokeResponse_empty(t *testing.T) {
	t.Parallel()

	resp := &http.Response{
		StatusCode: http.StatusNoContent,
		Header:     http.Header{},
	}
	result := newInvokeResponse(resp, nil)
	require.Equal(t, http.StatusNoContent, result.StatusCode)
	require.Equal(t, "", result.Body)
}

func TestOpenCommand(t *testing.T) {
	t.Parallel()

	cmd := openCommand("/tmp/test.json")
	require.Contains(t, cmd, "/tmp/test.json")
}
