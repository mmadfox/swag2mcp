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

func TestReader_Compress_NonJSONValueBodyFallback(t *testing.T) {
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

func TestReader_Compress_TruncateInvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":{not json}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath:  "data",
		Mode:      reader.CompressTruncateStrings,
		StringLen: 10,
	})
	require.Error(t, err)
}

func TestReader_Compress_KeysOnlyInvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":{not json}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "data",
		Mode:     reader.CompressKeysOnly,
	})
	require.Error(t, err)
}

func TestReader_Compress_SelectKeysInvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":{not json}`), 0o600))

	r := reader.New(dir)
	_, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "data",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"id"},
	})
	require.Error(t, err)
}

func TestReader_Outline_CountLinesAndStatError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":1}`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{})
	require.NoError(t, err)
	require.Equal(t, int64(7), outline.Size)
}

func TestReader_Outline_SkipValueErrors(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"1":1}]}`), 0o600))

	r := reader.New(dir)
	// Outline with max depth 1 triggers skipValue for object contents
	_, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 1})
	require.NoError(t, err)
}

func TestReader_SliceByJSONPath_LocateEmpty(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"a":""}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "a"})
	require.NoError(t, err)
	require.Equal(t, "value", slice.Context)
}
