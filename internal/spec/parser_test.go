package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParse_validSpecs(t *testing.T) {
	t.Parallel()
	entries, dirErr := os.ReadDir("testdata")
	if dirErr != nil {
		t.Fatal(dirErr)
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
			t.Parallel()
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
	t.Parallel()
	files := []string{
		"invalid_v_empty.yaml",
		"invalid_v_as_number.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
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
	t.Parallel()
	files := []string{
		"valid_v20_swagger.yaml",
		"valid_v311_openapi.yaml",
		"invalid_v_304.yaml",
		"invalid_v_conflict.yaml",
		"test_invalid_21_duplicate_tag_names.yaml",
		"test_invalid_22_undefined_tag_in_operation.yaml",
		"test_invalid_23_operation_without_responses.yaml",
		"test_invalid_24_empty_operation.yaml",
		"test_invalid_25_null_values.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
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

func TestParse_emptyDoc(t *testing.T) {
	t.Parallel()
	_, err := Parse([]byte{})
	if err == nil {
		t.Error("expected error for empty document")
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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

func TestToJSON_Empty(t *testing.T) {
	t.Parallel()

	_, err := toJSON([]byte{})
	if err == nil {
		t.Fatal("expected error for empty data")
	}
}

func TestToJSON_InvalidYAML(t *testing.T) {
	t.Parallel()

	_, err := toJSON([]byte("invalid: [yaml: broken"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestToJSON_JSON(t *testing.T) {
	t.Parallel()

	data, err := toJSON([]byte(`{"openapi":"3.0.0","info":{"title":"Test"}}`))
	if err != nil {
		t.Fatalf("toJSON() = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("data is empty")
	}
}

func TestToJSON_YAML(t *testing.T) {
	t.Parallel()

	data, err := toJSON([]byte("openapi: 3.0.0\ninfo:\n  title: Test\n"))
	if err != nil {
		t.Fatalf("toJSON() = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("data is empty")
	}
}

func TestPreprocessV3_InvalidJSON(t *testing.T) {
	t.Parallel()

	result := preprocessV3([]byte("{invalid json"))
	if string(result) != "{invalid json" {
		t.Error("expected original data to be returned unchanged")
	}
}

func TestPreprocessV3_ItemsFalse(t *testing.T) {
	t.Parallel()

	input := []byte(`{"openapi":"3.1.0","paths":{"/test":{"get":{"parameters":[{"schema":{"items":false}}]}}}}`)
	result := preprocessV3(input)
	if !strings.Contains(string(result), `"items":{}`) {
		t.Error("expected items:false to be replaced with items:{}")
	}
}

func TestPreprocessV3_ItemsTrue(t *testing.T) {
	t.Parallel()

	input := []byte(`{"openapi":"3.1.0","paths":{"/test":{"get":{"parameters":[{"schema":{"items":true}}]}}}}`)
	result := preprocessV3(input)
	if !strings.Contains(string(result), `"items":{}`) {
		t.Error("expected items:true to be replaced with items:{}")
	}
}
