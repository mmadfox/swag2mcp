package spec

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

//
//nolint:gocognit
func TestParse_postmanCollection(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "postman_petstore.json"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse(postman) failed: %v", err)
	}

	if doc.Version != "2.x" {
		t.Errorf("got version %q, want %q", doc.Version, "2.x")
	}
	if doc.Title != "Petstore API" {
		t.Errorf("got title %q, want %q", doc.Title, "Petstore API")
	}

	if len(doc.PathItems) != 4 {
		t.Fatalf("got %d path items, want 4", len(doc.PathItems))
	}

	var listPets, createPet, getPet, health bool
	for _, pi := range doc.PathItems {
		switch pi.Path {
		case "/v1/pets":
			switch pi.Method {
			case http.MethodGet:
				listPets = true
				op := pi.Operation
				if op.Summary != "List all pets" {
					t.Errorf("got summary %q, want %q", op.Summary, "List all pets")
				}
				if len(op.Parameters) == 0 {
					t.Fatal("expected query params on list pets")
				}
				var hasLimit bool
				for _, p := range op.Parameters {
					if p.Name == "limit" && p.In == "query" {
						hasLimit = true
					}
				}
				if !hasLimit {
					t.Error("expected limit query param")
				}
			case http.MethodPost:
				createPet = true
				op := pi.Operation
				if op.RequestBody == nil {
					t.Fatal("expected request body")
				}
				if op.RequestBody.Content["application/json"] == nil {
					t.Error("expected JSON content type")
				}
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
			if !hasID {
				t.Error("expected petId path param")
			}
		case "/health":
			health = true
		}
	}

	if !listPets {
		t.Error("GET /v1/pets not found")
	}
	if !createPet {
		t.Error("POST /v1/pets not found")
	}
	if !getPet {
		t.Error("GET /v1/pets/{petId} not found")
	}
	if !health {
		t.Error("GET /health not found")
	}
}

func TestParse_postmanHeaders(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "postman_petstore.json"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	for _, pi := range doc.PathItems {
		if pi.Path == "/v1/pets" && pi.Method == http.MethodGet {
			op := pi.Operation
			var hasAuth bool
			for _, p := range op.Parameters {
				if p.Name == "Authorization" && p.In == "header" {
					hasAuth = true
				}
			}
			if !hasAuth {
				t.Error("expected Authorization header")
			}
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
			if got != tt.want {
				t.Errorf("isPostman = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse_openapiServers(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v300_openapi.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Servers) == 0 {
		t.Fatal("expected at least 1 server from openapi servers")
	}
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
			if err != nil {
				t.Fatal(err)
			}
			doc, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse(%s) failed: %v", file, err)
			}
			if len(doc.PathItems) == 0 {
				t.Error("no path items")
			}
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
			if err != nil {
				t.Fatal(err)
			}
			doc, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse(%s) failed: %v", file, err)
			}
			if len(doc.PathItems) == 0 {
				t.Error("no path items")
			}
		})
	}
}

func TestPostmanTag_WithFolders(t *testing.T) {
	t.Parallel()

	tag := postmanTag("Get Pet", []string{"Pets", "Store"})
	if tag != "store" {
		t.Errorf("got %q, want %q", tag, "store")
	}
}

func TestPostmanTag_WithoutFolders(t *testing.T) {
	t.Parallel()

	tag := postmanTag("Get Pet", nil)
	if tag != "get-pet" {
		t.Errorf("got %q, want %q", tag, "get-pet")
	}
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
		if got != tt.want {
			t.Errorf("sanitizePostmanTag(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExtractPathFromURLString_NoScheme(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("example.com/api/v1/pets")
	if path != "/api/v1/pets" {
		t.Errorf("got %q, want %q", path, "/api/v1/pets")
	}
}

func TestExtractPathFromURLString_ColonParam(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("http://example.com/api/v1/pets/:petId")
	if path != "/api/v1/pets/{petId}" {
		t.Errorf("got %q, want %q", path, "/api/v1/pets/{petId}")
	}
}

func TestExtractPathFromURLString_EmptyPath(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("http://example.com")
	if path != "/" {
		t.Errorf("got %q, want %q", path, "/")
	}
}

func TestExtractPathFromURLString_InvalidURL(t *testing.T) {
	t.Parallel()

	path := extractPathFromURLString("http://[invalid]")
	if path != "http://[invalid]" {
		t.Errorf("got %q, want original", path)
	}
}

func TestExtractPostmanPath_NilURL(t *testing.T) {
	t.Parallel()

	path := extractPostmanPath(nil)
	if path != "/" {
		t.Errorf("got %q, want %q", path, "/")
	}
}

func TestExtractPostmanPath_StringURL(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal("http://example.com/api/v1/pets")
	path := extractPostmanPath(raw)
	if path != "/api/v1/pets" {
		t.Errorf("got %q, want %q", path, "/api/v1/pets")
	}
}

func TestExtractPostmanPath_StructuredPath(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal(postmanURL{
		Raw:  "http://example.com/api/v1/pets",
		Path: []json.RawMessage{json.RawMessage(`"api"`), json.RawMessage(`"v1"`), json.RawMessage(`"pets"`)},
	})
	path := extractPostmanPath(raw)
	if path != "/api/v1/pets" {
		t.Errorf("got %q, want %q", path, "/api/v1/pets")
	}
}

func TestExtractPostmanPath_PathVariable(t *testing.T) {
	t.Parallel()

	raw, _ := json.Marshal(postmanURL{
		Raw:  "http://example.com/pets/:petId",
		Path: []json.RawMessage{json.RawMessage(`"pets"`), json.RawMessage(`{"type":"string","value":"petId"}`)},
	})
	path := extractPostmanPath(raw)
	if path != "/pets/{petId}" {
		t.Errorf("got %q, want %q", path, "/pets/{petId}")
	}
}

func TestExtractPostmanPath_InvalidJSON(t *testing.T) {
	t.Parallel()

	path := extractPostmanPath(json.RawMessage("{invalid}"))
	if path != "/" {
		t.Errorf("got %q, want %q", path, "/")
	}
}

func TestAppendPostmanHeaders_Disabled(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	headers := []postmanHeader{
		{Key: "Authorization", Value: "Bearer token", Disabled: true},
	}
	appendPostmanHeaders(headers, op)
	if len(op.Parameters) != 0 {
		t.Errorf("Parameters = %d, want 0 (disabled header should be skipped)", len(op.Parameters))
	}
}

func TestAppendPostmanBody_Nil(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(nil, op, http.MethodPost)
	if op.RequestBody != nil {
		t.Fatal("RequestBody should be nil for nil body")
	}
}

func TestAppendPostmanBody_GetMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodGet)
	if op.RequestBody != nil {
		t.Fatal("RequestBody should be nil for GET")
	}
}

func TestAppendPostmanBody_HeadMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodHead)
	if op.RequestBody != nil {
		t.Fatal("RequestBody should be nil for HEAD")
	}
}

func TestAppendPostmanBody_DeleteMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodDelete)
	if op.RequestBody != nil {
		t.Fatal("RequestBody should be nil for DELETE")
	}
}

func TestAppendPostmanBody_OptionsMethod(t *testing.T) {
	t.Parallel()

	op := &Operation{}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodOptions)
	if op.RequestBody != nil {
		t.Fatal("RequestBody should be nil for OPTIONS")
	}
}

func TestAppendPostmanBody_RawJSON(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: `{"key":"value"}`}, op, http.MethodPost)
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}
	if op.RequestBody.Content["application/json"] == nil {
		t.Error("expected application/json content type")
	}
}

func TestAppendPostmanBody_RawXML(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: "<xml>data</xml>"}, op, http.MethodPost)
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}
	if op.RequestBody.Content["application/xml"] == nil {
		t.Error("expected application/xml content type")
	}
}

func TestAppendPostmanBody_RawText(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "raw", Raw: "plain text"}, op, http.MethodPost)
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}
	if op.RequestBody.Content["text/plain"] == nil {
		t.Error("expected text/plain content type")
	}
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
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}
	mt := op.RequestBody.Content["application/x-www-form-urlencoded"]
	if mt == nil {
		t.Fatal("expected urlencoded content type")
	}
	if mt.Schema.Properties["name"] == nil {
		t.Error("expected name property")
	}
	if mt.Schema.Properties["disabled_field"] != nil {
		t.Error("disabled field should not be present")
	}
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
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}
	mt := op.RequestBody.Content["multipart/form-data"]
	if mt == nil {
		t.Fatal("expected multipart/form-data content type")
	}
	if mt.Schema.Properties["file"].Type != "file" {
		t.Errorf("file type = %q, want %q", mt.Schema.Properties["file"].Type, "file")
	}
	if mt.Schema.Properties["name"].Type != "string" {
		t.Errorf("name type = %q, want %q", mt.Schema.Properties["name"].Type, "string")
	}
	if mt.Schema.Properties["disabled_field"] != nil {
		t.Error("disabled field should not be present")
	}
}

func TestAppendPostmanBody_GraphQL(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "graphql"}, op, http.MethodPost)
	if op.RequestBody == nil {
		t.Fatal("RequestBody is nil")
	}
	if op.RequestBody.Content["application/json"] == nil {
		t.Error("expected application/json content type")
	}
}

func TestAppendPostmanBody_UnknownMode(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanBody(&postmanBody{Mode: "unknown"}, op, http.MethodPost)
	if op.RequestBody != nil {
		t.Fatal("RequestBody should be nil for unknown mode")
	}
}

func TestGuessPostmanContentType_NilBody(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(nil)
	if ct != mediaTypeJSON {
		t.Errorf("got %q, want %q", ct, mediaTypeJSON)
	}
}

func TestGuessPostmanContentType_EmptyRaw(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: ""})
	if ct != mediaTypeJSON {
		t.Errorf("got %q, want %q", ct, mediaTypeJSON)
	}
}

func TestGuessPostmanContentType_Array(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: "[1,2,3]"})
	if ct != mediaTypeJSON {
		t.Errorf("got %q, want %q", ct, mediaTypeJSON)
	}
}

func TestGuessPostmanContentType_XML(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: "<root><item/></root>"})
	if ct != "application/xml" {
		t.Errorf("got %q, want %q", ct, "application/xml")
	}
}

func TestGuessPostmanContentType_Text(t *testing.T) {
	t.Parallel()

	ct := guessPostmanContentType(&postmanBody{Raw: "plain text"})
	if ct != "text/plain" {
		t.Errorf("got %q, want %q", ct, "text/plain")
	}
}

func TestParsePostman_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parsePostman([]byte("{invalid}"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParsePostman_EmptyCollection(t *testing.T) {
	t.Parallel()

	input := `{"info":{"name":"Test","schema":"https://schema.getpostman.com/collection/v2.1.0/collection.json"},"item":[]}`
	_, err := parsePostman([]byte(input))
	if err == nil {
		t.Fatal("expected error for empty collection")
	}
}

func TestIsPostman_InvalidJSON(t *testing.T) {
	t.Parallel()

	if isPostman([]byte("{invalid}")) {
		t.Error("expected false for invalid JSON")
	}
}

func TestIsPostman_NoPostmanSchema(t *testing.T) {
	t.Parallel()

	if isPostman([]byte(`{"info":{"schema":"https://example.com/schema.json"},"item":[{}]}`)) {
		t.Error("expected false for non-postman schema")
	}
}

func TestIsPostman_EmptyItems(t *testing.T) {
	t.Parallel()

	input := `{"info":{"schema":"https://schema.getpostman.com/collection/v2.1.0/collection.json"},"item":[]}`
	if isPostman([]byte(input)) {
		t.Error("expected false for empty items")
	}
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
	if err != nil {
		t.Fatalf("flattenPostmanItems() = %v", err)
	}
	if len(doc.PathItems) != 1 {
		t.Fatalf("PathItems = %d, want 1", len(doc.PathItems))
	}
	if doc.PathItems[0].Operation.Tags[0] != "folder1" {
		t.Errorf("Tag = %q, want %q", doc.PathItems[0].Operation.Tags[0], "folder1")
	}
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
	if pi.Method != http.MethodGet {
		t.Errorf("Method = %q, want %q", pi.Method, http.MethodGet)
	}
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
	if pi.Path != "/" {
		t.Errorf("Path = %q, want %q", pi.Path, "/")
	}
}

func TestAppendPostmanURLParams_InvalidURL(t *testing.T) {
	t.Parallel()

	op := &Operation{Parameters: make([]*Parameter, 0)}
	appendPostmanURLParams(json.RawMessage("{invalid}"), op)
	if len(op.Parameters) != 0 {
		t.Errorf("Parameters = %d, want 0", len(op.Parameters))
	}
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
	if len(op.Parameters) != 1 {
		t.Fatalf("Parameters = %d, want 1", len(op.Parameters))
	}
	if op.Parameters[0].Name != "petId" {
		t.Errorf("Name = %q, want %q", op.Parameters[0].Name, "petId")
	}
	if op.Parameters[0].In != "path" {
		t.Errorf("In = %q, want %q", op.Parameters[0].In, "path")
	}
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
	if len(op.Parameters) != 1 {
		t.Fatalf("Parameters = %d, want 1", len(op.Parameters))
	}
	if op.Parameters[0].Name != "limit" {
		t.Errorf("Name = %q, want %q", op.Parameters[0].Name, "limit")
	}
	if op.Parameters[0].In != "query" {
		t.Errorf("In = %q, want %q", op.Parameters[0].In, "query")
	}
}
