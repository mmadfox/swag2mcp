// Package workspace manages the workspace directory layout for swag2mcp.
//
// The workspace is the root directory (default ~/.swag2mcp) that contains
// subdirectories for cached remote specs, local spec files, invocation
// responses, and authentication scripts.
package workspace

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import "time"

const (
	// DefaultRootName is the default workspace directory name under the user's home.
	DefaultRootName = ".swag2mcp"

	// DirCache is the subdirectory for cached remote spec files.
	DirCache = "cache"

	// DirSpecs is the subdirectory for local spec files.
	DirSpecs = "specs"

	// DirResponses is the subdirectory for invocation response files.
	DirResponses = "responses"

	// DirAuthScripts is the subdirectory for authentication scripts.
	DirAuthScripts = "auth_scripts"

	// DefaultResponseMaxAge is the default age after which response files are cleaned up.
	DefaultResponseMaxAge = 48 * time.Hour

	// osWindows is the [runtime.GOOS] value for Windows.
	osWindows = "windows"
)
