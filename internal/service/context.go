package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"net/http"
	"sync/atomic"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/ratelimit"
)

// serviceContext holds all mutable, concurrency-safe fields shared across sub-services.
type serviceContext struct {
	httpClient         atomic.Value // *http.Client
	httpClientConfig   atomic.Value // httpclient.Config
	config             atomic.Value // *config.Config
	maxResponseSize    atomic.Int64
	startedAt          atomic.Int64 // UnixNano
	globalHeaders      atomic.Value // map[string]string
	globalUserAgent    atomic.Value // string
	globalCookies      atomic.Value // []httpclient.Cookie
	disableLLMAuth     atomic.Bool
	disableRateLimiter atomic.Bool
	rateLimiter        atomic.Value // ratelimit.Limiter
}

func newServiceContext() *serviceContext {
	return &serviceContext{}
}

func (c *serviceContext) loadHTTPClient() *http.Client {
	v := c.httpClient.Load()
	if v == nil {
		return http.DefaultClient
	}
	return v.(*http.Client)
}

func (c *serviceContext) storeHTTPClient(client *http.Client) {
	c.httpClient.Store(client)
}

func (c *serviceContext) loadHTTPClientConfig() httpclient.Config {
	v := c.httpClientConfig.Load()
	if v == nil {
		return httpclient.Config{}
	}
	return v.(httpclient.Config)
}

func (c *serviceContext) storeHTTPClientConfig(cfg httpclient.Config) {
	c.httpClientConfig.Store(cfg)
}

func (c *serviceContext) loadConfig() *config.Config {
	v := c.config.Load()
	if v == nil {
		return nil
	}
	return v.(*config.Config)
}

func (c *serviceContext) storeConfig(cfg *config.Config) {
	c.config.Store(cfg)
}

func (c *serviceContext) loadGlobalHeaders() map[string]string {
	v := c.globalHeaders.Load()
	if v == nil {
		return nil
	}
	return v.(map[string]string)
}

func (c *serviceContext) storeGlobalHeaders(headers map[string]string) {
	c.globalHeaders.Store(headers)
}

func (c *serviceContext) loadGlobalUserAgent() string {
	v := c.globalUserAgent.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (c *serviceContext) storeGlobalUserAgent(ua string) {
	c.globalUserAgent.Store(ua)
}

func (c *serviceContext) loadGlobalCookies() []httpclient.Cookie {
	v := c.globalCookies.Load()
	if v == nil {
		return nil
	}
	return v.([]httpclient.Cookie)
}

func (c *serviceContext) storeGlobalCookies(cookies []httpclient.Cookie) {
	c.globalCookies.Store(cookies)
}

// MaxResponseSize implements SettingsProvider.
func (c *serviceContext) MaxResponseSize() int {
	v := c.maxResponseSize.Load()
	if v <= 0 {
		return config.DefaultMaxResponseSize
	}
	return int(v)
}

// HTTPClientConfig implements settingsProvider.
func (c *serviceContext) HTTPClientConfig() httpclient.Config {
	return c.loadHTTPClientConfig()
}

// Config implements settingsProvider.
func (c *serviceContext) Config() *config.Config {
	return c.loadConfig()
}

func (c *serviceContext) loadRateLimiter() ratelimit.Limiter {
	v := c.rateLimiter.Load()
	if v == nil {
		return nil
	}
	return v.(ratelimit.Limiter)
}

func (c *serviceContext) storeRateLimiter(l ratelimit.Limiter) {
	c.rateLimiter.Store(l)
}
