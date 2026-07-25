package mockserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaForContentType_Nil(t *testing.T) {
	t.Parallel()

	result := schemaForContentType(nil)
	assert.Nil(t, result, "expected nil")
}

func TestSchemaForContentType_Empty(t *testing.T) {
	t.Parallel()

	result := schemaForContentType(make(map[string]*spec.MediaType))
	assert.Nil(t, result, "expected nil")
}

func TestSchemaForContentType_PrefersJSON(t *testing.T) {
	t.Parallel()

	jsonSchema := &spec.Schema{Type: "string"}
	xmlSchema := &spec.Schema{Type: "integer"}

	content := map[string]*spec.MediaType{
		"application/xml":  {Schema: xmlSchema},
		"application/json": {Schema: jsonSchema},
	}
	result := schemaForContentType(content)
	assert.Equal(t, jsonSchema, result, "expected application/json schema to be returned")
}

func TestSchemaForContentType_Fallback(t *testing.T) {
	t.Parallel()

	xmlSchema := &spec.Schema{Type: "integer"}
	content := map[string]*spec.MediaType{
		"application/xml": {Schema: xmlSchema},
	}
	result := schemaForContentType(content)
	assert.Equal(t, xmlSchema, result, "expected fallback schema to be returned")
}

func TestFindResponseSchema_NilOperation(t *testing.T) {
	t.Parallel()

	server := &apiMockServer{}
	result := server.findResponseSchema(nil)
	assert.Nil(t, result, "expected nil")
}

func TestFindResponseSchema_Prefers200(t *testing.T) {
	t.Parallel()

	server := &apiMockServer{}
	operation := &spec.Operation{
		Responses: map[string]*spec.Response{
			"200": {
				Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "string"}},
				},
			},
			"default": {
				Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "integer"}},
				},
			},
		},
	}
	result := server.findResponseSchema(operation)
	require.NotNil(t, result, "expected string schema from 200 response")
	assert.Equal(t, "string", result.Type, "expected string schema from 200 response")
}

func TestFindResponseSchema_FallbackToDefault(t *testing.T) {
	t.Parallel()

	server := &apiMockServer{}
	operation := &spec.Operation{
		Responses: map[string]*spec.Response{
			"default": {
				Content: map[string]*spec.MediaType{
					"application/json": {Schema: &spec.Schema{Type: "integer"}},
				},
			},
		},
	}
	result := server.findResponseSchema(operation)
	require.NotNil(t, result, "expected integer schema from default response")
	assert.Equal(t, "integer", result.Type, "expected integer schema from default response")
}

func TestCreateEndpointHandler_EmptyOperation(t *testing.T) {
	t.Parallel()

	server := &apiMockServer{
		doc: &spec.Doc{},
	}
	handler := server.createEndpointHandler(&spec.Operation{})
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "expected 200")
}

func TestCreateEndpointHandler_WithSchema(t *testing.T) {
	t.Parallel()

	server := &apiMockServer{
		doc: &spec.Doc{},
	}
	operation := &spec.Operation{
		Responses: map[string]*spec.Response{
			"200": {
				Content: map[string]*spec.MediaType{
					"application/json": {
						Schema: &spec.Schema{
							Type: "object",
							Properties: map[string]*spec.Schema{
								"name": {Type: "string"},
							},
						},
					},
				},
			},
		},
	}
	handler := server.createEndpointHandler(operation)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "expected 200")
}
