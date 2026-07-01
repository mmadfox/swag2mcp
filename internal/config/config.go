package config

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
)

const DefaultWorkspaceDir = ".swag2mcp"

type Config struct {
	WorkspaceDir string `yaml:"workspace_dir,omitempty"`
	Specs        []Spec `yaml:"specs"`
}

type Spec struct {
	Domain         string            `yaml:"domain" validate:"required,domain_format"`
	LLMTitle       string            `yaml:"llm_title,omitempty" validate:"required,min=20,max=120,title_format"`
	LLMInstruction string            `yaml:"llm_instruction,omitempty" validate:"omitempty,max=500,instruction_format"`
	Collections    []Collection      `yaml:"collections" validate:"required,min=1,max=30"`
	Disable        bool              `yaml:"disable"`
	Tags           []string          `yaml:"tags,omitempty"`
	BaseURL        string            `yaml:"base_url" validate:"required,url"`
	Headers        map[string]string `yaml:"headers,omitempty"`
}

type Collection struct {
	LLMTitle       string            `yaml:"llm_title,omitempty" json:"llm_title" validate:"omitempty,max=120,title_format"`
	LLMInstruction string            `yaml:"llm_instruction,omitempty" validate:"omitempty,max=360,instruction_format"`
	Title          string            `yaml:"title,omitempty"`
	Location       string            `yaml:"location" json:"location" validate:"required,min=5,max=250"`
	Disable        bool              `yaml:"disable" json:"disable"`
	Headers        map[string]string `yaml:"headers,omitempty"`
	BaseURL        string            `yaml:"base_url" validate:"omitempty,url"`
}

func (c *Config) SetDefaults() error {
	if c.WorkspaceDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		c.WorkspaceDir = filepath.Join(home, DefaultWorkspaceDir)
	}
	return nil
}

func (c *Config) Iterate(f *Filter) iter.Seq[*Spec] {
	return func(yield func(*Spec) bool) {
		for _, spec := range c.Specs {
			if spec.Disable {
				continue
			}
			if match := f.MatchSpec(spec.Tags...); !match {
				continue
			}
			if !yield(&spec) {
				break
			}
		}
	}
}

func (c *Config) Validate(f *Filter) error {
	if len(c.Specs) == 0 {
		return fmt.Errorf("no specs found")
	}

	specIndex := 1
	for spec := range c.Iterate(f) {
		if err := configValidator.Struct(spec); err != nil {
			return fmt.Errorf("failed to validate spec-%d: %w", specIndex, err)
		}

		for _, collection := range spec.Collections {
			if collection.Disable {
				continue
			}
			if err := configValidator.Struct(collection); err != nil {
				return fmt.Errorf("failed to validate collection %q: %w", collection.LLMTitle, err)
			}
		}

		specIndex++
	}

	if err := configValidator.Struct(c); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}
