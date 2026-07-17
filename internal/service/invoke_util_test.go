package service

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestResolveMaxResponseSize_NilConfig(t *testing.T) {
	t.Parallel()

	size := resolveMaxResponseSize(nil)
	require.Equal(t, defaultMaxResponseSize, size)
}

func TestResolveMaxResponseSize_NilField(t *testing.T) {
	t.Parallel()

	size := resolveMaxResponseSize(nil)
	require.Equal(t, defaultMaxResponseSize, size)
}

func TestResolveMaxResponseSize_Custom(t *testing.T) {
	t.Parallel()

	val := 4096
	size := resolveMaxResponseSize(&val)
	require.Equal(t, 4096, size)
}

func TestResolveMaxResponseSize_ExceedsMax(t *testing.T) {
	t.Parallel()

	val := 2 * 1024 * 1024
	size := resolveMaxResponseSize(&val)
	require.Equal(t, maxMaxResponseSize, size)
}

func TestResolveMaxResponseSize_Zero(t *testing.T) {
	t.Parallel()

	val := 0
	size := resolveMaxResponseSize(&val)
	require.Equal(t, defaultMaxResponseSize, size)
}

func TestResolveMaxResponseSize_Negative(t *testing.T) {
	t.Parallel()

	val := -100
	size := resolveMaxResponseSize(&val)
	require.Equal(t, defaultMaxResponseSize, size)
}

func TestOpenCommand_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("only on darwin")
	}

	cmd := openCommand("/tmp/test.json")
	require.Equal(t, "open /tmp/test.json", cmd)
}

func TestOpenCommand_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("only on linux")
	}

	cmd := openCommand("/tmp/test.json")
	require.Equal(t, "xdg-open /tmp/test.json", cmd)
}

func TestOpenCommand_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("only on windows")
	}

	cmd := openCommand("C:\\test.json")
	require.Equal(t, "start C:\\test.json", cmd)
}

func TestFormatSize_Bytes(t *testing.T) {
	t.Parallel()

	require.Equal(t, "500 B", formatSize(500))
}

func TestFormatSize_KB(t *testing.T) {
	t.Parallel()

	require.Equal(t, "2.0 KB", formatSize(2048))
}

func TestFormatSize_MB(t *testing.T) {
	t.Parallel()

	require.Equal(t, "1.0 MB", formatSize(1048576))
}

func TestFormatSize_GB(t *testing.T) {
	t.Parallel()

	require.Equal(t, "1.0 GB", formatSize(1073741824))
}

func TestFormatSize_Zero(t *testing.T) {
	t.Parallel()

	require.Equal(t, "0 B", formatSize(0))
}

func TestRandomSuffix_Length(t *testing.T) {
	t.Parallel()

	suffix := randomSuffix(6)
	require.Len(t, suffix, 6)
}

func TestRandomSuffix_HexChars(t *testing.T) {
	t.Parallel()

	suffix := randomSuffix(12)
	for _, c := range suffix {
		require.Contains(t, "0123456789abcdef", string(c))
	}
}

func TestRandomSuffix_Unique(t *testing.T) {
	t.Parallel()

	s1 := randomSuffix(6)
	s2 := randomSuffix(6)
	require.NotEqual(t, s1, s2)
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

	endpoint := &model.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 2048)
	require.NoError(t, err)
	require.NotNil(t, resp.FileRef)
	require.Equal(t, 10000, resp.FileRef.Size)
	require.NotEmpty(t, resp.FileRef.SizeHint)
	require.NotEmpty(t, resp.FileRef.MaxSizeHint)
	require.NotEmpty(t, resp.FileRef.Message)
	require.NotEmpty(t, resp.FileRef.OpenCmd)
	require.True(t, strings.HasPrefix(resp.FileRef.Path, svc.ws.ResponsesDir()))

	_, statErr := os.Stat(resp.FileRef.Path)
	require.False(t, os.IsNotExist(statErr))
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

	endpoint := &model.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 100)
	require.NoError(t, err)

	data, err := os.ReadFile(resp.FileRef.Path)
	require.NoError(t, err)
	require.Equal(t, body, data)
}

func TestValidateParameters_UnknownParameter(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, map[string]any{"unknown": "val"})
	require.Error(t, err)
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
	require.Error(t, err)
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
	require.NoError(t, err)
}

func TestValidateParameters_NoDeclaredParams(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{}
	err := validateParameters(op, map[string]any{"id": "123"})
	require.Error(t, err)
}

func TestValidateParameters_NilParams(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		Parameters: []*spec.Parameter{
			{Name: "id", In: "path", Required: true},
		},
	}
	err := validateParameters(op, nil)
	require.Error(t, err)
}

func TestValidateRequestBody_NilOperationBody(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{}
	err := validateRequestBody(op, map[string]any{"key": "val"})
	require.NoError(t, err)
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
	require.Error(t, err)
}

func TestValidateRequestBody_NotRequiredBodyNil(t *testing.T) {
	t.Parallel()

	op := &spec.Operation{
		RequestBody: &spec.RequestBody{
			Required: false,
		},
	}
	err := validateRequestBody(op, nil)
	require.NoError(t, err)
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
	require.NoError(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
}

func TestSchemaForContentType_NilContent(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(nil)
	require.Nil(t, schema)
}

func TestSchemaForContentType_EmptyContent(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(map[string]*spec.MediaType{})
	require.Nil(t, schema)
}

func TestSchemaForContentType_NonJSONContent(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(map[string]*spec.MediaType{
		"text/plain": {Schema: &spec.Schema{Type: "string"}},
	})
	require.Nil(t, schema)
}

func TestSchemaForContentType_JSONWithNilSchema(t *testing.T) {
	t.Parallel()

	schema := schemaForContentType(map[string]*spec.MediaType{
		"application/json": nil,
	})
	require.Nil(t, schema)
}

func TestSchemaForContentType_JSONWithSchema(t *testing.T) {
	t.Parallel()

	expected := &spec.Schema{Type: "object"}
	schema := schemaForContentType(map[string]*spec.MediaType{
		"application/json": {Schema: expected},
	})
	require.Equal(t, expected, schema)
}

func TestValidateSchemaValue_NilSchema(t *testing.T) {
	t.Parallel()

	err := validateSchemaValue(nil, "value", "$")
	require.NoError(t, err)
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
	require.Error(t, err)
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
	require.Error(t, err)
}

func TestValidateSchemaValue_UnknownType(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "string"}
	err := validateSchemaValue(schema, "hello", "$")
	require.NoError(t, err)
}

func TestValidateObjectSchema_ValueNotMap(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "object"}
	err := validateObjectSchema(schema, "not-a-map", "$")
	require.NoError(t, err)
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
	require.Error(t, err)
}

func TestValidateObjectSchema_UnknownField(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{
		Properties: map[string]*spec.Schema{
			"name": {Type: "string"},
		},
	}
	err := validateObjectSchema(schema, map[string]any{"name": "test", "unknown": "val"}, "$")
	require.Error(t, err)
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
	require.Error(t, err)
}

func TestValidateArraySchema_ValueNotSlice(t *testing.T) {
	t.Parallel()

	schema := &spec.Schema{Type: "array"}
	err := validateArraySchema(schema, "not-a-slice", "$")
	require.NoError(t, err)
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
	require.NoError(t, err)
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
	require.Error(t, err)
}

func TestNewHTTPClient_NilConfig(t *testing.T) {
	t.Parallel()

	client, err := httpclient.New(httpclient.Config{})
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewHTTPClient_WithTimeout(t *testing.T) {
	t.Parallel()

	client, err := httpclient.New(httpclient.Config{Timeout: 30 * time.Second})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, 30*time.Second, client.Timeout)
}

func TestNewHTTPClient_NoFollowRedirects(t *testing.T) {
	t.Parallel()

	follow := false
	client, err := httpclient.New(httpclient.Config{FollowRedirects: &follow})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.CheckRedirect)
}

func TestNewHTTPClient_MaxRedirects(t *testing.T) {
	t.Parallel()

	maxRedirects := 3
	client, err := httpclient.New(httpclient.Config{MaxRedirects: &maxRedirects})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.CheckRedirect)
}

func TestNewHTTPClient_TimeoutZero(t *testing.T) {
	t.Parallel()

	client, err := httpclient.New(httpclient.Config{Timeout: 0})
	require.NoError(t, err)
	require.NotZero(t, client.Timeout)
}

func TestNewHTTPClient_RedirectsNoop(t *testing.T) {
	t.Parallel()

	client, err := httpclient.New(httpclient.Config{})
	require.NoError(t, err)
	require.Nil(t, client.CheckRedirect)
}

func TestNewHTTPClient_FollowRedirectsFalse(t *testing.T) {
	t.Parallel()

	follow := false
	client, err := httpclient.New(httpclient.Config{FollowRedirects: &follow})
	require.NoError(t, err)
	require.NotNil(t, client.CheckRedirect)
}

func TestNewHTTPClient_MaxRedirectsSet(t *testing.T) {
	t.Parallel()

	maxRedirects := 5
	client, err := httpclient.New(httpclient.Config{MaxRedirects: &maxRedirects})
	require.NoError(t, err)
	require.NotNil(t, client.CheckRedirect)
}

func TestNewHTTPClient_BothSet(t *testing.T) {
	t.Parallel()

	follow := true
	maxRedirects := 5
	client, err := httpclient.New(httpclient.Config{FollowRedirects: &follow, MaxRedirects: &maxRedirects})
	require.NoError(t, err)
	require.NotNil(t, client.CheckRedirect)
}

func TestDumpRequest_EmptyDumpDir(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	svc.dumpRequest(req, "test-domain")
}

func TestDumpRequest_WithDumpDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	svc := newTestService(t, WithDumpDir(tmpDir))
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/api", nil)
	svc.dumpRequest(req, "test-domain")

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
}

func TestResolveBaseURL_CollectionOverride(t *testing.T) {
	t.Parallel()

	builder := &requestBuilder{
		spec:       &model.Spec{BaseURL: "https://spec.example.com"},
		collection: &model.Collection{BaseURL: "https://collection.example.com"},
	}
	url := builder.resolveBaseURL()
	require.Equal(t, "https://collection.example.com", url)
}

func TestResolveBaseURL_SpecFallback(t *testing.T) {
	t.Parallel()

	builder := &requestBuilder{
		spec:       &model.Spec{BaseURL: "https://spec.example.com"},
		collection: &model.Collection{},
	}
	url := builder.resolveBaseURL()
	require.Equal(t, "https://spec.example.com", url)
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

	endpoint := &model.Endpoint{
		Name: "GET",
		Path: "/test",
	}

	resp, err := svc.saveLargeResponse(response, body, t.Name(), endpoint, 100)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
