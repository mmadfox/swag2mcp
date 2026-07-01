package auth

import "net/http"

// APIKeyAuthClient holds an API key and describes where to place it in the request.
type APIKeyAuthClient struct {
	Key   string `yaml:"key"   validate:"required"`
	Value string `yaml:"value" validate:"required"`
	In    string `yaml:"in"    validate:"required,oneof=header query"`
}

func NewAPIKeyAuthClient(key, value, in string) *APIKeyAuthClient {
	return &APIKeyAuthClient{
		Key:   key,
		Value: value,
		In:    in,
	}
}

func (c *APIKeyAuthClient) Type() Type {
	return APIKeyAuth
}

func (c *APIKeyAuthClient) Apply(req *http.Request) error {
	switch c.In {
	case "query":
		q := req.URL.Query()
		q.Set(c.Key, c.Value)
		req.URL.RawQuery = q.Encode()
	default:
		req.Header.Set(c.Key, c.Value)
	}
	return nil
}

func (c *APIKeyAuthClient) Validate() error {
	return authValidator.Struct(c)
}
