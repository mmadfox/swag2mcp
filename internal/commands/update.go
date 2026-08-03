/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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
  swag2mcp update ./           — update ./
  swag2mcp update path/to      — update path/to`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}

			result, err := runUpdate(basePath)
			if err != nil {
				return err
			}

			for _, s := range result.specs {
				cmd.Printf("  %s (collections: %d)\n", s.domain, s.collections)
			}
			cmd.Println("  ------------------------------")
			cmd.Printf("✅ %d specs processed\n", result.total)
			return nil
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

type specSummary struct {
	domain      string
	collections int
}

type updateResult struct {
	total int
	specs []specSummary
}

func runUpdate(basePath string) (updateResult, error) {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return updateResult{}, fmt.Errorf("workspace: %w", err)
	}

	configPath := ws.ConfigPath()

	if ws.ConfigNotExists() {
		var err error
		configPath, err = ensureConfigExists(basePath)
		if err != nil {
			return updateResult{}, fmt.Errorf("configuration not found at %s: %w", configPath, err)
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return updateResult{}, fmt.Errorf("load config: %w", err)
	}

	wsDir := filepath.Dir(configPath)
	ca := cache.New(wsDir)

	validateOpts := config.ValidateOptions{
		Cache: ca,
	}
	if err := config.ValidateConfig(cfg, validateOpts); err != nil {
		return updateResult{}, fmt.Errorf("config validation failed:\n  %w", err)
	}

	if err := ws.Clean(); err != nil {
		return updateResult{}, fmt.Errorf("clean cache: %w", err)
	}

	if _, _, err := cacheSpecs(cfg, ca, ws); err != nil {
		return updateResult{}, err
	}

	if err := cleanOrphanAuthScripts(cfg, ws); err != nil {
		return updateResult{}, err
	}

	var specs []specSummary
	for s := range cfg.Iterate(nil) {
		count := 0
		for _, c := range s.Collections {
			if !c.Disable {
				count++
			}
		}
		specs = append(specs, specSummary{domain: s.Domain, collections: count})
	}

	return updateResult{total: len(specs), specs: specs}, nil
}

func cacheSpecs(cfg *config.Config, ca *cache.Cache, ws *workspace.Workspace) (int, int, error) {
	var remote, local int
	for spec := range cfg.Iterate(nil) {
		for _, col := range spec.Collections {
			if col.Disable {
				continue
			}
			if _, rErr := ca.Resolve(context.Background(), col.Location); rErr != nil {
				return 0, 0, fmt.Errorf("cache %s: %w", col.Location, rErr)
			}
			if isRemoteLocation(col.Location) {
				remote++
			} else {
				local++
			}
		}

		if spec.Auth.Client != nil && spec.Auth.Client.Type() == auth.ScriptAuth {
			if sErr := ws.EnsureAuthScript(spec.Domain); sErr != nil {
				return 0, 0, fmt.Errorf("auth script %s: %w", spec.Domain, sErr)
			}
		}
	}
	return remote, local, nil
}

func isRemoteLocation(location string) bool {
	return strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://")
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
