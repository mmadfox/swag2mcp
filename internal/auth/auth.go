package auth

import (
	"net/http"

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
