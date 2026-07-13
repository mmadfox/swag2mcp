package spec

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestParseV3_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseV3([]byte(`{"openapi":"3.0.0","paths":"not-an-object"}`))
	if err == nil {
		t.Fatal("expected error for invalid OpenAPI document")
	}
}

func TestOpenapi3OpToOp_NilValue(t *testing.T) {
	t.Parallel()

	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "testOp",
		Parameters: []*openapi3.ParameterRef{
			nil,
			{Value: nil},
		},
	})
	if op == nil {
		t.Fatal("op is nil")
	}
	if len(op.Parameters) != 0 {
		t.Errorf("Parameters = %d, want 0", len(op.Parameters))
	}
}

func TestOpenapi3OpToOp_NilDescription(t *testing.T) {
	t.Parallel()

	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: nil,
		},
	})

	op := openapi3OpToOp(&openapi3.Operation{
		OperationID: "testOp",
		Responses:   responses,
	})
	if op == nil {
		t.Fatal("op is nil")
	}
	resp, ok := op.Responses["200"]
	if !ok {
		t.Fatal("200 response not found")
	}
	if resp.Description != "" {
		t.Errorf("Description = %q, want empty", resp.Description)
	}
}

func TestSchemaRefToSchema_Nil(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(nil)
	if s != nil {
		t.Fatal("expected nil")
	}
}

func TestSchemaRefToSchema_MultipleTypes(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"string", "null"},
		},
	})
	if s == nil {
		t.Fatal("schema is nil")
	}
	if s.Type != "string" {
		t.Errorf("Type = %q, want %q", s.Type, "string")
	}
}

func TestSchemaRefToSchema_Items(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{"string"},
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

func TestSchemaRefToSchema_Composition(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			OneOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{"string"}}},
			},
			AnyOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{"integer"}}},
			},
			AllOf: []*openapi3.SchemaRef{
				{Value: &openapi3.Schema{Type: &openapi3.Types{"number"}}},
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
