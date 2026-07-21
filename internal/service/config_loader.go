package service

import (
	"fmt"
	"maps"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
)

const (
	defaultUserAgent      = "swag2mcp-global/1.0"
	defaultMockOAuth2Port = 9090
	defaultMockDigestPort = 9091
)

func buildGlobalHTTPConfig(global *config.GlobalHTTPClientConfig) httpclient.Config {
	if global == nil {
		return httpclient.Config{
			UserAgent: defaultUserAgent,
		}
	}

	cfg := httpclient.Config{
		Randomize:       global.Randomize,
		UserAgent:       global.UserAgent,
		Timeout:         global.Timeout,
		FollowRedirects: global.FollowRedirects,
		MaxRedirects:    global.MaxRedirects,
		MaxResponseSize: global.MaxResponseSize,
	}

	if cfg.UserAgent == "" && !cfg.Randomize {
		cfg.UserAgent = defaultUserAgent
	}
	if global.Headers != nil {
		cfg.Headers = make(map[string]string, len(global.Headers))
		maps.Copy(cfg.Headers, global.Headers)
	}
	if len(global.Cookies) > 0 {
		cfg.Cookies = make([]httpclient.Cookie, len(global.Cookies))
		for i, cookie := range global.Cookies {
			cfg.Cookies[i] = httpclient.Cookie{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Domain:   cookie.Domain,
				Path:     cookie.Path,
				Secure:   cookie.Secure,
				HTTPOnly: cookie.HTTPOnly,
			}
		}
	}
	if global.Proxy != nil {
		cfg.Proxy = &httpclient.ProxyConfig{
			URL:      global.Proxy.URL,
			Username: global.Proxy.Username,
			Password: global.Proxy.Password,
			Bypass:   append([]string{}, global.Proxy.Bypass...),
		}
	}
	return cfg
}

func applyMockAuthURLs(client auth.Authenticator, mockAuth *config.MockAuthConfig) {
	oauth2Port := defaultMockOAuth2Port
	digestPort := defaultMockDigestPort
	if mockAuth != nil {
		if mockAuth.OAuth2Port > 0 {
			oauth2Port = mockAuth.OAuth2Port
		}
		if mockAuth.DigestPort > 0 {
			digestPort = mockAuth.DigestPort
		}
	}
	if setter, ok := client.(auth.TokenURLSetter); ok {
		setter.SetTokenURL(fmt.Sprintf("http://127.0.0.1:%d/token", oauth2Port))
	}
	if setter, ok := client.(auth.MockBaseURLSetter); ok {
		setter.SetMockBaseURL(fmt.Sprintf("http://127.0.0.1:%d/", digestPort))
	}
}

func convertCookies(cookies []config.Cookie) []httpclient.Cookie {
	if len(cookies) == 0 {
		return nil
	}

	result := make([]httpclient.Cookie, len(cookies))
	for index, cookie := range cookies {
		result[index] = httpclient.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
		}
	}

	return result
}

// mergeHTTPClientConfig merges per-request HTTP configs with cascade:
// spec then collection. Collection overrides spec (last-wins).
func mergeHTTPClientConfig(
	spec, collection *config.HTTPClientConfig,
) *model.HTTPClientConfig {
	if spec == nil && collection == nil {
		return &model.HTTPClientConfig{}
	}

	result := &model.HTTPClientConfig{}

	levels := []*config.HTTPClientConfig{spec, collection}

	for _, level := range levels {
		if level == nil {
			continue
		}
		level.Resolve()
		mergeHeaders(result, level)
		mergeCookies(result, level)
	}

	return result
}

func mergeHeaders(result *model.HTTPClientConfig, level *config.HTTPClientConfig) {
	if len(level.Headers) == 0 {
		return
	}
	if result.Headers == nil {
		result.Headers = make(map[string]string, len(level.Headers))
	}
	maps.Copy(result.Headers, level.Headers)
}

func mergeCookies(result *model.HTTPClientConfig, level *config.HTTPClientConfig) {
	if len(level.Cookies) == 0 {
		return
	}
	result.Cookies = make([]httpclient.Cookie, len(level.Cookies))
	for index, cookie := range level.Cookies {
		result.Cookies[index] = httpclient.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
		}
	}
}

// configToModelHTTPClient converts a config-level HTTPClientConfig to a model-level one.
func configToModelHTTPClient(c *config.HTTPClientConfig) *model.HTTPClientConfig {
	if c == nil {
		return nil
	}
	m := &model.HTTPClientConfig{
		Randomize:       c.Randomize,
		Headers:         copyMap(c.Headers),
		Cookies:         convertCookies(c.Cookies),
		UserAgent:       c.UserAgent,
		Timeout:         c.Timeout,
		FollowRedirects: c.FollowRedirects,
		MaxRedirects:    c.MaxRedirects,
		MaxResponseSize: c.MaxResponseSize,
	}
	if c.Proxy != nil {
		m.Proxy = &httpclient.ProxyConfig{
			URL:      c.Proxy.URL,
			Username: c.Proxy.Username,
			Password: c.Proxy.Password,
			Bypass:   append([]string{}, c.Proxy.Bypass...),
		}
	}
	return m
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	r := make(map[K]V, len(m))
	maps.Copy(r, m)
	return r
}
