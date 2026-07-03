package workspace

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.want)
			}
		})
	}
}
