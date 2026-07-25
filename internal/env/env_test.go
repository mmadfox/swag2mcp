package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_NoMatch(t *testing.T) {
	t.Parallel()

	result := Parse("plain-text")
	assert.Equal(t, "plain-text", result)
}

func TestParse_EmptyParens(t *testing.T) {
	t.Parallel()

	result := Parse("$()")
	assert.Equal(t, "$()", result)
}

func TestParse_WithSpaces(t *testing.T) {
	t.Parallel()

	result := Parse("$( MY_VAR )")
	assert.Empty(t, result)
}

func TestParse_Resolved(t *testing.T) {
	t.Setenv("TEST_VAR", "resolved-value")
	result := Parse("$(TEST_VAR)")
	assert.Equal(t, "resolved-value", result)
}

func TestParse_Unset(t *testing.T) {
	result := Parse("$(NONEXISTENT_VAR_12345)")
	assert.Empty(t, result)
}

func TestParse_Trimmed(t *testing.T) {
	t.Setenv("TRIMMED", "ok")
	result := Parse("  $(TRIMMED)  ")
	assert.Equal(t, "ok", result)
}

func TestExpandTilde_Unix(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	result := ExpandTilde("~/config.yaml")
	expected := filepath.Join(home, "config.yaml")
	assert.Equal(t, expected, result)
}

func TestExpandTilde_NoMatch(t *testing.T) {
	t.Parallel()

	result := ExpandTilde("no-tilde")
	assert.Equal(t, "no-tilde", result)
}

func TestExpandTilde_AbsolutePath(t *testing.T) {
	t.Parallel()

	result := ExpandTilde("/absolute/path")
	assert.Equal(t, "/absolute/path", result)
}

func TestExpandTilde_RelativePath(t *testing.T) {
	t.Parallel()

	result := ExpandTilde("relative/path")
	assert.Equal(t, "relative/path", result)
}

func TestExpandTilde_WindowsBackslash(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	result := ExpandTilde("~\\config.yaml")
	expected := filepath.Join(home, "config.yaml")
	assert.Equal(t, expected, result)
}
