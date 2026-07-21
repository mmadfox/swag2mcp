package service

import (
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/stretchr/testify/require"
)

func TestServiceContext_HTTPClient(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Nil(t, c.loadHTTPClient())

	client := &http.Client{}
	c.storeHTTPClient(client)
	require.Same(t, client, c.loadHTTPClient())
}

func TestServiceContext_HTTPClientConfig(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Equal(t, httpclient.Config{}, c.loadHTTPClientConfig())

	cfg := httpclient.Config{UserAgent: "test"}
	c.storeHTTPClientConfig(cfg)
	require.Equal(t, "test", c.loadHTTPClientConfig().UserAgent)
}

func TestServiceContext_Config(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Nil(t, c.loadConfig())

	cfg := &config.Config{}
	c.storeConfig(cfg)
	require.Same(t, cfg, c.loadConfig())
}

func TestServiceContext_GlobalHeaders(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Nil(t, c.loadGlobalHeaders())

	h := map[string]string{"X-Custom": "val"}
	c.storeGlobalHeaders(h)
	require.Equal(t, "val", c.loadGlobalHeaders()["X-Custom"])
}

func TestServiceContext_GlobalUserAgent(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Empty(t, c.loadGlobalUserAgent())

	c.storeGlobalUserAgent("agent")
	require.Equal(t, "agent", c.loadGlobalUserAgent())
}

func TestServiceContext_GlobalCookies(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Nil(t, c.loadGlobalCookies())

	cookies := []httpclient.Cookie{{Name: "s", Value: "v"}}
	c.storeGlobalCookies(cookies)
	require.Len(t, c.loadGlobalCookies(), 1)
}

func TestServiceContext_MaxResponseSize(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	require.Equal(t, 0, c.MaxResponseSize())

	c.maxResponseSize.Store(2048)
	require.Equal(t, 2048, c.MaxResponseSize())
}

func TestServiceContext_HTTPClientConfigMethod(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	cfg := httpclient.Config{UserAgent: "ua"}
	c.storeHTTPClientConfig(cfg)
	require.Equal(t, "ua", c.HTTPClientConfig().UserAgent)
}

func TestServiceContext_ConfigMethod(t *testing.T) {
	t.Parallel()

	c := newServiceContext()
	cfg := &config.Config{}
	c.storeConfig(cfg)
	require.Same(t, cfg, c.Config())
}
