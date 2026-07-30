/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package auth

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

func TestAPIKeyAuthClient_Apply_EnvVars(t *testing.T) {
	t.Run("resolves value from env in header mode", func(t *testing.T) {
		t.Setenv("TEST_API_KEY_HDR", "env-key-value")

		client := &APIKeyAuthClient{
			Key:   "X-Api-Key",
			Value: "$(TEST_API_KEY_HDR)",
			In:    "header",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "env-key-value", req.Header.Get("X-Api-Key"))
		assert.Equal(t, "env-key-value", info.Headers["X-Api-Key"])
	})

	t.Run("resolves value from env in query mode", func(t *testing.T) {
		t.Setenv("TEST_API_KEY_QRY", "env-query-key")

		client := &APIKeyAuthClient{
			Key:   "api_key",
			Value: "$(TEST_API_KEY_QRY)",
			In:    "query",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "env-query-key", req.URL.Query().Get("api_key"))
		assert.Equal(t, "env-query-key", info.QueryParams["api_key"])
	})

	t.Run("resolves key name from env", func(t *testing.T) {
		t.Setenv("TEST_API_KEY_NAME", "X-Env-Key")
		t.Setenv("TEST_API_KEY_VAL", "env-val")

		client := &APIKeyAuthClient{
			Key:   "$(TEST_API_KEY_NAME)",
			Value: "$(TEST_API_KEY_VAL)",
			In:    "header",
		}
		require.NoError(t, client.New(), "New()")

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
		var info Info
		require.NoError(t, client.Apply(req, &info), "Apply()")

		assert.Equal(t, "env-val", req.Header.Get("X-Env-Key"))
		assert.Equal(t, "env-val", info.Headers["X-Env-Key"])
	})
}

func TestAPIKeyAuthClient_New(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{Key: "X-Key", Value: "val", In: "header"}
	require.NoError(t, client.New())
}

func TestAPIKeyAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *APIKeyAuthClient
		wantErr bool
	}{
		{name: "valid", client: &APIKeyAuthClient{Key: "X-Key", Value: "val", In: "header"}, wantErr: false},
		{name: "empty", client: &APIKeyAuthClient{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.client.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAPIKeyAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &APIKeyAuthClient{}
	assert.Equal(t, APIKeyAuth, client.Type())
}
