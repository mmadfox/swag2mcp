package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_SliceByJSONPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{
		"status": "ok",
		"data": [
			{"id": 1, "title": "first"},
			{"id": 2, "title": "second"}
		]
	}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "data.1"})
	require.NoError(t, err)
	require.Equal(t, "object", slice.Context)
	require.True(t, slice.IsComplete)
	require.Equal(t, "data.2", slice.NextPath)
	require.Equal(t, "data.0", slice.PrevPath)
	val, ok := slice.Value.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(2), val["id"])
}

func TestReader_SliceByLine(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte("{\n\t\"status\": \"ok\",\n\t\"data\": [1, 2, 3]\n}\n")
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Line: 1, Around: 1})
	require.NoError(t, err)
	require.Equal(t, 1, slice.Lines[0])
	require.Equal(t, 2, slice.Lines[1])
	require.NotEmpty(t, slice.Fragment)
	require.Equal(t, "object", slice.Context)
}

func TestReader_SlicePathNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{JSONPath: "missing"})
	require.Error(t, err)
}
