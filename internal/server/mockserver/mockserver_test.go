package mockserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

func TestGenerateFromSchema_NilSchema(t *testing.T) {
	t.Parallel()

	result := GenerateFromSchema(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestGenerateFromSchema_String(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "string"}
	result := GenerateFromSchema(schema)
	_, ok := result.(string)
	if !ok {
		t.Errorf("expected string, got %T", result)
	}
}

func TestGenerateFromSchema_Integer(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "integer"}
	result := GenerateFromSchema(schema)
	_, ok := result.(int)
	if !ok {
		t.Errorf("expected int, got %T", result)
	}
}

func TestGenerateFromSchema_Number(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "number"}
	result := GenerateFromSchema(schema)
	_, ok := result.(float64)
	if !ok {
		t.Errorf("expected float64, got %T", result)
	}
}

func TestGenerateFromSchema_Boolean(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "boolean"}
	result := GenerateFromSchema(schema)
	_, ok := result.(bool)
	if !ok {
		t.Errorf("expected bool, got %T", result)
	}
}

func TestGenerateFromSchema_Array(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type: "string",
		},
	}
	result := GenerateFromSchema(schema)
	array, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}
	if len(array) < 1 || len(array) > 3 {
		t.Errorf("expected array length 1-3, got %d", len(array))
	}
}

func TestGenerateFromSchema_Object(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
	}
	result := GenerateFromSchema(schema)
	obj, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if _, exists := obj["name"]; !exists {
		t.Error("expected 'name' key in result")
	}
	if _, exists := obj["age"]; !exists {
		t.Error("expected 'age' key in result")
	}
}

func TestGenerateFromSchema_Enum(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "string",
		Enum: []any{"red", "green", "blue"},
	}
	result := GenerateFromSchema(schema)
	value, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	valid := value == "red" || value == "green" || value == "blue"
	if !valid {
		t.Errorf("expected one of [red, green, blue], got %q", value)
	}
}

func TestGenerateFromSchema_Default(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:    "string",
		Default: "default-value",
	}
	result := GenerateFromSchema(schema)
	if result != "default-value" {
		t.Errorf("expected %q, got %v", "default-value", result)
	}
}

func TestGenerateFromSchema_Example(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:    "string",
		Example: "example-value",
	}
	result := GenerateFromSchema(schema)
	if result != "example-value" {
		t.Errorf("expected %q, got %v", "example-value", result)
	}
}

func TestGenerateFromSchema_AllOf(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		AllOf: []*spec.Schema{
			{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"name": {Type: "string"},
				},
			},
			{
				Type: "object",
				Properties: map[string]*spec.Schema{
					"age": {Type: "integer"},
				},
			},
		},
	}
	result := GenerateFromSchema(schema)
	obj, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if _, exists := obj["name"]; !exists {
		t.Error("expected 'name' key from allOf merge")
	}
	if _, exists := obj["age"]; !exists {
		t.Error("expected 'age' key from allOf merge")
	}
}

func TestGenerateFromSchema_OneOf(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		OneOf: []*spec.Schema{
			{Type: "string"},
			{Type: "integer"},
		},
	}
	result := GenerateFromSchema(schema)
	switch result.(type) {
	case string, int:
	default:
		t.Errorf("expected string or int, got %T", result)
	}
}

func TestGenerateFromSchema_AnyOf(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		AnyOf: []*spec.Schema{
			{Type: "string"},
			{Type: "integer"},
		},
	}
	result := GenerateFromSchema(schema)
	switch result.(type) {
	case string, int:
	default:
		t.Errorf("expected string or int, got %T", result)
	}
}

func TestGenerateFromSchema_StringFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format string
	}{
		{"date"},
		{"date-time"},
		{"email"},
		{"uri"},
		{"uuid"},
		{"ipv4"},
		{"hostname"},
		{"phone"},
		{"byte"},
		{"binary"},
		{"password"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			t.Parallel()
			schema := &spec.Schema{Type: "string", Format: tt.format}
			result := GenerateFromSchema(schema)
			_, ok := result.(string)
			if !ok {
				t.Errorf("expected string for format %q, got %T", tt.format, result)
			}
		})
	}
}

func TestGenerateFromSchema_IntegerFormat(t *testing.T) {
	t.Parallel()

	t.Run("int32", func(t *testing.T) {
		t.Parallel()
		schema := &spec.Schema{Type: "integer", Format: "int32"}
		result := GenerateFromSchema(schema)
		_, ok := result.(int32)
		if !ok {
			t.Errorf("expected int32, got %T", result)
		}
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		schema := &spec.Schema{Type: "integer", Format: "int64"}
		result := GenerateFromSchema(schema)
		_, ok := result.(int64)
		if !ok {
			t.Errorf("expected int64, got %T", result)
		}
	})
}

func TestGenerateFromSchema_ReadOnlySkipped(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "object",
		Properties: map[string]*spec.Schema{
			"readOnlyField":  {Type: "string", ReadOnly: true},
			"writeOnlyField": {Type: "string", WriteOnly: true},
			"normalField":    {Type: "string"},
		},
	}
	result := GenerateFromSchema(schema)
	obj, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if _, exists := obj["readOnlyField"]; exists {
		t.Error("readOnly field should be skipped")
	}
	if _, exists := obj["writeOnlyField"]; exists {
		t.Error("writeOnly field should be skipped")
	}
	if _, exists := obj["normalField"]; !exists {
		t.Error("normal field should be present")
	}
}

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

func TestAuthMockServer_handleDigest_NoChallenge(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerDigest, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	server.handleDigest(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", responseRecorder.Code)
	}
	wwwAuth := responseRecorder.Header().Get("WWW-Authenticate")
	if !strings.HasPrefix(wwwAuth, "Digest ") {
		t.Errorf("expected WWW-Authenticate to start with 'Digest ', got %q", wwwAuth)
	}
}

func TestAuthMockServer_handleDigest_ValidResponse(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerDigest, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", `Digest username="test", realm="test", nonce="abc", uri="/", response="def"`)

	server.handleDigest(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleOAuth2_InvalidMethod(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/token", nil)

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", responseRecorder.Code)
	}
}

func TestAuthMockServer_handleOAuth2_CC_ValidRequest(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	body := strings.NewReader("grant_type=client_credentials&client_id=test-client")
	request := httptest.NewRequest(http.MethodPost, "/token", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}

	var response map[string]any
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["access_token"] == "" {
		t.Error("expected access_token to be non-empty")
	}
}

func TestAuthMockServer_handleOAuth2_Password_ValidRequest(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	body := strings.NewReader("grant_type=password&username=alice&password=secret")
	request := httptest.NewRequest(http.MethodPost, "/token", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", responseRecorder.Code)
	}

	var response map[string]any
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response["access_token"] == "" {
		t.Error("expected access_token to be non-empty")
	}
}

func TestAuthMockServer_handleOAuth2_InvalidGrantType(t *testing.T) {
	t.Parallel()

	server := newAuthMockServer(authServerOAuth2, "127.0.0.1:0", nil, nil)
	responseRecorder := httptest.NewRecorder()
	body := strings.NewReader("grant_type=invalid_grant")
	request := httptest.NewRequest(http.MethodPost, "/token", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.handleOAuth2(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", responseRecorder.Code)
	}
}

func TestGenerateRandomToken(t *testing.T) {
	t.Parallel()

	token := generateRandomToken()
	if len(token) != authTokenLength*2 {
		t.Errorf("expected token length %d, got %d", authTokenLength*2, len(token))
	}
}

func TestParseDigestAuthorization(t *testing.T) {
	t.Parallel()

	server := &authMockServer{}
	authorization := `Digest username="test", realm="example", nonce="abc123", uri="/", response="def456", opaque="xyz789", algorithm=MD5, qop=auth`

	params := server.parseDigestAuthorization(authorization)

	expected := map[string]string{
		"username":  "test",
		"realm":     "example",
		"nonce":     "abc123",
		"uri":       "/",
		"response":  "def456",
		"opaque":    "xyz789",
		"algorithm": "MD5",
		"qop":       "auth",
	}

	for key, expectedValue := range expected {
		if params[key] != expectedValue {
			t.Errorf("expected %q=%q, got %q", key, expectedValue, params[key])
		}
	}
}

func TestNewMockServer_NoServers(t *testing.T) {
	t.Parallel()

	server := New(Options{
		Config: &config.Config{},
	})
	err := server.Start(context.Background())
	if err == nil {
		t.Error("expected error when mock_enabled is false")
	}
}

func TestNewMockServer_MockDisabled(t *testing.T) {
	t.Parallel()

	server := New(Options{
		Config: &config.Config{
			MockEnabled: false,
			Specs: []config.Spec{
				{
					Domain:   "test-api",
					LLMTitle: "Test API v1",
					BaseURL:  "https://api.example.com",
					Collections: []config.Collection{
						{
							LLMTitle:    "Main Collection",
							Location:    "https://example.com/spec.yaml",
							BaseMockURL: "localhost:8080",
						},
					},
				},
			},
		},
	})
	err := server.Start(context.Background())
	if err == nil {
		t.Error("expected error when mock_enabled is false")
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

	var response map[string]any
	if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, exists := response["name"]; !exists {
		t.Error("expected 'name' key in response")
	}
}

func TestExtractHostPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		addr string
		want string
	}{
		{"localhost:8080", "localhost:8080"},
		{"127.0.0.1:9000/v1/smev", "127.0.0.1:9000"},
		{"localhost:8080/api/v1", "localhost:8080"},
		{"0.0.0.0:3000/path/to/service", "0.0.0.0:3000"},
		{"127.0.0.1:80", "127.0.0.1:80"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			t.Parallel()
			got := extractHostPort(tt.addr)
			if got != tt.want {
				t.Errorf("extractHostPort(%q) = %q, want %q", tt.addr, got, tt.want)
			}
		})
	}
}
