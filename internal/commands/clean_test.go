/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestRunClean_EmptyWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}

	result, err := runClean(tmpDir)
	if err != nil {
		t.Fatalf("runClean() = %v", err)
	}
	if result.cache != cleanStatusCleared {
		t.Errorf("cache = %q, want 'cleared'", result.cache)
	}
	if result.responses != cleanStatusCleared {
		t.Errorf("responses = %q, want 'cleared'", result.responses)
	}
}

func TestRunClean_InvalidPath(t *testing.T) {
	result, err := runClean("/nonexistent/path")
	if err != nil {
		t.Fatalf("runClean() = %v", err)
	}
	if result.cache != cleanStatusCleared {
		t.Errorf("cache = %q, want 'cleared'", result.cache)
	}
	if result.responses != cleanStatusCleared {
		t.Errorf("responses = %q, want 'cleared'", result.responses)
	}
}
