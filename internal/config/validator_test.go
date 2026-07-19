package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestValidateConfig_Valid(t *testing.T) {
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err)
}

func TestValidateConfig_DuplicateDomain(t *testing.T) {
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
			{
				Domain:   "test-api",
				LLMTitle: "Test API v2",
				BaseURL:  "https://api2.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Second Collection",
						Location: "https://example.com/spec2.yaml",
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for duplicate domain")
}

func TestValidateConfig_WithTags(t *testing.T) {
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
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{Tags: []string{"public"}})
	require.NoError(t, err)
}

func TestValidateConfig_DisabledSpec(t *testing.T) {
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err, "disabled specs are skipped")
}

func TestValidateConfig_DisabledCollection(t *testing.T) {
	t.Parallel()

	cfg := &Config{
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
						LLMTitle: "Active Collection",
						Location: "https://example.com/spec2.yaml",
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err, "disabled collections are skipped")
}

func TestValidateConfig_AuthValidationError(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Auth: Auth{
					Client: &auth.BearerTokenAuthClient{Token: ""},
				},
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for invalid auth config")
}

func TestValidateConfig_WithCache(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	c := cache.New(dir)

	// Create a local spec file that the cache can resolve
	specFile := filepath.Join(dir, "spec.yaml")
	require.NoError(t, os.WriteFile(specFile, []byte("openapi: 3.0.0"), 0600))

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: specFile,
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{Cache: c})
	require.NoError(t, err)
}

func TestValidateConfig_WithCache_LocationError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	c := cache.New(dir)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: filepath.Join(dir, "nonexistent.yaml"),
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{Cache: c})
	require.Error(t, err, "expected error for nonexistent location with cache")
}

func TestValidationErrors_Empty(t *testing.T) {
	t.Parallel()

	var ve validationErrors
	assert.Equal(t, "no validation errors", ve.Error())
}

func TestValidationErrors_WithErrors(t *testing.T) {
	t.Parallel()

	ve := validationErrors{
		{field: "specs", message: "no specifications defined"},
	}
	assert.NotEmpty(t, ve.Error())
}

func TestValidationErrors_WithSpec(t *testing.T) {
	t.Parallel()

	ve := validationErrors{
		{
			field:   "specs[0].domain",
			message: "Domain is required",
			spec:    "test-api",
		},
	}
	assert.NotEmpty(t, ve.Error())
}

func TestValidationErrors_WithSpecAndCollection(t *testing.T) {
	t.Parallel()

	ve := validationErrors{
		{
			field:      "specs[0].collections[0].location",
			message:    "Location is required",
			spec:       "test-api",
			collection: "Main",
		},
	}
	assert.NotEmpty(t, ve.Error())
}

func TestHumanReadableError_Required(t *testing.T) {
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
	require.Error(t, err, "expected error for empty domain")
}

func TestHumanReadableError_Min(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "AB", // too short, min=5
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
	require.Error(t, err, "expected error for too short LLMTitle")
}

func TestHumanReadableError_Max(t *testing.T) {
	t.Parallel()

	longTitle := strings.Repeat("a", 130)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: longTitle, // too long, max=120
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
	require.Error(t, err, "expected error for too long LLMTitle")
}

func TestHumanReadableError_InvalidURL(t *testing.T) {
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
	require.Error(t, err, "expected error for invalid URL")
}

func TestHumanReadableError_LocationMin(t *testing.T) {
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
						Location: "ab", // too short, min=5
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for too short Location")
}

func TestHumanReadableError_LocationMax(t *testing.T) {
	t.Parallel()

	longLoc := strings.Repeat("a", 260)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: longLoc, // too long, max=250
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for too long Location")
}

func TestHumanReadableError_LLMInstructionMax(t *testing.T) {
	t.Parallel()

	longInstr := strings.Repeat("a", 510)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:         "test-api",
				LLMTitle:       "Test API v1",
				LLMInstruction: longInstr, // too long, max=500
				BaseURL:        "https://api.example.com",
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
	require.Error(t, err, "expected error for too long LLMInstruction")
}

func TestHumanReadableError_DomainFormat(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "invalid domain with spaces!",
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
	require.Error(t, err, "expected error for invalid domain format")
}

func TestHumanReadableError_TitleFormat(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test \x00 API", // null byte is invalid
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
	require.Error(t, err, "expected error for invalid title format")
}

func TestHumanReadableError_InstructionFormat(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:         "test-api",
				LLMTitle:       "Test API v1",
				LLMInstruction: "Valid instruction with \x00 null", // null byte is invalid
				BaseURL:        "https://api.example.com",
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
	require.Error(t, err, "expected error for invalid instruction format")
}

func TestHumanReadableError_CollectionURL(t *testing.T) {
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
						BaseURL:  "not-a-url",
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for invalid collection URL")
}

func TestHumanReadableError_CollectionLocationMin(t *testing.T) {
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
						Location: "ab", // too short, min=5
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for too short Location")
}

func TestHumanReadableError_CollectionLocationMax(t *testing.T) {
	t.Parallel()

	longLoc := strings.Repeat("a", 260)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main Collection",
						Location: longLoc, // too long, max=250
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for too long Location")
}

func TestHumanReadableError_CollectionLLMTitleMax(t *testing.T) {
	t.Parallel()

	longTitle := strings.Repeat("a", 130)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle: longTitle, // too long, max=120
						Location: "https://example.com/spec.yaml",
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for too long Collection LLMTitle")
}

func TestHumanReadableError_CollectionLLMInstructionMax(t *testing.T) {
	t.Parallel()

	longInstr := strings.Repeat("a", 370)

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "test-api",
				LLMTitle: "Test API v1",
				BaseURL:  "https://api.example.com",
				Collections: []Collection{
					{
						LLMTitle:       "Main Collection",
						LLMInstruction: longInstr, // too long, max=360
						Location:       "https://example.com/spec.yaml",
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for too long Collection LLMInstruction")
}

func TestConfig_Iterate_EarlyBreak(t *testing.T) {
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
		if count == 1 {
			break
		}
	}
	assert.Equal(t, 1, count)
}

func TestExampleSpecAddYAML(t *testing.T) {
	t.Parallel()

	data := ExampleSpecAddYAML()
	require.NotEmpty(t, data, "ExampleSpecAddYAML() returned empty data")

	// Verify it's valid YAML and contains expected fields
	var raw map[string]any
	require.NoError(t, yaml.Unmarshal(data, &raw))
	require.NotNil(t, raw["domain"], "domain is missing")
	require.NotNil(t, raw["base_url"], "base_url is missing")
	require.NotNil(t, raw["collections"], "collections is missing")
}

func TestExampleCollectionAddYAML(t *testing.T) {
	t.Parallel()

	data := ExampleCollectionAddYAML()
	require.NotEmpty(t, data, "ExampleCollectionAddYAML() returned empty data")

	var raw map[string]any
	require.NoError(t, yaml.Unmarshal(data, &raw))
	require.NotNil(t, raw["spec_domain"], "spec_domain is missing")
	require.NotNil(t, raw["location"], "location is missing")
}

func TestExampleMCPStdioYAML(t *testing.T) {
	t.Parallel()

	data := exampleMCPStdioYAML()
	require.NotEmpty(t, data, "exampleMCPStdioYAML() returned empty data")

	var raw map[string]any
	require.NoError(t, yaml.Unmarshal(data, &raw))
	mcp, ok := raw["mcp"].(map[string]any)
	require.True(t, ok, "mcp section is missing")
	assert.Equal(t, "stdio", mcp["transport"])
}

func TestExampleMCPSSEYAML(t *testing.T) {
	t.Parallel()

	data := exampleMCPSSEYAML()
	require.NotEmpty(t, data, "exampleMCPSSEYAML() returned empty data")

	var raw map[string]any
	require.NoError(t, yaml.Unmarshal(data, &raw))
	mcp, ok := raw["mcp"].(map[string]any)
	require.True(t, ok, "mcp section is missing")
	assert.Equal(t, "sse", mcp["transport"])
	assert.Equal(t, ":8080", mcp["addr"])
	assert.Equal(t, "/mcp", mcp["path"])
	auth, ok := mcp["auth"].(map[string]any)
	require.True(t, ok, "auth section is missing")
	assert.Equal(t, "your-secret-token", auth["token"])
}

func TestExampleMCPStreamableHTTPYAML(t *testing.T) {
	t.Parallel()

	data := exampleMCPStreamableHTTPYAML()
	require.NotEmpty(t, data, "exampleMCPStreamableHTTPYAML() returned empty data")

	var raw map[string]any
	require.NoError(t, yaml.Unmarshal(data, &raw))
	mcp, ok := raw["mcp"].(map[string]any)
	require.True(t, ok, "mcp section is missing")
	assert.Equal(t, "streamable-http", mcp["transport"])
	assert.Equal(t, ":9090", mcp["addr"])
	assert.Equal(t, "/api/mcp", mcp["path"])
	auth, ok := mcp["auth"].(map[string]any)
	require.True(t, ok, "auth section is missing")
	assert.Equal(t, "your-secret-token", auth["token"])
}

func TestValidateConfig_HTTPClient_Valid(t *testing.T) {
	t.Parallel()

	timeout := 15 * time.Second
	maxRedir := 5
	maxSize := 4096
	follow := true

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			UserAgent:       "test-agent/1.0",
			Timeout:         timeout,
			FollowRedirects: &follow,
			MaxRedirects:    &maxRedir,
			MaxResponseSize: &maxSize,
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err)
}

func TestValidateConfig_HTTPClient_Nil(t *testing.T) {
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err)
}

func TestValidateConfig_HTTPClient_InvalidProxyURL(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			Proxy: &ProxyConfig{
				URL: "ftp://proxy.example.com:8080",
			},
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for invalid proxy URL scheme")
}

func TestValidateConfig_HTTPClient_ValidProxyURL(t *testing.T) {
	t.Parallel()

	urls := []string{
		"http://proxy.example.com:8080",
		"https://proxy.example.com:8443",
		"socks5://127.0.0.1:1080",
		"socks5h://127.0.0.1:1080",
	}

	for _, u := range urls {
		t.Run(u, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{
				HTTPClient: &GlobalHTTPClientConfig{
					Proxy: &ProxyConfig{URL: u},
				},
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

			err := ValidateConfig(cfg, ValidateOptions{})
			require.NoError(t, err, "ValidateConfig() with proxy %q", u)
		})
	}
}

func TestValidateConfig_HTTPClient_MaxResponseSizeTooSmall(t *testing.T) {
	t.Parallel()

	val := 100
	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			MaxResponseSize: &val,
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for MaxResponseSize < 256")
}

func TestValidateConfig_HTTPClient_MaxResponseSizeTooLarge(t *testing.T) {
	t.Parallel()

	val := 20_000_000
	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			MaxResponseSize: &val,
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for MaxResponseSize > 10MB")
}

func TestValidateConfig_HTTPClient_MaxRedirectsTooLarge(t *testing.T) {
	t.Parallel()

	val := 100
	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			MaxRedirects: &val,
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for MaxRedirects > 50")
}

func TestValidateConfig_HTTPClient_TimeoutTooSmall(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			Timeout: 500 * time.Millisecond,
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for Timeout < 1s")
}

func TestValidateConfig_HTTPClient_TimeoutTooLarge(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{
			Timeout: 600 * time.Second,
		},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for Timeout > 300s")
}

func TestValidateConfig_HTTPClient_DefaultsApplied(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		HTTPClient: &GlobalHTTPClientConfig{},
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

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err)

	assert.Equal(t, "swag2mcp-global/1.0", cfg.HTTPClient.UserAgent)
	assert.Equal(t, 30*time.Second, cfg.HTTPClient.Timeout)
	require.NotNil(t, cfg.HTTPClient.MaxRedirects)
	assert.Equal(t, 10, *cfg.HTTPClient.MaxRedirects)
	require.NotNil(t, cfg.HTTPClient.MaxResponseSize)
	assert.Equal(t, 1048576, *cfg.HTTPClient.MaxResponseSize)
}

func TestConfig_Validate_UppercaseDomain(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Specs: []Spec{
			{
				Domain:   "PETSTORE",
				LLMTitle: "Petstore API",
				BaseURL:  "https://petstore.example.com",
				Collections: []Collection{
					{
						LLMTitle: "Main",
						Location: "https://petstore.example.com/openapi.yaml",
					},
				},
			},
		},
	}

	err := cfg.Validate(nil)
	require.Error(t, err, "expected error for uppercase domain")
	require.Contains(t, err.Error(), "lowercase")
}
