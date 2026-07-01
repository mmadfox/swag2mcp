package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize workspace directory and create configuration file",
		Long: `Initialize workspace directory and create configuration file.

This command creates:
- Workspace directory (default: ~/.swag2mcp)
- Configuration file (swag2mcp.yaml)
- Subdirectories for storing responses and cache`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return fmt.Errorf("cannot determine home directory: %w", homeErr)
			}

			workspaceDir := filepath.Join(home, ".swag2mcp")

			// Create workspace directory
			if err := os.MkdirAll(workspaceDir, 0750); err != nil {
				return fmt.Errorf("failed to create workspace directory %q: %w", workspaceDir, err)
			}

			// Create subdirectories
			subdirs := []string{"responses", "cache", "specs"}
			for _, subdir := range subdirs {
				dir := filepath.Join(workspaceDir, subdir)
				if mkErr := os.MkdirAll(dir, 0750); mkErr != nil {
					return fmt.Errorf("failed to create subdirectory %q: %w", dir, mkErr)
				}
			}

			// Create default config file
			configPath := filepath.Join(workspaceDir, "swag2mcp.yaml")
			if _, statErr := os.Stat(configPath); statErr == nil {
				// File exists
				cmd.Printf("Configuration file already exists at %s\n", configPath)
				return nil
			}

			configContent := `collections: []
`
			if writeErr := os.WriteFile(configPath, []byte(configContent), 0600); writeErr != nil {
				return fmt.Errorf("failed to write config file %q: %w", configPath, writeErr)
			}

			cmd.Printf("Workspace initialized at %s\n", workspaceDir)
			cmd.Printf("Configuration created at %s\n", configPath)
			cmd.Printf("Subdirectories created: %v\n", subdirs)

			return nil
		},
	}

	return cmd
}
