package spec

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseV3_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseV3([]byte(`{"openapi":"3.0.0","paths":"not-an-object"}`))
	require.Error(t, err, "expected error for invalid OpenAPI document")
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
	require.NotNil(t, op, "op is nil")
	assert.Empty(t, op.Parameters)
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
	require.NotNil(t, op, "op is nil")
	resp, ok := op.Responses["200"]
	require.True(t, ok, "200 response not found")
	assert.Empty(t, resp.Description)
}

func TestSchemaRefToSchema_Nil(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(nil)
	require.Nil(t, s, "expected nil")
}

func TestSchemaRefToSchema_MultipleTypes(t *testing.T) {
	t.Parallel()

	s := schemaRefToSchema(&openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: &openapi3.Types{"string", "null"},
		},
	})
	require.NotNil(t, s, "schema is nil")
	assert.Equal(t, "string", s.Type)
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
	require.NotNil(t, s, "schema is nil")
	require.NotNil(t, s.Items, "Items is nil")
	assert.Equal(t, "string", s.Items.Type)
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
	require.NotNil(t, s, "schema is nil")
	require.Len(t, s.OneOf, 1)
	assert.Equal(t, "string", s.OneOf[0].Type)
	require.Len(t, s.AnyOf, 1)
	assert.Equal(t, "integer", s.AnyOf[0].Type)
	require.Len(t, s.AllOf, 1)
	assert.Equal(t, "number", s.AllOf[0].Type)
}
