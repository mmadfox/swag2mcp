package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/tui"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a specification or collection from the configuration",
	}

	cmd.AddCommand(newDeleteSpecCmd())
	cmd.AddCommand(newDeleteCollectionCmd())

	return cmd
}

func newDeleteSpecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec [path]",
		Short: "Delete an API specification",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}

			configPath, err := ensureConfigExists(basePath)
			if err != nil {
				return err
			}
			if err := tui.DeleteSpecTUI(configPath); err != nil {
				return fmt.Errorf("delete spec: %w", err)
			}
			return nil
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func newDeleteCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection [path]",
		Short: "Delete a collection from a specification",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			basePath := ""
			if len(args) > 0 {
				basePath = args[0]
			}

			configPath, err := ensureConfigExists(basePath)
			if err != nil {
				return err
			}
			if err := tui.DeleteCollectionTUI(configPath); err != nil {
				return fmt.Errorf("delete collection: %w", err)
			}
			return nil
		},
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
