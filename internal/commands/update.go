package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [path]",
		Short: "Update cache from configuration",
		Long: `Validate configuration, clear cache, and re-cache all spec files.

  swag2mcp update              — update ~/.swag2mcp
  swag2mcp update ./           — update ./.swag2mcp
  swag2mcp update path/to      — update path/to/.swag2mcp`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}

			total, err := runUpdate(basePath)
			if err != nil {
				return err
			}

			cmd.Printf("✅ Cache updated (%d specs cached)\n", total)
			return nil
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runUpdate(basePath string) (int, error) {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return 0, fmt.Errorf("workspace: %w", err)
	}

	configPath := ws.ConfigPath()

	if ws.ConfigNotExists() {
		return 0, fmt.Errorf("configuration not found at %s", configPath)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return 0, fmt.Errorf("load config: %w", err)
	}

	wsDir := filepath.Dir(configPath)
	ca := cache.New(wsDir)

	validateOpts := config.ValidateOptions{
		Cache: ca,
	}
	if err := config.ValidateConfig(cfg, validateOpts); err != nil {
		return 0, fmt.Errorf("config validation failed:\n  %w", err)
	}

	if err := ws.Clean(); err != nil {
		return 0, fmt.Errorf("clean cache: %w", err)
	}

	total, err := cacheSpecs(cfg, ca, ws)
	if err != nil {
		return 0, err
	}

	if err := cleanOrphanAuthScripts(cfg, ws); err != nil {
		return 0, err
	}

	return total, nil
}

func cacheSpecs(cfg *config.Config, ca *cache.Cache, ws *workspace.Workspace) (int, error) {
	var total int
	for spec := range cfg.Iterate(nil) {
		for _, col := range spec.Collections {
			if col.Disable {
				continue
			}
			if _, rErr := ca.Resolve(col.Location); rErr != nil {
				return 0, fmt.Errorf("cache %s: %w", col.Location, rErr)
			}
			total++
		}

		if spec.Auth.Client != nil && spec.Auth.Client.Type() == auth.ScriptAuth {
			if sErr := ws.EnsureAuthScript(spec.Domain); sErr != nil {
				return 0, fmt.Errorf("auth script %s: %w", spec.Domain, sErr)
			}
		}
	}
	return total, nil
}

func cleanOrphanAuthScripts(cfg *config.Config, ws *workspace.Workspace) error {
	var activeDomains []string
	for spec := range cfg.Iterate(nil) {
		activeDomains = append(activeDomains, spec.Domain)
	}
	if oErr := ws.RemoveOrphanAuthScripts(activeDomains); oErr != nil {
		return fmt.Errorf("remove orphan auth scripts: %w", oErr)
	}
	return nil
}
