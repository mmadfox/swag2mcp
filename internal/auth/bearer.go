package auth

import "net/http"

// BearerTokenAuthClient holds a bearer token for HTTP Bearer authentication.
type BearerTokenAuthClient struct {
	Token string `yaml:"token" validate:"required"`
}

func NewBearerTokenAuthClient(token string) *BearerTokenAuthClient {
	return &BearerTokenAuthClient{
		Token: token,
	}
}

func (c *BearerTokenAuthClient) Type() Type {
	return BearerTokenAuth
}

func (c *BearerTokenAuthClient) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	return nil
}

func (c *BearerTokenAuthClient) Validate() error {
	return authValidator.Struct(c)
}
