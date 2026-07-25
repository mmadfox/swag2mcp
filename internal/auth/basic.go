package auth

import "net/http"

// BasicAuthClient holds credentials for HTTP Basic authentication.
type BasicAuthClient struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

func (c *BasicAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	return nil
}

func (c *BasicAuthClient) Type() Type {
	return BasicAuth
}

func (c *BasicAuthClient) Apply(req *http.Request, out *Info) error {
	req.SetBasicAuth(c.Username, c.Password)
	if out != nil {
		val := req.Header.Get("Authorization")
		if out.Headers == nil {
			out.Headers = make(map[string]string)
		}
		out.Headers["Authorization"] = val
	}
	return nil
}

func (c *BasicAuthClient) Validate() error {
	return authValidator.Struct(c)
}
