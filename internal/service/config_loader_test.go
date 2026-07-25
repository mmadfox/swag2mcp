package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
)

func TestBuildGlobalHTTPConfig_nil(t *testing.T) {
	t.Parallel()

	cfg := BuildGlobalHTTPConfig(nil)
	require.Equal(t, "swag2mcp-global/1.0", cfg.UserAgent)
}

func TestBuildGlobalHTTPConfig_withValues(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		UserAgent: "custom-agent",
		Timeout:   60 * time.Second,
		Headers:   map[string]string{"X-Custom": "val"},
		Cookies: []config.Cookie{
			{Name: "session", Value: "abc"},
		},
		Proxy: &config.ProxyConfig{
			URL: "http://proxy:8080",
		},
	}

	cfg := BuildGlobalHTTPConfig(global)
	require.Equal(t, "custom-agent", cfg.UserAgent)
	require.Equal(t, 60*time.Second, cfg.Timeout)
	require.Equal(t, "val", cfg.Headers["X-Custom"])
	require.Len(t, cfg.Cookies, 1)
	require.NotNil(t, cfg.Proxy)
	require.Equal(t, "http://proxy:8080", cfg.Proxy.URL)
}

func TestConvertCookies_nil(t *testing.T) {
	t.Parallel()

	result := convertCookies(nil)
	require.Nil(t, result)
}

func TestConvertCookies_withValues(t *testing.T) {
	t.Parallel()

	cookies := []config.Cookie{
		{Name: "a", Value: "1", Domain: "example.com", Path: "/", Secure: true, HTTPOnly: true},
	}
	result := convertCookies(cookies)
	require.Len(t, result, 1)
	require.Equal(t, "a", result[0].Name)
	require.Equal(t, "1", result[0].Value)
	require.True(t, result[0].Secure)
}

func TestMergeHTTPClientConfig_bothNil(t *testing.T) {
	t.Parallel()

	result := mergeHTTPClientConfig(nil, nil)
	require.NotNil(t, result)
}

func TestMergeHTTPClientConfig_specOnly(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Api": "key"},
	}
	result := mergeHTTPClientConfig(spec, nil)
	require.Equal(t, "key", result.Headers["X-Api"])
}

func TestMergeHTTPClientConfig_collectionOverrides(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Api": "spec-value"},
	}
	coll := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Api": "coll-value"},
	}
	result := mergeHTTPClientConfig(spec, coll)
	require.Equal(t, "coll-value", result.Headers["X-Api"])
}

func TestConfigToModelHTTPClient_nil(t *testing.T) {
	t.Parallel()

	result := configToModelHTTPClient(nil)
	require.Nil(t, result)
}

func TestConfigToModelHTTPClient_withProxy(t *testing.T) {
	t.Parallel()

	c := &config.HTTPClientConfig{
		UserAgent: "ua",
		Proxy: &config.ProxyConfig{
			URL: "http://proxy:9090",
		},
	}
	result := configToModelHTTPClient(c)
	require.NotNil(t, result)
	require.Equal(t, "ua", result.UserAgent)
	require.NotNil(t, result.Proxy)
	require.Equal(t, "http://proxy:9090", result.Proxy.URL)
}

func TestCopyMap_nil(t *testing.T) {
	t.Parallel()

	result := copyMap[string, string](nil)
	require.Nil(t, result)
}

func TestCopyMap_withValues(t *testing.T) {
	t.Parallel()

	m := map[string]string{"a": "1", "b": "2"}
	result := copyMap(m)
	require.Equal(t, m, result)
	result["a"] = "changed"
	require.Equal(t, "1", m["a"]) // original unchanged
}

func TestMergeCookies_withValues(t *testing.T) {
	t.Parallel()

	result := &model.HTTPClientConfig{}
	level := &config.HTTPClientConfig{
		Cookies: []config.Cookie{
			{Name: "session", Value: "abc123"},
		},
	}
	mergeCookies(result, level)
	require.Len(t, result.Cookies, 1)
	require.Equal(t, "session", result.Cookies[0].Name)
	require.Equal(t, "abc123", result.Cookies[0].Value)
}

func TestMergeCookies_nil(t *testing.T) {
	t.Parallel()

	result := &model.HTTPClientConfig{}
	mergeCookies(result, &config.HTTPClientConfig{})
	require.Nil(t, result.Cookies)
}

func TestApplyMockAuthURLs_nilConfig(t *testing.T) {
	t.Parallel()

	client := auth.NewNoAuthClient()
	applyMockAuthURLs(client, nil)
	// No panic, no error - just verifies it handles nil config
}

func TestApplyMockAuthURLs_withPorts(t *testing.T) {
	t.Parallel()

	client := auth.NewNoAuthClient()
	applyMockAuthURLs(client, &config.MockAuthConfig{
		OAuth2Port: 9095,
		DigestPort: 9096,
	})
	// No panic, no error - NoAuthClient doesn't implement TokenURLSetter or MockBaseURLSetter
}
