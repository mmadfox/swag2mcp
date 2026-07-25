package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import "net/http"

// BasicAuthClient holds credentials for HTTP Basic authentication.
type BasicAuthClient struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

// New resolves environment variables in Username and Password and returns nil.
func (c *BasicAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	return nil
}

// Type returns the authentication type for HTTP Basic auth.
func (c *BasicAuthClient) Type() Type {
	return BasicAuth
}

// Apply sets the Basic Authorization header on the request using the configured credentials.
func (c *BasicAuthClient) Apply(req *http.Request, out *Info) error {
	if c.Username == "" || c.Password == "" {
		return nil
	}
	req.SetBasicAuth(c.Username, c.Password)
	if out == nil {
		return nil
	}
	val := req.Header.Get(headerAuthorization)
	if out.Headers == nil {
		out.Headers = make(map[string]string)
	}
	out.Headers[headerAuthorization] = val
	return nil
}

// Validate checks that the Username and Password fields are present and valid.
func (c *BasicAuthClient) Validate() error {
	return authValidator.Struct(c)
}
