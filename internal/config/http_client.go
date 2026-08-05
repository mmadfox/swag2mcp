/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"time"

	"github.com/mmadfox/swag2mcp/internal/env"
)

// Cookie represents an HTTP cookie for configuration.
type Cookie struct {
	Name     string `yaml:"name"     validate:"required"`
	Value    string `yaml:"value"    validate:"required"`
	Domain   string `yaml:"domain,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Secure   bool   `yaml:"secure,omitempty"`
	HTTPOnly bool   `yaml:"http_only,omitempty"`
}

// ProxyConfig holds proxy connection settings.
type ProxyConfig struct {
	URL      string   `yaml:"url"                  validate:"omitempty,proxy_url_format"`
	Username string   `yaml:"username,omitempty"`
	Password string   `yaml:"password,omitempty"`
	Bypass   []string `yaml:"bypass,omitempty"`
}

// HTTPClientConfig holds per-request HTTP settings for a spec or collection.
type HTTPClientConfig struct {
	Randomize       bool              `yaml:"random,omitempty"`
	Proxy           *ProxyConfig      `yaml:"proxy,omitempty"              validate:"omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Cookies         []Cookie          `yaml:"cookies,omitempty"            validate:"omitempty,dive"`
	UserAgent       string            `yaml:"user_agent,omitempty"`
	Timeout         time.Duration     `yaml:"timeout,omitempty"            validate:"omitempty,min=1000000000,max=300000000000"`
	FollowRedirects *bool             `yaml:"follow_redirects,omitempty"`
	MaxRedirects    *int              `yaml:"max_redirects,omitempty"      validate:"omitempty,min=0,max=50"`
	MaxResponseSize *int              `yaml:"max_response_size,omitempty"   validate:"omitempty,min=256,max=10485760"`
}

// GlobalHTTPClientConfig holds global HTTP client settings.
type GlobalHTTPClientConfig struct {
	Randomize       bool              `yaml:"random,omitempty"`
	Proxy           *ProxyConfig      `yaml:"proxy,omitempty"              validate:"omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Cookies         []Cookie          `yaml:"cookies,omitempty"            validate:"omitempty,dive"`
	UserAgent       string            `yaml:"user_agent,omitempty"`
	Timeout         time.Duration     `yaml:"timeout,omitempty"            validate:"omitempty,min=1000000000,max=300000000000"`
	FollowRedirects *bool             `yaml:"follow_redirects,omitempty"`
	MaxRedirects    *int              `yaml:"max_redirects,omitempty"      validate:"omitempty,min=0,max=50"`
	MaxResponseSize *int              `yaml:"max_response_size,omitempty"   validate:"omitempty,min=256,max=10485760"`
}

// SetDefaults fills nil/zero fields with sensible defaults.
func (c *GlobalHTTPClientConfig) SetDefaults() {
	if c == nil {
		return
	}
	if c.UserAgent == "" && !c.Randomize {
		c.UserAgent = defaultUserAgent
	}
	if c.Timeout <= 0 {
		c.Timeout = defaultTimeout
	}
	if c.FollowRedirects == nil {
		v := true
		c.FollowRedirects = &v
	}
	if c.MaxRedirects == nil {
		v := defaultMaxRedirects
		c.MaxRedirects = &v
	}
	if c.MaxResponseSize == nil {
		v := defaultMaxResponseSize
		c.MaxResponseSize = &v
	}
}

// Resolve resolves environment variable references in proxy fields.
func (c *ProxyConfig) Resolve() {
	if c == nil {
		return
	}
	c.URL = env.Parse(c.URL)
	c.Username = env.Parse(c.Username)
	c.Password = env.Parse(c.Password)
}

// Resolve resolves environment variable references in headers and cookie values.
func (c *HTTPClientConfig) Resolve() {
	if c == nil {
		return
	}
	c.Proxy.Resolve()
	for k, v := range c.Headers {
		c.Headers[k] = env.Parse(v)
	}
	for i := range c.Cookies {
		c.Cookies[i].Value = env.Parse(c.Cookies[i].Value)
	}
}

// Resolve resolves environment variable references in headers and cookie values.
func (c *GlobalHTTPClientConfig) Resolve() {
	if c == nil {
		return
	}
	c.Proxy.Resolve()
	for k, v := range c.Headers {
		c.Headers[k] = env.Parse(v)
	}
	for i := range c.Cookies {
		c.Cookies[i].Value = env.Parse(c.Cookies[i].Value)
	}
}
