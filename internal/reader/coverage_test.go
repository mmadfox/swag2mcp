/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package reader_test

import (
	"os"
	"path/filepath"
	"strings"
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

// TestReader_Outline_LargeSingleLineFile ensures countLines handles files
// with lines exceeding the default [bufio.Scanner] 64KB limit.
func TestReader_Outline_LargeSingleLineFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "large.json")

	// Create a 1MB single-line JSON file (no newlines).
	payload := strings.Repeat("x", 1024*1024)
	content := `{"data":"` + payload + `"}`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, 1, outline.LineCount)
	require.Equal(t, "object", outline.Type)
}

// TestReader_SliceByLine_LargeFile ensures sliceByLines handles files
// with lines exceeding the default [bufio.Scanner] 64KB limit.
func TestReader_SliceByLine_LargeLineFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "large.json")

	// Create a file with a single 1MB line.
	payload := strings.Repeat("x", 1024*1024)
	content := `{"data":"` + payload + `"}`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-1"})
	require.NoError(t, err)
	require.Equal(t, [2]int{1, 1}, slice.Lines)
	require.Contains(t, slice.Fragment, `"data"`)
}

// TestReader_SliceByLine_LargeFile_MultipleLongLines ensures sliceByLines
// handles a file with multiple lines each exceeding 64KB.
func TestReader_SliceByLine_MultipleLongLines(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "large.json")

	var b strings.Builder
	for i := range 5 {
		b.WriteString(`{"id":` + string(rune('0'+i)) + `,"data":"`)
		b.WriteString(strings.Repeat("x", 200*1024))
		b.WriteString(`"}` + "\n")
	}
	require.NoError(t, os.WriteFile(path, []byte(b.String()), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "2-4"})
	require.NoError(t, err)
	require.Equal(t, [2]int{2, 4}, slice.Lines)
	require.Contains(t, slice.Fragment, `"data"`)
}

// TestReader_SliceByLine_SingleLongLine_Range ensures sliceByLines works
// on a range within a single huge line.
func TestReader_SliceByLine_SingleHugeLineRange(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "huge.json")

	// Single line, 10MB of data, no newlines.
	payload := strings.Repeat("z", 10*1024*1024)
	require.NoError(t, os.WriteFile(path, []byte(payload), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{Range: "1-1"})
	require.NoError(t, err)
	require.Equal(t, [2]int{1, 1}, slice.Lines)
}
