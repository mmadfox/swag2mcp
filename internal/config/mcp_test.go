/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMCPAuthConfig_Resolve_Nil(t *testing.T) {
	t.Parallel()

	var c *MCPAuthConfig
	c.Resolve() // should not panic
}

func TestMCPAuthConfig_Resolve_NoEnv(t *testing.T) {
	t.Parallel()

	c := &MCPAuthConfig{Token: "static-token"}
	c.Resolve()
	assert.Equal(t, "static-token", c.Token)
}

func TestMCPAuthConfig_Resolve_WithEnv(t *testing.T) {
	t.Setenv("MCP_TOKEN", "resolved-token")
	c := &MCPAuthConfig{Token: "$(MCP_TOKEN)"}
	c.Resolve()
	assert.Equal(t, "resolved-token", c.Token)
}

func TestMCPConfig_Defaults(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MCP: &MCPConfig{
			Transport: "sse",
			Addr:      ":9090",
			Path:      "/api/mcp",
		},
	}
	assert.Equal(t, "sse", cfg.MCP.Transport)
	assert.Equal(t, ":9090", cfg.MCP.Addr)
	assert.Equal(t, "/api/mcp", cfg.MCP.Path)
}
