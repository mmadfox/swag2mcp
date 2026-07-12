package service

import (
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
)

func TestBuildGlobalHTTPConfig_Nil(t *testing.T) {
	t.Parallel()

	cfg := buildGlobalHTTPConfig(nil)

	if cfg.Randomize {
		t.Error("Randomize = true, want false")
	}
	if cfg.Proxy != nil {
		t.Error("Proxy != nil, want nil")
	}
}

func TestBuildGlobalHTTPConfig_Full(t *testing.T) {
	t.Parallel()

	timeout := 30 * time.Second
	follow := false
	maxRedir := 5
	maxSize := 4096

	global := &config.GlobalHTTPClientConfig{
		Randomize:       true,
		UserAgent:       "test-agent",
		Timeout:         timeout,
		FollowRedirects: &follow,
		MaxRedirects:    &maxRedir,
		MaxResponseSize: &maxSize,
		Headers:         map[string]string{"X-Custom": "value"},
		Cookies: []config.Cookie{
			{Name: "session", Value: "abc", Domain: ".example.com", Path: "/", Secure: true, HTTPOnly: true},
		},
		Proxy: &config.ProxyConfig{
			URL:      "http://proxy:8080",
			Username: "user",
			Password: "pass",
			Bypass:   []string{"localhost"},
		},
	}

	cfg := buildGlobalHTTPConfig(global)

	if !cfg.Randomize {
		t.Error("Randomize = false, want true")
	}
	if cfg.UserAgent != "test-agent" {
		t.Errorf("UserAgent = %q, want %q", cfg.UserAgent, "test-agent")
	}
	if cfg.Timeout != timeout {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, timeout)
	}
	if cfg.FollowRedirects == nil || *cfg.FollowRedirects != false {
		t.Error("FollowRedirects = nil or true, want false")
	}
	if cfg.MaxRedirects == nil || *cfg.MaxRedirects != 5 {
		t.Error("MaxRedirects = nil or wrong, want 5")
	}
	if cfg.MaxResponseSize == nil || *cfg.MaxResponseSize != 4096 {
		t.Error("MaxResponseSize = nil or wrong, want 4096")
	}
	if cfg.Headers["X-Custom"] != "value" {
		t.Errorf("Headers[X-Custom] = %q, want %q", cfg.Headers["X-Custom"], "value")
	}
	if len(cfg.Cookies) != 1 || cfg.Cookies[0].Name != "session" {
		t.Error("Cookies missing session")
	}
	if cfg.Proxy == nil || cfg.Proxy.URL != "http://proxy:8080" {
		t.Error("Proxy missing or wrong URL")
	}
	if cfg.Proxy.Username != "user" {
		t.Errorf("Proxy.Username = %q, want %q", cfg.Proxy.Username, "user")
	}
	if len(cfg.Proxy.Bypass) != 1 || cfg.Proxy.Bypass[0] != "localhost" {
		t.Error("Proxy.Bypass missing localhost")
	}
}

func TestBuildGlobalHTTPConfig_NoHeaders(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{Randomize: true}

	cfg := buildGlobalHTTPConfig(global)

	if !cfg.Randomize {
		t.Error("Randomize = false, want true")
	}
	if cfg.Headers != nil {
		t.Error("Headers should be nil when not set")
	}
	if cfg.Cookies != nil {
		t.Error("Cookies should be nil when not set")
	}
	if cfg.Proxy != nil {
		t.Error("Proxy should be nil when not set")
	}
}

func TestBuildGlobalHTTPConfig_NoProxy(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		Headers: map[string]string{"X-Test": "val"},
	}

	cfg := buildGlobalHTTPConfig(global)

	if cfg.Proxy != nil {
		t.Error("Proxy should be nil")
	}
	if cfg.Headers["X-Test"] != "val" {
		t.Errorf("Headers[X-Test] = %q, want %q", cfg.Headers["X-Test"], "val")
	}
}

func TestBuildGlobalHTTPConfig_NoCookies(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{Cookies: []config.Cookie{}}

	cfg := buildGlobalHTTPConfig(global)

	if cfg.Cookies != nil {
		t.Error("Cookies should be nil for empty slice")
	}
}

func TestBuildGlobalHTTPConfig_ProxyNoAuth(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		Proxy: &config.ProxyConfig{
			URL:    "http://proxy:8080",
			Bypass: []string{"localhost", "127.0.0.1"},
		},
	}

	cfg := buildGlobalHTTPConfig(global)

	if cfg.Proxy == nil {
		t.Fatal("Proxy is nil")
	}
	if cfg.Proxy.URL != "http://proxy:8080" {
		t.Errorf("Proxy.URL = %q, want %q", cfg.Proxy.URL, "http://proxy:8080")
	}
	if cfg.Proxy.Username != "" {
		t.Errorf("Proxy.Username = %q, want empty", cfg.Proxy.Username)
	}
	if len(cfg.Proxy.Bypass) != 2 {
		t.Errorf("Proxy.Bypass length = %d, want 2", len(cfg.Proxy.Bypass))
	}
}

func TestBuildGlobalHTTPConfig_ConvertCookies(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		Cookies: []config.Cookie{
			{Name: "c1", Value: "v1"},
			{Name: "c2", Value: "v2", Domain: ".example.com", Path: "/", Secure: true, HTTPOnly: true},
		},
	}

	cfg := buildGlobalHTTPConfig(global)

	if len(cfg.Cookies) != 2 {
		t.Fatalf("Cookies length = %d, want 2", len(cfg.Cookies))
	}
	if cfg.Cookies[0].Name != "c1" || cfg.Cookies[0].Value != "v1" {
		t.Error("Cookies[0] wrong")
	}
	if cfg.Cookies[1].Name != "c2" || !cfg.Cookies[1].Secure || !cfg.Cookies[1].HTTPOnly {
		t.Error("Cookies[1] wrong")
	}
}

func TestBuildGlobalHTTPConfig_ZeroValues(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{}

	cfg := buildGlobalHTTPConfig(global)

	if cfg.Randomize {
		t.Error("Randomize = true, want false")
	}
	if cfg.UserAgent != "" {
		t.Errorf("UserAgent = %q, want empty", cfg.UserAgent)
	}
	if cfg.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0", cfg.Timeout)
	}
	if cfg.FollowRedirects != nil {
		t.Error("FollowRedirects should be nil")
	}
	if cfg.MaxRedirects != nil {
		t.Error("MaxRedirects should be nil")
	}
	if cfg.MaxResponseSize != nil {
		t.Error("MaxResponseSize should be nil")
	}
}

func TestBuildGlobalHTTPConfig_WithProxyPassword(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		Proxy: &config.ProxyConfig{
			URL:      "http://proxy:8080",
			Username: "user",
			Password: "secret",
		},
	}

	cfg := buildGlobalHTTPConfig(global)

	if cfg.Proxy == nil {
		t.Fatal("Proxy is nil")
	}
	if cfg.Proxy.Password != "secret" {
		t.Errorf("Proxy.Password = %q, want %q", cfg.Proxy.Password, "secret")
	}
}

func TestBuildGlobalHTTPConfig_WithBypass(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		Proxy: &config.ProxyConfig{
			URL:    "http://proxy:8080",
			Bypass: []string{"localhost", "*.internal", "*.local"},
		},
	}

	cfg := buildGlobalHTTPConfig(global)

	if len(cfg.Proxy.Bypass) != 3 {
		t.Fatalf("Bypass length = %d, want 3", len(cfg.Proxy.Bypass))
	}
	if cfg.Proxy.Bypass[0] != "localhost" {
		t.Errorf("Bypass[0] = %q, want %q", cfg.Proxy.Bypass[0], "localhost")
	}
}

func TestBuildGlobalHTTPConfig_EmptyProxy(t *testing.T) {
	t.Parallel()

	global := &config.GlobalHTTPClientConfig{
		Proxy: &config.ProxyConfig{},
	}

	cfg := buildGlobalHTTPConfig(global)

	if cfg.Proxy == nil {
		t.Fatal("Proxy is nil")
	}
	if cfg.Proxy.URL != "" {
		t.Errorf("Proxy.URL = %q, want empty", cfg.Proxy.URL)
	}
}

func TestBuildSpecInfo_WithAuth(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
		Auth: config.Auth{
			Client: &auth.BasicAuthClient{Username: "user", Password: "pass"},
		},
	}

	spec, err := svc.buildSpecInfo(specConfig, false, nil)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.Domain != "test-api" {
		t.Errorf("Domain = %q, want %q", spec.Domain, "test-api")
	}
	if spec.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q, want %q", spec.BaseURL, "https://api.example.com")
	}
	if spec.Auth == nil {
		t.Fatal("Auth is nil")
	}
}

func TestBuildSpecInfo_WithHTTPClient(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
		HTTPClient: &config.HTTPClientConfig{
			Headers: map[string]string{"X-Spec": "value"},
		},
	}

	spec, err := svc.buildSpecInfo(specConfig, false, nil)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.HTTPClient == nil {
		t.Fatal("HTTPClient is nil")
	}
	if spec.HTTPClient.Headers["X-Spec"] != "value" {
		t.Errorf("Headers[X-Spec] = %q, want %q", spec.HTTPClient.Headers["X-Spec"], "value")
	}
}

func TestBuildSpecInfo_NoAuth(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
	}

	spec, err := svc.buildSpecInfo(specConfig, false, nil)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.Auth != nil {
		t.Error("Auth should be nil when not configured")
	}
}

func TestBuildSpecInfo_NoHTTPClient(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
	}

	spec, err := svc.buildSpecInfo(specConfig, false, nil)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.HTTPClient != nil {
		t.Error("HTTPClient should be nil when not configured")
	}
}

func TestBuildSpecInfo_WithCookies(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
		HTTPClient: &config.HTTPClientConfig{
			Cookies: []config.Cookie{
				{Name: "session", Value: "abc"},
			},
		},
	}

	spec, err := svc.buildSpecInfo(specConfig, false, nil)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.HTTPClient == nil {
		t.Fatal("HTTPClient is nil")
	}
	if len(spec.HTTPClient.Cookies) != 1 || spec.HTTPClient.Cookies[0].Name != "session" {
		t.Error("Cookies missing session")
	}
}

func TestBuildSpecInfo_WithMockEnabled(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
		Auth: config.Auth{
			Client: &auth.BasicAuthClient{Username: "user", Password: "pass"},
		},
	}

	spec, err := svc.buildSpecInfo(specConfig, true, nil)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.Auth == nil {
		t.Fatal("Auth is nil")
	}
}

func TestBuildSpecInfo_WithMockAuthConfig(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)

	specConfig := &config.Spec{
		Domain:   "test-api",
		LLMTitle: "Test API",
		BaseURL:  "https://api.example.com",
		Auth: config.Auth{
			Client: &auth.BasicAuthClient{Username: "user", Password: "pass"},
		},
	}

	mockAuth := &config.MockAuthConfig{
		OAuth2Port: 9099,
		DigestPort: 9098,
	}

	spec, err := svc.buildSpecInfo(specConfig, true, mockAuth)
	if err != nil {
		t.Fatalf("buildSpecInfo() = %v", err)
	}

	if spec.Auth == nil {
		t.Fatal("Auth is nil")
	}
}
