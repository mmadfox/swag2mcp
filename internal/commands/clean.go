/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean [path]",
		Short: "Remove temporary data (cache and responses)",
		Long: `Remove cached remote specs and invocation responses.

  swag2mcp clean              — clean ~/.swag2mcp/{cache,responses}
  swag2mcp clean ./           — clean ./{cache,responses}
  swag2mcp clean path/to      — clean path/to/{cache,responses}`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			result, err := runClean(basePath)
			if err != nil {
				return err
			}
			cmd.Printf("  cache: %s (%d files)\n", result.cache, result.cacheFiles)
			cmd.Printf("  responses: %s (%d files)\n", result.responses, result.responsesFiles)
			if result.orphanScripts > 0 {
				cmd.Printf("  orphan auth scripts: %d removed\n", result.orphanScripts)
			}
			cmd.Println("  ------------------------------")
			cmd.Println("✅ Clean completed")
			return nil
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

const cleanStatusCleared = "cleared"

type cleanResult struct {
	cache          string
	cacheFiles     int
	responses      string
	responsesFiles int
	orphanScripts  int
}

func runClean(basePath string) (cleanResult, error) {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return cleanResult{}, fmt.Errorf("workspace: %w", err)
	}

	cacheFiles := countFiles(ws.CacheDir())
	responsesFiles := countFiles(ws.ResponsesDir())

	if err := ws.Clean(); err != nil {
		return cleanResult{}, fmt.Errorf("clean: %w", err)
	}

	result := cleanResult{
		cache:          cleanStatusCleared,
		cacheFiles:     cacheFiles,
		responses:      cleanStatusCleared,
		responsesFiles: responsesFiles,
	}

	if ws.ConfigExists() {
		cfg, loadErr := config.Load(ws.ConfigPath())
		if loadErr == nil {
			var activeDomains []string
			for spec := range cfg.Iterate(nil) {
				activeDomains = append(activeDomains, spec.Domain)
			}
			before := countScripts(ws.AuthScriptsDir())
			if oErr := ws.RemoveOrphanAuthScripts(activeDomains); oErr != nil {
				return cleanResult{}, fmt.Errorf("remove orphan auth scripts: %w", oErr)
			}
			after := countScripts(ws.AuthScriptsDir())
			result.orphanScripts = before - after
		}
	}

	return result, nil
}

func countFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count
}

func countScripts(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".sh") || strings.HasSuffix(name, ".bat") {
			count++
		}
	}
	return count
}
