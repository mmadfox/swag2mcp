package spec

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_swaggerHost(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	require.NotEmpty(t, doc.Servers, "expected at least 1 server from swagger host")

	assert.Equal(t, "https://api.example.com/v1", doc.Servers[0].URL)
}

func TestParse_operationMetadata(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodGet {
			found = true
			op := pi.Operation
			assert.Equal(t, "Список пользователей", op.Summary)
			require.NotEmpty(t, op.Parameters, "expected parameters")
			assert.Equal(t, "limit", op.Parameters[0].Name)
			assert.Equal(t, "query", op.Parameters[0].In)
			break
		}
	}
	require.True(t, found, "GET /users not found in parsed doc")
}

func TestParse_requestBody(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodPost {
			found = true
			op := pi.Operation
			require.NotNil(t, op.RequestBody, "expected request body")
			assert.True(t, op.RequestBody.Required, "expected required request body")
			require.NotNil(t, op.RequestBody.Content, "expected request body content")
			break
		}
	}
	require.True(t, found, "POST /users not found")
}

func TestParse_responses(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == http.MethodGet {
			found = true
			op := pi.Operation
			require.NotEmpty(t, op.Responses, "expected responses")
			resp, ok := op.Responses["200"]
			require.True(t, ok, "expected 200 response")
			assert.Equal(t, "OK", resp.Description)
			break
		}
	}
	require.True(t, found, "GET /users not found")
}

func TestParse_swaggerFileUpload(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/files/upload" && pi.Method == http.MethodPost {
			found = true
			op := pi.Operation
			require.NotEmpty(t, op.Parameters, "expected parameters")
			assert.Equal(t, "file", op.Parameters[0].Name)
			assert.Equal(t, "formData", op.Parameters[0].In)
			break
		}
	}
	require.True(t, found, "POST /files/upload not found")
}

func TestParseV2_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseV2([]byte("{invalid}"))
	require.Error(t, err, "expected error for invalid JSON")
}

func TestParseV2_EmptyHost(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{"swagger":"2.0","info":{"title":"Test","version":"1.0"},"paths":{}}`)
	result, err := parseV2(jsonData)
	require.NoError(t, err, "parseV2() failed")
	assert.Empty(t, result.Servers)
}

func TestSwaggerSchemaToSchema_Nil(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(nil)
	require.Nil(t, s, "expected nil")
}

func TestSwaggerSchemaToSchema_Ref(t *testing.T) {
	t.Parallel()

	s := swaggerSchemaToSchema(&spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/Pet"),
		},
	})
	require.NotNil(t, s, "schema is nil")
	assert.Equal(t, "#/definitions/Pet", s.Ref)
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
	require.NotNil(t, s, "schema is nil")
	require.NotNil(t, s.Items, "Items is nil")
	assert.Equal(t, "string", s.Items.Type)
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
	require.NotNil(t, s, "schema is nil")
	require.Len(t, s.Properties, 1)
	assert.Equal(t, "string", s.Properties["name"].Type)
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
	require.NotNil(t, s, "schema is nil")
	require.Len(t, s.OneOf, 1)
	assert.Equal(t, "string", s.OneOf[0].Type)
	require.Len(t, s.AnyOf, 1)
	assert.Equal(t, "integer", s.AnyOf[0].Type)
	require.Len(t, s.AllOf, 1)
	assert.Equal(t, "number", s.AllOf[0].Type)
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
	require.NotNil(t, op, "op is nil")
	resp, ok := op.Responses["default"]
	require.True(t, ok, "default response not found")
	assert.Equal(t, "Default response", resp.Description)
}
