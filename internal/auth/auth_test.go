package auth

import (
	"testing"
)

// Compile-time checks: all auth clients satisfy the Authenticator interface.
var (
	_ Authenticator = (*NoAuthClient)(nil)
	_ Authenticator = (*BasicAuthClient)(nil)
	_ Authenticator = (*BearerTokenAuthClient)(nil)
	_ Authenticator = (*DigestAuthClient)(nil)
	_ Authenticator = (*OAuth2ClientCredentialsAuthClient)(nil)
	_ Authenticator = (*OAuth2PasswordAuthClient)(nil)
	_ Authenticator = (*APIKeyAuthClient)(nil)
	_ Authenticator = (*ScriptAuthClient)(nil)
)

func TestType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  Type
		want string
	}{
		{"none", NoAuth, "none"},
		{"basic", BasicAuth, "basic"},
		{"bearer", BearerTokenAuth, "bearer"},
		{"digest", DigestAuth, "digest"},
		{"oauth2-cc", OAuth2ClientCredentials, "oauth2-cc"},
		{"oauth2-pwd", OAuth2Password, "oauth2-pwd"},
		{"api-key", APIKeyAuth, "api-key"},
		{"script", ScriptAuth, "script"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.typ.String(); got != tt.want {
				t.Errorf("Type.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

//nolint:tparallel // testdata is not parallel-safe
func TestResolveEnv(t *testing.T) {
	t.Run("returns value when env var is set", func(t *testing.T) {
		t.Setenv("TEST_SWAG_MYVAR", "hello")
		got := resolveEnv("$(TEST_SWAG_MYVAR)")
		if got != "hello" {
			t.Errorf("resolveEnv = %q, want %q", got, "hello")
		}
	})

	t.Run("trims leading whitespace inside parens", func(t *testing.T) {
		t.Setenv("TEST_SWAG_MYVAR", "hello")
		got := resolveEnv("$(  TEST_SWAG_MYVAR)")
		if got != "hello" {
			t.Errorf("resolveEnv = %q, want %q", got, "hello")
		}
	})

	t.Run("trims trailing whitespace inside parens", func(t *testing.T) {
		t.Setenv("TEST_SWAG_MYVAR", "hello")
		got := resolveEnv("$(TEST_SWAG_MYVAR  )")
		if got != "hello" {
			t.Errorf("resolveEnv = %q, want %q", got, "hello")
		}
	})

	t.Run("trims whitespace on both sides inside parens", func(t *testing.T) {
		t.Setenv("TEST_SWAG_MYVAR", "hello")
		got := resolveEnv("$(  TEST_SWAG_MYVAR  )")
		if got != "hello" {
			t.Errorf("resolveEnv = %q, want %q", got, "hello")
		}
	})

	t.Run("returns empty string when env var is not set", func(t *testing.T) {
		t.Parallel()
		got := resolveEnv("$(TEST_SWAG_UNSET_VAR)")
		if got != "" {
			t.Errorf("resolveEnv = %q, want %q", got, "")
		}
	})

	t.Run("returns original string when no pattern matches", func(t *testing.T) {
		t.Parallel()
		got := resolveEnv("plaintext")
		if got != "plaintext" {
			t.Errorf("resolveEnv = %q, want %q", got, "plaintext")
		}
	})

	t.Run("returns original string when parens are empty", func(t *testing.T) {
		t.Parallel()
		got := resolveEnv("$(  )")
		if got != "$(  )" {
			t.Errorf("resolveEnv = %q, want %q", got, "$(  )")
		}
	})

	t.Run("returns original string when only open paren", func(t *testing.T) {
		t.Parallel()
		got := resolveEnv("$(")
		if got != "$(" {
			t.Errorf("resolveEnv = %q, want %q", got, "$(")
		}
	})

	t.Run("returns empty string for empty input", func(t *testing.T) {
		t.Parallel()
		got := resolveEnv("")
		if got != "" {
			t.Errorf("resolveEnv = %q, want %q", got, "")
		}
	})

	t.Run("trims whitespace around input", func(t *testing.T) {
		t.Setenv("TEST_SWAG_MYVAR", "hello")
		got := resolveEnv("  $(TEST_SWAG_MYVAR)  ")
		if got != "hello" {
			t.Errorf("resolveEnv = %q, want %q", got, "hello")
		}
	})
}
