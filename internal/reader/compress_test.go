package reader_test

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

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

func TestReader_CompressFirstOfArray_Empty(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"data":[]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "data",
		Mode:     reader.CompressFirstOfArray,
	})
	require.NoError(t, err)
	require.False(t, result.TooLarge)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Empty(t, arr)
}

func TestReader_CompressFirstOfArray_NotArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":{"id":1}}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "data",
		Mode:     reader.CompressFirstOfArray,
	})
	require.Error(t, err)
}

func TestReader_CompressFirstOfArray_InvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{not json`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{Mode: reader.CompressFirstOfArray})
	require.Error(t, err)
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

func TestReader_CompressSampleArray_Small(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"items":[{"id":1},{"id":2}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "items",
		Mode:      reader.CompressSampleArray,
		ArrayHead: 2,
		ArrayTail: 2,
	})
	require.NoError(t, err)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	compressed, ok := body["compressed"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(2), compressed["length"])
	require.Equal(t, float64(0), compressed["skipped"])
}

func TestReader_CompressSampleArray_NotArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":{"id":1}}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "items",
		Mode:     reader.CompressSampleArray,
	})
	require.Error(t, err)
}

func TestReader_CompressTruncateStrings(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"name":"Very long string value here","nested":{"desc":"Another long text"}}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "nested",
		Mode:      reader.CompressTruncateStrings,
		StringLen: 5,
	})
	require.NoError(t, err)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Anoth...", body["desc"])
}

func TestReader_CompressTruncateStrings_RootArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`["short","a very long string here"]`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode:      reader.CompressTruncateStrings,
		StringLen: 5,
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Len(t, arr, 2)
	require.Equal(t, "short", arr[0])
	require.Equal(t, "a ver...", arr[1])
}

func TestReader_CompressTruncateStrings_ShortString(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`"short"`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode:      reader.CompressTruncateStrings,
		StringLen: 100,
	})
	require.NoError(t, err)
	require.Equal(t, "short", result.Body)
}

func TestReader_CompressKeysOnly(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"user":{"id":1,"name":"Alice","tags":["a","b"]}}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "user",
		Mode:     reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "number", body["id"])
	require.Equal(t, "string", body["name"])
	require.Equal(t, "array", body["tags"])
}

func TestReader_CompressKeysOnly_RootArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[{"id":1,"name":"a"},{"id":2,"name":"b"}]`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode: reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.Equal(t, "array", result.Body)
}

func TestReader_CompressKeysOnly_NonStringKeyError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{Mode: reader.CompressKeysOnly})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"a": "number"}, result.Body)
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

func TestReader_CompressSelectKeys_EmptyKeys(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"users":[{"id":1}]}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "users",
		Mode:     reader.CompressSelectKeys,
	})
	require.ErrorIs(t, err, reader.ErrSelectKeysRequired)
}

func TestReader_Compress_InvalidMode(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		Mode: reader.CompressMode("unknown"),
	})
	require.ErrorIs(t, err, reader.ErrInvalidCompressMode)
}

func TestReader_Compress_JSONPathNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "missing",
		Mode:     reader.CompressFirstOfArray,
	})
	require.ErrorIs(t, err, reader.ErrPathNotFound)
}

func TestReader_Compress_TooLarge(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	large := make([]byte, 500)
	for i := range large {
		large[i] = 'A'
	}
	data := []byte(`{"data":[{"id":1,"text":"` + string(large) + `"}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "data",
		Mode:     reader.CompressFirstOfArray,
		Limit:    50,
	})
	require.NoError(t, err)
	require.True(t, result.TooLarge)
	require.Empty(t, result.Body)
}

func TestReader_Compress_RootWithoutJSONPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"id":1,"name":"short"}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode:      reader.CompressTruncateStrings,
		StringLen: 100,
	})
	require.NoError(t, err)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), body["id"])
}

func TestReader_Compress_BodyAsString(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":"not json"}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "data",
		Mode:      reader.CompressTruncateStrings,
		StringLen: 100,
	})
	require.NoError(t, err)
	require.Equal(t, "not json", result.Body)
}

func TestReader_CompressKeysOnly_Array(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[{"id":1,"name":"a"},{"id":2,"name":"b"}]`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode: reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.Equal(t, "array", result.Body)
}

func TestReader_CompressKeysOnly_ScalarValues(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"s":"x","n":1,"b":true,"nu":null}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode: reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "string", body["s"])
	require.Equal(t, "number", body["n"])
	require.Equal(t, "bool", body["b"])
	require.Equal(t, "null", body["nu"])
}

func TestReader_Compress_ReadFileError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	_, err := r.Compress(filepath.Join(dir, "missing.json"), reader.CompressOptions{
		Mode: reader.CompressFirstOfArray,
	})
	require.ErrorIs(t, err, reader.ErrFileNotFound)
}

func TestReader_CompressSelectKeys_NonObjectArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[1,2,3]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "items",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id"},
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Equal(t, []any{float64(1), float64(2), float64(3)}, arr)
}

func TestReader_CompressSelectKeys_RootObject(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"id":1,"name":"a","secret":"x"}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id", "name"},
	})
	require.NoError(t, err)
	body, ok := result.Body.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), body["id"])
	require.Equal(t, "a", body["name"])
	require.Nil(t, body["secret"])
}

func TestReader_CompressSelectKeys_SkipNestedObject(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"id":1,"meta":{"x":"y"}}]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "items",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id"},
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	first, ok := arr[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), first["id"])
	require.Nil(t, first["meta"])
}

func TestReader_CompressSelectKeys_SkipNestedArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"id":1,"tags":["a","b"]}]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "items",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id"},
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	first, ok := arr[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), first["id"])
	require.Nil(t, first["tags"])
}

func TestReader_CompressTruncateStrings_EmptyObjectArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "items",
		Mode:      reader.CompressTruncateStrings,
		StringLen: 10,
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Empty(t, arr)
}
