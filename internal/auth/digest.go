package auth

import "net/http"

// DigestAuthClient holds credentials for HTTP Digest authentication.
type DigestAuthClient struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

func NewDigestAuthClient(username, password string) *DigestAuthClient {
	return &DigestAuthClient{
		Username: username,
		Password: password,
	}
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
