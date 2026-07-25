package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_Outline(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{
		"status": "ok",
		"data": [
			{"id": 1, "title": "first"},
			{"id": 2, "title": "second"},
			{"id": 3, "title": "third"}
		]
	}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 3, MaxArrayItems: 2})
	require.NoError(t, err)
	require.Equal(t, "object", outline.Type)
	require.Equal(t, int64(len(data)), outline.Size)
	require.Equal(t, []string{"status", "data"}, outline.Keys)
	require.Equal(t, 3, outline.Structure.Structure["data"].ItemCount)
	require.Equal(t, "array", outline.Structure.Structure["data"].Type)

	require.NotEmpty(t, outline.CompressionHints)
}

func TestReader_Outline_DefaultOptions(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
}

func TestReader_ValidatePath_FileNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	_, err := r.Outline(filepath.Join(dir, "missing.json"), reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrFileNotFound)
}

func TestReader_ValidatePath_Directory(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	_, err := r.Outline(dir, reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrFileNotFound)
}

func TestReader_ValidatePath_OutsideResponsesDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	outside := t.TempDir()
	path := filepath.Join(outside, "outside.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrPathNotAllowed)
}

func TestReader_ValidatePath_TraversalEscape(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	parent := filepath.Dir(dir)
	path := filepath.Join(parent, "escape.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrPathNotAllowed)
}

func TestReader_ValidatePath_AbsoluteInput(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	_, err := r.Outline("/etc/passwd", reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrPathNotAllowed)
}

func TestReader_Outline_InvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(path, []byte(`{not json`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrNotJSON)
}

func TestReader_Outline_ScalarTypes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		content  string
		expected string
	}{
		{"string", `"hello"`, "string"},
		{"number", `42`, "number"},
		{"true", `true`, "true"},
		{"false", `false`, "false"},
		{"null", `null`, "null"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "scalar.json")
			require.NoError(t, os.WriteFile(path, []byte(tc.content), 0o600))

			r := reader.New(dir)
			outline, err := r.Outline(path, reader.OutlineOptions{})
			require.NoError(t, err)
			require.Equal(t, tc.expected, outline.Type)
			require.NotEmpty(t, outline.SchemaHint)
		})
	}
}

func TestReader_Outline_TopLevelArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"status":"ok","employees":[{"id":1},{"id":2}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 3, MaxArrayItems: 2})
	require.NoError(t, err)
	require.Equal(t, "object", outline.Type)
	require.NotEmpty(t, outline.CompressionHints)
	require.NotEmpty(t, outline.NavigationHints.ArrayPaths)
	require.Equal(t, "employees", outline.NavigationHints.ArrayPaths[0].Path)
}

func TestReader_Outline_HomogeneousArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[1, 2, 3, 4, 5]`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, "array", outline.Type)
	require.Equal(t, "number", outline.ItemType)
	require.Equal(t, 5, outline.ItemCount)
}

func TestReader_Outline_HeterogeneousArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[1, "two", true]`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, "array", outline.Type)
	require.Empty(t, outline.ItemType)
	require.Len(t, outline.Structure.SampleItems, 3)
}

func TestReader_Outline_EmptyLineCountFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	require.NoError(t, os.WriteFile(path, []byte(`{}`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, outline.LineCount)
}

func TestReader_Outline_NestedObjectArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"departments":[{"name":"eng","employees":[{"id":1}]}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 6, MaxArrayItems: 2})
	require.NoError(t, err)
	require.Equal(t, "object", outline.Type)
	require.NotEmpty(t, outline.NavigationHints.ArrayPaths)
	require.GreaterOrEqual(t, len(outline.NavigationHints.ArrayPaths), 1)
}

func TestReader_Outline_ObjectAtMaxDepth(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"a":{"b":{"c":1}}}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 2})
	require.NoError(t, err)
	require.Equal(t, 2, outline.Depth)
	require.Empty(t, outline.Structure.Structure["a"].Structure)
}

func TestReader_Outline_LongStringValue(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	long := make([]byte, 200)
	for i := range long {
		long[i] = 'x'
	}
	require.NoError(t, os.WriteFile(path, []byte(`{"value":"`+string(long)+`"}`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, "string", outline.Structure.Structure["value"].Type)
	require.Contains(t, outline.Structure.Structure["value"].Value, "...")
}

func TestReader_Outline_TopLevelArrayDepth(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[[1,2],[3,4]]`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, "array", outline.Type)
	require.Equal(t, 2, outline.ItemCount)
}

func TestReader_Outline_InvalidFileKey(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(path, []byte(`{1:"x"}`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.Error(t, err)
}

func TestReader_Outline_NestedArraySkippedByMaxDepth(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"items":[{"id":1,"tags":["a","b"]}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 2})
	require.NoError(t, err)
	require.Empty(t, outline.Structure.Structure["items"].Structure)
}

func TestReader_Outline_ArraySampleLargerThanMax(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"items":[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5},{"id":6}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 3, MaxArrayItems: 2})
	require.NoError(t, err)
	require.Equal(t, 6, outline.Structure.Structure["items"].ItemCount)
	require.Len(t, outline.Structure.Structure["items"].SampleItems, 2)
}

func TestReader_Outline_NavigationForNestedArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"items":[{"id":1,"tags":["a","b"]}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 3, MaxArrayItems: 2})
	require.NoError(t, err)
	require.NotEmpty(t, outline.NavigationHints.ArrayPaths)
}

func TestReader_Slice_ReadFileError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	_, err := r.Slice(filepath.Join(dir, "missing.json"), reader.SliceOptions{
		JSONPath: "x",
	})
	require.ErrorIs(t, err, reader.ErrFileNotFound)
}

func TestReader_Slice_JSONPathNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{JSONPath: "missing"})
	require.ErrorIs(t, err, reader.ErrPathNotFound)
}

func TestReader_ValidatePath_RelativePathInsideDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	require.NoError(t, os.MkdirAll(subdir, 0o700))
	path := filepath.Join(subdir, "file.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
}
