/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"context"
	"testing"
)

func TestRunMockServer_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	opts := &mockRootCmdOptions{}
	err := runMockServer(tmpDir, opts, context.Background())
	if err == nil {
		t.Fatal("runMockServer() expected error, got nil")
	}
}
