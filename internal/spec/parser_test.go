package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse_validSpecs(t *testing.T) {
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()

		if strings.HasPrefix(name, "test_invalid_") {
			continue
		}
		if strings.HasPrefix(name, "invalid_") {
			continue
		}

		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", name))
			if err != nil {
				t.Fatal(err)
			}

			doc, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse(%s) failed: %v", name, err)
			}

			if doc.Version == "" {
				t.Error("version is empty")
			}
		})
	}
}

func TestParse_invalidSpecs_structural(t *testing.T) {
	// These files are structurally malformed and should fail to parse.
	files := []string{
		"invalid_v_empty.yaml",
		"invalid_v_as_number.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", file))
			if err != nil {
				t.Fatal(err)
			}
			_, err = Parse(data)
			if err == nil {
				t.Error("expected parse error, got nil")
			}
		})
	}
}

func TestParse_invalidSpecs_semantic(t *testing.T) {
	// These files are structurally valid YAML/JSON but semantically invalid.
	// Our parser is lenient and only parses, so they should succeed.
	files := []string{
		"valid_v20_swagger.yaml",  // YAML version of valid spec
		"valid_v311_openapi.yaml", // 3.1.1 with items: false
		"invalid_v_304.yaml",      // openapi 3.0.4 (valid YAML, non-standard minor)
		"invalid_v_conflict.yaml", // both swagger and openapi
		"test_invalid_21_duplicate_tag_names.yaml",
		"test_invalid_22_undefined_tag_in_operation.yaml",
		"test_invalid_23_operation_without_responses.yaml",
		"test_invalid_24_empty_operation.yaml",
		"test_invalid_25_null_values.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", file))
			if err != nil {
				t.Fatal(err)
			}
			_, err = Parse(data)
			if err != nil {
				t.Fatalf("Parse(%s) should have succeeded (lenient parser): %v", file, err)
			}
		})
	}
}

func TestParse_versionDetection(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantVer string
	}{
		{"swagger 2.0", "valid_v20_swagger.yaml", "2.0"},
		{"openapi 3.0.0", "valid_v300_openapi.yaml", "3.0.0"},
		{"openapi 3.0.1", "valid_v301_openapi.yaml", "3.0.1"},
		{"openapi 3.0.2", "valid_v302_openapi.yaml", "3.0.2"},
		{"openapi 3.0.3", "valid_v303_openapi.yaml", "3.0.3"},
		{"openapi 3.1.0", "valid_v310_openapi.yaml", "3.1.0"},
		{"openapi 3.1.1", "valid_v311_openapi.yaml", "3.1.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tt.file))
			if err != nil {
				t.Fatal(err)
			}

			doc, err := Parse(data)
			if err != nil {
				t.Fatalf("Parse(%s) failed: %v", tt.file, err)
			}

			if doc.Version != tt.wantVer {
				t.Errorf("got version %q, want %q", doc.Version, tt.wantVer)
			}
		})
	}
}

func TestParse_swaggerHost(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	if len(doc.Servers) == 0 {
		t.Fatal("expected at least 1 server from swagger host")
	}

	if doc.Servers[0].URL != "https://api.example.com/v1" {
		t.Errorf("got server URL %q, want %q", doc.Servers[0].URL, "https://api.example.com/v1")
	}
}

func TestParse_openapiServers(t *testing.T) {
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

func TestParse_operationMetadata(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == "GET" {
			found = true
			op := pi.Operation
			if op.Summary != "Список пользователей" {
				t.Errorf("got summary %q, want %q", op.Summary, "List users")
			}
			if len(op.Parameters) == 0 {
				t.Fatal("expected parameters")
			}
			if op.Parameters[0].Name != "limit" {
				t.Errorf("got param name %q, want %q", op.Parameters[0].Name, "limit")
			}
			if op.Parameters[0].In != "query" {
				t.Errorf("got param in %q, want %q", op.Parameters[0].In, "query")
			}
			break
		}
	}
	if !found {
		t.Fatal("GET /users not found in parsed doc")
	}
}

func TestParse_requestBody(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == "POST" {
			found = true
			op := pi.Operation
			if op.RequestBody == nil {
				t.Fatal("expected request body")
			}
			if !op.RequestBody.Required {
				t.Error("expected required request body")
			}
			if op.RequestBody.Content == nil {
				t.Fatal("expected request body content")
			}
			break
		}
	}
	if !found {
		t.Fatal("POST /users not found")
	}
}

func TestParse_responses(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/users" && pi.Method == "GET" {
			found = true
			op := pi.Operation
			if len(op.Responses) == 0 {
				t.Fatal("expected responses")
			}
			resp, ok := op.Responses["200"]
			if !ok {
				t.Fatal("expected 200 response")
			}
			if resp.Description != "OK" {
				t.Errorf("got description %q, want %q", resp.Description, "OK")
			}
			break
		}
	}
	if !found {
		t.Fatal("GET /users not found")
	}
}

func TestParse_tagHierarchies(t *testing.T) {
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

func TestParse_emptyDoc(t *testing.T) {
	_, err := Parse([]byte{})
	if err == nil {
		t.Error("expected error for empty document")
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		file    string
		wantPre string
	}{
		{"valid_v20_swagger.yaml", "2."},
		{"valid_v300_openapi.yaml", "3."},
		{"valid_v301_openapi.yaml", "3."},
		{"valid_v302_openapi.yaml", "3."},
		{"valid_v303_openapi.yaml", "3."},
		{"valid_v310_openapi.yaml", "3."},
		{"valid_v311_openapi.yaml", "3."},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tt.file))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := Parse(data)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(doc.Version, tt.wantPre) {
				t.Errorf("got version %q, want %q prefix", doc.Version, tt.wantPre)
			}
		})
	}
}

func TestParse_swaggerFileUpload(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "valid_v20_swagger.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, pi := range doc.PathItems {
		if pi.Path == "/files/upload" && pi.Method == "POST" {
			found = true
			op := pi.Operation
			if len(op.Parameters) == 0 {
				t.Fatal("expected parameters")
			}
			if op.Parameters[0].Name != "file" {
				t.Errorf("got param name %q, want %q", op.Parameters[0].Name, "file")
			}
			if op.Parameters[0].In != "formData" {
				t.Errorf("got param in %q, want %q", op.Parameters[0].In, "formData")
			}
			break
		}
	}
	if !found {
		t.Fatal("POST /files/upload not found")
	}
}

func TestParse_postmanCollection(t *testing.T) {
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

	// Check specific endpoints
	var listPets, createPet, getPet, health bool
	for _, pi := range doc.PathItems {
		switch pi.Path {
		case "/v1/pets":
			switch pi.Method {
			case "GET":
				listPets = true
				op := pi.Operation
				if op.Summary != "List all pets" {
					t.Errorf("got summary %q, want %q", op.Summary, "List all pets")
				}
				// query params
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
			case "POST":
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
	data, err := os.ReadFile(filepath.Join("testdata", "postman_petstore.json"))
	if err != nil {
		t.Fatal(err)
	}

	doc, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}

	for _, pi := range doc.PathItems {
		if pi.Path == "/v1/pets" && pi.Method == "GET" {
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
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "postman collection",
			data: []byte(`{"info":{"schema":"https://schema.getpostman.com/collection/v2.1.0/collection.json"},"item":[{"name":"test","request":{"method":"GET","url":"http://example.com"}}]}`),
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
			got := isPostman(tt.data)
			if got != tt.want {
				t.Errorf("isPostman = %v, want %v", got, tt.want)
			}
		})
	}
}
