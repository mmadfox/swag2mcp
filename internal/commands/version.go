/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time (e.g. -X ...=v1.1.1).
// Defaults to "dev" for local development.
var Version = "dev"

const versionUse = "version"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   versionUse,
		Short: "Print the swag2mcp version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "swag2mcp %s\n", Version)
			return err
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}

func newMockVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   versionUse,
		Short: "Print the swag2mcp-mock version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "swag2mcp-mock %s\n", Version)
			return err
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}
