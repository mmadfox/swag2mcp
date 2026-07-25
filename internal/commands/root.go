package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root swag2mcp command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "swag2mcp",
		Short:        "swag2mcp - MCP server for OpenAPI/Swagger specifications",
		Long:         `swag2mcp provides LLM agents with tools to work with Swagger/OpenAPI 3+ specifications.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if showVersion, _ := cmd.Flags().GetBool("version"); showVersion {
				fmt.Fprintf(os.Stdout, "swag2mcp %s\n", Version)
				return nil
			}
			return cmd.Help()
		},
	}

	cmd.Flags().Bool("version", false, "Print the swag2mcp version")
	cmd.Flags().BoolP("help", "h", false, "Print help")

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newLsCmd())
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newCleanCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newMCPCmd(Version))
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newImportCmd())
	cmd.AddCommand(newExportCmd())

	return cmd
}
