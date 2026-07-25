package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_CompressKeysOnly_NonEmptyArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"id":1},{"id":2}]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "items",
		Mode:     reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.Equal(t, "array", result.Body)
}

func TestReader_SliceByJSONPath_LocateNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "a"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
}

func TestReader_Outline_DecodeObjectEndError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.Error(t, err)
}

func TestReader_Outline_DecodeArrayEndError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[1,2`), 0o600))

	r := reader.New(dir)
	_, err := r.Outline(path, reader.OutlineOptions{})
	require.Error(t, err)
}

func TestReader_Compress_TruncateEmpty(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":""}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "a",
		Mode:      reader.CompressTruncateStrings,
		StringLen: 10,
	})
	require.NoError(t, err)
	require.Equal(t, "", result.Body)
}

func TestReader_Compress_SelectKeysEmptyArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "items",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id"},
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Empty(t, arr)
}
