/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

// Package reader provides streaming and path-based access to large JSON
// response files stored in the workspace. It is used by the response_* MCP
// tools to build outlines, compress arrays, and extract fragments without
// loading entire files into memory.
package reader
