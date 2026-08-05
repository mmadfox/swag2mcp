/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import "github.com/mmadfox/swag2mcp/internal/env"

// MCPConfig holds the MCP server configuration.
type MCPConfig struct {
	Transport string         `yaml:"transport,omitempty" validate:"omitempty,oneof=stdio sse streamable-http"`
	Addr      string         `yaml:"addr,omitempty"`
	Path      string         `yaml:"path,omitempty"`
	Auth      *MCPAuthConfig `yaml:"auth,omitempty"`
}

// MCPAuthConfig holds the MCP server authentication configuration.
type MCPAuthConfig struct {
	Token            string `yaml:"token,omitempty"`
	Type             string `yaml:"type,omitempty"              validate:"omitempty,oneof=jwks introspection oidc"`
	JWKSURL          string `yaml:"jwks_url,omitempty"`
	Issuer           string `yaml:"issuer,omitempty"`
	Audience         string `yaml:"audience,omitempty"`
	IntrospectionURL string `yaml:"introspection_url,omitempty"`
	ClientID         string `yaml:"client_id,omitempty"`
	ClientSecret     string `yaml:"client_secret,omitempty"`
}

// Resolve resolves environment variable references in the token and JWT fields.
func (c *MCPAuthConfig) Resolve() {
	if c == nil {
		return
	}
	c.Token = env.Parse(c.Token)
	c.JWKSURL = env.Parse(c.JWKSURL)
	c.Issuer = env.Parse(c.Issuer)
	c.Audience = env.Parse(c.Audience)
	c.IntrospectionURL = env.Parse(c.IntrospectionURL)
	c.ClientID = env.Parse(c.ClientID)
	c.ClientSecret = env.Parse(c.ClientSecret)
}
