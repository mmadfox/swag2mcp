package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/tui"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

func newLsCmd() *cobra.Command {
	opts := struct {
		Tags string
	}{}

	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List specifications and collections",
		Long: `List all specifications and their collections from the configuration.

  swag2mcp ls              — list ~/.swag2mcp/swag2mcp.yaml
  swag2mcp ls ./           — list ./swag2mcp.yaml
  swag2mcp ls path/to      — list path/to/swag2mcp.yaml
  swag2mcp ls --tags=public,internal`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}
			return runLs(basePath, opts.Tags, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&opts.Tags, "tags", "t", "", "Filter by tags (comma-separated)")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runLs(basePath, tagsFilter string, w io.Writer) error {
	ws, err := workspace.NewFromBase(basePath)
	if err != nil {
		return err
	}

	configPath := ws.ConfigPath()

	if ws.ConfigNotExists() {
		configPath, err = ensureConfigExists(basePath)
		if err != nil {
			return err
		}
	}

	var tags []string
	if tagsFilter != "" {
		tags = strings.Split(tagsFilter, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	output, err := tui.ListConfig(configPath, tags)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, output)
	return err
}
