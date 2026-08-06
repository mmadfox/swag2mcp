/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

type authService struct {
	index           IndexReader
	llmAuthDisabled func() bool
	log             *slog.Logger
}

func newAuthService(index IndexReader, llmAuthDisabled func() bool, log *slog.Logger) *authService {
	return &authService{index: index, llmAuthDisabled: llmAuthDisabled, log: log}
}

func (as *authService) Auth(ctx context.Context, rq AuthRequest) (AuthResponse, error) {
	if as.llmAuthDisabled() {
		return AuthResponse{}, nil
	}

	sp, err := as.index.SpecByID(rq.SpecID)
	if err != nil {
		as.log.ErrorContext(ctx, "auth failed: spec not found", "spec_id", rq.SpecID, "error", err)
		return AuthResponse{}, NewSpecNotFoundError(rq.SpecID, err)
	}

	if sp.Auth == nil {
		return AuthResponse{}, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		as.log.ErrorContext(ctx, "auth failed: create request", "spec_id", rq.SpecID, "error", err)
		return AuthResponse{}, NewAuthError(
			"Failed to prepare the authentication request. This is an internal error.",
			err,
		)
	}

	var info auth.Info
	if err := sp.Auth.Apply(req, &info); err != nil {
		as.log.ErrorContext(ctx, "auth failed: apply auth", "spec_id", rq.SpecID, "spec", sp.Domain, "error", err)
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
