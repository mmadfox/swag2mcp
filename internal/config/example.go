package config

import (
	"github.com/mmadfox/swag2mcp/internal/auth"
	"go.yaml.in/yaml/v3"
)

// SpecAddRequest wraps Spec for the add spec command.
// It uses the same Spec type so validation and YAML tags are shared.
type SpecAddRequest struct {
	Spec `yaml:",inline"`
}

// CollectionAddRequest wraps Collection with a spec domain reference.
type CollectionAddRequest struct {
	Collection `yaml:",inline"`
	SpecDomain string `yaml:"spec_domain" validate:"required"`
}

// ExampleSpecAddYAML returns a pretty-printed YAML example for add spec.
func ExampleSpecAddYAML() []byte {
	example := SpecAddRequest{
		Spec: Spec{
			Domain:         "meteo",
			LLMTitle:       "Open-Meteo API",
			LLMInstruction: "Use this API to manage pets, orders, and users.",
			BaseURL:        "https://meteo.swagger.io/v2",
			Tags:           []string{"public", "demo"},
			Auth: Auth{
				Client: &auth.BearerTokenAuthClient{Token: "your-token-here"},
			},
			Collections: []Collection{
				{
					LLMTitle: "Open-Meteo Swagger",
					Location: "https://meteo.swagger.io/v2/swagger.json",
				},
			},
		},
	}
	data, _ := yaml.Marshal(example)
	return data
}

// ExampleCollectionAddYAML returns a pretty-printed YAML example for add collection.
func ExampleCollectionAddYAML() []byte {
	example := CollectionAddRequest{
		SpecDomain: "meteo",
		Collection: Collection{
			LLMTitle: "Orders Collection",
			Location: "https://meteo.example.com/orders.json",
		},
	}
	data, _ := yaml.Marshal(example)
	return data
}

// ExampleMCPStdioYAML returns a YAML example for MCP with stdio transport.
func exampleMCPStdioYAML() []byte {
	data, _ := yaml.Marshal(map[string]any{
		"mcp": map[string]any{
			"transport": "stdio",
		},
	})
	return data
}

// ExampleMCPSSEYAML returns a YAML example for MCP with SSE transport and auth.
func exampleMCPSSEYAML() []byte {
	data, _ := yaml.Marshal(map[string]any{
		"mcp": map[string]any{
			"transport": "sse",
			"addr":      ":8080",
			"path":      "/mcp",
			"auth": map[string]any{
				"token": "your-secret-token",
			},
		},
	})
	return data
}

// ExampleMCPStreamableHTTPYAML returns a YAML example for MCP with streamable HTTP transport and auth.
func exampleMCPStreamableHTTPYAML() []byte {
	data, _ := yaml.Marshal(map[string]any{
		"mcp": map[string]any{
			"transport": "streamable-http",
			"addr":      ":9090",
			"path":      "/api/mcp",
			"auth": map[string]any{
				"token": "your-secret-token",
			},
		},
	})
	return data
}

// ExampleMockConfigYAML returns a YAML example for mock server configuration.
func exampleMockConfigYAML() []byte {
	data, _ := yaml.Marshal(map[string]any{
		"mock_enabled": true,
		"specs": []map[string]any{
			{
				"domain":    "meteo",
				"llm_title": "Open-Meteo API",
				"base_url":  "https://meteo.swagger.io/v2",
				"collections": []map[string]any{
					{
						"llm_title":     "Open-Meteo Swagger",
						"location":      "specs/meteo.json",
						"base_mock_url": "localhost:8080",
					},
				},
			},
		},
	})
	return data
}
