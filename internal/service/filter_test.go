package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeFilter_empty(t *testing.T) {
	t.Parallel()

	f := makeFilter(nil)
	require.True(t, f.match("anything"))
}

func TestMakeFilter_match(t *testing.T) {
	t.Parallel()

	f := makeFilter([]string{"alpha", "beta"})
	require.True(t, f.match("alpha"))
	require.True(t, f.match("beta"))
	require.False(t, f.match("gamma"))
}

func TestSpecFileName_withTitle(t *testing.T) {
	t.Parallel()

	name := specFileName("mydomain", "My Collection", "https://example.com/spec.yaml")
	require.Contains(t, name, "mydomain")
	require.Contains(t, name, "my-collection")
}

func TestSpecFileName_withoutTitle(t *testing.T) {
	t.Parallel()

	name := specFileName("mydomain", "", "https://example.com/spec.yaml")
	require.Contains(t, name, "mydomain")
	require.Contains(t, name, ".yaml")
}

func TestSpecFileNameBase_url(t *testing.T) {
	t.Parallel()

	base := specFileNameBase("https://example.com/api/v3/swagger.json")
	require.Equal(t, "swagger.json", base)
}

func TestSpecFileNameBase_local(t *testing.T) {
	t.Parallel()

	base := specFileNameBase("/path/to/spec.yaml")
	require.Equal(t, "spec.yaml", base)
}

func TestSpecFileNameBase_default(t *testing.T) {
	t.Parallel()

	base := specFileNameBase("https://example.com/")
	require.Equal(t, "spec", base)
}

func TestRemoveDiacritics(t *testing.T) {
	t.Parallel()

	result := removeDiacritics("café naïve")
	require.Equal(t, "cafe naive", result)
}
