package httpclient

import (
	"net/http"
	"testing"
	"time"
)

func TestNew_DefaultTimeout(t *testing.T) {
	client, err := New(Config{})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.Timeout != defaultTimeout {
		t.Errorf("Timeout = %v, want %v", client.Timeout, defaultTimeout)
	}
}

func TestNew_CustomTimeout(t *testing.T) {
	client, err := New(Config{Timeout: 15 * time.Second})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.Timeout != 15*time.Second {
		t.Errorf("Timeout = %v, want %v", client.Timeout, 15*time.Second)
	}
}

func TestNew_NoFollowRedirects(t *testing.T) {
	follow := false
	client, err := New(Config{FollowRedirects: &follow})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil, expected ErrUseLastResponse")
	}
}

func TestNew_MaxRedirects(t *testing.T) {
	maxRedirects := 3
	client, err := New(Config{MaxRedirects: &maxRedirects})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil")
	}
}

func TestNew_DefaultRedirects(t *testing.T) {
	client, err := New(Config{})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client.CheckRedirect != nil {
		t.Fatal("CheckRedirect should be nil by default")
	}
}

func TestNew_NoConfig(t *testing.T) {
	client, err := New(Config{})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNew_Randomize(t *testing.T) {
	client, err := New(Config{
		Randomize: true,
		UserAgent: "test-agent",
		Headers:   map[string]string{"Accept": "text/plain"},
	})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNew_ProxyHTTP(t *testing.T) {
	client, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "http://127.0.0.1:8080",
		},
	})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNew_ProxySOCKS5(t *testing.T) {
	client, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "socks5://127.0.0.1:1080",
		},
	})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNew_ProxySOCKS5h(t *testing.T) {
	client, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "socks5h://127.0.0.1:1080",
		},
	})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNew_ProxyInvalidScheme(t *testing.T) {
	_, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "ftp://127.0.0.1:21",
		},
	})
	if err == nil {
		t.Fatal("expected error for unsupported proxy scheme")
	}
}

func TestNew_ProxyInvalidURL(t *testing.T) {
	_, err := New(Config{
		Proxy: &ProxyConfig{
			URL: "://invalid",
		},
	})
	if err == nil {
		t.Fatal("expected error for invalid proxy URL")
	}
}

func TestNewDefault_NoGlobalConfig(t *testing.T) {
	oldCfg := globalConfig
	globalConfig = Config{}
	t.Cleanup(func() { globalConfig = oldCfg })

	client, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() = %v", err)
	}
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestSetGlobalConfig(t *testing.T) {
	oldCfg := globalConfig
	t.Cleanup(func() { globalConfig = oldCfg })

	SetGlobalConfig(Config{
		Timeout: 10 * time.Second,
	})

	client, err := NewDefault()
	if err != nil {
		t.Fatalf("NewDefault() = %v", err)
	}
	if client.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want %v", client.Timeout, 10*time.Second)
	}
}

func TestMatchBypass_Exact(t *testing.T) {
	if !matchBypass("example.com", []string{"example.com"}) {
		t.Error("expected match for exact host")
	}
}

func TestMatchBypass_Wildcard(t *testing.T) {
	if !matchBypass("api.example.com", []string{"*.example.com"}) {
		t.Error("expected match for wildcard")
	}
}

func TestMatchBypass_NoMatch(t *testing.T) {
	if matchBypass("other.com", []string{"example.com"}) {
		t.Error("expected no match")
	}
}

func TestMatchBypass_Empty(t *testing.T) {
	if matchBypass("example.com", nil) {
		t.Error("expected no match for empty bypass list")
	}
}

func TestRandomizeConfig_Nil(_ *testing.T) {
	RandomizeConfig(nil)
}

func TestRandomizeConfig_FillsEmpty(t *testing.T) {
	cfg := &Config{}
	RandomizeConfig(cfg)

	if cfg.UserAgent == "" {
		t.Error("UserAgent should be filled")
	}
	if cfg.Headers == nil {
		t.Fatal("Headers should be initialized")
	}
	if cfg.Headers["Accept"] == "" {
		t.Error("Accept should be filled")
	}
	if cfg.Headers["Accept-Language"] == "" {
		t.Error("Accept-Language should be filled")
	}
	if cfg.Headers["Accept-Encoding"] == "" {
		t.Error("Accept-Encoding should be filled")
	}
	if cfg.Headers["Referer"] == "" {
		t.Error("Referer should be filled")
	}
	if cfg.Headers["Sec-Ch-Ua"] == "" {
		t.Error("Sec-Ch-Ua should be filled")
	}
}

func TestRandomizeConfig_DoesNotOverwrite(t *testing.T) {
	cfg := &Config{
		UserAgent: "MyBot/1.0",
		Headers: map[string]string{
			"Accept": "application/custom",
		},
	}
	RandomizeConfig(cfg)

	if cfg.UserAgent != "MyBot/1.0" {
		t.Errorf("UserAgent was overwritten: %q", cfg.UserAgent)
	}
	if cfg.Headers["Accept"] != "application/custom" {
		t.Errorf("Accept was overwritten: %q", cfg.Headers["Accept"])
	}
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

	if req.Header.Get("User-Agent") != "test-agent" {
		t.Errorf("User-Agent = %q, want test-agent", req.Header.Get("User-Agent"))
	}
	if req.Header.Get("X-Custom") != "value" {
		t.Errorf("X-Custom = %q, want value", req.Header.Get("X-Custom"))
	}
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

	if req.Header.Get("User-Agent") != "existing-agent" {
		t.Errorf("User-Agent was overwritten: %q", req.Header.Get("User-Agent"))
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("Accept was overwritten: %q", req.Header.Get("Accept"))
	}
}

type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header)}, nil
}
