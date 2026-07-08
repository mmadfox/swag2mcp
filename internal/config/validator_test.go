package config

import (
	"testing"
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
	if err != nil {
		t.Fatalf("ValidateConfig() = %v, want nil", err)
	}
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
	if err == nil {
		t.Fatal("expected error for duplicate domain")
	}
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
	if err != nil {
		t.Fatalf("ValidateConfig() = %v, want nil", err)
	}
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
	if err != nil {
		t.Fatalf("ValidateConfig() = %v, want nil (disabled specs are skipped)", err)
	}
}

func TestValidationErrors_Empty(t *testing.T) {
	t.Parallel()

	var ve validationErrors
	if ve.Error() != "no validation errors" {
		t.Errorf("Error() = %q, want %q", ve.Error(), "no validation errors")
	}
}

func TestValidationErrors_WithErrors(t *testing.T) {
	t.Parallel()

	ve := validationErrors{
		{field: "specs", message: "no specifications defined"},
	}
	errStr := ve.Error()
	if errStr == "" {
		t.Fatal("Error() returned empty string")
	}
}
