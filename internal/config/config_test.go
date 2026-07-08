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
		HTTPClient: &HTTPClientConfig{
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
		HTTPClient: &HTTPClientConfig{},
	}
	if cfg.HTTPClient.MaxResponseSize != nil {
		t.Error("MaxResponseSize should be nil by default")
	}
}
