package model

import (
	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

// HTTPClientConfig holds per-request HTTP settings for a spec or collection.
type HTTPClientConfig struct {
	Headers map[string]string   `json:"headers,omitempty"`
	Cookies []httpclient.Cookie `json:"cookies,omitempty"`
}

// Spec represents an API specification with its domain, metadata, authentication, and HTTP client config.
type Spec struct {
	ID             string             `json:"id"`
	Domain         string             `json:"domain"`
	LLMTitle       string             `json:"llmtitle"`
	LLMInstruction string             `json:"llminstruction"`
	BaseURL        string             `json:"baseurl"`
	HTTPClient     *HTTPClientConfig  `json:"httpClient,omitempty"`
	Auth           auth.Authenticator `json:"auth"`
	Stats          struct {
		Collections int `json:"collections"`
		Tags        int `json:"tags"`
		Methods     int `json:"methods"`
	}
}

// InitAuthenticator initializes the spec's authenticator, if one is configured.
func (s *Spec) InitAuthenticator() error {
	if s.Auth == nil {
		return nil
	}
	return s.Auth.New()
}

// Collection represents a group of endpoints within a spec.
type Collection struct {
	ID             string            `json:"id"`
	SpecID         string            `json:"specId"`
	LLMTitle       string            `json:"llmtitle"`
	LLMInstruction string            `json:"llminstruction"`
	Title          string            `json:"title"`
	BaseURL        string            `json:"baseurl,omitempty"`
	BaseMockURL    string            `json:"base_mock_url,omitempty"`
	HTTPClient     *HTTPClientConfig `json:"httpClient,omitempty"`
	Stats          struct {
		Tags    int `json:"tags"`
		Methods int `json:"methods"`
	}
}

// Tag represents a logical grouping of endpoints within a collection.
type Tag struct {
	ID           string `json:"id"`
	CollectionID string `json:"collectionId"`
	SpecID       string `json:"specId"`
	Name         string `json:"name"`
	Stats        struct {
		Methods int `json:"methods"`
	}
}

// Endpoint represents a single API endpoint with its method, path, tag, and operation details.
type Endpoint struct {
	ID           string          `json:"id"`
	TagID        string          `json:"tagId"`
	CollectionID string          `json:"collectionId"`
	SpecID       string          `json:"specId"`
	Name         string          `json:"method"`
	Path         string          `json:"path"`
	Tag          string          `json:"tag"`
	Operation    *spec.Operation `json:"operation"`
}

// SummaryOrFallback returns the endpoint summary, falling back to description, then "Method /path".
func (e *Endpoint) SummaryOrFallback() string {
	if e.Operation.Summary != "" {
		return e.Operation.Summary
	}
	if e.Operation.Description != "" {
		return e.Operation.Description
	}
	return e.Name + " " + e.Path
}
