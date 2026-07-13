package spec

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/spec"
)

func TestParse_swaggerHost(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Servers) == 0 {
		t.Fatal("expected at least 1 server from swagger host")
	}

	if doc.Servers[0].URL != "https://api.example.com/v1" {
		t.Errorf("got server URL %q, want %q", doc.Servers[0].URL, "https://api.example.com/v1")
	}
}

func TestParse_operationMetadata(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodGet {
			found = true
			op := pi.Operation
			if op.Summary != "Список пользователей" {
				t.Errorf("got summary %q, want %q", op.Summary, "List users")
			}
			if len(op.Parameters) == 0 {
				t.Fatal("expected parameters")
			}
			if op.Parameters[0].Name != "limit" {
				t.Errorf("got param name %q, want %q", op.Parameters[0].Name, "limit")
			}
			if op.Parameters[0].In != "query" {
				t.Errorf("got param in %q, want %q", op.Parameters[0].In, "query")
			}
			break
		}
	}
	if !found {
		t.Fatal("GET /users not found in parsed doc")
	}
}

func TestParse_requestBody(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodPost {
			found = true
			op := pi.Operation
			if op.RequestBody == nil {
				t.Fatal("expected request body")
			}
			if !op.RequestBody.Required {
				t.Error("expected required request body")
			}
			if op.RequestBody.Content == nil {
				t.Fatal("expected request body content")
			}
			break
		}
	}
	if !found {
		t.Fatal("POST /users not found")
	}
}

func TestParse_responses(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodGet {
			found = true
			op := pi.Operation
			if len(op.Responses) == 0 {
				t.Fatal("expected responses")
			}
			resp, ok := op.Responses["200"]
			if !ok {
				t.Fatal("expected 200 response")
			}
			if resp.Description != "OK" {
				t.Errorf("got description %q, want %q", resp.Description, "OK")
			}
			break
		}
	}
	if !found {
		t.Fatal("GET /users not found")
	}
}

func TestParse_swaggerFileUpload(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/files/upload" && pi.Method == http.MethodPost {
			found = true
			op := pi.Operation
			if len(op.Parameters) == 0 {
				t.Fatal("expected parameters")
			}
			if op.Parameters[0].Name != "file" {
				t.Errorf("got param name %q, want %q", op.Parameters[0].Name, "file")
			}
			if op.Parameters[0].In != "formData" {
				t.Errorf("got param in %q, want %q", op.Parameters[0].In, "formData")
			}
			break
		}
	}
	if !found {
		t.Fatal("POST /files/upload not found")
	}
}

func TestParseV2_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseV2([]byte("{invalid}"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseV2_EmptyHost(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{"swagger":"2.0","info":{"title":"Test","version":"1.0"},"paths":{}}`)
	result, err := parseV2(jsonData)
	if err != nil {
		t.Fatalf("parseV2() = %v", err)
	}
	if len(result.Servers) != 0 {
		t.Errorf("Servers = %d, want 0", len(result.Servers))
	}
}

func TestSwaggerSchemaToSchema_Nil(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(nil)
	if s != nil {
		t.Fatal("expected nil")
	}
}

func TestSwaggerSchemaToSchema_Ref(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/Pet"),
		},
	})
	if s == nil {
		t.Fatal("schema is nil")
	}
	if s.Ref != "#/definitions/Pet" {
		t.Errorf("Ref = %q, want %q", s.Ref, "#/definitions/Pet")
	}
}

func TestSwaggerSchemaToSchema_Items(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Items: &spec.SchemaOrArray{
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"string"},
					},
				},
			},
		},
	})
	if s == nil {
		t.Fatal("schema is nil")
	}
	if s.Items == nil {
		t.Fatal("Items is nil")
	}
	if s.Items.Type != "string" {
		t.Errorf("Items.Type = %q, want %q", s.Items.Type, "string")
	}
}

func TestSwaggerSchemaToSchema_Properties(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Properties: map[string]spec.Schema{
				"name": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
			},
		},
	})
	if s == nil {
		t.Fatal("schema is nil")
	}
	if len(s.Properties) != 1 {
		t.Fatalf("Properties = %d, want 1", len(s.Properties))
	}
	if s.Properties["name"].Type != "string" {
		t.Errorf("Properties[name].Type = %q, want %q", s.Properties["name"].Type, "string")
	}
}

func TestSwaggerSchemaToSchema_OneOfAnyOfAllOf(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			OneOf: []spec.Schema{
				{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
			},
			AnyOf: []spec.Schema{
				{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
			},
			AllOf: []spec.Schema{
				{SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"number"}}},
			},
		},
	})
	if s == nil {
		t.Fatal("schema is nil")
	}
	if len(s.OneOf) != 1 || s.OneOf[0].Type != "string" {
		t.Error("OneOf not preserved")
	}
	if len(s.AnyOf) != 1 || s.AnyOf[0].Type != "integer" {
		t.Error("AnyOf not preserved")
	}
	if len(s.AllOf) != 1 || s.AllOf[0].Type != "number" {
		t.Error("AllOf not preserved")
	}
}

func TestSwaggerOpToOp_DefaultResponse(t *testing.T) {
	t.Parallel()

	op := swaggerOpToOp(&spec.Operation{
		OperationProps: spec.OperationProps{
			ID: "testOp",
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					Default: &spec.Response{
						ResponseProps: spec.ResponseProps{
							Description: "Default response",
						},
					},
				},
			},
		},
	})
	if op == nil {
		t.Fatal("op is nil")
	}
	resp, ok := op.Responses["default"]
	if !ok {
		t.Fatal("default response not found")
	}
	if resp.Description != "Default response" {
		t.Errorf("Description = %q, want %q", resp.Description, "Default response")
	}
}
