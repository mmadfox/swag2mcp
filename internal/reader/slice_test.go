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

func TestReader_SliceByJSONPath_LimitExceeded(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"data":[{"id":1,"text":"` + string(make([]byte, 500)) + `"}]}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{
		JSONPath: "data.0",
		Limit:    50,
	})
	require.ErrorIs(t, err, reader.ErrPathNotFound)
}

func TestReader_SliceByJSONPath_AdjacentPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":[1,2,3]}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "data.0"})
	require.NoError(t, err)
	require.Equal(t, "data.1", slice.NextPath)
	require.Equal(t, "data.-1", slice.PrevPath)
}

func TestReader_SliceByJSONPath_AdjacentMiddle(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[10,20,30]}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "items.1"})
	require.NoError(t, err)
	require.Equal(t, "items.2", slice.NextPath)
	require.Equal(t, "items.0", slice.PrevPath)
}

func TestReader_SliceByJSONPath_AdjacentInvalid(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":"b"}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "a"})
	require.NoError(t, err)
	require.Empty(t, slice.NextPath)
	require.Empty(t, slice.PrevPath)
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

func TestReader_SliceByRange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte("[1, 2, 3]\n"), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-1"})
	require.NoError(t, err)
	require.Equal(t, 1, slice.Lines[0])
	require.Equal(t, 1, slice.Lines[1])
	require.Equal(t, "array", slice.Context)
	require.True(t, slice.IsComplete)
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

func TestReader_Slice_InvalidRange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{Range: "invalid"})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)

	_, err = r.Slice(path, reader.SliceOptions{Range: "5-1"})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_Slice_NoSelector(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_Slice_OutOfRange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{Line: 100, Around: 5})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_SliceByJSONPath_Scalar(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"count":42}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "count"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
	require.Equal(t, float64(42), slice.Value)
	require.Empty(t, slice.NextPath)
	require.Empty(t, slice.PrevPath)
}

func TestReader_SliceByJSONPath_NestedArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"nested":[{"id":1}]}]}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "items.0.nested.0.id"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
	require.Equal(t, float64(1), slice.Value)
}

func TestReader_SliceByLine_LimitExceeded(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte(`{"text":"` + string(make([]byte, 500)) + `"}`)
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{Range: "1-1", Limit: 50})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_SliceByLine_RangeWithInvalidStart(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{Range: "x-1"})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_SliceByLine_RangeWithInvalidEnd(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{Range: "1-x"})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_SliceByLine_RangeZeroStart(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{Range: "0-1"})
	require.ErrorIs(t, err, reader.ErrInvalidLineRange)
}

func TestReader_SliceByLine_WhitespaceContext(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`   {"a":1}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-1"})
	require.NoError(t, err)
	require.Equal(t, "object", slice.Context)
}

func TestReader_SliceByLine_JSONFragment(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	data := []byte("[1,\n2,\n3]\n")
	require.NoError(t, os.WriteFile(path, data, 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-3"})
	require.NoError(t, err)
	require.True(t, slice.IsComplete)
	arr, ok := slice.Value.([]any)
	require.True(t, ok)
	require.Len(t, arr, 3)
}

func TestReader_SliceByLine_InvalidJSONFragment(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte("[1,\n"), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-2"})
	require.NoError(t, err)
	require.False(t, slice.IsComplete)
	require.Nil(t, slice.Value)
}

func TestReader_SliceByJSONPath_EmptyResult(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":""}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "a"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
	require.Equal(t, "", slice.Value)
}

func TestReader_SliceByLine_OpenFileError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	_, err := r.Slice(filepath.Join(dir, "missing.json"), reader.SliceOptions{Range: "1-1"})
	require.ErrorIs(t, err, reader.ErrFileNotFound)
}
