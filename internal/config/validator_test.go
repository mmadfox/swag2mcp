package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"gopkg.in/yaml.v3"
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
	if err != nil {
		t.Fatalf("ValidateConfig() = %v, want nil (disabled collections are skipped)", err)
	}
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
	if err == nil {
		t.Fatal("expected error for invalid auth config")
	}
}

func TestValidateConfig_WithCache(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	c := cache.New(dir)

	// Create a local spec file that the cache can resolve
	specFile := filepath.Join(dir, "spec.yaml")
	if err := os.WriteFile(specFile, []byte("openapi: 3.0.0"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

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
	if err != nil {
		t.Fatalf("ValidateConfig() with cache = %v, want nil", err)
	}
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
	if err == nil {
		t.Fatal("expected error for nonexistent location with cache")
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

func TestValidationErrors_WithSpec(t *testing.T) {
	t.Parallel()

	ve := validationErrors{
		{
			field:   "specs[0].domain",
			message: "Domain is required",
			spec:    "test-api",
		},
	}
	errStr := ve.Error()
	if errStr == "" {
		t.Fatal("Error() returned empty string")
	}
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
	errStr := ve.Error()
	if errStr == "" {
		t.Fatal("Error() returned empty string")
	}
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
	if err == nil {
		t.Fatal("expected error for empty domain")
	}
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
	if err == nil {
		t.Fatal("expected error for too short LLMTitle")
	}
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
	if err == nil {
		t.Fatal("expected error for too long LLMTitle")
	}
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
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
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
	if err == nil {
		t.Fatal("expected error for too short Location")
	}
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
	if err == nil {
		t.Fatal("expected error for too long Location")
	}
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
	if err == nil {
		t.Fatal("expected error for too long LLMInstruction")
	}
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
	if err == nil {
		t.Fatal("expected error for invalid domain format")
	}
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
	if err == nil {
		t.Fatal("expected error for invalid title format")
	}
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
	if err == nil {
		t.Fatal("expected error for invalid instruction format")
	}
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
	if err == nil {
		t.Fatal("expected error for invalid collection URL")
	}
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
	if err == nil {
		t.Fatal("expected error for too short Location")
	}
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
	if err == nil {
		t.Fatal("expected error for too long Location")
	}
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
	if err == nil {
		t.Fatal("expected error for too long Collection LLMTitle")
	}
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
	if err == nil {
		t.Fatal("expected error for too long Collection LLMInstruction")
	}
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
	if count != 1 {
		t.Errorf("Iterate early break count = %d, want 1", count)
	}
}

func TestExampleSpecAddYAML(t *testing.T) {
	t.Parallel()

	data := ExampleSpecAddYAML()
	if len(data) == 0 {
		t.Fatal("ExampleSpecAddYAML() returned empty data")
	}

	// Verify it's valid YAML and contains expected fields
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal example YAML: %v", err)
	}
	if raw["domain"] == nil {
		t.Error("domain is missing")
	}
	if raw["base_url"] == nil {
		t.Error("base_url is missing")
	}
	if raw["collections"] == nil {
		t.Error("collections is missing")
	}
}

func TestExampleCollectionAddYAML(t *testing.T) {
	t.Parallel()

	data := ExampleCollectionAddYAML()
	if len(data) == 0 {
		t.Fatal("ExampleCollectionAddYAML() returned empty data")
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal example YAML: %v", err)
	}
	if raw["spec_domain"] == nil {
		t.Error("spec_domain is missing")
	}
	if raw["location"] == nil {
		t.Error("location is missing")
	}
}
