package auth

import (
	"testing"
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
