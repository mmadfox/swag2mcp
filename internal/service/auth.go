package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

// AuthRequest contains the parameters needed to retrieve authentication
// information for a given spec.
type (
	AuthRequest struct {
		SpecID string `json:"specId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the spec/domain to get an auth token for"`
	}

	// AuthResponse contains the authentication token, headers, and query
	// parameters returned after applying the spec's auth configuration.
	AuthResponse struct {
		Token       string            `json:"token"`
		Headers     map[string]string `json:"headers,omitempty"`
		QueryParams map[string]string `json:"queryParams,omitempty"`
	}
)

// Auth retrieves authentication information for the spec identified by the
// request SpecID. It applies the spec's auth configuration and returns the
// resulting token, headers, and query parameters.
func (s *Service) Auth(ctx context.Context, rq AuthRequest) (AuthResponse, error) {
	if s.disableLLMAuth.Load() {
		return AuthResponse{}, nil
	}

	sp, err := s.index.SpecByID(rq.SpecID)
	if err != nil {
		return AuthResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", rq.SpecID), err)
	}

	if sp.Auth == nil {
		return AuthResponse{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to create dummy request: %w", err)
	}

	var info auth.Info
	if err := sp.Auth.Apply(req, &info); err != nil {
		return AuthResponse{}, fmt.Errorf("failed to apply auth: %w", err)
	}

	return AuthResponse{
		Token:       info.Headers["Authorization"],
		Headers:     info.Headers,
		QueryParams: info.QueryParams,
	}, nil
}
