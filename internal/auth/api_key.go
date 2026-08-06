/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

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

// New initializes the APIKeyAuthClient. It always returns nil.
func (c *APIKeyAuthClient) New() error {
	c.Key = resolveEnv(c.Key)
	c.Value = resolveEnv(c.Value)
	return nil
}

// Type returns the authentication type for API key auth.
func (c *APIKeyAuthClient) Type() Type {
	return APIKeyAuth
}

// Apply places the API key into the request as either a header or query parameter based on the In field.
func (c *APIKeyAuthClient) Apply(req *http.Request, out *Info) error {
	switch c.In {
	case paramInQuery:
		q := req.URL.Query()
		q.Del(c.Key)
		q.Set(c.Key, c.Value)
		req.URL.RawQuery = q.Encode()
		if out != nil {
			if out.QueryParams == nil {
				out.QueryParams = make(map[string]string)
			}
			out.QueryParams[c.Key] = c.Value
		}
	default:
		setAuthHeader(req, out, c.Key, c.Value)
	}
	return nil
}

// Validate checks that the Key, Value, and In fields are present and valid.
func (c *APIKeyAuthClient) Validate() error {
	return authValidator.Struct(c)
}
