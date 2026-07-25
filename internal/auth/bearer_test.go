package auth

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
