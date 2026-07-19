package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_CompressFirstOfArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"data":[{"id":1,"title":"first"},{"id":2,"title":"second"}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "data",
		Mode:     reader.CompressFirstOfArray,
	})
	require.NoError(t, err)
	require.False(t, result.TooLarge)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	compressed, ok := body["compressed"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "array", compressed["type"])
	require.Equal(t, float64(1), compressed["length"])
}

func TestReader_CompressSampleArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"items":[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "items",
		Mode:      reader.CompressSampleArray,
		ArrayHead: 2,
		ArrayTail: 1,
	})
	require.NoError(t, err)
	require.False(t, result.TooLarge)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	compressed, ok := body["compressed"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(3), compressed["length"])
	require.Equal(t, float64(5), compressed["original"])
}

func TestReader_CompressSelectKeys(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"users":[` +
		`{"id":1,"name":"Alice","email":"a@example.com"},` +
		`{"id":2,"name":"Bob","email":"b@example.com"}` +
		`]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "users",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id", "name"},
	})
	require.NoError(t, err)
	require.False(t, result.TooLarge)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Len(t, arr, 2)
	first, ok := arr[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), first["id"])
	require.Equal(t, "Alice", first["name"])
	require.Nil(t, first["email"])
}
