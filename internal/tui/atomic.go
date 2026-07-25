package tui

import (
	"fmt"
	"os"

	"github.com/mmadfox/swag2mcp/internal/config"
	"go.yaml.in/yaml/v3"
)

// AtomicWriteConfig reads the config from configPath, calls fn to modify it,
// then writes the result atomically via a temporary file and [os.Rename].
func AtomicWriteConfig(configPath string, fn func(*config.Config) error) error {
	cfg, loadErr := config.Load(configPath)
	if loadErr != nil {
		return fmt.Errorf("load config: %w", loadErr)
	}

	if fnErr := fn(cfg); fnErr != nil {
		return fnErr
	}

	if valErr := cfg.Validate(config.NewFilter(nil)); valErr != nil {
		return fmt.Errorf("config validation failed:\n  %w\n  File: %s", valErr, configPath)
	}

	data, marshalErr := yaml.Marshal(cfg)
	if marshalErr != nil {
		return fmt.Errorf("marshal config: %w", marshalErr)
	}

	tmpPath := configPath + ".tmp"
	if writeErr := os.WriteFile(tmpPath, data, 0600); writeErr != nil {
		return fmt.Errorf("write temp file: %w", writeErr)
	}

	if renameErr := os.Rename(tmpPath, configPath); renameErr != nil {
		return fmt.Errorf("rename temp file: %w", renameErr)
	}

	return nil
}
