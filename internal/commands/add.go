package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/initmcp"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a specification or collection to the configuration",
	}

	cmd.AddCommand(newAddSpecCmd())
	cmd.AddCommand(newAddCollectionCmd())

	return cmd
}

func newAddSpecCmd() *cobra.Command {
	opts := struct {
		ConfigPath string
		YAML       string
		Example    bool
	}{}

	cmd := &cobra.Command{
		Use:   "spec",
		Short: "Add a new API specification",
		Long: `Add a new API specification to the configuration.

Interactive mode (default):
  swag2mcp add spec

Non-interactive mode with YAML:
  swag2mcp add spec --yaml 'domain: petstore
  llm_title: Petstore API
  base_url: https://petstore.swagger.io/v2'

  Or pipe from stdin:
  cat spec.yaml | swag2mcp add spec --yaml -

Show YAML example:
  swag2mcp add spec --example`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.Example {
				fmt.Print(string(config.ExampleSpecAddYAML()))
				return nil
			}

			configPath, err := ensureConfigExists(opts.ConfigPath)
			if err != nil {
				return err
			}

			if opts.YAML != "" {
				var data []byte
				if opts.YAML == "-" {
					d, err := io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("read stdin: %w", err)
					}
					data = d
				} else {
					data = []byte(opts.YAML)
				}
				if err := initmcp.AddSpecFromYAML(configPath, data); err != nil {
					return fmt.Errorf("add spec: %w", err)
				}
				return nil
			}
			if err := initmcp.AddSpecTUI(configPath); err != nil {
				return fmt.Errorf("add spec: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&opts.YAML, "yaml", "y", "", "YAML input (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Example, "example", "e", false, "Show YAML example and exit")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func newAddCollectionCmd() *cobra.Command {
	opts := struct {
		ConfigPath string
		YAML       string
		Example    bool
	}{}

	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Add a new collection to an existing specification",
		Long: `Add a new collection to an existing specification.

Interactive mode (default):
  swag2mcp add collection

Non-interactive mode with YAML:
  swag2mcp add collection --yaml 'spec_domain: petstore
  llm_title: Orders Collection
  location: https://petstore.example.com/orders.json'

  Or pipe from stdin:
  cat collection.yaml | swag2mcp add collection --yaml -

Show YAML example:
  swag2mcp add collection --example`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.Example {
				fmt.Print(string(config.ExampleCollectionAddYAML()))
				return nil
			}

			configPath, err := ensureConfigExists(opts.ConfigPath)
			if err != nil {
				return err
			}

			if opts.YAML != "" {
				var data []byte
				if opts.YAML == "-" {
					d, err := io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("read stdin: %w", err)
					}
					data = d
				} else {
					data = []byte(opts.YAML)
				}
				if err := initmcp.AddCollectionFromYAML(configPath, data); err != nil {
					return fmt.Errorf("add collection: %w", err)
				}
				return nil
			}
			if err := initmcp.AddCollectionTUI(configPath); err != nil {
				return fmt.Errorf("add collection: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.ConfigPath, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&opts.YAML, "yaml", "y", "", "YAML input (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Example, "example", "e", false, "Show YAML example and exit")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}
