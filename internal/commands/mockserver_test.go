package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"
)

func TestRunMockServer_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	opts := &mockServerCmdOptions{}
	err := runMockServer(tmpDir, opts, context.Background())
	if err == nil {
		t.Fatal("runMockServer() expected error, got nil")
	}
}
