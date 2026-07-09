package config

import (
	"github.com/mmadfox/swag2mcp/internal/auth"
	"gopkg.in/yaml.v3"
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
			Domain:         "petstore",
			LLMTitle:       "Petstore API",
			LLMInstruction: "Use this API to manage pets, orders, and users.",
			BaseURL:        "https://petstore.swagger.io/v2",
			Tags:           []string{"public", "demo"},
			Auth: Auth{
				Client: &auth.BearerTokenAuthClient{Token: "your-token-here"},
			},
			Collections: []Collection{
				{
					LLMTitle: "Petstore Swagger",
					Location: "https://petstore.swagger.io/v2/swagger.json",
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
		SpecDomain: "petstore",
		Collection: Collection{
			LLMTitle: "Orders Collection",
			Location: "https://petstore.example.com/orders.json",
		},
	}
	data, _ := yaml.Marshal(example)
	return data
}

// ExampleMCPStdioYAML returns a YAML example for MCP with stdio transport.
func ExampleMCPStdioYAML() []byte {
	data, _ := yaml.Marshal(map[string]any{
		"mcp": map[string]any{
			"transport": "stdio",
		},
	})
	return data
}

// ExampleMCPSSEYAML returns a YAML example for MCP with SSE transport and auth.
func ExampleMCPSSEYAML() []byte {
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
func ExampleMCPStreamableHTTPYAML() []byte {
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
