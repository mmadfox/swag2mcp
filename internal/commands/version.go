package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

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
