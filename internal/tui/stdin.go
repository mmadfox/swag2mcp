/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package tui

import (
	"bufio"
	"os"
)

// stdinReader is a shared buffered reader for stdin.
// Used after Bubbletea TUI exits to drain any remaining \r\n on Windows.
//
//nolint:gochecknoglobals // shared stdin reader for Windows compatibility
var stdinReader = bufio.NewReader(os.Stdin)
