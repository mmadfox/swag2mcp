package config

import (
	"testing"

	"gopkg.in/yaml.v3"
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
	if err == nil {
		t.Fatal("expected error for duplicate mock ports")
	}
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
	if err == nil {
		t.Fatal("expected error for duplicate mock ports in collections")
	}
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
	if err != nil {
		t.Fatalf("ValidateConfig() = %v, want nil", err)
	}
}

func TestExampleMockConfigYAML(t *testing.T) {
	t.Parallel()

	data := ExampleMockConfigYAML()
	if len(data) == 0 {
		t.Fatal("ExampleMockConfigYAML() returned empty data")
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal example YAML: %v", err)
	}
	if raw["mock_enabled"] != true {
		t.Error("mock_enabled should be true")
	}
	specs, ok := raw["specs"].([]any)
	if !ok {
		t.Fatal("specs section is missing")
	}
	if len(specs) == 0 {
		t.Fatal("specs is empty")
	}
	spec, ok := specs[0].(map[string]any)
	if !ok {
		t.Fatal("first spec is not a map")
	}
	collections, ok := spec["collections"].([]any)
	if !ok {
		t.Fatal("collections section is missing in spec")
	}
	if len(collections) == 0 {
		t.Fatal("collections is empty")
	}
	collection, ok := collections[0].(map[string]any)
	if !ok {
		t.Fatal("first collection is not a map")
	}
	if collection["base_mock_url"] == nil {
		t.Error("base_mock_url is missing in collection")
	}
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
			if got != tt.want {
				t.Errorf("extractPort(%q) = %d, want %d", tt.addr, got, tt.want)
			}
		})
	}
}
