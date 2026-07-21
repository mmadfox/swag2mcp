package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIKeyAuthClient_Apply(t *testing.T) {
	t.Parallel()

	t.Run("sets header when In is header", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "X-Api-Key",
			Value: "my-api-key-value",
			In:    "header",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "my-api-key-value", req.Header.Get("X-Api-Key"))
		assert.Equal(t, "my-api-key-value", info.Headers["X-Api-Key"])
	})

	t.Run("sets query param when In is query", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "api_key",
			Value: "query-key-value",
			In:    "query",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "query-key-value", req.URL.Query().Get("api_key"))
		assert.Equal(t, "query-key-value", info.QueryParams["api_key"])
	})

	t.Run("defaults to header when In is unknown", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "X-Auth",
			Value: "fallback-value",
			In:    "unknown",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		require.NoError(t, client.Apply(req, nil), "Apply()")

		assert.Equal(t, "fallback-value", req.Header.Get("X-Auth"))
	})

	t.Run("does not set empty value", func(t *testing.T) {
		t.Parallel()

		client := &APIKeyAuthClient{
			Key:   "X-Key",
			Value: "",
			In:    "header",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Empty(t, req.Header.Get("X-Key"))
		assert.Nil(t, info.Headers)
	})
}
