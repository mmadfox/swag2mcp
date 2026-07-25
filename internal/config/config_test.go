package config

import (
	"testing"
)

func TestConfig_Validate_NoSpecs(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	err := cfg.Validate(nil)
	if err == nil {
		t.Fatal("expected error for empty config")
	}
}

func TestConfig_Validate_ValidSpec(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestConfig_Validate_DisabledSpec(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "disabled-api",
				LLMTitle: "Disabled API v1",
				BaseURL:  "https://api.example.com",
				Disable:  true,
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err != nil {
		t.Fatalf("Validate() = %v, want nil (disabled specs are skipped)", err)
	}
}

func TestConfig_Validate_InvalidDomain(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestConfig_Validate_InvalidBaseURL(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "not-a-url",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestConfig_Validate_WithFilter(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "public-api",
				LLMTitle: "Public API v1",
				BaseURL:  "https://api.example.com",
				Tags:     []string{"public"},
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
			{
				Domain:   "internal-api",
				LLMTitle: "Internal API v1",
				BaseURL:  "https://internal.example.com",
				Tags:     []string{"internal"},
				Collections: []Collection{
					{
						LLMTitle: "Internal Collection",
						Location: "https://internal.example.com/spec.yaml",
					},
				},
			},
		},
	}
	filter := NewFilter([]string{"public"})
	err := cfg.Validate(filter)
	if err != nil {
		t.Fatalf("Validate() with filter = %v, want nil", err)
	}
}

func TestConfig_Iterate_All(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "api-1",
				LLMTitle: "API One",
				BaseURL:  "https://api1.example.com",
				Collections: []Collection{
					{LLMTitle: "Coll", Location: "https://example.com/spec.yaml"},
				},
			},
			{
				Domain:   "api-2",
				LLMTitle: "API Two",
				BaseURL:  "https://api2.example.com",
				Collections: []Collection{
					{LLMTitle: "Coll", Location: "https://example.com/spec.yaml"},
				},
			},
		},
	}
	count := 0
	for range cfg.Iterate(nil) {
		count++
	}
	if count != 2 {
		t.Errorf("Iterate count = %d, want 2", count)
	}
}

func TestConfig_Iterate_SkipDisabled(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "api-1",
				LLMTitle: "API One",
				BaseURL:  "https://api1.example.com",
				Disable:  true,
				Collections: []Collection{
					{LLMTitle: "Coll", Location: "https://example.com/spec.yaml"},
				},
			},
			{
				Domain:   "api-2",
				LLMTitle: "API Two",
				BaseURL:  "https://api2.example.com",
				Collections: []Collection{
					{LLMTitle: "Coll", Location: "https://example.com/spec.yaml"},
				},
			},
		},
	}
	count := 0
	for range cfg.Iterate(nil) {
		count++
	}
	if count != 1 {
		t.Errorf("Iterate count = %d, want 1", count)
	}
}

func TestConfig_Iterate_WithFilter(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "public-api",
				LLMTitle: "Public API",
				BaseURL:  "https://api.example.com",
				Tags:     []string{"public"},
				Collections: []Collection{
					{LLMTitle: "Coll", Location: "https://example.com/spec.yaml"},
				},
			},
			{
				Domain:   "internal-api",
				LLMTitle: "Internal API",
				BaseURL:  "https://internal.example.com",
				Tags:     []string{"internal"},
				Collections: []Collection{
					{LLMTitle: "Coll", Location: "https://internal.example.com/spec.yaml"},
				},
			},
		},
	}
	filter := NewFilter([]string{"public"})
	count := 0
	for range cfg.Iterate(filter) {
		count++
	}
	if count != 1 {
		t.Errorf("Iterate count = %d, want 1", count)
	}
}

func TestConfig_Iterate_Empty(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	count := 0
	for range cfg.Iterate(nil) {
		count++
	}
	if count != 0 {
		t.Errorf("Iterate count = %d, want 0", count)
	}
}

func TestHTTPClientConfig_MaxResponseSize(t *testing.T) {
	t.Parallel()

	val := 4096
	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			MaxResponseSize: &val,
		},
	}

	if cfg.HTTPClient.MaxResponseSize == nil {
		t.Fatal("MaxResponseSize is nil")
	}
	if *cfg.HTTPClient.MaxResponseSize != 4096 {
		t.Errorf("MaxResponseSize = %d, want %d", *cfg.HTTPClient.MaxResponseSize, 4096)
	}
}

func TestHTTPClientConfig_MaxResponseSize_Nil(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{},
	}
	if cfg.HTTPClient.MaxResponseSize != nil {
		t.Error("MaxResponseSize should be nil by default")
	}
}

func TestMCPAuthConfig_Resolve_Nil(t *testing.T) {
	t.Parallel()

	var c *MCPAuthConfig
	c.Resolve() // should not panic
}

func TestMCPAuthConfig_Resolve_NoEnv(t *testing.T) {
	t.Parallel()

	c := &MCPAuthConfig{Token: "static-token"}
	c.Resolve()
	if c.Token != "static-token" {
		t.Errorf("Token = %q, want %q", c.Token, "static-token")
	}
}

func TestMCPAuthConfig_Resolve_WithEnv(t *testing.T) {
	t.Setenv("MCP_TOKEN", "resolved-token")
	c := &MCPAuthConfig{Token: "$(MCP_TOKEN)"}
	c.Resolve()
	if c.Token != "resolved-token" {
		t.Errorf("Token = %q, want %q", c.Token, "resolved-token")
	}
}

func TestMCPConfig_Defaults(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MCP: &MCPConfig{
			Transport: "sse",
			Addr:      ":9090",
			Path:      "/api/mcp",
		},
	}
	if cfg.MCP.Transport != "sse" {
		t.Errorf("Transport = %q, want %q", cfg.MCP.Transport, "sse")
	}
	if cfg.MCP.Addr != ":9090" {
		t.Errorf("Addr = %q, want %q", cfg.MCP.Addr, ":9090")
	}
	if cfg.MCP.Path != "/api/mcp" {
		t.Errorf("Path = %q, want %q", cfg.MCP.Path, "/api/mcp")
	}
}

func TestConfig_Validate_MockEnabled_RequiresBaseMockURL(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err == nil {
		t.Fatal("expected error when mock_enabled is true but collection BaseMockURL is empty")
	}
}

func TestConfig_Validate_MockEnabled_Valid(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Main Collection",
						Location:    "https://example.com/spec.yaml",
						BaseMockURL: "localhost:8080",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestConfig_Validate_MockEnabled_CollectionRequiresBaseMockURL(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err == nil {
		t.Fatal("expected error when mock_enabled is true but collection BaseMockURL is empty")
	}
}

func TestConfig_Validate_BaseMockURL_InvalidFormat(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Main Collection",
						Location:    "https://example.com/spec.yaml",
						BaseMockURL: "invalid-addr",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err == nil {
		t.Fatal("expected error for invalid BaseMockURL format")
	}
}

func TestConfig_Validate_BaseMockURL_ValidFormats(t *testing.T) {
	t.Parallel()

	formats := []string{
		"localhost:8080",
		"127.0.0.1:8080",
		"0.0.0.0:8080",
		"localhost:80",
		"127.0.0.1:65535",
		"localhost:8080/api/v1",
		"127.0.0.1:9000/v1/smev",
		"0.0.0.0:3000/path/to/service",
	}

	for _, addr := range formats {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()
			cfg := &Config{
				Specs: []Spec{
					{
						Domain:   "test-api",
						LLMTitle: "Test API v1",
						BaseURL:  "https://api.example.com",
						Collections: []Collection{
							{
								LLMTitle:    "Main Collection",
								Location:    "https://example.com/spec.yaml",
								BaseMockURL: addr,
							},
						},
					},
				},
			}
			err := cfg.Validate(nil)
			if err != nil {
				t.Errorf("Validate() with addr %q = %v, want nil", addr, err)
			}
		})
	}
}

func TestConfig_Validate_BaseMockURL_InvalidFormats(t *testing.T) {
	t.Parallel()

	formats := []string{
		"example.com:8080",
		"192.168.1.1:8080",
		"localhost",
		":8080",
		"localhost:",
		"localhost:abc",
		"example.com:8080/api/v1",
		"192.168.1.1:9000/v1/smev",
	}

	for _, addr := range formats {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()
			cfg := &Config{
				Specs: []Spec{
					{
						Domain:   "test-api",
						LLMTitle: "Test API v1",
						BaseURL:  "https://api.example.com",
						Collections: []Collection{
							{
								LLMTitle:    "Main Collection",
								Location:    "https://example.com/spec.yaml",
								BaseMockURL: addr,
							},
						},
					},
				},
			}
			err := cfg.Validate(nil)
			if err == nil {
				t.Errorf("expected error for invalid addr %q", addr)
			}
		})
	}
}

func TestConfig_Validate_MockEnabled_DisabledSpecSkipped(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "disabled-api",
				LLMTitle: "Disabled API v1",
				BaseURL:  "https://api.example.com",
				Disable:  true,
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err != nil {
		t.Fatalf("Validate() = %v, want nil (disabled specs are skipped)", err)
	}
}

func TestConfig_Validate_MockEnabled_DisabledCollectionSkipped(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Disabled Collection",
						Location: "https://example.com/spec.yaml",
						Disable:  true,
					},
					{
						LLMTitle:    "Active Collection",
						Location:    "https://example.com/spec2.yaml",
						BaseMockURL: "localhost:8080",
					},
				},
			},
		},
	}
	err := cfg.Validate(nil)
	if err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}
