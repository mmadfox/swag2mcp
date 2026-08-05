/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientConfig_MaxResponseSize(t *testing.T) {
	t.Parallel()

	val := 4096
	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			MaxResponseSize: &val,
		},
	}

	require.NotNil(t, cfg.HTTPClient.MaxResponseSize)
	assert.Equal(t, 4096, *cfg.HTTPClient.MaxResponseSize)
}

func TestHTTPClientConfig_MaxResponseSize_Nil(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{},
	}
	assert.Nil(t, cfg.HTTPClient.MaxResponseSize, "MaxResponseSize should be nil by default")
}

func TestGlobalHTTPClientConfig_SetDefaults(t *testing.T) {
	t.Parallel()

	cfg := &GlobalHTTPClientConfig{}
	cfg.SetDefaults()

	assert.Equal(t, "swag2mcp-global/1.0", cfg.UserAgent)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
	require.NotNil(t, cfg.FollowRedirects)
	assert.True(t, *cfg.FollowRedirects)
	require.NotNil(t, cfg.MaxRedirects)
	assert.Equal(t, 10, *cfg.MaxRedirects)
	require.NotNil(t, cfg.MaxResponseSize)
	assert.Equal(t, 1048576, *cfg.MaxResponseSize)
}

func TestGlobalHTTPClientConfig_SetDefaults_Nil(t *testing.T) {
	t.Parallel()

	var cfg *GlobalHTTPClientConfig
	cfg.SetDefaults() // should not panic
}

func TestGlobalHTTPClientConfig_SetDefaults_DoesNotOverwrite(t *testing.T) {
	t.Parallel()

	timeout := 10 * time.Second
	follow := false
	maxRedir := 5
	maxSize := 4096

	cfg := &GlobalHTTPClientConfig{
		UserAgent:       "custom-agent/1.0",
		Timeout:         timeout,
		FollowRedirects: &follow,
		MaxRedirects:    &maxRedir,
		MaxResponseSize: &maxSize,
	}
	cfg.SetDefaults()

	assert.Equal(t, "custom-agent/1.0", cfg.UserAgent)
	assert.Equal(t, timeout, cfg.Timeout)
	assert.Equal(t, follow, *cfg.FollowRedirects)
	assert.Equal(t, maxRedir, *cfg.MaxRedirects)
	assert.Equal(t, maxSize, *cfg.MaxResponseSize)
}

func TestGlobalHTTPClientConfig_SetDefaults_WithRandomize(t *testing.T) {
	t.Parallel()

	cfg := &GlobalHTTPClientConfig{
		Randomize: true,
	}
	cfg.SetDefaults()

	assert.Empty(t, cfg.UserAgent, "UserAgent should be empty when Randomize is true")
}

func TestHTTPClientConfig_Resolve(t *testing.T) {
	t.Setenv("MY_HEADER", "resolved-header")
	t.Setenv("MY_COOKIE", "resolved-cookie")

	cfg := &HTTPClientConfig{
		Headers: map[string]string{
			"X-Custom": "$(MY_HEADER)",
			"X-Static": "static-value",
		},
		Cookies: []Cookie{
			{Name: "session", Value: "$(MY_COOKIE)"},
			{Name: "static", Value: "static-cookie"},
		},
	}

	cfg.Resolve()

	assert.Equal(t, "resolved-header", cfg.Headers["X-Custom"])
	assert.Equal(t, "static-value", cfg.Headers["X-Static"])
	assert.Equal(t, "resolved-cookie", cfg.Cookies[0].Value)
	assert.Equal(t, "static-cookie", cfg.Cookies[1].Value)
}

func TestHTTPClientConfig_Resolve_Nil(_ *testing.T) {
	var cfg *HTTPClientConfig
	cfg.Resolve() // should not panic
}

func TestGlobalHTTPClientConfig_Resolve(t *testing.T) {
	t.Setenv("GLOBAL_HEADER", "global-val")

	cfg := &GlobalHTTPClientConfig{
		Headers: map[string]string{
			"X-Global": "$(GLOBAL_HEADER)",
		},
		Cookies: []Cookie{
			{Name: "gc", Value: "$(GLOBAL_HEADER)"},
		},
	}

	cfg.Resolve()

	assert.Equal(t, "global-val", cfg.Headers["X-Global"])
	assert.Equal(t, "global-val", cfg.Cookies[0].Value)
}

func TestHTTPClientConfig_Validate_TransportFields(t *testing.T) {
	t.Parallel()

	timeout := 5 * time.Second
	follow := false
	maxRedir := 3
	maxSize := 4096

	cfg := &HTTPClientConfig{
		Timeout:         timeout,
		FollowRedirects: &follow,
		MaxRedirects:    &maxRedir,
		MaxResponseSize: &maxSize,
		Randomize:       true,
		UserAgent:       "spec-agent",
		Proxy: &ProxyConfig{
			URL: "http://proxy:8080",
		},
	}

	errs := collectStructErrors("http_client", *cfg)
	assert.Empty(t, errs, "expected no validation errors")
}

func TestHTTPClientConfig_Validate_InvalidTimeout(t *testing.T) {
	t.Parallel()

	cfg := &HTTPClientConfig{
		Timeout: 100 * time.Millisecond, // below 1s minimum
	}

	errs := collectStructErrors("http_client", *cfg)
	assert.NotEmpty(t, errs, "expected validation error for too-short timeout")
}

func TestProxyConfig_Resolve(t *testing.T) {
	t.Setenv("PROXY_URL", "http://proxy.example.com:8080")
	t.Setenv("PROXY_USER", "user1")
	t.Setenv("PROXY_PASS", "pass1")

	cfg := &ProxyConfig{
		URL:      "$(PROXY_URL)",
		Username: "$(PROXY_USER)",
		Password: "$(PROXY_PASS)",
	}
	cfg.Resolve()

	assert.Equal(t, "http://proxy.example.com:8080", cfg.URL)
	assert.Equal(t, "user1", cfg.Username)
	assert.Equal(t, "pass1", cfg.Password)
}

func TestProxyConfig_Resolve_nil(t *testing.T) {
	t.Parallel()

	var cfg *ProxyConfig
	cfg.Resolve() // should not panic
}

func TestHTTPClientConfig_Resolve_Proxy(t *testing.T) {
	t.Setenv("PROXY_URL", "http://proxy:3128")

	cfg := &HTTPClientConfig{
		Proxy: &ProxyConfig{
			URL: "$(PROXY_URL)",
		},
	}
	cfg.Resolve()

	assert.Equal(t, "http://proxy:3128", cfg.Proxy.URL)
}

func TestGlobalHTTPClientConfig_Resolve_Proxy(t *testing.T) {
	t.Setenv("PROXY_URL", "socks5://127.0.0.1:1080")

	cfg := &GlobalHTTPClientConfig{
		Proxy: &ProxyConfig{
			URL: "$(PROXY_URL)",
		},
	}
	cfg.Resolve()

	assert.Equal(t, "socks5://127.0.0.1:1080", cfg.Proxy.URL)
}
