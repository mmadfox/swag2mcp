package spec

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_validSpecs(t *testing.T) {
	t.Parallel()
	entries, dirErr := os.ReadDir("testdata")
	require.NoError(t, dirErr)

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
			require.NoError(t, err)

			doc, err := Parse(data)
			require.NoError(t, err, "Parse(%s) failed", name)

			assert.NotEmpty(t, doc.Version, "version is empty")
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
			require.NoError(t, err)
			_, err = Parse(data)
			require.Error(t, err, "expected parse error, got nil")
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
			require.NoError(t, err)
			_, err = Parse(data)
			require.NoError(t, err, "Parse(%s) should have succeeded (lenient parser)", file)
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
			require.NoError(t, err)

			doc, err := Parse(data)
			require.NoError(t, err, "Parse(%s) failed", tt.file)

			assert.Equal(t, tt.wantVer, doc.Version)
		})
	}
}

func TestParse_emptyDoc(t *testing.T) {
	t.Parallel()
	_, err := Parse([]byte{})
	require.Error(t, err, "expected error for empty document")
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
			require.NoError(t, err)
			doc, err := Parse(data)
			require.NoError(t, err)
			assert.True(t,
				strings.HasPrefix(doc.Version, tt.wantPre),
				"got version %q, want %q prefix", doc.Version, tt.wantPre)
		})
	}
}

func TestToJSON_Empty(t *testing.T) {
	t.Parallel()

	_, err := toJSON([]byte{})
	require.Error(t, err, "expected error for empty data")
}

func TestToJSON_InvalidYAML(t *testing.T) {
	t.Parallel()

	_, err := toJSON([]byte("invalid: [yaml: broken"))
	require.Error(t, err, "expected error for invalid YAML")
}

func TestToJSON_JSON(t *testing.T) {
	t.Parallel()

	data, err := toJSON([]byte(`{"openapi":"3.0.0","info":{"title":"Test"}}`))
	require.NoError(t, err, "toJSON() failed")
	require.NotEmpty(t, data, "data is empty")
}

func TestToJSON_YAML(t *testing.T) {
	t.Parallel()

	data, err := toJSON([]byte("openapi: 3.0.0\ninfo:\n  title: Test\n"))
	require.NoError(t, err, "toJSON() failed")
	require.NotEmpty(t, data, "data is empty")
}

func TestPreprocessV3_InvalidJSON(t *testing.T) {
	t.Parallel()

	result := preprocessV3([]byte("{invalid json"))
	assert.Equal(t, "{invalid json", string(result), "expected original data to be returned unchanged")
}

func TestPreprocessV3_ItemsFalse(t *testing.T) {
	t.Parallel()

	input := []byte(`{"openapi":"3.1.0","paths":{"/test":{"get":{"parameters":[{"schema":{"items":false}}]}}}}`)
	result := preprocessV3(input)
	assert.Contains(t, string(result), `"items":{}`, "expected items:false to be replaced with items:{}")
}

func TestPreprocessV3_ItemsTrue(t *testing.T) {
	t.Parallel()

	input := []byte(`{"openapi":"3.1.0","paths":{"/test":{"get":{"parameters":[{"schema":{"items":true}}]}}}}`)
	result := preprocessV3(input)
	assert.Contains(t, string(result), `"items":{}`, "expected items:true to be replaced with items:{}")
}
