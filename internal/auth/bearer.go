package auth

import "net/http"

// BearerTokenAuthClient holds a bearer token for HTTP Bearer authentication.
type BearerTokenAuthClient struct {
	Token string `yaml:"token" validate:"required"`
}

func (c *BearerTokenAuthClient) New() error {
	c.Token = resolveEnv(c.Token)
	return nil
}

func (c *BearerTokenAuthClient) Type() Type {
	return BearerTokenAuth
}

func (c *BearerTokenAuthClient) Apply(req *http.Request, out *Info) error {
	if c.Token != "" {
		setAuthHeader(req, out, "Authorization", "Bearer "+c.Token)
	}
	return nil
}

func (c *BearerTokenAuthClient) Validate() error {
	return authValidator.Struct(c)
}
