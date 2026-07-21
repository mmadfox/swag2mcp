package commands

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/tui"
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

func resolveBasePath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func readYAMLInput(yaml string) ([]byte, error) {
	if yaml == "-" {
		d, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return d, nil
	}
	return []byte(yaml), nil
}

func newAddSpecCmd() *cobra.Command {
	opts := struct {
		YAML    string
		Example bool
	}{}

	cmd := &cobra.Command{
		Use:   "spec [path]",
		Short: "Add a new API specification",
		Long: `Add a new API specification to the configuration.

Interactive mode (default):
  swag2mcp add spec

Non-interactive mode with YAML:
  swag2mcp add spec --yaml 'domain: meteo
  llm_title: Open-Meteo API
  base_url: https://meteo.swagger.io/v2'

  Or pipe from stdin:
  cat spec.yaml | swag2mcp add spec --yaml -

Show YAML example:
  swag2mcp add spec --example`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddSpec(resolveBasePath(args), opts.YAML, opts.Example)
		},
	}

	cmd.Flags().StringVarP(&opts.YAML, "yaml", "y", "", "YAML input (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Example, "example", "e", false, "Show YAML example and exit")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runAddSpec(basePath, yaml string, example bool) error {
	if example {
		fmt.Print(string(config.ExampleSpecAddYAML()))
		return nil
	}

	configPath, err := ensureConfigExists(basePath)
	if err != nil {
		return err
	}

	if yaml != "" {
		data, err := readYAMLInput(yaml)
		if err != nil {
			return err
		}
		return tui.AddSpecFromYAML(configPath, data)
	}
	return tui.AddSpecTUI(configPath)
}

func newAddCollectionCmd() *cobra.Command {
	opts := struct {
		YAML    string
		Example bool
	}{}

	cmd := &cobra.Command{
		Use:   "collection [path]",
		Short: "Add a new collection to an existing specification",
		Long: `Add a new collection to an existing specification.

Interactive mode (default):
  swag2mcp add collection

Non-interactive mode with YAML:
  swag2mcp add collection --yaml 'spec_domain: meteo
  llm_title: Orders Collection
  location: https://meteo.example.com/orders.json'

  Or pipe from stdin:
  cat collection.yaml | swag2mcp add collection --yaml -

Show YAML example:
  swag2mcp add collection --example`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddCollection(resolveBasePath(args), opts.YAML, opts.Example)
		},
	}

	cmd.Flags().StringVarP(&opts.YAML, "yaml", "y", "", "YAML input (use - for stdin)")
	cmd.Flags().BoolVarP(&opts.Example, "example", "e", false, "Show YAML example and exit")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

func runAddCollection(basePath, yaml string, example bool) error {
	if example {
		fmt.Print(string(config.ExampleCollectionAddYAML()))
		return nil
	}

	configPath, err := ensureConfigExists(basePath)
	if err != nil {
		return err
	}

	if yaml != "" {
		data, err := readYAMLInput(yaml)
		if err != nil {
			return err
		}
		return tui.AddCollectionFromYAML(configPath, data)
	}
	return tui.AddCollectionTUI(configPath)
}
