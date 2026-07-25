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

func TestBearerTokenAuthClient_Apply(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: "my-bearer-token"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Equal(t, "Bearer my-bearer-token", req.Header.Get(headerAuthorization))
	assert.Equal(t, "Bearer my-bearer-token", info.Headers[headerAuthorization])
}

func TestBearerTokenAuthClient_Apply_EmptyToken(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{Token: ""}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Empty(t, req.Header.Get(headerAuthorization))
	assert.Nil(t, info.Headers)
}

func TestBearerTokenAuthClient_Apply_EnvVars(t *testing.T) {
	t.Setenv("TEST_BEARER_TOKEN", "env-bearer-token")

	client := &BearerTokenAuthClient{Token: "$(TEST_BEARER_TOKEN)"}
	require.NoError(t, client.New(), "New()")

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	require.NoError(t, client.Apply(req, nil), "Apply()")

	assert.Equal(t, "Bearer env-bearer-token", req.Header.Get(headerAuthorization))
}

func TestBearerTokenAuthClient_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		client  *BearerTokenAuthClient
		wantErr bool
	}{
		{name: "valid", client: &BearerTokenAuthClient{Token: "valid-token"}, wantErr: false},
		{name: "empty token", client: &BearerTokenAuthClient{Token: ""}, wantErr: true},
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

func TestBearerTokenAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := &BearerTokenAuthClient{}
	assert.Equal(t, BearerTokenAuth, client.Type())
}
