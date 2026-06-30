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
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("cannot determine home directory: %w", err)
			}

			workspaceDir := filepath.Join(home, ".swag2mcp")

			// Create workspace directory
			if err := os.MkdirAll(workspaceDir, 0755); err != nil {
				return fmt.Errorf("failed to create workspace directory %q: %w", workspaceDir, err)
			}

			// Create subdirectories
			subdirs := []string{"responses", "cache", "specs"}
			for _, subdir := range subdirs {
				dir := filepath.Join(workspaceDir, subdir)
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create subdirectory %q: %w", dir, err)
				}
			}

			// Create default config file
			configPath := filepath.Join(workspaceDir, "swag2mcp.yaml")
			if _, err := os.Stat(configPath); err == nil {
				// File exists
				fmt.Printf("Configuration file already exists at %s\n", configPath)
				return nil
			}

			configContent := `collections: []
`
			if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
				return fmt.Errorf("failed to write config file %q: %w", configPath, err)
			}

			fmt.Printf("Workspace initialized at %s\n", workspaceDir)
			fmt.Printf("Configuration created at %s\n", configPath)
			fmt.Printf("Subdirectories created: %v\n", subdirs)

			return nil
		},
	}

	return cmd
}
