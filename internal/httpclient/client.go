package httpclient

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

const (
	defaultTimeout      = 30 * time.Second
	maxIdleConnections  = 100
	maxIdleConnsPerHost = 10
	idleConnTimeout     = 90 * time.Second
)

var globalConfig Config

// SetGlobalConfig stores the configuration used by NewDefault.
func SetGlobalConfig(cfg Config) {
	globalConfig = cfg
}

// NewDefault creates an HTTP client from the global configuration.
func NewDefault() (*http.Client, error) {
	return New(globalConfig)
}

// New creates an HTTP client with the given configuration.
func New(cfg Config) (*http.Client, error) {
	transport, err := newTransport(cfg)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}
	if client.Timeout == 0 {
		client.Timeout = defaultTimeout
	}

	applyRedirects(client, cfg)

	return client, nil
}

func newTransport(cfg Config) (http.RoundTripper, error) {
	base, err := newBaseTransport(cfg.Proxy)
	if err != nil {
		return nil, err
	}

	var transport http.RoundTripper = base

	if cfg.Randomize {
		transport = &randomizingTransport{
			Base:      transport,
			UserAgent: cfg.UserAgent,
			Headers:   cfg.Headers,
			Cookies:   cfg.Cookies,
		}
	}

	return transport, nil
}

func newBaseTransport(proxyCfg *ProxyConfig) (*http.Transport, error) {
	t := &http.Transport{
		MaxIdleConns:        maxIdleConnections,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     idleConnTimeout,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	if proxyCfg == nil || proxyCfg.URL == "" {
		return t, nil
	}

	proxyURL, err := url.Parse(proxyCfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL %q: %w", proxyCfg.URL, err)
	}

	switch proxyURL.Scheme {
	case "http", "https":
		t.Proxy = http.ProxyURL(proxyURL)
		if proxyCfg.Username != "" || proxyCfg.Password != "" {
			proxyURL.User = url.UserPassword(proxyCfg.Username, proxyCfg.Password)
			t.Proxy = http.ProxyURL(proxyURL)
		}
		if len(proxyCfg.Bypass) > 0 {
			t.Proxy = bypassProxy(proxyURL, proxyCfg.Bypass)
		}

	case "socks5", "socks5h":
		auth := &proxy.Auth{}
		if proxyCfg.Username != "" {
			auth.User = proxyCfg.Username
			auth.Password = proxyCfg.Password
		}

		var dialer proxy.Dialer
		dialer, err = proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("socks5 dialer: %w", err)
		}

		ctxDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			return nil, errors.New("socks5 dialer does not support context")
		}

		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if len(proxyCfg.Bypass) > 0 && matchBypass(addr, proxyCfg.Bypass) {
				return (&net.Dialer{}).DialContext(ctx, network, addr)
			}
			return ctxDialer.DialContext(ctx, network, addr)
		}

	default:
		return nil, fmt.Errorf("unsupported proxy scheme %q (use http, https, socks5, or socks5h)", proxyURL.Scheme)
	}

	return t, nil
}

func bypassProxy(proxyURL *url.URL, bypass []string) func(*http.Request) (*url.URL, error) {
	return func(req *http.Request) (*url.URL, error) {
		host := req.URL.Hostname()
		if matchBypass(host, bypass) {
			return nil, nil
		}
		return proxyURL, nil
	}
}

func matchBypass(host string, bypass []string) bool {
	for _, b := range bypass {
		switch {
		case strings.HasPrefix(b, "*."):
			suffix := b[1:]
			if strings.HasSuffix(host, suffix) {
				return true
			}
		case b == host:
			return true
		case strings.HasPrefix(b, "/") && strings.HasSuffix(b, "/"):
			pattern := b[1 : len(b)-1]
			if strings.Contains(host, pattern) {
				return true
			}
		}
	}
	return false
}

func applyRedirects(client *http.Client, cfg Config) {
	if cfg.FollowRedirects != nil && !*cfg.FollowRedirects {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if cfg.MaxRedirects != nil {
		maxRedirects := *cfg.MaxRedirects
		client.CheckRedirect = func(_ *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("too many redirects (max %d)", maxRedirects)
			}
			return nil
		}
	}
}
