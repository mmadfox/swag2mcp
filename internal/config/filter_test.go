package config

import (
	"testing"
)

func TestFilter_MatchSpec_Empty(t *testing.T) {
	t.Parallel()

	f := NewFilter(nil)
	if !f.MatchSpec("anything") {
		t.Error("MatchSpec() = false, want true for empty filter")
	}
}

func TestFilter_MatchSpec_Match(t *testing.T) {
	t.Parallel()

	f := NewFilter([]string{"public", "demo"})
	if !f.MatchSpec("public") {
		t.Error("MatchSpec(public) = false, want true")
	}
}

func TestFilter_MatchSpec_NoMatch(t *testing.T) {
	t.Parallel()

	f := NewFilter([]string{"public"})
	if f.MatchSpec("internal") {
		t.Error("MatchSpec(internal) = true, want false")
	}
}

func TestFilter_MatchSpec_MultipleSpecTags(t *testing.T) {
	t.Parallel()

	f := NewFilter([]string{"public"})
	if !f.MatchSpec("internal", "public") {
		t.Error("MatchSpec(internal, public) = false, want true")
	}
}
