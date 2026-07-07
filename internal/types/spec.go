package types

import (
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

// Cookie represents an HTTP cookie.
type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain,omitempty"`
	Path     string `json:"path,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	HTTPOnly bool   `json:"httpOnly,omitempty"`
}

// HTTPClientConfig holds HTTP client settings.
type HTTPClientConfig struct {
	Headers         map[string]string `json:"headers,omitempty"`
	Cookies         []Cookie          `json:"cookies,omitempty"`
	UserAgent       string            `json:"userAgent,omitempty"`
	Timeout         time.Duration     `json:"timeout,omitempty"`
	FollowRedirects *bool             `json:"followRedirects,omitempty"`
	MaxRedirects    *int              `json:"maxRedirects,omitempty"`
}

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

func (s *Spec) InitAuthenticator() error {
	if s.Auth == nil {
		return nil
	}
	return s.Auth.New()
}

type Collection struct {
	ID             string            `json:"id"`
	SpecID         string            `json:"specId"`
	LLMTitle       string            `json:"llmtitle"`
	LLMInstruction string            `json:"llminstruction"`
	Title          string            `json:"title"`
	BaseURL        string            `json:"baseurl,omitempty"`
	HTTPClient     *HTTPClientConfig `json:"httpClient,omitempty"`
	Stats          struct {
		Tags    int `json:"tags"`
		Methods int `json:"methods"`
	}
}

type Tag struct {
	ID           string `json:"id"`
	CollectionID string `json:"collectionId"`
	SpecID       string `json:"specId"`
	Name         string `json:"name"`
	Stats        struct {
		Methods int `json:"methods"`
	}
}

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

func (e *Endpoint) SummaryOrFallback() string {
	if e.Operation.Summary != "" {
		return e.Operation.Summary
	}
	if e.Operation.Description != "" {
		return e.Operation.Description
	}
	return e.Name + " " + e.Path
}
