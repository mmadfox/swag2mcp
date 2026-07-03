package initpkg

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmadfox/swag2mcp/internal/workspace"
)

//go:embed init.swag2mcp.yaml
var exampleConfig []byte

//go:embed config.tmpl
var configTemplate string

// Setup creates a workspace and writes the example configuration file.
// Use this for non-interactive / scripted initialization.
func Setup(configPath, workspaceDir string) error {
	ws, err := workspace.New(workspaceDir)
	if err != nil {
		return fmt.Errorf("workspace: %w", err)
	}
	if err := ws.Init(); err != nil {
		return fmt.Errorf("init workspace: %w", err)
	}
	if err := os.WriteFile(configPath, exampleConfig, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// ExampleConfig returns the embedded example configuration.
func ExampleConfig() []byte {
	return exampleConfig
}

// WriteConfig writes the given YAML content to the config file path.
func WriteConfig(configPath string, data []byte) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(configPath, data, 0600)
}
