package auth

import "net/http"

// BasicAuthClient holds credentials for HTTP Basic authentication.
type BasicAuthClient struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

func NewBasicAuthClient(username, password string) *BasicAuthClient {
	return &BasicAuthClient{
		Username: username,
		Password: password,
	}
}

func (c *BasicAuthClient) Type() Type {
	return BasicAuth
}

func (c *BasicAuthClient) Apply(req *http.Request) error {
	req.SetBasicAuth(c.Username, c.Password)
	return nil
}

func (c *BasicAuthClient) Validate() error {
	return authValidator.Struct(c)
}
