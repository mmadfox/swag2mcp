package mockserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/spec"
)

func TestSchemaForContentType_Nil(t *testing.T) {
	t.Parallel()

	result := schemaForContentType(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestSchemaForContentType_Empty(t *testing.T) {
	t.Parallel()

	result := schemaForContentType(make(map[string]*spec.MediaType))
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
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
	if result != jsonSchema {
		t.Error("expected application/json schema to be returned")
	}
}

func TestSchemaForContentType_Fallback(t *testing.T) {
	t.Parallel()

	xmlSchema := &spec.Schema{Type: "integer"}
	content := map[string]*spec.MediaType{
		"application/xml": {Schema: xmlSchema},
	}
	result := schemaForContentType(content)
	if result != xmlSchema {
		t.Error("expected fallback schema to be returned")
	}
}

func TestFindResponseSchema_NilOperation(t *testing.T) {
	t.Parallel()

	server := &apiMockServer{}
	result := server.findResponseSchema(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
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
	if result == nil || result.Type != "string" {
		t.Errorf("expected string schema from 200 response, got %v", result)
	}
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
	if result == nil || result.Type != "integer" {
		t.Errorf("expected integer schema from default response, got %v", result)
	}
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

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}
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

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}
}
