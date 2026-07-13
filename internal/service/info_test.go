package service

import (
	"context"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
)

func TestInfo_WithConfig(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, WithVersion("v1.0.0"))

	cfg := &config.Config{
		MockEnabled: true,
		HTTPClient: &config.GlobalHTTPClientConfig{
			Randomize:       true,
			UserAgent:       "test-agent",
			Timeout:         30 * time.Second,
			FollowRedirects: new(bool),
			MaxRedirects:    new(int),
			MaxResponseSize: new(int),
			Headers:         map[string]string{"X-Custom": "value"},
			Cookies: []config.Cookie{
				{Name: "session", Value: "abc", Domain: ".example.com", Path: "/", Secure: true, HTTPOnly: true},
			},
			Proxy: &config.ProxyConfig{
				URL:      "http://proxy:8080",
				Username: "user",
				Bypass:   []string{"localhost"},
			},
		},
		MCP: &config.MCPConfig{
			Transport: "sse",
			Addr:      ":9090",
			Path:      "/mcp",
			Auth:      &config.MCPAuthConfig{Token: "secret"},
		},
		Specs: []config.Spec{
			{Domain: "active-spec", LLMTitle: "Active Spec", BaseURL: "https://api.example.com", Collections: []config.Collection{{Location: "/dev/null", Title: "c1"}}},
			{Domain: "disabled-spec", LLMTitle: "Disabled Spec", BaseURL: "https://api.example.com", Disable: true, Collections: []config.Collection{{Location: "/dev/null", Title: "c2"}}},
		},
	}
	*cfg.HTTPClient.FollowRedirects = true
	*cfg.HTTPClient.MaxRedirects = 5
	*cfg.HTTPClient.MaxResponseSize = 4096

	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Version != "v1.0.0" {
		t.Errorf("Version = %q, want %q", info.Version, "v1.0.0")
	}
	if info.Specs.Total != 2 {
		t.Errorf("Specs.Total = %d, want 2", info.Specs.Total)
	}
	if info.Specs.Active != 1 {
		t.Errorf("Specs.Active = %d, want 1", info.Specs.Active)
	}
	if info.Specs.Disabled != 1 {
		t.Errorf("Specs.Disabled = %d, want 1", info.Specs.Disabled)
	}
	if info.Specs.Collections != 0 {
		t.Errorf("Specs.Collections = %d, want 0 (no indexed specs)", info.Specs.Collections)
	}
	if info.Specs.Endpoints != 0 {
		t.Errorf("Specs.Endpoints = %d, want 0 (no indexed specs)", info.Specs.Endpoints)
	}
	if !info.HTTPClient.Randomize {
		t.Error("HTTPClient.Randomize = false, want true")
	}
	if info.HTTPClient.UserAgent != "test-agent" {
		t.Errorf("HTTPClient.UserAgent = %q, want %q", info.HTTPClient.UserAgent, "test-agent")
	}
	if info.HTTPClient.Timeout != "30s" {
		t.Errorf("HTTPClient.Timeout = %q, want %q", info.HTTPClient.Timeout, "30s")
	}
	if info.HTTPClient.MaxResponseSize != 4096 {
		t.Errorf("HTTPClient.MaxResponseSize = %d, want 4096", info.HTTPClient.MaxResponseSize)
	}
	if info.HTTPClient.Headers == nil || info.HTTPClient.Headers["X-Custom"] != "value" {
		t.Error("HTTPClient.Headers missing X-Custom")
	}
	if len(info.HTTPClient.Cookies) != 1 || info.HTTPClient.Cookies[0].Name != "session" {
		t.Error("HTTPClient.Cookies missing session cookie")
	}
	if info.HTTPClient.Proxy == nil || info.HTTPClient.Proxy.URL != "http://proxy:8080" {
		t.Error("HTTPClient.Proxy missing or wrong URL")
	}
	if info.MCP.Transport != "sse" {
		t.Errorf("MCP.Transport = %q, want %q", info.MCP.Transport, "sse")
	}
	if info.MCP.Addr != ":9090" {
		t.Errorf("MCP.Addr = %q, want %q", info.MCP.Addr, ":9090")
	}
	if !info.MCP.AuthEnabled {
		t.Error("MCP.AuthEnabled = false, want true")
	}
	if !info.Mock.Enabled {
		t.Error("Mock.Enabled = false, want true")
	}
}

func TestInfo_WithNilConfig(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, WithVersion("v1.0.0"))
	seedTestData(t, svc, t.Name())

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Version != "v1.0.0" {
		t.Errorf("Version = %q, want %q", info.Version, "v1.0.0")
	}
	if info.Specs.Total != 0 {
		t.Errorf("Specs.Total = %d, want 0", info.Specs.Total)
	}
}

func TestInfo_WithNilHTTPClient(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		HTTPClient: nil,
		Specs:      []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.HTTPClient.MaxResponseSize != defaultMaxResponseSize {
		t.Errorf("MaxResponseSize = %d, want %d", info.HTTPClient.MaxResponseSize, defaultMaxResponseSize)
	}
}

func TestInfo_WithNilMCP(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		MCP:   nil,
		Specs: []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.MCP.Transport != "stdio" {
		t.Errorf("MCP.Transport = %q, want %q", info.MCP.Transport, "stdio")
	}
	if info.MCP.AuthEnabled {
		t.Error("MCP.AuthEnabled = true, want false")
	}
}

func TestInfo_WithAuthMethods(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:  "basic-spec",
				BaseURL: "https://api.example.com",
				Auth: config.Auth{
					Client: &auth.BasicAuthClient{Username: "user", Password: "pass"},
				},
				Collections: []config.Collection{{Location: "/dev/null", Title: "c1"}},
			},
			{
				Domain:  "bearer-spec",
				BaseURL: "https://api.example.com",
				Auth: config.Auth{
					Client: &auth.BearerTokenAuthClient{Token: "tok"},
				},
				Collections: []config.Collection{{Location: "/dev/null", Title: "c2"}},
			},
			{
				Domain:  "no-auth-spec",
				BaseURL: "https://api.example.com",
				Auth: config.Auth{
					Client: auth.NewNoAuthClient(),
				},
				Collections: []config.Collection{{Location: "/dev/null", Title: "c3"}},
			},
			{
				Domain:  "disabled-spec",
				BaseURL: "https://api.example.com",
				Disable: true,
				Auth: config.Auth{
					Client: &auth.DigestAuthClient{Username: "u", Password: "p"},
				},
				Collections: []config.Collection{{Location: "/dev/null", Title: "c4"}},
			},
		},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if len(info.Auth.Methods) != 2 {
		t.Fatalf("Auth.Methods = %v, want 2 methods", info.Auth.Methods)
	}

	seen := make(map[string]bool)
	for _, m := range info.Auth.Methods {
		seen[m] = true
	}
	if !seen["basic"] {
		t.Error("missing basic auth method")
	}
	if !seen["bearer"] {
		t.Error("missing bearer auth method")
	}
}

func TestInfo_WithNoAuth(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:  "no-auth",
				BaseURL: "https://api.example.com",
				Auth: config.Auth{
					Client: auth.NewNoAuthClient(),
				},
				Collections: []config.Collection{{Location: "/dev/null", Title: "c1"}},
			},
		},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if len(info.Auth.Methods) != 0 {
		t.Errorf("Auth.Methods = %v, want empty", info.Auth.Methods)
	}
}

func TestInfo_WithEmptySpecs(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		Specs: []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Specs.Total != 0 {
		t.Errorf("Specs.Total = %d, want 0", info.Specs.Total)
	}
	if info.Specs.Active != 0 {
		t.Errorf("Specs.Active = %d, want 0", info.Specs.Active)
	}
	if info.Specs.Disabled != 0 {
		t.Errorf("Specs.Disabled = %d, want 0", info.Specs.Disabled)
	}
}

func TestInfo_WithAllDisabledSpecs(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		Specs: []config.Spec{
			{
				Domain:      "disabled-1",
				BaseURL:     "https://api.example.com",
				Disable:     true,
				Collections: []config.Collection{{Location: "/dev/null", Title: "c1"}},
			},
			{
				Domain:      "disabled-2",
				BaseURL:     "https://api.example.com",
				Disable:     true,
				Collections: []config.Collection{{Location: "/dev/null", Title: "c2"}},
			},
		},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Specs.Total != 2 {
		t.Errorf("Specs.Total = %d, want 2", info.Specs.Total)
	}
	if info.Specs.Active != 0 {
		t.Errorf("Specs.Active = %d, want 0", info.Specs.Active)
	}
	if info.Specs.Disabled != 2 {
		t.Errorf("Specs.Disabled = %d, want 2", info.Specs.Disabled)
	}
}

func TestInfo_WithFullHTTPClientConfig(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	follow := true
	maxRedir := 10
	maxSize := 8192

	cfg := &config.Config{
		HTTPClient: &config.GlobalHTTPClientConfig{
			Randomize:       true,
			UserAgent:       "my-ua",
			Timeout:         5 * time.Second,
			FollowRedirects: &follow,
			MaxRedirects:    &maxRedir,
			MaxResponseSize: &maxSize,
			Headers:         map[string]string{"X-Api-Key": "123"},
			Cookies: []config.Cookie{
				{Name: "test-cookie", Value: "val", Domain: ".test.com", Path: "/api", Secure: true, HTTPOnly: true},
			},
			Proxy: &config.ProxyConfig{
				URL:      "http://proxy:3128",
				Username: "proxy-user",
				Password: "proxy-pass",
				Bypass:   []string{"localhost", "*.internal"},
			},
		},
		Specs: []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if !info.HTTPClient.Randomize {
		t.Error("Randomize = false, want true")
	}
	if info.HTTPClient.UserAgent != "my-ua" {
		t.Errorf("UserAgent = %q, want %q", info.HTTPClient.UserAgent, "my-ua")
	}
	if info.HTTPClient.Timeout != "5s" {
		t.Errorf("Timeout = %q, want %q", info.HTTPClient.Timeout, "5s")
	}
	if info.HTTPClient.FollowRedirects == nil || *info.HTTPClient.FollowRedirects != true {
		t.Error("FollowRedirects = nil or false, want true")
	}
	if info.HTTPClient.MaxRedirects == nil || *info.HTTPClient.MaxRedirects != 10 {
		t.Error("MaxRedirects = nil or wrong value, want 10")
	}
	if info.HTTPClient.MaxResponseSize != 8192 {
		t.Errorf("MaxResponseSize = %d, want 8192", info.HTTPClient.MaxResponseSize)
	}
	if info.HTTPClient.Headers["X-Api-Key"] != "123" {
		t.Errorf("Headers missing X-Api-Key")
	}
	if len(info.HTTPClient.Cookies) != 1 || info.HTTPClient.Cookies[0].Name != "test-cookie" {
		t.Error("Cookies missing test-cookie")
	}
	if info.HTTPClient.Proxy == nil || info.HTTPClient.Proxy.URL != "http://proxy:3128" {
		t.Error("Proxy missing or wrong URL")
	}
	if info.HTTPClient.Proxy.Username != "proxy-user" {
		t.Errorf("Proxy.Username = %q, want %q", info.HTTPClient.Proxy.Username, "proxy-user")
	}
	if len(info.HTTPClient.Proxy.Bypass) != 2 {
		t.Errorf("Proxy.Bypass length = %d, want 2", len(info.HTTPClient.Proxy.Bypass))
	}
}

func TestInfo_WithMCPTransports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mcp       *config.MCPConfig
		wantTrans string
		wantAuth  bool
	}{
		{
			name:      "stdio",
			mcp:       &config.MCPConfig{Transport: "stdio"},
			wantTrans: "stdio",
			wantAuth:  false,
		},
		{
			name:      "sse",
			mcp:       &config.MCPConfig{Transport: "sse", Addr: ":8080", Path: "/mcp"},
			wantTrans: "sse",
			wantAuth:  false,
		},
		{
			name:      "streamable-http",
			mcp:       &config.MCPConfig{Transport: "streamable-http", Addr: ":9090"},
			wantTrans: "streamable-http",
			wantAuth:  false,
		},
		{
			name:      "with-auth",
			mcp:       &config.MCPConfig{Transport: "sse", Auth: &config.MCPAuthConfig{Token: "tok"}},
			wantTrans: "sse",
			wantAuth:  true,
		},
		{
			name:      "default-transport",
			mcp:       &config.MCPConfig{},
			wantTrans: "stdio",
			wantAuth:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTestService(t)
			seedTestData(t, svc, t.Name())

			cfg := &config.Config{
				MCP:   tt.mcp,
				Specs: []config.Spec{},
			}
			svc.config = cfg
			svc.buildSnapshot()

			info, err := svc.Info(context.Background())
			if err != nil {
				t.Fatalf("Info() = %v", err)
			}

			if info.MCP.Transport != tt.wantTrans {
				t.Errorf("Transport = %q, want %q", info.MCP.Transport, tt.wantTrans)
			}
			if info.MCP.AuthEnabled != tt.wantAuth {
				t.Errorf("AuthEnabled = %v, want %v", info.MCP.AuthEnabled, tt.wantAuth)
			}
		})
	}
}

func TestInfo_Uptime(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Uptime == "" {
		t.Error("Uptime is empty")
	}
}

func TestInfo_Version(t *testing.T) {
	t.Parallel()

	svc := newTestService(t, WithVersion("v2.0.0-rc1"))
	seedTestData(t, svc, t.Name())

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Version != "v2.0.0-rc1" {
		t.Errorf("Version = %q, want %q", info.Version, "v2.0.0-rc1")
	}
}

func TestInfo_Workspace(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Workspace == "" {
		t.Error("Workspace is empty")
	}
}

func TestInfo_MockDisabled(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		MockEnabled: false,
		Specs:       []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.Mock.Enabled {
		t.Error("Mock.Enabled = true, want false")
	}
}

func TestInfo_WithProxyNoAuth(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		HTTPClient: &config.GlobalHTTPClientConfig{
			Proxy: &config.ProxyConfig{
				URL:    "http://proxy:8080",
				Bypass: []string{"localhost", "127.0.0.1"},
			},
		},
		Specs: []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if info.HTTPClient.Proxy == nil {
		t.Fatal("Proxy is nil")
	}
	if info.HTTPClient.Proxy.URL != "http://proxy:8080" {
		t.Errorf("Proxy.URL = %q, want %q", info.HTTPClient.Proxy.URL, "http://proxy:8080")
	}
	if info.HTTPClient.Proxy.Username != "" {
		t.Errorf("Proxy.Username = %q, want empty", info.HTTPClient.Proxy.Username)
	}
	if len(info.HTTPClient.Proxy.Bypass) != 2 {
		t.Errorf("Proxy.Bypass length = %d, want 2", len(info.HTTPClient.Proxy.Bypass))
	}
}

func TestInfo_WithCookies(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		HTTPClient: &config.GlobalHTTPClientConfig{
			Cookies: []config.Cookie{
				{Name: "c1", Value: "v1"},
				{Name: "c2", Value: "v2", Domain: ".example.com", Path: "/", Secure: true, HTTPOnly: true},
			},
		},
		Specs: []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if len(info.HTTPClient.Cookies) != 2 {
		t.Fatalf("Cookies length = %d, want 2", len(info.HTTPClient.Cookies))
	}
	if info.HTTPClient.Cookies[0].Name != "c1" {
		t.Errorf("Cookies[0].Name = %q, want %q", info.HTTPClient.Cookies[0].Name, "c1")
	}
	if info.HTTPClient.Cookies[1].Domain != ".example.com" {
		t.Errorf("Cookies[1].Domain = %q, want %q", info.HTTPClient.Cookies[1].Domain, ".example.com")
	}
	if !info.HTTPClient.Cookies[1].Secure {
		t.Error("Cookies[1].Secure = false, want true")
	}
	if !info.HTTPClient.Cookies[1].HTTPOnly {
		t.Error("Cookies[1].HTTPOnly = false, want true")
	}
}

func TestInfo_WithHeaders(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	cfg := &config.Config{
		HTTPClient: &config.GlobalHTTPClientConfig{
			Headers: map[string]string{
				"Authorization": "Bearer test",
				"X-Request-Id":  "12345",
			},
		},
		Specs: []config.Spec{},
	}
	svc.config = cfg
	svc.buildSnapshot()

	info, err := svc.Info(context.Background())
	if err != nil {
		t.Fatalf("Info() = %v", err)
	}

	if len(info.HTTPClient.Headers) != 2 {
		t.Fatalf("Headers length = %d, want 2", len(info.HTTPClient.Headers))
	}
	if info.HTTPClient.Headers["Authorization"] != "Bearer test" {
		t.Errorf("Headers[Authorization] = %q, want %q", info.HTTPClient.Headers["Authorization"], "Bearer test")
	}
	if info.HTTPClient.Headers["X-Request-Id"] != "12345" {
		t.Errorf("Headers[X-Request-Id] = %q, want %q", info.HTTPClient.Headers["X-Request-Id"], "12345")
	}
}
