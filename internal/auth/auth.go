package auth

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/env"
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

// Info holds the authentication details extracted during Apply.
type Info struct {
	Headers     map[string]string
	QueryParams map[string]string
}

// Authenticator is an interface for authenticating requests.
type Authenticator interface {
	New() error
	Type() Type
	Apply(req *http.Request, out *Info) error
	Validate() error
}

// TokenURLSetter is an optional interface for auth clients that have a configurable token URL.
type TokenURLSetter interface {
	SetTokenURL(url string)
}

// MockBaseURLSetter is an optional interface for auth clients that need a mock base URL.
type MockBaseURLSetter interface {
	SetMockBaseURL(url string)
}

func setAuthHeader(req *http.Request, out *Info, key, value string) {
	if value == "" {
		return
	}
	req.Header.Set(key, value)
	if out != nil {
		if out.Headers == nil {
			out.Headers = make(map[string]string)
		}
		out.Headers[key] = value
	}
}

func setAuthQuery(req *http.Request, out *Info, key, value string) {
	if value == "" {
		return
	}
	q := req.URL.Query()
	q.Set(key, value)
	req.URL.RawQuery = q.Encode()
	if out != nil {
		if out.QueryParams == nil {
			out.QueryParams = make(map[string]string)
		}
		out.QueryParams[key] = value
	}
}

// resolveEnv resolves $(VAR_NAME) patterns to environment variable values.
func resolveEnv(s string) string {
	return env.Parse(s)
}
