package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Type is the type of authentication used.
type Type string

// String returns the string representation of the auth type.
func (t Type) String() string {
	return string(t)
}

const (
	NoAuth                  Type = "none"
	BasicAuth               Type = "basic"
	BearerTokenAuth         Type = "bearer"
	DigestAuth              Type = "digest"
	OAuth2ClientCredentials Type = "oauth2-cc"  //nolint:gosec // not a credential, type name
	OAuth2Password          Type = "oauth2-pwd" //nolint:gosec // not a credential, type name
	APIKeyAuth              Type = "api-key"
	ScriptAuth              Type = "script"
)

// authValidator validates auth client structs using struct tags.
//
//nolint:gochecknoglobals // validator is stateless and safe to reuse.
var authValidator = validator.New(validator.WithRequiredStructEnabled())

// Authenticator is an interface for authenticating requests.
type Authenticator interface {
	New() error
	Type() Type
	Apply(req *http.Request) error
	Validate() error
}

// resolveEnv checks if s matches the pattern $(VARNAME) with optional
// whitespace inside the parentheses. If it matches, the variable name is
// extracted and looked up via [os.Getenv]. Otherwise s is returned unchanged.
func resolveEnv(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "$(") || !strings.HasSuffix(s, ")") {
		return s
	}
	inner := s[2 : len(s)-1]
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return s
	}
	return os.Getenv(inner)
}
