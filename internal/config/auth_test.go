package config

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

func TestAuth_UnmarshalYAML_none(t *testing.T) {
	t.Parallel()

	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte("type: none"), &a))
	require.NotNil(t, a.Client)
	_, ok := a.Client.(*auth.NoAuthClient)
	assert.True(t, ok, "expected *auth.NoAuthClient")
}

func TestAuth_UnmarshalYAML_none_empty(t *testing.T) {
	t.Parallel()

	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte("type:"), &a))
	require.NotNil(t, a.Client)
	_, ok := a.Client.(*auth.NoAuthClient)
	assert.True(t, ok, "expected *auth.NoAuthClient")
}

func TestAuth_UnmarshalYAML_basic(t *testing.T) {
	t.Parallel()

	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte("type: basic\nconfig:\n  username: myuser\n  password: secret"), &a))
	c, ok := a.Client.(*auth.BasicAuthClient)
	require.True(t, ok, "expected *auth.BasicAuthClient")
	assert.Equal(t, "myuser", c.Username)
	assert.Equal(t, "secret", c.Password)
}

func TestAuth_UnmarshalYAML_bearer(t *testing.T) {
	t.Parallel()

	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte("type: bearer\nconfig:\n  token: mytoken"), &a))
	c, ok := a.Client.(*auth.BearerTokenAuthClient)
	require.True(t, ok, "expected *auth.BearerTokenAuthClient")
	assert.Equal(t, "mytoken", c.Token)
}

func TestAuth_UnmarshalYAML_digest(t *testing.T) {
	t.Parallel()

	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte("type: digest\nconfig:\n  username: u\n  password: p"), &a))
	c, ok := a.Client.(*auth.DigestAuthClient)
	require.True(t, ok, "expected *auth.DigestAuthClient")
	assert.Equal(t, "u", c.Username)
	assert.Equal(t, "p", c.Password)
}

func TestAuth_UnmarshalYAML_oauth2_cc(t *testing.T) {
	t.Parallel()

	y := "type: oauth2-cc\nconfig:\n  client_id: cid\n  client_secret: cs\n  token_url: https://example.com/token\n  scopes:\n    - read\n    - write"
	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte(y), &a))
	c, ok := a.Client.(*auth.OAuth2ClientCredentialsAuthClient)
	require.True(t, ok, "expected *auth.OAuth2ClientCredentialsAuthClient")
	assert.Equal(t, "cid", c.ClientID)
	assert.Equal(t, "cs", c.ClientSecret)
	assert.Equal(t, "https://example.com/token", c.TokenURL)
	assert.Equal(t, []string{"read", "write"}, c.Scopes)
}

func TestAuth_UnmarshalYAML_oauth2_pwd(t *testing.T) {
	t.Parallel()

	y := "type: oauth2-pwd\nconfig:\n  username: u\n  password: p\n  client_id: cid\n  client_secret: cs\n  token_url: https://example.com/token"
	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte(y), &a))
	c, ok := a.Client.(*auth.OAuth2PasswordAuthClient)
	require.True(t, ok, "expected *auth.OAuth2PasswordAuthClient")
	assert.Equal(t, "u", c.Username)
	assert.Equal(t, "p", c.Password)
	assert.Equal(t, "cid", c.ClientID)
	assert.Equal(t, "cs", c.ClientSecret)
	assert.Equal(t, "https://example.com/token", c.TokenURL)
	assert.Nil(t, c.Scopes)
}

func TestAuth_UnmarshalYAML_api_key(t *testing.T) {
	t.Parallel()

	var a Auth
	y := "type: api-key\nconfig:\n  key: X-API-Key\n  value: abc123\n  in: header"
	require.NoError(t, yaml.Unmarshal([]byte(y), &a))
	c, ok := a.Client.(*auth.APIKeyAuthClient)
	require.True(t, ok, "expected *auth.APIKeyAuthClient")
	assert.Equal(t, "X-API-Key", c.Key)
	assert.Equal(t, "abc123", c.Value)
	assert.Equal(t, "header", c.In)
}

func TestAuth_UnmarshalYAML_script(t *testing.T) {
	t.Parallel()

	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte("type: script\nconfig:\n  domain: my-api"), &a))
	c, ok := a.Client.(*auth.ScriptAuthClient)
	require.True(t, ok, "expected *auth.ScriptAuthClient")
	assert.Equal(t, "my-api", c.Domain)
}

func TestAuth_UnmarshalYAML_hmac(t *testing.T) {
	t.Parallel()

	var a Auth
	y := "type: hmac\nconfig:\n  api_key: my-api-key\n  secret_key: my-secret-key"
	require.NoError(t, yaml.Unmarshal([]byte(y), &a))
	c, ok := a.Client.(*auth.HMACAuthClient)
	require.True(t, ok, "expected *auth.HMACAuthClient")
	assert.Equal(t, "my-api-key", c.APIKey)
	assert.Equal(t, "my-secret-key", c.SecretKey)
}

func TestAuth_UnmarshalYAML_underscore_type(t *testing.T) {
	t.Parallel()

	y := "type: oauth2_cc\nconfig:\n  client_id: cid\n  client_secret: cs\n  token_url: https://example.com/token"
	var a Auth
	require.NoError(t, yaml.Unmarshal([]byte(y), &a))
	_, ok := a.Client.(*auth.OAuth2ClientCredentialsAuthClient)
	assert.True(t, ok, "expected *auth.OAuth2ClientCredentialsAuthClient")
}

func TestAuth_UnmarshalYAML_empty_config_error(t *testing.T) {
	t.Parallel()

	var a Auth
	err := yaml.Unmarshal([]byte("type: basic"), &a)
	require.Error(t, err)
	assert.Equal(t, "auth config is required for this type", err.Error())
}

func TestAuth_UnmarshalYAML_unknown_type(t *testing.T) {
	t.Parallel()

	var a Auth
	err := yaml.Unmarshal([]byte("type: unknown\nconfig: {}"), &a)
	require.Error(t, err)
	assert.Equal(t, `unsupported auth type "unknown"`, err.Error())
}

func TestAuth_UnmarshalYAML_decode_error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		yaml string
	}{
		{"basic", "type: basic\nconfig: 123"},
		{"bearer", "type: bearer\nconfig: 123"},
		{"digest", "type: digest\nconfig: 123"},
		{"oauth2-cc", "type: oauth2-cc\nconfig: 123"},
		{"oauth2-pwd", "type: oauth2-pwd\nconfig: 123"},
		{"api-key", "type: api-key\nconfig: 123"},
		{"script", "type: script\nconfig: 123"},
		{"hmac", "type: hmac\nconfig: 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var a Auth
			err := yaml.Unmarshal([]byte(tt.yaml), &a)
			require.Error(t, err, "expected decode error")
		})
	}
}

func TestAuth_UnmarshalYAML_top_level_decode_error(t *testing.T) {
	t.Parallel()

	var a Auth
	err := yaml.Unmarshal([]byte("type:\n  - broken"), &a)
	require.Error(t, err, "expected error for malformed YAML")
}
