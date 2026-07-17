package spec

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_postmanCollection(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "postman_petstore.json"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err, "Parse(postman) failed")

	assert.Equal(t, "2.x", doc.Version)
	assert.Equal(t, "Petstore API", doc.Title)

	require.Len(t, doc.PathItems, 4)

	var listPets, createPet, getPet, health bool
	for _, pi := range doc.PathItems {
		switch pi.Path {
		case "/v1/pets":
			switch pi.Method {
			case http.MethodGet:
				listPets = true
				op := pi.Operation
				assert.Equal(t, "List all pets", op.Summary)
				require.NotEmpty(t, op.Parameters, "expected query params on list pets")
				var hasLimit bool
				for _, p := range op.Parameters {
					if p.Name == "limit" && p.In == "query" {
						hasLimit = true
					}
				}
				assert.True(t, hasLimit, "expected limit query param")
			case http.MethodPost:
				createPet = true
				op := pi.Operation
				require.NotNil(t, op.RequestBody, "expected request body")
				require.NotNil(t, op.RequestBody.Content["application/json"], "expected JSON content type")
			}
		case "/v1/pets/{petId}":
			getPet = true
			op := pi.Operation
			var hasID bool
			for _, p := range op.Parameters {
				if p.Name == "petId" && p.In == "path" {
					hasID = true
				}
			}
			assert.True(t, hasID, "expected petId path param")
		case "/health":
			health = true
		}
	}

	assert.True(t, listPets, "GET /v1/pets not found")
	assert.True(t, createPet, "POST /v1/pets not found")
	assert.True(t, getPet, "GET /v1/pets/{petId} not found")
	assert.True(t, health, "GET /health not found")
}

func TestParse_postmanHeaders(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "postman_petstore.json"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	for _, pi := range doc.PathItems {
		if pi.Path == "/v1/pets" && pi.Method == http.MethodGet {
			op := pi.Operation
			var hasAuth bool
			for _, p := range op.Parameters {
				if p.Name == "Authorization" && p.In == "header" {
					hasAuth = true
				}
			}
			assert.True(t, hasAuth, "expected Authorization header")
			break
		}
	}
}

func TestParse_isPostman(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "postman collection",
			data: []byte(
				`{"info":{"schema":"https://schema.getpostman.com/collection/v2.1.0/collection.json"},"item":[{"name":"test","request":{"method":"GET","url":"http://example.com"}}]}`,
			),
			want: true,
		},
		{
			name: "openapi 3",
			data: []byte(`{"openapi":"3.0.0","info":{"title":"Test"},"paths":{}}`),
			want: false,
		},
		{
			name: "swagger 2",
			data: []byte(`{"swagger":"2.0","info":{"title":"Test"},"paths":{}}`),
			want: false,
		},
		{
			name: "empty",
			data: []byte(`{}`),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isPostman(tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParse_openapiServers(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v300_openapi.yaml"))
	require.NoError(t, err)

	doc, err := Parse(data)
	require.NoError(t, err)

	require.NotEmpty(t, doc.Servers, "expected at least 1 server from openapi servers")
}

func TestParse_tagHierarchies(t *testing.T) {
	t.Parallel()
	files := []string{
		"test_tags_01_flat.yaml",
		"test_tags_02_slash_hierarchy.yaml",
		"test_tags_03_dot_hierarchy.yaml",
		"test_tags_04_double_colon.yaml",
		"test_tags_05_dash_hierarchy.yaml",
		"test_tags_07_deep_hierarchy.yaml",
		"test_tags_14_multiple_hierarchical.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile(filepath.Join("testdata", file))
			require.NoError(t, err)
			doc, err := Parse(data)
			require.NoError(t, err, "Parse(%s) failed", file)
			assert.NotEmpty(t, doc.PathItems, "no path items")
		})
	}
}

func TestParse_multiTags(t *testing.T) {
	t.Parallel()
	files := []string{
		"test_multi_tags_01_cross_domain.yaml",
		"test_multi_tags_02_roles.yaml",
		"test_multi_tags_03_versions.yaml",
		"test_multi_tags_05_microservices.yaml",
		"test_multi_tags_09_mixed_hierarchy.yaml",
		"test_multi_tags_10_many_tags.yaml",
		"test_multi_tags_11_special_chars.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile(filepath.Join("testdata", file))
			require.NoError(t, err)
			doc, err := Parse(data)
			require.NoError(t, err, "Parse(%s) failed", file)
			assert.NotEmpty(t, doc.PathItems, "no path items")
		})
	}
}

func TestPostmanTag_WithFolders(t *testing.T) {
	t.Parallel()

	tag := postmanTag("Get Pet", []string{"Pets", "Store"})
	assert.Equal(t, "store", tag)
}

func TestPostmanTag_WithoutFolders(t *testing.T) {
	t.Parallel()

	tag := postmanTag("Get Pet", nil)
	assert.Equal(t, "get-pet", tag)
}

func TestSanitizePostmanTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"Get Pet", "get-pet"},
		{"Hello World!", "hello-world"},
		{"  spaces  ", "spaces"},
		{"UPPERCASE", "uppercase"},
		{"special_chars!@#", "special-chars"},
	}
	for _, tt := range tests {
		got := sanitizePostmanTag(tt.input)
		assert.Equal(t, tt.want, got, "sanitizePostmanTag(%q)", tt.input)
	}
}

func TestExtractPathFromURLString_NoScheme(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("example.com/api/v1/pets")
	assert.Equal(t, "/api/v1/pets", path)
}

func TestExtractPathFromURLString_ColonParam(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("http://example.com/api/v1/pets/:petId")
	assert.Equal(t, "/api/v1/pets/{petId}", path)
}

func TestExtractPathFromURLString_EmptyPath(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("http://example.com")
	assert.Equal(t, "/", path)
}

func TestExtractPathFromURLString_InvalidURL(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("http://[invalid]")
	assert.Equal(t, "http://[invalid]", path)
}

func TestExtractPostmanPath_NilURL(t *testing.T) {
	t.Parallel()

	path := extractPostmanPath(nil)
	assert.Equal(t, "/", path)
}

func TestExtractPostmanPath_StringURL(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal("http://example.com/api/v1/pets")
	path := extractPostmanPath(raw)
	assert.Equal(t, "/api/v1/pets", path)
}

func TestExtractPostmanPath_StructuredPath(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal(postmanURL{
		Raw:  "http://example.com/api/v1/pets",
		Path: []json.RawMessage{json.RawMessage(`"api"`), json.RawMessage(`"v1"`), json.RawMessage(`"pets"`)},
	})
	path := extractPostmanPath(raw)
	assert.Equal(t, "/api/v1/pets", path)
}

func TestExtractPostmanPath_PathVariable(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal(postmanURL{
		Raw:  "http://example.com/pets/:petId",
		Path: []json.RawMessage{json.RawMessage(`"pets"`), json.RawMessage(`{"type":"string","value":"petId"}`)},
	})
	path := extractPostmanPath(raw)
	assert.Equal(t, "/pets/{petId}", path)
}

func TestExtractPostmanPath_InvalidJSON(t *testing.T) {
	t.Parallel()

	path := extractPostmanPath(json.RawMessage("{invalid}"))
	assert.Equal(t, "/", path)
}

func TestAppendPostmanHeaders_Disabled(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	headers := []postmanHeader{
		{Key: "Authorization", Value: "Bearer token", Disabled: true},
	}
	appendPostmanHeaders(headers, op)
	assert.Empty(t, op.Parameters, "disabled header should be skipped")
}

func TestAppendPostmanBody_Nil(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(nil, op, http.MethodPost)
	require.Nil(t, op.RequestBody, "RequestBody should be nil for nil body")
}

func TestAppendPostmanBody_GetMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodGet)
	require.Nil(t, op.RequestBody, "RequestBody should be nil for GET")
}

func TestAppendPostmanBody_HeadMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodHead)
	require.Nil(t, op.RequestBody, "RequestBody should be nil for HEAD")
}

func TestAppendPostmanBody_DeleteMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodDelete)
	require.Nil(t, op.RequestBody, "RequestBody should be nil for DELETE")
}

func TestAppendPostmanBody_OptionsMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodOptions)
	require.Nil(t, op.RequestBody, "RequestBody should be nil for OPTIONS")
}

func TestAppendPostmanBody_RawJSON(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodPost)
	require.NotNil(t, op.RequestBody, "RequestBody is nil")
	require.NotNil(t, op.RequestBody.Content["application/json"], "expected application/json content type")
}

func TestAppendPostmanBody_RawXML(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: "<xml>data</xml>"}, op, http.MethodPost)
	require.NotNil(t, op.RequestBody, "RequestBody is nil")
	require.NotNil(t, op.RequestBody.Content["application/xml"], "expected application/xml content type")
}

func TestAppendPostmanBody_RawText(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: "plain text"}, op, http.MethodPost)
	require.NotNil(t, op.RequestBody, "RequestBody is nil")
	require.NotNil(t, op.RequestBody.Content["text/plain"], "expected text/plain content type")
}

func TestAppendPostmanBody_URLEncoded(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{
		Mode: "urlencoded",
		URLEnc: []postmanURLEncoded{
			{Key: "name", Value: "test"},
			{Key: "disabled_field", Value: "skip", Disabled: true},
		},
	}, op, http.MethodPost)
	require.NotNil(t, op.RequestBody, "RequestBody is nil")
	mt := op.RequestBody.Content["application/x-www-form-urlencoded"]
	require.NotNil(t, mt, "expected urlencoded content type")
	require.NotNil(t, mt.Schema.Properties["name"], "expected name property")
	require.Nil(t, mt.Schema.Properties["disabled_field"], "disabled field should not be present")
}

func TestAppendPostmanBody_FormData(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{
		Mode: "formdata",
		FormData: []postmanFormData{
			{Key: "file", Value: "test.txt", Type: "file"},
			{Key: "name", Value: "test"},
			{Key: "disabled_field", Value: "skip", Disabled: true},
		},
	}, op, http.MethodPost)
	require.NotNil(t, op.RequestBody, "RequestBody is nil")
	mt := op.RequestBody.Content["multipart/form-data"]
	require.NotNil(t, mt, "expected multipart/form-data content type")
	assert.Equal(t, "file", mt.Schema.Properties["file"].Type)
	assert.Equal(t, "string", mt.Schema.Properties["name"].Type)
	require.Nil(t, mt.Schema.Properties["disabled_field"], "disabled field should not be present")
}

func TestAppendPostmanBody_GraphQL(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "graphql"}, op, http.MethodPost)
	require.NotNil(t, op.RequestBody, "RequestBody is nil")
	require.NotNil(t, op.RequestBody.Content["application/json"], "expected application/json content type")
}

func TestAppendPostmanBody_UnknownMode(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "unknown"}, op, http.MethodPost)
	require.Nil(t, op.RequestBody, "RequestBody should be nil for unknown mode")
}

func TestGuessPostmanContentType_NilBody(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(nil)
	assert.Equal(t, "application/json", ct)
}

func TestGuessPostmanContentType_EmptyRaw(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: ""})
	assert.Equal(t, "application/json", ct)
}

func TestGuessPostmanContentType_Array(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: "[1,2,3]"})
	assert.Equal(t, "application/json", ct)
}

func TestGuessPostmanContentType_XML(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: "<root><item/></root>"})
	assert.Equal(t, "application/xml", ct)
}

func TestGuessPostmanContentType_Text(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: "plain text"})
	assert.Equal(t, "text/plain", ct)
}

func TestParsePostman_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parsePostman([]byte("{invalid}"))
	require.Error(t, err, "expected error for invalid JSON")
}

func TestParsePostman_EmptyCollection(t *testing.T) {
	t.Parallel()

	input := `{"info":{"name":"Test","schema":"https://schema.getpostman.com/collection/v2.1.0/collection.json"},"item":[]}`
	_, err := parsePostman([]byte(input))
	require.Error(t, err, "expected error for empty collection")
}

func TestIsPostman_InvalidJSON(t *testing.T) {
	t.Parallel()

	assert.False(t, isPostman([]byte("{invalid}")), "expected false for invalid JSON")
}

func TestIsPostman_NoPostmanSchema(t *testing.T) {
	t.Parallel()

	assert.False(t,
		isPostman([]byte(`{"info":{"schema":"https://example.com/schema.json"},"item":[{}]}`)),
		"expected false for non-postman schema")
}

func TestIsPostman_EmptyItems(t *testing.T) {
	t.Parallel()

	input := `{"info":{"schema":"https://schema.getpostman.com/collection/v2.1.0/collection.json"},"item":[]}`
	assert.False(t, isPostman([]byte(input)), "expected false for empty items")
}

func TestFlattenPostmanItems_NestedFolders(t *testing.T) {
	t.Parallel()

	doc := &Doc{PathItems: make([]*PathItem, 0)}
	items := []postmanItem{
		{
			Name: "Folder1",
			Item: []postmanItem{
				{
					Name: "SubItem",
					Request: &postmanRequest{
						Method: "GET",
						URL:    json.RawMessage(`"http://example.com/api"`),
					},
				},
			},
		},
	}
	err := flattenPostmanItems(nil, items, doc)
	require.NoError(t, err, "flattenPostmanItems() failed")
	require.Len(t, doc.PathItems, 1)
	assert.Equal(t, "folder1", doc.PathItems[0].Operation.Tags[0])
}

func TestPostmanItemToPathItem_EmptyMethod(t *testing.T) {
	t.Parallel()

	pi := postmanItemToPathItem(postmanItem{
		Name: "Test Item",
		Request: &postmanRequest{
			Method: "",
			URL:    json.RawMessage(`"http://example.com/api"`),
		},
	}, nil)
	assert.Equal(t, http.MethodGet, pi.Method)
}

func TestPostmanItemToPathItem_EmptyPath(t *testing.T) {
	t.Parallel()

	pi := postmanItemToPathItem(postmanItem{
		Name: "Test Item",
		Request: &postmanRequest{
			Method: "POST",
			URL:    nil,
		},
	}, nil)
	assert.Equal(t, "/", pi.Path)
}

func TestAppendPostmanURLParams_InvalidURL(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanURLParams(json.RawMessage("{invalid}"), op)
	assert.Empty(t, op.Parameters)
}

func TestAppendPostmanURLParams_WithVariables(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal(postmanURL{
		Raw: "http://example.com/pets/:petId",
		Variable: []postmanVariable{
			{Key: "petId", Value: "123"},
		},
	})
	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanURLParams(raw, op)
	require.Len(t, op.Parameters, 1)
	assert.Equal(t, "petId", op.Parameters[0].Name)
	assert.Equal(t, "path", op.Parameters[0].In)
}

func TestAppendPostmanURLParams_WithQuery(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal(postmanURL{
		Raw: "http://example.com/api?limit=10",
		Query: []postmanQueryVar{
			{Key: "limit", Value: "10"},
			{Key: "offset", Value: "0", Disabled: true},
		},
	})
	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanURLParams(raw, op)
	require.Len(t, op.Parameters, 1)
	assert.Equal(t, "limit", op.Parameters[0].Name)
	assert.Equal(t, "query", op.Parameters[0].In)
}
