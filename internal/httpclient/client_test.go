package httpclient

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultTimeout(t *testing.T) {
	client, err := New(Config{})
	require.NoError(t, err)
	assert.Equal(t, defaultTimeout, client.Timeout)
}

func TestNew_CustomTimeout(t *testing.T) {
	client, err := New(Config{Timeout: 15 * time.Second})
	require.NoError(t, err)
	assert.Equal(t, 15*time.Second, client.Timeout)
}

func TestNew_NoFollowRedirects(t *testing.T) {
	follow := false
	client, err := New(Config{FollowRedirects: &follow})
	require.NoError(t, err)
	require.NotNil(t, client.CheckRedirect)
}

func TestNew_MaxRedirects(t *testing.T) {
	maxRedirects := 3
	client, err := New(Config{MaxRedirects: &maxRedirects})
	require.NoError(t, err)
	require.NotNil(t, client.CheckRedirect)
}

func TestNew_DefaultRedirects(t *testing.T) {
	client, err := New(Config{})
	require.NoError(t, err)
	assert.Nil(t, client.CheckRedirect)
}

func TestNew_NoConfig(t *testing.T) {
	client, err := New(Config{})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNew_Randomize(t *testing.T) {
	client, err := New(Config{
		Randomize: true,
		UserAgent: "test-agent",
		Headers:   map[string]string{"Accept": "text/plain"},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNew_ProxyHTTP(t *testing.T) {
	client, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "http://127.0.0.1:8080",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNew_ProxySOCKS5(t *testing.T) {
	client, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "socks5://127.0.0.1:1080",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNew_ProxySOCKS5h(t *testing.T) {
	client, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "socks5h://127.0.0.1:1080",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNew_ProxyInvalidScheme(t *testing.T) {
	_, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "ftp://127.0.0.1:21",
		},
	})
	require.Error(t, err)
}

func TestNew_ProxyInvalidURL(t *testing.T) {
	_, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "://invalid",
		},
	})
	require.Error(t, err)
}

func TestNewDefault_NoGlobalConfig(t *testing.T) {
	globalConfig = atomic.Value{}
	globalConfig.Store(Config{})

	client, err := NewDefault()
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestSetGlobalConfig(t *testing.T) {
	globalConfig = atomic.Value{}

	SetGlobalConfig(Config{
		Timeout: 10 * time.Second,
	})

	client, err := NewDefault()
	require.NoError(t, err)
	assert.Equal(t, 10*time.Second, client.Timeout)
}

func TestMatchBypass_Exact(t *testing.T) {
	assert.True(t, matchBypass("example.com", []string{"example.com"}))
}

func TestMatchBypass_Wildcard(t *testing.T) {
	assert.True(t, matchBypass("api.example.com", []string{"*.example.com"}))
}

func TestMatchBypass_NoMatch(t *testing.T) {
	assert.False(t, matchBypass("other.com", []string{"example.com"}))
}

func TestMatchBypass_Empty(t *testing.T) {
	assert.False(t, matchBypass("example.com", nil))
}

func TestBypassProxy_Matches(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy:8080")
	bypassFn := bypassProxy(proxyURL, []string{"example.com"})

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/path", nil)
	result, err := bypassFn(req)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestBypassProxy_NoMatch(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy:8080")
	bypassFn := bypassProxy(proxyURL, []string{"example.com"})

	req, _ := http.NewRequest(http.MethodGet, "http://other.com/path", nil)
	result, err := bypassFn(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "http://proxy:8080", result.String())
}

func TestBypassProxy_EmptyBypass(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy:8080")
	bypassFn := bypassProxy(proxyURL, nil)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/path", nil)
	result, err := bypassFn(req)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestBypassProxy_Wildcard(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy:8080")
	bypassFn := bypassProxy(proxyURL, []string{"*.example.com"})

	req, _ := http.NewRequest(http.MethodGet, "http://api.example.com/path", nil)
	result, err := bypassFn(req)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestBypassProxy_RegexPattern(t *testing.T) {
	proxyURL, _ := url.Parse("http://proxy:8080")
	bypassFn := bypassProxy(proxyURL, []string{"/internal/"})

	req, _ := http.NewRequest(http.MethodGet, "http://api.internal.example.com/path", nil)
	result, err := bypassFn(req)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestApplyRedirects_NoFollow(t *testing.T) {
	client := &http.Client{}
	follow := false
	applyRedirects(client, Config{FollowRedirects: &follow})

	require.NotNil(t, client.CheckRedirect)
	err := client.CheckRedirect(nil, nil)
	assert.True(t, errors.Is(err, http.ErrUseLastResponse))
}

func TestApplyRedirects_MaxRedirects(t *testing.T) {
	client := &http.Client{}
	maxRedir := 3
	applyRedirects(client, Config{MaxRedirects: &maxRedir})

	require.NotNil(t, client.CheckRedirect)

	req1, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req2, _ := http.NewRequest(http.MethodGet, "http://example.com/2", nil)
	req3, _ := http.NewRequest(http.MethodGet, "http://example.com/3", nil)

	err := client.CheckRedirect(req1, []*http.Request{req1, req2, req3})
	require.Error(t, err)
}

func TestApplyRedirects_UnderMaxRedirects(t *testing.T) {
	client := &http.Client{}
	maxRedir := 5
	applyRedirects(client, Config{MaxRedirects: &maxRedir})

	require.NotNil(t, client.CheckRedirect)

	req1, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	err := client.CheckRedirect(req1, []*http.Request{req1})
	require.NoError(t, err)
}

func TestApplyRedirects_Default(t *testing.T) {
	client := &http.Client{}
	applyRedirects(client, Config{})
	assert.Nil(t, client.CheckRedirect)
}

func TestApplyRedirects_FollowTrue(t *testing.T) {
	client := &http.Client{}
	follow := true
	applyRedirects(client, Config{FollowRedirects: &follow})
	assert.Nil(t, client.CheckRedirect)
}

func TestMatchBypass_Regex(t *testing.T) {
	assert.True(t, matchBypass("api.internal.example.com", []string{"/internal/"}))
}

func TestMatchBypass_RegexNoMatch(t *testing.T) {
	assert.False(t, matchBypass("api.example.com", []string{"/internal/"}))
}

func TestMatchBypass_WildcardNoMatch(t *testing.T) {
	assert.False(t, matchBypass("example.com", []string{"*.example.com"}))
}

func TestMatchBypass_Multiple(t *testing.T) {
	assert.True(t, matchBypass("test.local", []string{"example.com", "*.local", "other.com"}))
}

func TestMatchBypass_NoMatchMultiple(t *testing.T) {
	assert.False(t, matchBypass("remote.com", []string{"example.com", "*.local"}))
}

func TestRandomizeConfig_Nil(_ *testing.T) {
	RandomizeConfig(nil)
}

func TestRandomizeConfig_FillsEmpty(t *testing.T) {
	cfg := &Config{}
	RandomizeConfig(cfg)

	require.NotEmpty(t, cfg.UserAgent)
	require.NotNil(t, cfg.Headers)
	assert.NotEmpty(t, cfg.Headers["Accept"])
	assert.NotEmpty(t, cfg.Headers["Accept-Language"])
	assert.NotEmpty(t, cfg.Headers["Accept-Encoding"])
	assert.NotEmpty(t, cfg.Headers["Referer"])
	assert.NotEmpty(t, cfg.Headers["Sec-Ch-Ua"])
}

func TestRandomizeConfig_DoesNotOverwrite(t *testing.T) {
	cfg := &Config{
		UserAgent: "MyBot/1.0",
		Headers: map[string]string{
			"Accept": "application/custom",
		},
	}
	RandomizeConfig(cfg)

	assert.Equal(t, "MyBot/1.0", cfg.UserAgent)
	assert.Equal(t, "application/custom", cfg.Headers["Accept"])
}

func TestRandomizingTransport_SetsHeaders(t *testing.T) {
	base := &mockRoundTripper{}
	rt := &randomizingTransport{
		Base:      base,
		UserAgent: "test-agent",
		Headers:   map[string]string{"X-Custom": "value"},
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	_, _ = rt.RoundTrip(req)

	assert.Equal(t, "test-agent", req.Header.Get("User-Agent"))
	assert.Equal(t, "value", req.Header.Get("X-Custom"))
}

func TestRandomizingTransport_DoesNotOverwrite(t *testing.T) {
	base := &mockRoundTripper{}
	rt := &randomizingTransport{
		Base:      base,
		UserAgent: "test-agent",
		Headers:   map[string]string{"Accept": "text/html"},
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Set("User-Agent", "existing-agent")
	req.Header.Set("Accept", "application/json")
	_, _ = rt.RoundTrip(req)

	assert.Equal(t, "existing-agent", req.Header.Get("User-Agent"))
	assert.Equal(t, "application/json", req.Header.Get("Accept"))
}

type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}, nil
}
