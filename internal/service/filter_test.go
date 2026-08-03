/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

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

func TestExtFromLocation_urlWithExt(t *testing.T) {
	t.Parallel()

	ext := extFromLocation("https://example.com/api/v3/doc.json")
	require.Equal(t, ".json", ext)
}

func TestExtFromLocation_urlNoExt(t *testing.T) {
	t.Parallel()

	ext := extFromLocation("https://example.com/")
	require.Empty(t, ext)
}

func TestExtFromLocation_local(t *testing.T) {
	t.Parallel()

	ext := extFromLocation("/path/to/spec.yaml")
	require.Empty(t, ext)
}

func TestSpecFileName_titleNoExt_urlHasExt(t *testing.T) {
	t.Parallel()

	name := specFileName("nspd", "nspd", "https://example.com/api/actors/swagger/v2/doc.json", "")
	require.Contains(t, name, "nspd")
	require.Contains(t, name, ".json")
}

func TestSpecFileName_titleNoExt_urlNoExt(t *testing.T) {
	t.Parallel()

	name := specFileName("dom", "myspec", "https://example.com/", "")
	require.NotContains(t, name, ".yaml")
	require.NotContains(t, name, ".json")
}

func TestSpecFileName_titleNoExt_localPathHasExt(t *testing.T) {
	t.Parallel()

	name := specFileName("nspd1", "#2", "/Users/mmadfox/tmp/test.json", "")
	require.Contains(t, name, "nspd1")
	require.Contains(t, name, ".json")
	require.Contains(t, name, "#2")
}

func TestSpecFileName_withTitle(t *testing.T) {
	t.Parallel()

	name := specFileName("mydomain", "My Collection", "https://example.com/spec.yaml", "")
	require.Contains(t, name, "mydomain")
	require.Contains(t, name, "my-collection")
}

func TestSpecFileName_withoutTitle(t *testing.T) {
	t.Parallel()

	name := specFileName("mydomain", "", "https://example.com/spec.yaml", "")
	require.Contains(t, name, "mydomain")
	require.Contains(t, name, ".yaml")
}

func TestSpecFileName_withPathPart(t *testing.T) {
	t.Parallel()

	name := specFileName("nspd", "Actors", "https://example.com/api/actors/swagger/v2/doc.json", "api-actors-swagger-v2")
	require.Contains(t, name, "nspd")
	require.Contains(t, name, "actors")
	require.Contains(t, name, "api-actors-swagger-v2")
}

func TestSpecFileName_withPathPartNoTitle(t *testing.T) {
	t.Parallel()

	name := specFileName("nspd", "", "https://example.com/api/actors/swagger/v2/doc.json", "api-actors-swagger-v2")
	require.Contains(t, name, "nspd")
	require.Contains(t, name, "doc")
	require.Contains(t, name, "api-actors-swagger-v2")
}

func TestSpecFileName_truncatesLongTitle(t *testing.T) {
	t.Parallel()

	name := specFileName("dom", "LLM Title API Agreements Swagger V1", "https://example.com/spec.yaml", "")
	require.Contains(t, name, "dom")
	require.Contains(t, name, "llm-title-api")
	require.NotContains(t, name, "agreements")
}

func TestSpecFileName_truncatesLongPathPart(t *testing.T) {
	t.Parallel()

	name := specFileName("dom", "col", "https://example.com/api/agreements/swagger/v1/doc.json", "api-agreements-swagger-v1-extra-long-path-segment")
	require.Contains(t, name, "dom")
	require.Contains(t, name, "col")
	require.Contains(t, name, ".json")
	require.NotContains(t, name, "segment")
}

func TestTruncateByLastHyphen_underLimit(t *testing.T) {
	t.Parallel()

	result := truncateByLastHyphen("short", 20)
	require.Equal(t, "short", result)
}

func TestTruncateByLastHyphen_overLimit(t *testing.T) {
	t.Parallel()

	result := truncateByLastHyphen("llm-title-api-agreements-swagger-v1", 20)
	require.Equal(t, "llm-title-api", result)
}

func TestTruncateByLastHyphen_noHyphen(t *testing.T) {
	t.Parallel()

	result := truncateByLastHyphen("abcdefghijklmnopqrstuvwxyz", 10)
	require.Equal(t, "abcdefghij", result)
}

func TestPathPartFromLocation_url(t *testing.T) {
	t.Parallel()

	part := pathPartFromLocation("https://example.com/api/actors/swagger/v2/doc.json")
	require.Equal(t, "api-actors-swagger-v2", part)
}

func TestPathPartFromLocation_root(t *testing.T) {
	t.Parallel()

	part := pathPartFromLocation("https://example.com/spec.yaml")
	require.Empty(t, part)
}

func TestPathPartFromLocation_local(t *testing.T) {
	t.Parallel()

	part := pathPartFromLocation("/path/to/spec.yaml")
	require.Empty(t, part)
}

func TestPathPartFromLocation_empty(t *testing.T) {
	t.Parallel()

	part := pathPartFromLocation("")
	require.Empty(t, part)
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
