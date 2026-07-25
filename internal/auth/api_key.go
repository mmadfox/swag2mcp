package auth

import (
	"net/http"
)

// APIKeyAuthClient holds an API key and describes where to place it in the request.
type APIKeyAuthClient struct {
	Key   string `yaml:"key"   validate:"required"`
	Value string `yaml:"value" validate:"required"`
	In    string `yaml:"in"    validate:"required,oneof=header query"`
}

func (c *APIKeyAuthClient) New() error {
	return nil
}

func (c *APIKeyAuthClient) Type() Type {
	return APIKeyAuth
}

func (c *APIKeyAuthClient) Apply(req *http.Request, out *Info) error {
	switch c.In {
	case "query":
		setAuthQuery(req, out, c.Key, c.Value)
	default:
		setAuthHeader(req, out, c.Key, c.Value)
	}
	return nil
}

func (c *APIKeyAuthClient) Validate() error {
	return authValidator.Struct(c)
}
