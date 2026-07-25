package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func TestRunLs_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	var buf strings.Builder
	err := runLs(tmpDir, "", &buf)
	if err != nil {
		t.Fatalf("runLs() = %v", err)
	}
	if buf.Len() == 0 {
		t.Error("runLs() produced no output")
	}
}

func TestRunLs_WithTags(t *testing.T) {
	tmpDir := t.TempDir()
	ws, _ := workspace.New(tmpDir)
	if err := ws.Init(); err != nil {
		t.Fatalf("Init() = %v", err)
	}
	if err := os.WriteFile(ws.ConfigPath(), []byte("specs: []"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	var buf strings.Builder
	err := runLs(tmpDir, "public,internal", &buf)
	if err != nil {
		t.Fatalf("runLs() = %v", err)
	}
	if buf.Len() == 0 {
		t.Error("runLs() produced no output")
	}
}

func TestRunLs_InvalidPath(t *testing.T) {
	var buf strings.Builder
	err := runLs("/nonexistent/path", "", &buf)
	if err == nil {
		t.Fatal("runLs() expected error, got nil")
	}
}
