package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"DefaultRootName", DefaultRootName, ".swag2mcp"},
		{"DirCache", DirCache, "cache"},
		{"DirSpecs", DirSpecs, "specs"},
		{"DirResponses", DirResponses, "responses"},
		{"DirAuthScripts", DirAuthScripts, "auth_scripts"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.value, "%s = %q, want %q", tt.name, tt.value, tt.want)
		})
	}
}

func TestDefaultResponseMaxAge(t *testing.T) {
	assert.Equal(t, 48*time.Hour, DefaultResponseMaxAge,
		"DefaultResponseMaxAge = %v, want %v", DefaultResponseMaxAge, 48*time.Hour)
}
