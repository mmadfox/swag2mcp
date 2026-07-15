package config

import (
	"testing"

	"go.yaml.in/yaml/v3"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

func TestAuth_UnmarshalYAML_none(t *testing.T) {
	t.Parallel()

	var a Auth
	if err := yaml.Unmarshal([]byte("type: none"), &a); err != nil {
		t.Fatal(err)
	}
	if a.Client == nil {
		t.Fatal("Client is nil")
	}
	if _, ok := a.Client.(auth.NoAuthClient); !ok {
		t.Fatalf("expected auth.NoAuthClient, got %T", a.Client)
	}
}

func TestAuth_UnmarshalYAML_none_empty(t *testing.T) {
	t.Parallel()

	var a Auth
	if err := yaml.Unmarshal([]byte("type:"), &a); err != nil {
		t.Fatal(err)
	}
	if a.Client == nil {
		t.Fatal("Client is nil")
	}
	if _, ok := a.Client.(auth.NoAuthClient); !ok {
		t.Fatalf("expected auth.NoAuthClient, got %T", a.Client)
	}
}

func TestAuth_UnmarshalYAML_basic(t *testing.T) {
	t.Parallel()

	var a Auth
	if err := yaml.Unmarshal([]byte("type: basic\nconfig:\n  username: myuser\n  password: secret"), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.BasicAuthClient)
	if !ok {
		t.Fatalf("expected *auth.BasicAuthClient, got %T", a.Client)
	}
	if c.Username != "myuser" {
		t.Errorf("Username = %q, want %q", c.Username, "myuser")
	}
	if c.Password != "secret" {
		t.Errorf("Password = %q, want %q", c.Password, "secret")
	}
}

func TestAuth_UnmarshalYAML_bearer(t *testing.T) {
	t.Parallel()

	var a Auth
	if err := yaml.Unmarshal([]byte("type: bearer\nconfig:\n  token: mytoken"), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.BearerTokenAuthClient)
	if !ok {
		t.Fatalf("expected *auth.BearerTokenAuthClient, got %T", a.Client)
	}
	if c.Token != "mytoken" {
		t.Errorf("Token = %q, want %q", c.Token, "mytoken")
	}
}

func TestAuth_UnmarshalYAML_digest(t *testing.T) {
	t.Parallel()

	var a Auth
	if err := yaml.Unmarshal([]byte("type: digest\nconfig:\n  username: u\n  password: p"), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.DigestAuthClient)
	if !ok {
		t.Fatalf("expected *auth.DigestAuthClient, got %T", a.Client)
	}
	if c.Username != "u" {
		t.Errorf("Username = %q, want %q", c.Username, "u")
	}
	if c.Password != "p" {
		t.Errorf("Password = %q, want %q", c.Password, "p")
	}
}

func TestAuth_UnmarshalYAML_oauth2_cc(t *testing.T) {
	t.Parallel()

	y := "type: oauth2-cc\nconfig:\n  client_id: cid\n  client_secret: cs\n  token_url: https://example.com/token\n  scopes:\n    - read\n    - write"
	var a Auth
	if err := yaml.Unmarshal([]byte(y), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.OAuth2ClientCredentialsAuthClient)
	if !ok {
		t.Fatalf("expected *auth.OAuth2ClientCredentialsAuthClient, got %T", a.Client)
	}
	if c.ClientID != "cid" {
		t.Errorf("ClientID = %q, want %q", c.ClientID, "cid")
	}
	if c.ClientSecret != "cs" {
		t.Errorf("ClientSecret = %q, want %q", c.ClientSecret, "cs")
	}
	if c.TokenURL != "https://example.com/token" {
		t.Errorf("TokenURL = %q, want %q", c.TokenURL, "https://example.com/token")
	}
	if len(c.Scopes) != 2 || c.Scopes[0] != "read" || c.Scopes[1] != "write" {
		t.Errorf("Scopes = %v, want [read write]", c.Scopes)
	}
}

func TestAuth_UnmarshalYAML_oauth2_pwd(t *testing.T) {
	t.Parallel()

	y := "type: oauth2-pwd\nconfig:\n  username: u\n  password: p\n  client_id: cid\n  client_secret: cs\n  token_url: https://example.com/token"
	var a Auth
	if err := yaml.Unmarshal([]byte(y), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.OAuth2PasswordAuthClient)
	if !ok {
		t.Fatalf("expected *auth.OAuth2PasswordAuthClient, got %T", a.Client)
	}
	if c.Username != "u" {
		t.Errorf("Username = %q, want %q", c.Username, "u")
	}
	if c.Password != "p" {
		t.Errorf("Password = %q, want %q", c.Password, "p")
	}
	if c.ClientID != "cid" {
		t.Errorf("ClientID = %q, want %q", c.ClientID, "cid")
	}
	if c.ClientSecret != "cs" {
		t.Errorf("ClientSecret = %q, want %q", c.ClientSecret, "cs")
	}
	if c.TokenURL != "https://example.com/token" {
		t.Errorf("TokenURL = %q, want %q", c.TokenURL, "https://example.com/token")
	}
	if c.Scopes != nil {
		t.Errorf("Scopes = %v, want nil", c.Scopes)
	}
}

func TestAuth_UnmarshalYAML_api_key(t *testing.T) {
	t.Parallel()

	var a Auth
	y := "type: api-key\nconfig:\n  key: X-API-Key\n  value: abc123\n  in: header"
	if err := yaml.Unmarshal([]byte(y), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.APIKeyAuthClient)
	if !ok {
		t.Fatalf("expected *auth.APIKeyAuthClient, got %T", a.Client)
	}
	if c.Key != "X-API-Key" {
		t.Errorf("Key = %q, want %q", c.Key, "X-API-Key")
	}
	if c.Value != "abc123" {
		t.Errorf("Value = %q, want %q", c.Value, "abc123")
	}
	if c.In != "header" {
		t.Errorf("In = %q, want %q", c.In, "header")
	}
}

func TestAuth_UnmarshalYAML_script(t *testing.T) {
	t.Parallel()

	var a Auth
	if err := yaml.Unmarshal([]byte("type: script\nconfig:\n  domain: my-api"), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.ScriptAuthClient)
	if !ok {
		t.Fatalf("expected *auth.ScriptAuthClient, got %T", a.Client)
	}
	if c.Domain != "my-api" {
		t.Errorf("Domain = %q, want %q", c.Domain, "my-api")
	}
}

func TestAuth_UnmarshalYAML_hmac(t *testing.T) {
	t.Parallel()

	var a Auth
	y := "type: hmac\nconfig:\n  api_key: my-api-key\n  secret_key: my-secret-key"
	if err := yaml.Unmarshal([]byte(y), &a); err != nil {
		t.Fatal(err)
	}
	c, ok := a.Client.(*auth.HMACAuthClient)
	if !ok {
		t.Fatalf("expected *auth.HMACAuthClient, got %T", a.Client)
	}
	if c.APIKey != "my-api-key" {
		t.Errorf("APIKey = %q, want %q", c.APIKey, "my-api-key")
	}
	if c.SecretKey != "my-secret-key" {
		t.Errorf("SecretKey = %q, want %q", c.SecretKey, "my-secret-key")
	}
}

func TestAuth_UnmarshalYAML_underscore_type(t *testing.T) {
	t.Parallel()

	y := "type: oauth2_cc\nconfig:\n  client_id: cid\n  client_secret: cs\n  token_url: https://example.com/token"
	var a Auth
	if err := yaml.Unmarshal([]byte(y), &a); err != nil {
		t.Fatal(err)
	}
	if _, ok := a.Client.(*auth.OAuth2ClientCredentialsAuthClient); !ok {
		t.Fatalf("expected *auth.OAuth2ClientCredentialsAuthClient, got %T", a.Client)
	}
}

func TestAuth_UnmarshalYAML_empty_config_error(t *testing.T) {
	t.Parallel()

	var a Auth
	err := yaml.Unmarshal([]byte("type: basic"), &a)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "auth config is required for this type" {
		t.Errorf("error = %q, want %q", err.Error(), "auth config is required for this type")
	}
}

func TestAuth_UnmarshalYAML_unknown_type(t *testing.T) {
	t.Parallel()

	var a Auth
	err := yaml.Unmarshal([]byte("type: unknown\nconfig: {}"), &a)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != `unsupported auth type "unknown"` {
		t.Errorf("error = %q, want %q", err.Error(), `unsupported auth type "unknown"`)
	}
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
			if err == nil {
				t.Fatal("expected decode error, got nil")
			}
		})
	}
}

func TestAuth_UnmarshalYAML_top_level_decode_error(t *testing.T) {
	t.Parallel()

	var a Auth
	err := yaml.Unmarshal([]byte("type:\n  - broken"), &a)
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}
