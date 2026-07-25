package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

type authService struct {
	index           IndexReader
	llmAuthDisabled func() bool
}

func newAuthService(index IndexReader, llmAuthDisabled func() bool) *authService {
	return &authService{index: index, llmAuthDisabled: llmAuthDisabled}
}

// Auth retrieves authentication information for the spec identified by the
// request SpecID. It applies the spec's auth configuration and returns the
// resulting token, headers, and query parameters.
func (as *authService) Auth(ctx context.Context, rq AuthRequest) (AuthResponse, error) {
	if as.llmAuthDisabled() {
		return AuthResponse{}, nil
	}

	sp, err := as.index.SpecByID(rq.SpecID)
	if err != nil {
		return AuthResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. Use spec_list to find the correct spec ID.", rq.SpecID),
			err,
		)
	}

	if sp.Auth == nil {
		return AuthResponse{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		return AuthResponse{}, NewAuthError(
			"Failed to prepare the authentication request. This is an internal error.",
			err,
		)
	}

	var info auth.Info
	if err := sp.Auth.Apply(req, &info); err != nil {
		return AuthResponse{}, NewAuthError(
			"Failed to apply the authentication configuration for this spec. "+
				"Check that the auth credentials are valid and the auth server is reachable.",
			err,
		)
	}

	return AuthResponse{
		Token:       info.Headers["Authorization"],
		Headers:     info.Headers,
		QueryParams: info.QueryParams,
	}, nil
}
