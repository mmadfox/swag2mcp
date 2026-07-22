package tui

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
)

func promptSelection(label string, maxVal int) (int, error) {
	fmt.Printf("  Select %s [1-%d] > ", label, maxVal)
	var idx int
	fmt.Scanln(&idx)
	if idx < 1 || idx > maxVal {
		return 0, fmt.Errorf("invalid number")
	}
	return idx, nil
}

func confirmAction(prompt string) bool {
	fmt.Printf("  %s (y/n) > ", prompt)
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// DeleteSpecTUI runs an interactive wizard to delete a specification.
func DeleteSpecTUI(configPath string) error {
	configPath = resolveConfigPath(configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Specs) == 0 {
		return fmt.Errorf("no specs found")
	}

	fmt.Printf("\n  Current specs: %d\n", len(cfg.Specs))
	for i, sp := range cfg.Specs {
		fmt.Printf("    %d. %s (%s)\n", i+1, sp.LLMTitle, sp.Domain)
	}
	fmt.Println()

	specIdx, err := promptSelection("a spec to delete", len(cfg.Specs))
	if err != nil {
		return err
	}
	spec := cfg.Specs[specIdx-1]

	if !confirmAction(fmt.Sprintf("Are you sure you want to delete \"%s\" (%s)?", spec.LLMTitle, spec.Domain)) {
		fmt.Println("  Cancelled.")
		return nil
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		cfg.Specs = append(cfg.Specs[:specIdx-1], cfg.Specs[specIdx:]...)
		return nil
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println("  Specification deleted.")
	return nil
}

// DeleteCollectionTUI runs an interactive wizard to delete a collection.
func DeleteCollectionTUI(configPath string) error {
	configPath = resolveConfigPath(configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Specs) == 0 {
		return fmt.Errorf("no specs found")
	}

	fmt.Printf("\n  Current specs: %d\n", len(cfg.Specs))
	for i, sp := range cfg.Specs {
		fmt.Printf("    %d. %s (%s) — %d collections\n", i+1, sp.LLMTitle, sp.Domain, len(sp.Collections))
	}
	fmt.Println()

	specIdx, err := promptSelection("a spec", len(cfg.Specs))
	if err != nil {
		return err
	}
	spec := cfg.Specs[specIdx-1]

	if len(spec.Collections) == 0 {
		return fmt.Errorf("spec \"%s\" has no collections", spec.LLMTitle)
	}

	fmt.Printf("\n  Collections for \"%s\":\n", spec.LLMTitle)
	for i, col := range spec.Collections {
		fmt.Printf("    %d. %s → %s\n", i+1, col.LLMTitle, col.Location)
	}
	fmt.Println()

	colIdx, err := promptSelection("a collection to delete", len(spec.Collections))
	if err != nil {
		return err
	}
	col := spec.Collections[colIdx-1]

	if !confirmAction(fmt.Sprintf("Are you sure you want to delete \"%s\"?", col.LLMTitle)) {
		fmt.Println("  Cancelled.")
		return nil
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		if specIdx-1 < len(cfg.Specs) {
			colls := cfg.Specs[specIdx-1].Collections
			cfg.Specs[specIdx-1].Collections = append(colls[:colIdx-1], colls[colIdx:]...)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println("  Collection deleted.")
	return nil
}
