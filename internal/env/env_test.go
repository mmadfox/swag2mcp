package env

import (
	"testing"
)

func TestParse_NoMatch(t *testing.T) {
	t.Parallel()

	result := Parse("plain-text")
	if result != "plain-text" {
		t.Errorf("got %q, want %q", result, "plain-text")
	}
}

func TestParse_EmptyParens(t *testing.T) {
	t.Parallel()

	result := Parse("$()")
	if result != "$()" {
		t.Errorf("got %q, want %q", result, "$()")
	}
}

func TestParse_WithSpaces(t *testing.T) {
	t.Parallel()

	result := Parse("$( MY_VAR )")
	if result != "" {
		t.Errorf("got %q, want empty", result)
	}
}

func TestParse_Resolved(t *testing.T) {
	t.Setenv("TEST_VAR", "resolved-value")
	result := Parse("$(TEST_VAR)")
	if result != "resolved-value" {
		t.Errorf("got %q, want %q", result, "resolved-value")
	}
}

func TestParse_Unset(t *testing.T) {
	result := Parse("$(NONEXISTENT_VAR_12345)")
	if result != "" {
		t.Errorf("got %q, want empty", result)
	}
}

func TestParse_Trimmed(t *testing.T) {
	t.Setenv("TRIMMED", "ok")
	result := Parse("  $(TRIMMED)  ")
	if result != "ok" {
		t.Errorf("got %q, want %q", result, "ok")
	}
}
