// Package env provides environment variable resolution for $(VAR_NAME) patterns.
package env

import (
	"os"
	"strings"
)

// Parse checks if s matches the pattern $(VARNAME) with optional
// whitespace inside the parentheses. If it matches, the variable name is
// extracted and looked up via [os.Getenv]. Otherwise s is returned unchanged.
func Parse(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "$(") || !strings.HasSuffix(s, ")") {
		return s
	}
	inner := s[2 : len(s)-1]
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return s
	}
	return os.Getenv(inner)
}
