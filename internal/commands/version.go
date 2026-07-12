package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time (e.g. -X ...=v1.1.1).
// Defaults to "dev" for local development.
var Version = "dev"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
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
