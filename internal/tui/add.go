package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"gopkg.in/yaml.v3"
)

func resolveConfigPath(configPath string) string {
	if configPath == "" {
		configPath = workspace.DefaultConfigPath()
	}
	if info, err := os.Stat(configPath); err == nil && info.IsDir() {
		configPath = filepath.Join(configPath, "swag2mcp.yaml")
	}
	return configPath
}

// AddSpecFromYAML adds a specification from a YAML string (non-interactive).
func AddSpecFromYAML(configPath string, data []byte) error {
	configPath = resolveConfigPath(configPath)

	var input config.SpecAddRequest
	if err := yaml.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	if input.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if input.LLMTitle == "" {
		return fmt.Errorf("llm_title is required")
	}
	if input.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		cfg.Specs = append(cfg.Specs, input.Spec)
		return nil
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("  ✅ Specification \"%s\" added.\n", input.Domain)
	return nil
}

// AddCollectionFromYAML adds a collection from a YAML string (non-interactive).
func AddCollectionFromYAML(configPath string, data []byte) error {
	configPath = resolveConfigPath(configPath)

	var input config.CollectionAddRequest
	if err := yaml.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}

	if input.SpecDomain == "" {
		return fmt.Errorf("spec_domain is required")
	}
	if input.LLMTitle == "" {
		return fmt.Errorf("llm_title is required")
	}
	if input.Location == "" {
		return fmt.Errorf("location is required")
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		for i := range cfg.Specs {
			if cfg.Specs[i].Domain == input.SpecDomain {
				cfg.Specs[i].Collections = append(cfg.Specs[i].Collections, input.Collection)
				return nil
			}
		}
		return fmt.Errorf("spec with domain %q not found", input.SpecDomain)
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("  ✅ Collection \"%s\" added to spec \"%s\".\n", input.LLMTitle, input.SpecDomain)
	return nil
}

// AddSpecTUI runs an interactive wizard to add a specification to an existing config.
func AddSpecTUI(configPath string) error {
	configPath = resolveConfigPath(configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	fmt.Printf("\n  Current specs: %d\n", len(cfg.Specs))
	for i, sp := range cfg.Specs {
		fmt.Printf("    %d. %s (%s)\n", i+1, sp.LLMTitle, sp.Domain)
	}
	fmt.Println()

	spec, err := collectSpec(len(cfg.Specs) + 1)
	if err != nil {
		return err
	}

	fmt.Printf("\n  Review:\n")
	fmt.Printf("  ────────\n")
	fmt.Printf("  Spec: %s (%s)\n", spec.LLMTitle, spec.Domain)
	fmt.Printf("  Base URL: %s\n", spec.BaseURL)
	if len(spec.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(spec.Tags, ", "))
	}
	if spec.AuthType != "" && spec.AuthType != "none" {
		fmt.Printf("  Auth: %s\n", spec.AuthType)
	}
	fmt.Printf("  Collections: %d\n", len(spec.Collections))
	for j, col := range spec.Collections {
		note := ""
		if strings.HasPrefix(col.Location, "http://") || strings.HasPrefix(col.Location, "https://") {
			note = " (cached)"
		}
		fmt.Printf("    %d. %s → %s%s\n", j+1, col.Title, col.Location, note)
	}
	fmt.Println()

	fmt.Print("  Write changes? (y/n) > ")
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("  Cancelled.")
		return nil
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		newSpec := config.Spec{
			Domain:         spec.Domain,
			LLMTitle:       spec.LLMTitle,
			LLMInstruction: spec.Instruction,
			BaseURL:        spec.BaseURL,
			Tags:           spec.Tags,
		}
		if spec.AuthType != "" && spec.AuthType != "none" {
			newSpec.Auth = config.Auth{}
		}
		for _, col := range spec.Collections {
			newSpec.Collections = append(newSpec.Collections, config.Collection{
				LLMTitle: col.Title,
				Location: col.Location,
			})
		}
		cfg.Specs = append(cfg.Specs, newSpec)
		return nil
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println("  ✅ Specification added.")
	return nil
}

// AddCollectionTUI runs an interactive wizard to add a collection to an existing spec.
func AddCollectionTUI(configPath string) error {
	configPath = resolveConfigPath(configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if len(cfg.Specs) == 0 {
		return fmt.Errorf("no specs found. Run 'swag2mcp add spec' first")
	}

	fmt.Printf("\n  Current specs: %d\n", len(cfg.Specs))
	for i, sp := range cfg.Specs {
		fmt.Printf("    %d. %s (%s)\n", i+1, sp.LLMTitle, sp.Domain)
	}
	fmt.Println()

	fmt.Printf("  Select a spec [1-%d] > ", len(cfg.Specs))
	var specIdx int
	fmt.Scanf("%d", &specIdx)
	if specIdx < 1 || specIdx > len(cfg.Specs) {
		return fmt.Errorf("invalid spec number")
	}
	spec := cfg.Specs[specIdx-1]

	col, err := collectCollection(specIdx, len(spec.Collections)+1, spec.Domain)
	if err != nil {
		return err
	}

	note := ""
	if strings.HasPrefix(col.Location, "http://") || strings.HasPrefix(col.Location, "https://") {
		note = " (cached)"
	}

	fmt.Printf("\n  Review:\n")
	fmt.Printf("  ────────\n")
	fmt.Printf("  Spec: %s (%s)\n", spec.LLMTitle, spec.Domain)
	fmt.Printf("  New collection: %s → %s%s\n", col.Title, col.Location, note)
	fmt.Println()

	fmt.Print("  Write changes? (y/n) > ")
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("  Cancelled.")
		return nil
	}

	if err := AtomicWriteConfig(configPath, func(cfg *config.Config) error {
		if specIdx-1 < len(cfg.Specs) {
			cfg.Specs[specIdx-1].Collections = append(cfg.Specs[specIdx-1].Collections, config.Collection{
				LLMTitle: col.Title,
				Location: col.Location,
			})
		}
		return nil
	}); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println("  ✅ Collection added.")
	return nil
}
