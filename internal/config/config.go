/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"
	"iter"
	"time"
)

// Config is the top-level swag2mcp configuration.
type Config struct {
	MockEnabled        bool                    `yaml:"mock_enabled,omitempty"`
	MockAuth           *MockAuthConfig         `yaml:"mock_auth,omitempty"`
	HTTPClient         *GlobalHTTPClientConfig `yaml:"http_client,omitempty"`
	MCP                *MCPConfig              `yaml:"mcp,omitempty"`
	DisableRateLimiter bool                    `yaml:"disable_ratelimiter,omitempty"`
	RateLimitInterval  time.Duration           `yaml:"rate_limit_interval,omitempty"`
	GlobalRateLimit    int                     `yaml:"global_rate_limit,omitempty"`
	Specs              []Spec                  `yaml:"specs"`
}

// Spec defines a single API specification group.
type Spec struct {
	Domain         string            `yaml:"domain"                    validate:"required,domain_format"`
	LLMTitle       string            `yaml:"llm_title,omitempty"       validate:"required,min=5,max=120,title_format"`
	LLMInstruction string            `yaml:"llm_instruction,omitempty" validate:"omitempty,max=500,instruction_format"`
	Collections    []Collection      `yaml:"collections,omitempty"     validate:"required,min=1,max=30"`
	Disable        bool              `yaml:"disable,omitempty"`
	Tags           []string          `yaml:"tags,omitempty"`
	BaseURL        string            `yaml:"base_url,omitempty"        validate:"required,url"`
	HTTPClient     *HTTPClientConfig `yaml:"http_client,omitempty"`
	Auth           Auth              `yaml:"auth,omitempty"`
}

// Collection defines a single spec file (Swagger/OpenAPI) within a Spec.
type Collection struct {
	LLMTitle       string            `yaml:"llm_title,omitempty"       json:"llm_title" validate:"omitempty,max=120,title_format"`
	LLMInstruction string            `yaml:"llm_instruction,omitempty"                  validate:"omitempty,max=360,instruction_format"`
	Title          string            `yaml:"title,omitempty"`
	Location       string            `yaml:"location"                  json:"location"  validate:"required,min=5,max=250,spec_location"`
	Disable        bool              `yaml:"disable,omitempty"          json:"disable"`
	HTTPClient     *HTTPClientConfig `yaml:"http_client,omitempty"`
	BaseURL        string            `yaml:"base_url,omitempty"                          validate:"omitempty,url"`
	BaseMockURL    string            `yaml:"base_mock_url,omitempty"                      validate:"omitempty,mock_addr_format"`
}

// SetDefaults fills zero fields with sensible defaults.
func (c *Config) SetDefaults() {
	if c == nil {
		return
	}
	if c.RateLimitInterval <= 0 {
		c.RateLimitInterval = defaultRateLimitInterval
	}
}

// Iterate returns an iterator over non-disabled specs that match the given filter.
func (c *Config) Iterate(f *Filter) iter.Seq[*Spec] {
	return func(yield func(*Spec) bool) {
		for _, spec := range c.Specs {
			if spec.Disable {
				continue
			}
			if f != nil {
				if match := f.MatchSpec(spec.Tags...); !match {
					continue
				}
			}
			if !yield(&spec) {
				break
			}
		}
	}
}

// Validate validates the configuration against struct tags and business rules.
// It skips disabled specs and applies the optional filter. Returns nil if valid.
func (c *Config) Validate(f *Filter) error {
	var errs validationErrors

	if len(c.Specs) == 0 {
		errs = append(errs, validationError{
			field:   "specs",
			message: "no specifications defined",
		})
	}

	for i := range c.Specs {
		spec := &c.Specs[i]
		if spec.Disable {
			continue
		}
		if f != nil && !f.MatchSpec(spec.Tags...) {
			continue
		}
		specPrefix := fmt.Sprintf("specs[%d]", i)
		errs = append(errs, c.validateSpec(specPrefix, spec)...)
	}

	errs = append(errs, collectStructErrors("config", c)...)

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (c *Config) validateSpec(prefix string, spec *Spec) []validationError {
	var errs validationErrors

	specErrs := collectStructErrors(prefix, *spec)
	for k := range specErrs {
		specErrs[k].spec = spec.Domain
	}
	errs = append(errs, specErrs...)

	if spec.Auth.Client != nil {
		if verr := spec.Auth.Client.Validate(); verr != nil {
			errs = append(errs, validationError{
				field:   prefix + ".auth",
				spec:    spec.Domain,
				message: fmt.Sprintf("auth client validation failed: %s", verr),
			})
		}
	}

	for j := range spec.Collections {
		coll := &spec.Collections[j]
		if coll.Disable {
			continue
		}
		collPrefix := fmt.Sprintf("%s.collections[%d]", prefix, j)
		errs = append(errs, c.validateCollection(collPrefix, spec.Domain, j, coll)...)
	}

	return errs
}

func (c *Config) validateCollection(prefix, specDomain string, collIdx int, coll *Collection) []validationError {
	var errs validationErrors

	collErrs := collectStructErrors(prefix, *coll)
	collTitle := coll.LLMTitle
	if collTitle == "" {
		collTitle = fmt.Sprintf("#%d", collIdx)
	}
	for k := range collErrs {
		collErrs[k].spec = specDomain
		collErrs[k].collection = collTitle
	}
	errs = append(errs, collErrs...)

	if c.MockEnabled && coll.BaseMockURL == "" {
		errs = append(errs, validationError{
			field:      prefix + ".base_mock_url",
			spec:       specDomain,
			collection: collTitle,
			message:    "BaseMockURL is required when mock_enabled is true",
		})
	}

	return errs
}
