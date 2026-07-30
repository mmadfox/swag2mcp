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

func TestBasicAuthClient_Apply(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{
		Username: "alice",
		Password: "secret123",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	user, pass, ok := req.BasicAuth()
	require.True(t, ok, "expected BasicAuth to be set")
	assert.Equal(t, "alice", user)
	assert.Equal(t, "secret123", pass)
	assert.NotEmpty(t, info.Headers[headerAuthorization])
}

func TestBasicAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_BASIC_USER", "bob")
	t.Setenv("TEST_BASIC_PASS", "bobpass")

	client := &BasicAuthClient{
		Username: "$(TEST_BASIC_USER)",
		Password: "$(TEST_BASIC_PASS)",
	}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	user, pass, ok := req.BasicAuth()
	require.True(t, ok, "expected BasicAuth to be set")
	assert.Equal(t, "bob", user)
	assert.Equal(t, "bobpass", pass)
}

func TestBasicAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &BasicAuthClient{}
	assert.Equal(t, BasicAuth, client.Type())
}

func TestBasicAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *BasicAuthClient
		wantErr bool
	}{
		{name: "valid", client: &BasicAuthClient{Username: "u", Password: "p"}, wantErr: false},
		{name: "empty", client: &BasicAuthClient{}, wantErr: true},
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
