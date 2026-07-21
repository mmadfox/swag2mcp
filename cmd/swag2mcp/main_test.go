package main

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"os"
	"testing"
)

func TestRunMain_Help(t *testing.T) {
	os.Args = []string{"swag2mcp", "--help"} //nolint:reassign // test needs to set args
	err := runMain()
	if err != nil {
		t.Fatalf("runMain() with --help returned error: %v", err)
	}
}

func TestRunMain_UnknownFlag(t *testing.T) {
	os.Args = []string{"swag2mcp", "--unknown-flag"} //nolint:reassign // test needs to set args
	err := runMain()
	if err == nil {
		t.Fatal("runMain() expected error for unknown flag, got nil")
	}
}

func TestRunMain_Init(t *testing.T) {
	tmpDir := t.TempDir()
	os.Args = []string{"swag2mcp", "init", tmpDir} //nolint:reassign // test needs to set args
	err := runMain()
	if err != nil {
		t.Fatalf("runMain() with init returned error: %v", err)
	}
}
