package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate_NoSpecs(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for empty config")
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
	require.NoError(t, err)
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
	require.NoError(t, err, "disabled specs are skipped")
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
	require.Error(t, err, "expected validation error")
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
	require.Error(t, err, "expected validation error")
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
	require.NoError(t, err)
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
	assert.Equal(t, 2, count)
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
	assert.Equal(t, 1, count)
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
	assert.Equal(t, 1, count)
}

func TestConfig_Iterate_Empty(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	count := 0
	for range cfg.Iterate(nil) {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestHTTPClientConfig_MaxResponseSize(t *testing.T) {
	t.Parallel()

	val := 4096
	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			MaxResponseSize: &val,
		},
	}

	require.NotNil(t, cfg.HTTPClient.MaxResponseSize)
	assert.Equal(t, 4096, *cfg.HTTPClient.MaxResponseSize)
}

func TestHTTPClientConfig_MaxResponseSize_Nil(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{},
	}
	assert.Nil(t, cfg.HTTPClient.MaxResponseSize, "MaxResponseSize should be nil by default")
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
	assert.Equal(t, "static-token", c.Token)
}

func TestMCPAuthConfig_Resolve_WithEnv(t *testing.T) {
	t.Setenv("MCP_TOKEN", "resolved-token")
	c := &MCPAuthConfig{Token: "$(MCP_TOKEN)"}
	c.Resolve()
	assert.Equal(t, "resolved-token", c.Token)
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
	assert.Equal(t, "sse", cfg.MCP.Transport)
	assert.Equal(t, ":9090", cfg.MCP.Addr)
	assert.Equal(t, "/api/mcp", cfg.MCP.Path)
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
	require.Error(t, err, "expected error when mock_enabled is true but collection BaseMockURL is empty")
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
	require.NoError(t, err)
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
	require.Error(t, err, "expected error when mock_enabled is true but collection BaseMockURL is empty")
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
	require.Error(t, err, "expected error for invalid BaseMockURL format")
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
			require.NoError(t, err, "Validate() with addr %q", addr)
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
			require.Error(t, err, "expected error for invalid addr %q", addr)
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
	require.NoError(t, err, "disabled specs are skipped")
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
	require.NoError(t, err)
}

func TestGlobalHTTPClientConfig_SetDefaults(t *testing.T) {
	t.Parallel()

	cfg := &GlobalHTTPClientConfig{}
	cfg.SetDefaults()

	assert.Equal(t, "swag2mcp-global/1.0", cfg.UserAgent)
	assert.Equal(t, 30*time.Second, cfg.Timeout)
	require.NotNil(t, cfg.FollowRedirects)
	assert.True(t, *cfg.FollowRedirects)
	require.NotNil(t, cfg.MaxRedirects)
	assert.Equal(t, 10, *cfg.MaxRedirects)
	require.NotNil(t, cfg.MaxResponseSize)
	assert.Equal(t, 2048, *cfg.MaxResponseSize)
}

func TestGlobalHTTPClientConfig_SetDefaults_Nil(t *testing.T) {
	t.Parallel()

	var cfg *GlobalHTTPClientConfig
	cfg.SetDefaults() // should not panic
}

func TestGlobalHTTPClientConfig_SetDefaults_DoesNotOverwrite(t *testing.T) {
	t.Parallel()

	timeout := 10 * time.Second
	follow := false
	maxRedir := 5
	maxSize := 4096

	cfg := &GlobalHTTPClientConfig{
		UserAgent:       "custom-agent/1.0",
		Timeout:         timeout,
		FollowRedirects: &follow,
		MaxRedirects:    &maxRedir,
		MaxResponseSize: &maxSize,
	}
	cfg.SetDefaults()

	assert.Equal(t, "custom-agent/1.0", cfg.UserAgent)
	assert.Equal(t, timeout, cfg.Timeout)
	assert.Equal(t, follow, *cfg.FollowRedirects)
	assert.Equal(t, maxRedir, *cfg.MaxRedirects)
	assert.Equal(t, maxSize, *cfg.MaxResponseSize)
}

func TestGlobalHTTPClientConfig_SetDefaults_WithRandomize(t *testing.T) {
	t.Parallel()

	cfg := &GlobalHTTPClientConfig{
		Randomize: true,
	}
	cfg.SetDefaults()

	assert.Empty(t, cfg.UserAgent, "UserAgent should be empty when Randomize is true")
}
