/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package reader_test

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

func TestNormalizeJSONPath_rootArrayLast(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"id":1},{"id":2},{"id":3}]`)
	require.Equal(t, "2", reader.NormalizeJSONPath(data, "-"))
}

func TestNormalizeJSONPath_nestedArrayLast(t *testing.T) {
	t.Parallel()

	data := []byte(`{"data":[10,20,30]}`)
	require.Equal(t, "data.2", reader.NormalizeJSONPath(data, "data.-"))
}

func TestNormalizeJSONPath_noTrailingDash(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"id":1},{"id":2}]`)
	require.Equal(t, "0", reader.NormalizeJSONPath(data, "0"))
	require.Equal(t, "data.0", reader.NormalizeJSONPath(data, "data.0"))
	require.Equal(t, "data.name", reader.NormalizeJSONPath(data, "data.name"))
}

func TestNormalizeJSONPath_empty(t *testing.T) {
	t.Parallel()

	data := []byte(`[1,2,3]`)
	require.Equal(t, "", reader.NormalizeJSONPath(data, ""))
}

func TestNormalizeJSONPath_notArray(t *testing.T) {
	t.Parallel()

	data := []byte(`{"a":1}`)
	require.Equal(t, "-", reader.NormalizeJSONPath(data, "-"))
	require.Equal(t, "a.-", reader.NormalizeJSONPath(data, "a.-"))
}

func TestNormalizeJSONPath_emptyArray(t *testing.T) {
	t.Parallel()

	data := []byte(`[]`)
	require.Equal(t, "-", reader.NormalizeJSONPath(data, "-"))
}

func TestNormalizeJSONPath_middleSegment(t *testing.T) {
	t.Parallel()

	data := []byte(`{"items":[{"name":"a"},{"name":"b"}]}`)
	require.Equal(t, "items.1.name", reader.NormalizeJSONPath(data, "items.-.name"))
}

func TestNormalizeJSONPath_multipleSegments(t *testing.T) {
	t.Parallel()

	data := []byte(`{"a":[{"b":[10,20,30]}]}`)
	require.Equal(t, "a.0.b.2", reader.NormalizeJSONPath(data, "a.-.b.-"))
}

func TestNormalizeJSONPath_middleSegmentNotArray(t *testing.T) {
	t.Parallel()

	data := []byte(`{"items":[{"name":"a"}]}`)
	// "items.0" is an object, not an array, so "-" after it cannot be resolved
	// and the path is left unchanged.
	require.Equal(t, "items.0.-", reader.NormalizeJSONPath(data, "items.0.-"))
}
