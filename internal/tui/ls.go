package tui

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

		fmt.Fprintf(w, "  %s\t%s\t%s\n", spec.Domain, spec.LLMTitle, spec.BaseURL)

		if len(spec.Tags) > 0 {
			fmt.Fprintf(w, "    Tags:\t%s\n", strings.Join(spec.Tags, ", "))
		}

		if spec.Auth.Client != nil && spec.Auth.Client.Type() != "none" {
			fmt.Fprintf(w, "    Auth:\t%s\n", spec.Auth.Client.Type())
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
