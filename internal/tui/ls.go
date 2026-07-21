package tui

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/mmadfox/swag2mcp/internal/config"
)

// ListConfig loads the config and returns a formatted listing of specs and collections.
func ListConfig(configPath string, tags []string) (string, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return "", fmt.Errorf("load config: %w", err)
	}

	filter := config.NewFilter(tags)
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 3, ' ', 0) //nolint:mnd // tab padding width

	fmt.Fprintln(w, "Specifications:")

	for _, spec := range cfg.Specs {
		if spec.Disable {
			continue
		}
		if !filter.MatchSpec(spec.Tags...) {
			continue
		}

		authType := "-"
		if spec.Auth.Client != nil && spec.Auth.Client.Type() != "none" {
			authType = string(spec.Auth.Client.Type())
		}

		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", spec.Domain, spec.LLMTitle, spec.BaseURL, authType)

		if len(spec.Tags) > 0 {
			fmt.Fprintf(w, "    Tags:\t%s\n", strings.Join(spec.Tags, ", "))
		}

		if len(spec.Collections) > 0 {
			fmt.Fprintf(w, "    Collections (%d):\n", len(spec.Collections))
			for _, col := range spec.Collections {
				if col.Disable {
					continue
				}
				fmt.Fprintf(w, "      %s\t%s\n", col.LLMTitle, col.Location)
			}
		}

		fmt.Fprintln(w)
	}

	w.Flush()
	return b.String(), nil
}
