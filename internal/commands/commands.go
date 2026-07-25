package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"path/filepath"

	"github.com/mmadfox/swag2mcp/internal/tui"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// ensureConfigExists checks if the config file exists.
// If not, it runs the init wizard to create one.
func ensureConfigExists(basePath string) (string, error) {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return "", fmt.Errorf("workspace: %w", err)
	}

	configPath := ws.ConfigPath()

	if ws.ConfigNotExists() {
		fmt.Printf("\n  Configuration file not found at %s\n", configPath)
		fmt.Println("  Let's create one first.")

		wsDir := filepath.Dir(configPath)
		if err := tui.Setup(configPath, wsDir); err != nil {
			return "", fmt.Errorf("setup: %w", err)
		}

		fmt.Printf("\n  ✅ Configuration written to: %s\n", configPath)
		fmt.Printf("  ✅ Workspace initialized at: %s\n", wsDir)
		fmt.Println()

		return configPath, nil
	}

	return configPath, nil
}
