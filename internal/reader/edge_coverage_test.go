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

func TestReader_SliceByJSONPath_AdjacentBadSuffix(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":[1,2]}`), 0o600))

	r := reader.New(dir)
	// Valid sibling path generation also covers numeric conversion.
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "data.1"})
	require.NoError(t, err)
	require.Equal(t, "data.2", slice.NextPath)
}

func TestReader_SliceByJSONPath_InferContextArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"data":[1,2]}`), 0o600))

	r := reader.New(dir)
	slice, err := r.Slice(path, reader.SliceOptions{JSONPath: "data"})
	require.NoError(t, err)
	require.Equal(t, "array", slice.Context)
}

func TestReader_CompressKeysOnly_DeepArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[[{"id":1}]]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath: "items",
		Mode:     reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.Equal(t, "array", result.Body)
}

func TestReader_CompressSelectKeys_DeepObject(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"meta":{"id":1}}]}`), 0o600))

	r := reader.New(dir)
	result, err := r.Compress(path, reader.CompressOptions{
		JSONPath:   "items",
		Mode:       reader.CompressSelectKeys,
		SelectKeys: []string{"meta"},
	})
	require.NoError(t, err)
	arr, ok := result.Body.([]any)
	require.True(t, ok)
	require.Len(t, arr, 1)
}

func TestReader_Outline_SkipNestedObject(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"tags":["a","b"]}]}`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 2})
	require.NoError(t, err)
	require.Equal(t, "array", outline.Structure.Structure["items"].Type)
}

func TestReader_Outline_SkipNestedArray(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	require.NoError(t, os.WriteFile(path, []byte(`{"items":[{"nested":[{"id":1}]}]}`), 0o600))

	r := reader.New(dir)
	outline, err := r.Outline(path, reader.OutlineOptions{MaxDepth: 2})
	require.NoError(t, err)
	require.Empty(t, outline.Structure.Structure["items"].Structure)
}
