package commands

import (
	"fmt"
	"os"

	"github.com/mmadfox/swag2mcp/internal/initmcp"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// ensureConfigExists checks if the config file exists.
// If not, it runs the init wizard to create one.
func ensureConfigExists(configPath string) (string, error) {
	if configPath == "" {
		configPath = workspace.DefaultConfigPath()
	}
	if info, statErr := os.Stat(configPath); statErr == nil && info.IsDir() {
		configPath = workspace.DefaultConfigPath()
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("\n  Configuration file not found at %s\n", configPath)
		fmt.Println("  Let's create one first.")

		cfgPath, wsDir, _, err := initmcp.RunTUI()
		if err != nil {
			return "", fmt.Errorf("init wizard: %w", err)
		}

		fmt.Printf("\n  ✅ Configuration written to: %s\n", cfgPath)
		fmt.Printf("  ✅ Workspace initialized at: %s\n", wsDir)
		fmt.Println()

		return cfgPath, nil
	}

	return configPath, nil
}
