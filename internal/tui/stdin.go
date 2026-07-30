package tui

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bufio"
	"os"
)

// stdinReader is a shared buffered reader for stdin.
// Used after Bubbletea TUI exits to drain any remaining \r\n on Windows.
//
//nolint:gochecknoglobals // shared stdin reader for Windows compatibility
var stdinReader = bufio.NewReader(os.Stdin)
