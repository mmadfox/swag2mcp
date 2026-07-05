package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// DeleteSpecTUI runs an interactive wizard to delete a specification.
func DeleteSpecTUI(configPath string) error {
	if configPath == "" {
		configPath = workspace.DefaultConfigPath()
	}
	if info, statErr := os.Stat(configPath); statErr == nil && info.IsDir() {
		configPath = filepath.Join(configPath, "swag2mcp.yaml")
	}

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

	fmt.Print("  Select a spec to delete [1-", len(cfg.Specs), "] > ")
	var specIdx int
	fmt.Scanf("%d", &specIdx)
	if specIdx < 1 || specIdx > len(cfg.Specs) {
		return fmt.Errorf("invalid spec number")
	}
	spec := cfg.Specs[specIdx-1]

	fmt.Printf("\n  Are you sure you want to delete \"%s\" (%s)? (y/n) > ", spec.LLMTitle, spec.Domain)
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("  Cancelled.")
		return nil
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		cfg.Specs = append(cfg.Specs[:specIdx-1], cfg.Specs[specIdx:]...)
		return nil
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println("  ✅ Specification deleted.")
	return nil
}

// DeleteCollectionTUI runs an interactive wizard to delete a collection.
func DeleteCollectionTUI(configPath string) error {
	if configPath == "" {
		configPath = workspace.DefaultConfigPath()
	}
	if info, statErr := os.Stat(configPath); statErr == nil && info.IsDir() {
		configPath = filepath.Join(configPath, "swag2mcp.yaml")
	}

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

	fmt.Print("  Select a spec [1-", len(cfg.Specs), "] > ")
	var specIdx int
	fmt.Scanf("%d", &specIdx)
	if specIdx < 1 || specIdx > len(cfg.Specs) {
		return fmt.Errorf("invalid spec number")
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

	fmt.Print("  Select a collection to delete [1-", len(spec.Collections), "] > ")
	var colIdx int
	fmt.Scanf("%d", &colIdx)
	if colIdx < 1 || colIdx > len(spec.Collections) {
		return fmt.Errorf("invalid collection number")
	}
	col := spec.Collections[colIdx-1]

	fmt.Printf("\n  Are you sure you want to delete \"%s\"? (y/n) > ", col.LLMTitle)
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
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

	fmt.Println("  ✅ Collection deleted.")
	return nil
}
