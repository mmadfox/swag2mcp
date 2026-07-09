package service

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

func TestResolveMaxResponseSize_NilConfig(t *testing.T) {
	t.Parallel()

	size := resolveMaxResponseSize(nil)
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_NilField(t *testing.T) {
	t.Parallel()

	size := resolveMaxResponseSize(&types.HTTPClientConfig{})
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_Custom(t *testing.T) {
	t.Parallel()

	val := 4096
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != 4096 {
		t.Errorf("got %d, want %d", size, 4096)
	}
}

func TestResolveMaxResponseSize_ExceedsMax(t *testing.T) {
	t.Parallel()

	val := 2 * 1024 * 1024 // 2 MB
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != maxMaxResponseSize {
		t.Errorf("got %d, want %d", size, maxMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_Zero(t *testing.T) {
	t.Parallel()

	val := 0
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestResolveMaxResponseSize_Negative(t *testing.T) {
	t.Parallel()

	val := -100
	size := resolveMaxResponseSize(&types.HTTPClientConfig{MaxResponseSize: &val})
	if size != defaultMaxResponseSize {
		t.Errorf("got %d, want %d", size, defaultMaxResponseSize)
	}
}

func TestOpenCommand_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("only on darwin")
	}

	cmd := openCommand("/tmp/test.json")
	if cmd != "open /tmp/test.json" {
		t.Errorf("got %q, want %q", cmd, "open /tmp/test.json")
	}
}

func TestOpenCommand_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("only on linux")
	}

	cmd := openCommand("/tmp/test.json")
	if cmd != "xdg-open /tmp/test.json" {
		t.Errorf("got %q, want %q", cmd, "xdg-open /tmp/test.json")
	}
}

func TestOpenCommand_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("only on windows")
	}

	cmd := openCommand("C:\\test.json")
	if cmd != "start C:\\test.json" {
		t.Errorf("got %q, want %q", cmd, "start C:\\test.json")
	}
}

func TestFormatSize_Bytes(t *testing.T) {
	t.Parallel()

	if s := formatSize(500); s != "500 B" {
		t.Errorf("got %q, want %q", s, "500 B")
	}
}

func TestFormatSize_KB(t *testing.T) {
	t.Parallel()

	if s := formatSize(2048); s != "2.0 KB" {
		t.Errorf("got %q, want %q", s, "2.0 KB")
	}
}

func TestFormatSize_MB(t *testing.T) {
	t.Parallel()

	if s := formatSize(1048576); s != "1.0 MB" {
		t.Errorf("got %q, want %q", s, "1.0 MB")
	}
}

func TestFormatSize_GB(t *testing.T) {
	t.Parallel()

	if s := formatSize(1073741824); s != "1.0 GB" {
		t.Errorf("got %q, want %q", s, "1.0 GB")
	}
}

func TestFormatSize_Zero(t *testing.T) {
	t.Parallel()

	if s := formatSize(0); s != "0 B" {
		t.Errorf("got %q, want %q", s, "0 B")
	}
}

func TestRandomSuffix_Length(t *testing.T) {
	t.Parallel()

	suffix := randomSuffix(6)
	if len(suffix) != 6 {
		t.Errorf("len = %d, want %d", len(suffix), 6)
	}
}

func TestRandomSuffix_HexChars(t *testing.T) {
	t.Parallel()

	suffix := randomSuffix(12)
	for _, c := range suffix {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("unexpected char %c in suffix %q", c, suffix)
		}
	}
}

func TestRandomSuffix_Unique(t *testing.T) {
	t.Parallel()

	s1 := randomSuffix(6)
	s2 := randomSuffix(6)
	if s1 == s2 {
		t.Error("two random suffixes are identical")
	}
}

func TestSaveLargeResponse(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	body := make([]byte, 10000)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	endpoint := &types.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 2048)
	if err != nil {
		t.Fatalf("saveLargeResponse() = %v", err)
	}

	if resp.FileRef == nil {
		t.Fatal("FileRef is nil")
	}
	if resp.FileRef.Size != 10000 {
		t.Errorf("Size = %d, want %d", resp.FileRef.Size, 10000)
	}
	if resp.FileRef.SizeHint == "" {
		t.Error("SizeHint is empty")
	}
	if resp.FileRef.MaxSizeHint == "" {
		t.Error("MaxSizeHint is empty")
	}
	if resp.FileRef.Message == "" {
		t.Error("Message is empty")
	}
	if resp.FileRef.OpenCmd == "" {
		t.Error("OpenCmd is empty")
	}
	if !strings.HasPrefix(resp.FileRef.Path, svc.ws.ResponsesDir()) {
		t.Errorf("Path %q not in responses dir %q", resp.FileRef.Path, svc.ws.ResponsesDir())
	}

	if _, statErr := os.Stat(resp.FileRef.Path); os.IsNotExist(statErr) {
		t.Error("response file was not created on disk")
	}
}

func TestSaveLargeResponse_FileContent(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	body := []byte(`{"key": "value"}`)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
	}

	endpoint := &types.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 100)
	if err != nil {
		t.Fatalf("saveLargeResponse() = %v", err)
	}

	data, err := os.ReadFile(resp.FileRef.Path)
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(data) != string(body) {
		t.Errorf("file content = %q, want %q", string(data), string(body))
	}
}

// --- validateParameters tests ---

func TestValidateParameters_UnknownParameter(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, map[string]any{"unknown": "val"})
	if err == nil {
		t.Fatal("expected error for unknown parameter")
	}
}

func TestValidateParameters_MissingRequired(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
			{Name: "name", In: "query", Required: false},
		},
	}
	err := validateParameters(op, map[string]any{"name": "test"})
	if err == nil {
		t.Fatal("expected error for missing required parameter")
	}
}

func TestValidateParameters_AllValid(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
			{Name: "name", In: "query", Required: false},
		},
	}
	err := validateParameters(op, map[string]any{"id": "123", "name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateParameters_NoDeclaredParams(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{}
	err := validateParameters(op, map[string]any{"id": "123"})
	if err == nil {
		t.Fatal("expected error for unknown parameter with no declared params")
	}
}

func TestValidateParameters_NilParams(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, nil)
	if err == nil {
		t.Fatal("expected error for nil params with required parameter")
	}
}

// --- validateRequestBody tests ---

func TestValidateRequestBody_NilOperationBody(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{}
	err := validateRequestBody(op, map[string]any{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRequestBody_RequiredBodyNil(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content:  map[string]*spec.MediaType{"application/json": {Schema: &spec.Schema{Type: "object"}}},
		},
	}
	err := validateRequestBody(op, nil)
	if err == nil {
		t.Fatal("expected error for required body with nil")
	}
}

func TestValidateRequestBody_NotRequiredBodyNil(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: false,
		},
	}
	err := validateRequestBody(op, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRequestBody_ValidBody(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
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
	}
	err := validateRequestBody(op, map[string]any{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRequestBody_MissingRequiredField(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type:     "object",
						Required: []string{"name"},
						Properties: map[string]*spec.Schema{
							"name": {Type: "string"},
						},
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestValidateRequestBody_UnknownField(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
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
	}
	err := validateRequestBody(op, map[string]any{"name": "test", "unknown": "val"})
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestValidateRequestBody_NestedObject(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"address": {
								Type:     "object",
								Required: []string{"city"},
								Properties: map[string]*spec.Schema{
									"city": {Type: "string"},
								},
							},
						},
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{"address": map[string]any{}})
	if err == nil {
		t.Fatal("expected error for missing nested required field")
	}
}

func TestValidateRequestBody_NestedArray(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: true,
			Content: map[string]*spec.MediaType{
				"application/json": {
					Schema: &spec.Schema{
						Type: "object",
						Properties: map[string]*spec.Schema{
							"items": {
								Type: "array",
								Items: &spec.Schema{
									Type:     "object",
									Required: []string{"id"},
									Properties: map[string]*spec.Schema{
										"id": {Type: "string"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := validateRequestBody(op, map[string]any{"items": []any{map[string]any{}}})
	if err == nil {
		t.Fatal("expected error for missing nested array item required field")
	}
}

// --- schemaForContentType tests ---

func TestSchemaForContentType_NilContent(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(nil)
	if schema != nil {
		t.Fatal("expected nil")
	}
}

func TestSchemaForContentType_EmptyContent(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(map[string]*spec.MediaType{})
	if schema != nil {
		t.Fatal("expected nil")
	}
}

func TestSchemaForContentType_NonJSONContent(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(map[string]*spec.MediaType{
		"text/plain": {Schema: &spec.Schema{Type: "string"}},
	})
	if schema != nil {
		t.Fatal("expected nil for non-json content type")
	}
}

func TestSchemaForContentType_JSONWithNilSchema(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(map[string]*spec.MediaType{
		"application/json": nil,
	})
	if schema != nil {
		t.Fatal("expected nil when media type is nil")
	}
}

func TestSchemaForContentType_JSONWithSchema(t *testing.T) {
	t.Parallel()

	expected := &spec.Schema{Type: "object"}
	schema := schemaForContentType(map[string]*spec.MediaType{
		"application/json": {Schema: expected},
	})
	if schema != expected {
		t.Fatal("expected the schema to be returned")
	}
}

// --- validateSchemaValue tests ---

func TestValidateSchemaValue_NilSchema(t *testing.T) {
	t.Parallel()

	err := validateSchemaValue(nil, "value", "$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSchemaValue_ObjectType(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type:     "object",
		Required: []string{"name"},
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
		},
	}
	err := validateSchemaValue(schema, map[string]any{}, "$")
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestValidateSchemaValue_ArrayType(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type:     "object",
			Required: []string{"id"},
			Properties: map[string]*spec.Schema{
				"id": {Type: "string"},
			},
		},
	}
	err := validateSchemaValue(schema, []any{map[string]any{}}, "$")
	if err == nil {
		t.Fatal("expected error for missing required field in array item")
	}
}

func TestValidateSchemaValue_UnknownType(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "string"}
	err := validateSchemaValue(schema, "hello", "$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- validateObjectSchema tests ---

func TestValidateObjectSchema_ValueNotMap(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "object"}
	err := validateObjectSchema(schema, "not-a-map", "$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateObjectSchema_MissingRequired(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Required: []string{"name"},
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
		},
	}
	err := validateObjectSchema(schema, map[string]any{}, "$")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateObjectSchema_UnknownField(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
		},
	}
	err := validateObjectSchema(schema, map[string]any{"name": "test", "unknown": "val"}, "$")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateObjectSchema_NestedValidation(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"child": {
				Type:     "object",
				Required: []string{"id"},
				Properties: map[string]*spec.Schema{
					"id": {Type: "string"},
				},
			},
		},
	}
	err := validateObjectSchema(schema, map[string]any{"child": map[string]any{}}, "$")
	if err == nil {
		t.Fatal("expected error for nested missing required field")
	}
}

// --- validateArraySchema tests ---

func TestValidateArraySchema_ValueNotSlice(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "array"}
	err := validateArraySchema(schema, "not-a-slice", "$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateArraySchema_ItemValidationSuccess(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type: "object",
			Properties: map[string]*spec.Schema{
				"id": {Type: "string"},
			},
		},
	}
	err := validateArraySchema(schema, []any{map[string]any{"id": "1"}}, "$")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateArraySchema_ItemValidationFailure(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Type: "array",
		Items: &spec.Schema{
			Type:     "object",
			Required: []string{"id"},
			Properties: map[string]*spec.Schema{
				"id": {Type: "string"},
			},
		},
	}
	err := validateArraySchema(schema, []any{map[string]any{}}, "$")
	if err == nil {
		t.Fatal("expected error for missing required field in item")
	}
}

// --- newAuthHTTPClient tests ---

func TestNewAuthHTTPClient_NilSpec(t *testing.T) {
	t.Parallel()

	client := newAuthHTTPClient(nil, nil)
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNewAuthHTTPClient_WithAuth(t *testing.T) {
	t.Parallel()

	authenticator := &auth.BearerTokenAuthClient{Token: "test-token"}
	if err := authenticator.New(); err != nil {
		t.Fatalf("failed to init authenticator: %v", err)
	}

	spec := &types.Spec{Auth: authenticator}
	client := newAuthHTTPClient(spec, nil)
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestNewAuthHTTPClient_WithTimeout(t *testing.T) {
	t.Parallel()

	config := &types.HTTPClientConfig{Timeout: 30}
	client := newAuthHTTPClient(nil, config)
	if client == nil {
		t.Fatal("client is nil")
	}
	if client.Timeout != 30 {
		t.Errorf("Timeout = %v, want %v", client.Timeout, 30)
	}
}

func TestNewAuthHTTPClient_NoFollowRedirects(t *testing.T) {
	t.Parallel()

	follow := false
	config := &types.HTTPClientConfig{FollowRedirects: &follow}
	client := newAuthHTTPClient(nil, config)
	if client == nil {
		t.Fatal("client is nil")
	}
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil, expected ErrUseLastResponse")
	}
}

func TestNewAuthHTTPClient_MaxRedirects(t *testing.T) {
	t.Parallel()

	maxRedirects := 3
	config := &types.HTTPClientConfig{MaxRedirects: &maxRedirects}
	client := newAuthHTTPClient(nil, config)
	if client == nil {
		t.Fatal("client is nil")
	}
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil")
	}
}

// --- applyHTTPClientTimeout tests ---

func TestApplyHTTPClientTimeout_Zero(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	applyHTTPClientTimeout(client, &types.HTTPClientConfig{Timeout: 0})
	if client.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0", client.Timeout)
	}
}

func TestApplyHTTPClientTimeout_Positive(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	applyHTTPClientTimeout(client, &types.HTTPClientConfig{Timeout: 15})
	if client.Timeout != 15 {
		t.Errorf("Timeout = %v, want %v", client.Timeout, 15)
	}
}

// --- applyHTTPClientRedirects tests ---

func TestApplyHTTPClientRedirects_Noop(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	config := &types.HTTPClientConfig{}
	applyHTTPClientRedirects(client, config)
	if client.CheckRedirect != nil {
		t.Fatal("CheckRedirect should be nil when both fields are nil")
	}
}

func TestApplyHTTPClientRedirects_NilConfig(t *testing.T) {
	t.Parallel()

	// applyHTTPClientRedirects panics on nil config, so we test via newAuthHTTPClient
	result := newAuthHTTPClient(&types.Spec{}, nil)
	if result.CheckRedirect != nil {
		t.Fatal("CheckRedirect should be nil")
	}
}

func TestApplyHTTPClientRedirects_FollowRedirectsFalse(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	follow := false
	applyHTTPClientRedirects(client, &types.HTTPClientConfig{FollowRedirects: &follow})
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil")
	}
}

func TestApplyHTTPClientRedirects_MaxRedirects(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	maxRedirects := 5
	applyHTTPClientRedirects(client, &types.HTTPClientConfig{MaxRedirects: &maxRedirects})
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil")
	}
}

func TestApplyHTTPClientRedirects_BothSet(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	follow := true
	maxRedirects := 5
	applyHTTPClientRedirects(client, &types.HTTPClientConfig{FollowRedirects: &follow, MaxRedirects: &maxRedirects})
	if client.CheckRedirect == nil {
		t.Fatal("CheckRedirect is nil")
	}
}

// --- dumpRequest tests ---

func TestDumpRequest_EmptyDumpDir(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	svc.dumpRequest(req, "test-domain")
	// Should not panic or write anything
}

func TestDumpRequest_WithDumpDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t, WithDumpDir(tmpDir))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	svc.dumpRequest(req, "test-domain")

	// Check that a file was created in the dump dir
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir() = %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no dump files created")
	}
}

func TestResolveBaseURL_CollectionOverride(t *testing.T) {
	t.Parallel()

	builder := &requestBuilder{
		spec:       &types.Spec{BaseURL: "https://spec.example.com"},
		collection: &types.Collection{BaseURL: "https://collection.example.com"},
	}
	url := builder.resolveBaseURL()
	if url != "https://collection.example.com" {
		t.Errorf("got %q, want %q", url, "https://collection.example.com")
	}
}

func TestResolveBaseURL_SpecFallback(t *testing.T) {
	t.Parallel()

	builder := &requestBuilder{
		spec:       &types.Spec{BaseURL: "https://spec.example.com"},
		collection: &types.Collection{},
	}
	url := builder.resolveBaseURL()
	if url != "https://spec.example.com" {
		t.Errorf("got %q, want %q", url, "https://spec.example.com")
	}
}

func TestSaveLargeResponse_StatusCode(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	seedTestData(t, svc, t.Name())

	body := []byte("test")
	response := &http.Response{
		StatusCode: http.StatusNotFound,
		Header:     http.Header{},
	}

	endpoint := &types.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 100)
	if err != nil {
		t.Fatalf("saveLargeResponse() = %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
