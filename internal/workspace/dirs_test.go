package workspace

import (
	"testing"
	"time"
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
			if tt.value != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.want)
			}
		})
	}
}

func TestDefaultResponseMaxAge(t *testing.T) {
	if DefaultResponseMaxAge != 48*time.Hour {
		t.Errorf("DefaultResponseMaxAge = %v, want %v", DefaultResponseMaxAge, 48*time.Hour)
	}
}
