package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestReader_SliceByJSONPath_AdjacentRootArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`[{"id":1},{"id":2}]`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "1"})
	require.NoError(t, err)
	require.Equal(t, "2", slice.NextPath)
	require.Equal(t, "0", slice.PrevPath)
}

func TestReader_CompressKeysOnly_EmptyArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "items",
		Mode:     reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.Equal(t, "array", result.Body)
}

func TestReader_Outline_ScalarRootSchemaHint(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`42`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, "number value", outline.SchemaHint)
}

func TestReader_ValidatePath_StatError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	r := reader.New(dir)
	// Path inside dir but missing file triggers stat not-exist branch already.
	_, err := r.Outline(filepath.Join(dir, "missing.json"), reader.OutlineOptions{})
	require.ErrorIs(t, err, reader.ErrFileNotFound)
}

func TestReader_SliceByJSONPath_EmptyRawContext(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":""}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "a"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
}

func TestReader_SliceByLine_EmptyFragmentContext(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte("\n\n"), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-1"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
}

func TestReader_SliceByJSONPath_LocateLinesNotFound(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	_, err := r.Slice(path, reader.SliceOptions{JSONPath: "a"})
	require.NoError(t, err)
}
