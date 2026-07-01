package auth

import "net/http"

// DigestAuthClient holds credentials for HTTP Digest authentication.
type DigestAuthClient struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

func (c *DigestAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	return nil
}

func (c *DigestAuthClient) Type() Type {
	return DigestAuth
}

func (c *DigestAuthClient) Apply(_ *http.Request) error {
	return nil
}

func (c *DigestAuthClient) Validate() error {
	return authValidator.Struct(c)
}
