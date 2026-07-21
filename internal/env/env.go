// Package env provides environment variable resolution for $(VAR_NAME) patterns
// and tilde expansion for file paths.
package env

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"path/filepath"
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

// ExpandTilde replaces a leading ~/ or ~\ prefix with the user's home directory.
// Works on both Unix and Windows. Returns the path unchanged if no tilde is found
// or if the home directory cannot be determined.
func ExpandTilde(path string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
