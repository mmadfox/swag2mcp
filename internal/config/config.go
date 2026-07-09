package config

import (
	"fmt"
	"iter"
	"time"

	"github.com/mmadfox/swag2mcp/internal/env"
)

// Cookie represents an HTTP cookie for configuration.
type Cookie struct {
	Name     string `yaml:"name"     validate:"required"`
	Value    string `yaml:"value"    validate:"required"`
	Domain   string `yaml:"domain,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Secure   bool   `yaml:"secure,omitempty"`
	HTTPOnly bool   `yaml:"http_only,omitempty"`
}

// ProxyConfig holds proxy connection settings.
type ProxyConfig struct {
	URL      string   `yaml:"url"`
	Username string   `yaml:"username,omitempty"`
	Password string   `yaml:"password,omitempty"`
	Bypass   []string `yaml:"bypass,omitempty"`
}

// HTTPClientConfig holds per-request HTTP settings for a spec or collection.
// These values are applied to each request at invocation time.
type HTTPClientConfig struct {
	Headers map[string]string `yaml:"headers,omitempty"`
	Cookies []Cookie          `yaml:"cookies,omitempty"`
}

// GlobalHTTPClientConfig holds global HTTP client settings.
type GlobalHTTPClientConfig struct {
	Randomize       bool              `yaml:"random,omitempty"`
	Proxy           *ProxyConfig      `yaml:"proxy,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Cookies         []Cookie          `yaml:"cookies,omitempty"`
	UserAgent       string            `yaml:"user_agent,omitempty"`
	Timeout         time.Duration     `yaml:"timeout,omitempty"`
	FollowRedirects *bool             `yaml:"follow_redirects,omitempty"`
	MaxRedirects    *int              `yaml:"max_redirects,omitempty"`
	MaxResponseSize *int              `yaml:"max_response_size,omitempty"`
}

// Config is the top-level swag2mcp configuration.
//
// Validation rules:
//   - Specs: at least one spec must be defined.
type Config struct {
	HTTPClient *GlobalHTTPClientConfig `yaml:"http_client,omitempty"`
	MCP        *MCPConfig              `yaml:"mcp,omitempty"`
	Specs      []Spec                  `yaml:"specs"`
}

// MCPConfig holds the MCP server configuration.
type MCPConfig struct {
	Transport string         `yaml:"transport,omitempty" validate:"omitempty,oneof=stdio sse streamable-http"`
	Addr      string         `yaml:"addr,omitempty"`
	Path      string         `yaml:"path,omitempty"`
	Auth      *MCPAuthConfig `yaml:"auth,omitempty"`
}

// MCPAuthConfig holds the MCP server authentication configuration.
type MCPAuthConfig struct {
	Token string `yaml:"token,omitempty"`
}

// Resolve resolves environment variable references in the token.
func (c *MCPAuthConfig) Resolve() {
	if c == nil {
		return
	}
	c.Token = env.Parse(c.Token)
}

// Spec defines a single API specification group.
//
// Validation rules:
//   - Domain: required, 1-60 chars, letters/digits/underscore/hyphen only.
//   - LLMTitle: required, 20-120 chars, allows letters/digits/punctuation.
//   - LLMInstruction: optional, max 500 chars, allows letters/digits/punctuation.
//   - Collections: required, 1-30 collections per spec.
//   - BaseURL: required, must be a valid URL.
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
//
// Validation rules:
//   - LLMTitle: optional, max 120 chars, allows letters/digits/punctuation.
//   - LLMInstruction: optional, max 360 chars, allows letters/digits/punctuation.
//   - Location: required, 5-250 chars (path or URL to the spec file).
//   - BaseURL: optional, must be a valid URL if set.
type Collection struct {
	LLMTitle       string            `yaml:"llm_title,omitempty"       json:"llm_title" validate:"omitempty,max=120,title_format"`
	LLMInstruction string            `yaml:"llm_instruction,omitempty"                  validate:"omitempty,max=360,instruction_format"`
	Title          string            `yaml:"title,omitempty"`
	Location       string            `yaml:"location"                  json:"location"  validate:"required,min=5,max=250"`
	Disable        bool              `yaml:"disable,omitempty"          json:"disable"`
	HTTPClient     *HTTPClientConfig `yaml:"http_client,omitempty"`
	BaseURL        string            `yaml:"base_url,omitempty"                          validate:"omitempty,url"`
}

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

func (c *Config) Validate(f *Filter) error {
	var errs validationErrors

	if len(c.Specs) == 0 {
		errs = append(errs, validationError{
			field:   "specs",
			message: "no specifications defined",
		})
	}

	for i, spec := range c.Specs {
		if spec.Disable {
			continue
		}
		if f != nil && !f.MatchSpec(spec.Tags...) {
			continue
		}

		specPrefix := fmt.Sprintf("specs[%d]", i)
		errs = append(errs, collectStructErrors(specPrefix, spec)...)

		if spec.Auth.Client != nil {
			if verr := spec.Auth.Client.Validate(); verr != nil {
				errs = append(errs, validationError{
					field:   specPrefix + ".auth",
					message: fmt.Sprintf("auth client validation failed: %s", verr),
				})
			}
		}

		for j, collection := range spec.Collections {
			if collection.Disable {
				continue
			}
			collPrefix := fmt.Sprintf("%s.collections[%d]", specPrefix, j)
			errs = append(errs, collectStructErrors(collPrefix, collection)...)
		}
	}

	errs = append(errs, collectStructErrors("config", c)...)

	if len(errs) == 0 {
		return nil
	}
	return errs
}
