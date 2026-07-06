package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

type (
	AuthRequest struct {
		DomainID string `json:"domainId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the spec/domain to get an auth token for"`
	}

	AuthResponse struct {
		Token       string            `json:"token"`
		Headers     map[string]string `json:"headers,omitempty"`
		QueryParams map[string]string `json:"queryParams,omitempty"`
	}
)

func (s *Service) Auth(ctx context.Context, req AuthRequest) (AuthResponse, error) {
	if s.disableLLMAuth.Load() {
		return AuthResponse{}, nil
	}

	spec, err := s.index.SpecByID(req.DomainID)
	if err != nil {
		return AuthResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", req.DomainID), err)
	}

	if spec.Auth == nil {
		return AuthResponse{}, nil
	}

	dummyReq, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if reqErr != nil {
		return AuthResponse{}, fmt.Errorf("failed to create dummy request: %w", reqErr)
	}

	var info auth.Info
	if aErr := spec.Auth.Apply(dummyReq, &info); aErr != nil {
		return AuthResponse{}, fmt.Errorf("failed to apply auth: %w", aErr)
	}

	return AuthResponse{
		Token:       info.Headers["Authorization"],
		Headers:     info.Headers,
		QueryParams: info.QueryParams,
	}, nil
}
