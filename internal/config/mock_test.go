package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

func TestValidateConfig_MockEnabled_DuplicatePorts(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "api-one",
				LLMTitle: "API One",
				BaseURL:  "https://api1.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Collection One",
						Location:    "https://example.com/spec.yaml",
						BaseMockURL: "localhost:8080",
					},
				},
			},
			{
				Domain:   "api-two",
				LLMTitle: "API Two",
				BaseURL:  "https://api2.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Collection Two",
						Location:    "https://example.com/spec2.yaml",
						BaseMockURL: "localhost:8080",
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for duplicate mock ports")
}

func TestValidateConfig_MockEnabled_DuplicatePortsCollection(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "api-one",
				LLMTitle: "API One",
				BaseURL:  "https://api1.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Collection One",
						Location:    "https://example.com/spec.yaml",
						BaseMockURL: "localhost:8080",
					},
					{
						LLMTitle:    "Collection Two",
						Location:    "https://example.com/spec2.yaml",
						BaseMockURL: "localhost:8080",
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{})
	require.Error(t, err, "expected error for duplicate mock ports in collections")
}

func TestValidateConfig_MockEnabled_UniquePorts(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		MockEnabled: true,
		Specs: []Spec{
			{
				Domain:   "api-one",
				LLMTitle: "API One",
				BaseURL:  "https://api1.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Collection One",
						Location:    "https://example.com/spec.yaml",
						BaseMockURL: "localhost:8080",
					},
				},
			},
			{
				Domain:   "api-two",
				LLMTitle: "API Two",
				BaseURL:  "https://api2.example.com",
				Collections: []Collection{
					{
						LLMTitle:    "Collection Two",
						Location:    "https://example.com/spec2.yaml",
						BaseMockURL: "localhost:8081",
					},
				},
			},
		},
	}

	err := ValidateConfig(cfg, ValidateOptions{})
	require.NoError(t, err)
}

func TestExampleMockConfigYAML(t *testing.T) {
	t.Parallel()

	data := exampleMockConfigYAML()
	require.NotEmpty(t, data, "exampleMockConfigYAML() returned empty data")

	var raw map[string]any
	require.NoError(t, yaml.Unmarshal(data, &raw))
	assert.Equal(t, true, raw["mock_enabled"])

	specs, ok := raw["specs"].([]any)
	require.True(t, ok, "specs section is missing")
	require.NotEmpty(t, specs, "specs is empty")

	spec, ok := specs[0].(map[string]any)
	require.True(t, ok, "first spec is not a map")

	collections, ok := spec["collections"].([]any)
	require.True(t, ok, "collections section is missing in spec")
	require.NotEmpty(t, collections, "collections is empty")

	collection, ok := collections[0].(map[string]any)
	require.True(t, ok, "first collection is not a map")
	require.NotNil(t, collection["base_mock_url"], "base_mock_url is missing in collection")
}

func TestExtractPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		addr string
		want int
	}{
		{"localhost:8080", 8080},
		{"127.0.0.1:80", 80},
		{"0.0.0.0:65535", 65535},
		{"localhost:", 0},
		{"invalid", 0},
		{"", 0},
		{"127.0.0.1:9000/v1/smev", 9000},
		{"localhost:8080/api/v1", 8080},
		{"0.0.0.0:3000/path/to/service", 3000},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			t.Parallel()
			got := extractPort(tt.addr)
			assert.Equal(t, tt.want, got, "extractPort(%q)", tt.addr)
		})
	}
}
