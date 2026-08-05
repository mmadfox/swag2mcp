/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import "time"

const (
	// DefaultGlobalRateLimit is the default maximum number of invoke requests per second.
	DefaultGlobalRateLimit = 5
	// MinResponseSize is the minimum allowed response size in bytes (10 KB).
	MinResponseSize = 10240
	// MaxAllowedResponseSize is the maximum allowed response size in bytes (10 MB).
	MaxAllowedResponseSize = 10485760
	// DefaultMaxResponseSize is the default maximum response size in bytes (1 MB).
	DefaultMaxResponseSize = 1048576
	// DefaultRateLimitInterval is the default per-endpoint rate limit interval.
	DefaultRateLimitInterval = 10 * time.Second
	// RandSuffixLen is the length of the random suffix appended to randomized user agents.
	RandSuffixLen = 6
)

const (
	defaultUserAgent         = "swag2mcp-global/1.0"
	defaultTimeout           = 30 * time.Second
	defaultMaxRedirects      = 10
	defaultRateLimitInterval = 10 * time.Second
	defaultMaxResponseSize   = DefaultMaxResponseSize
	// DefaultMCPTransport is the default MCP transport protocol.
	DefaultMCPTransport = "stdio"
)
