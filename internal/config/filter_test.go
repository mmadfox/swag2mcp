package config

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_MatchSpec_Empty(t *testing.T) {
	t.Parallel()

	f := NewFilter(nil)
	assert.True(t, f.MatchSpec("anything"), "MatchSpec() = false, want true for empty filter")
}

func TestFilter_MatchSpec_Match(t *testing.T) {
	t.Parallel()

	f := NewFilter([]string{"public", "demo"})
	assert.True(t, f.MatchSpec("public"))
}

func TestFilter_MatchSpec_NoMatch(t *testing.T) {
	t.Parallel()

	f := NewFilter([]string{"public"})
	assert.False(t, f.MatchSpec("internal"))
}

func TestFilter_MatchSpec_MultipleSpecTags(t *testing.T) {
	t.Parallel()

	f := NewFilter([]string{"public"})
	assert.True(t, f.MatchSpec("internal", "public"))
}
