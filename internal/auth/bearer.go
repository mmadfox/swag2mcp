package auth

import "net/http"

// BearerTokenAuthClient holds a bearer token for HTTP Bearer authentication.
type BearerTokenAuthClient struct {
	Token string `yaml:"token" validate:"required"`
}

// New resolves environment variables in the Token and returns nil.
func (c *BearerTokenAuthClient) New() error {
	c.Token = resolveEnv(c.Token)
	return nil
}

// Type returns the authentication type for Bearer token auth.
func (c *BearerTokenAuthClient) Type() Type {
	return BearerTokenAuth
}

// Apply sets the Bearer Authorization header on the request using the configured token.
func (c *BearerTokenAuthClient) Apply(req *http.Request, out *Info) error {
	if c.Token == "" {
		return nil
	}
	setAuthHeader(req, out, headerAuthorization, bearerToken(c.Token))
	return nil
}

// Validate checks that the Token field is present and valid.
func (c *BearerTokenAuthClient) Validate() error {
	return authValidator.Struct(c)
}
