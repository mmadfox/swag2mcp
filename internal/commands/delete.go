package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/initmcp"
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
	opts := struct {
		ConfigPath string
	}{}

	cmd := &cobra.Command{
		Use:   "spec",
		Short: "Delete an API specification",
		RunE: func(cmd *cobra.Command, _ []string) error {
			configPath, err := ensureConfigExists(opts.ConfigPath)
			if err != nil {
				return err
			}
			if err := initmcp.DeleteSpecTUI(configPath); err != nil {
				return fmt.Errorf("delete spec: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to configuration file")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func newDeleteCollectionCmd() *cobra.Command {
	opts := struct {
		ConfigPath string
	}{}

	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Delete a collection from a specification",
		RunE: func(cmd *cobra.Command, _ []string) error {
			configPath, err := ensureConfigExists(opts.ConfigPath)
			if err != nil {
				return err
			}
			if err := initmcp.DeleteCollectionTUI(configPath); err != nil {
				return fmt.Errorf("delete collection: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to configuration file")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
