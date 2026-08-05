/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/auth"
)

const (
	// authTypeJWKS is the JWT auth type for JWKS (JSON Web Key Set) verification.
	// Used by newJWTVerifier to dispatch to newJWKSVerifier.
	authTypeJWKS = "jwks"
	// authTypeOIDC is the JWT auth type for OIDC Discovery verification.
	// Used by newJWTVerifier to dispatch to newOIDCVerifier.
	authTypeOIDC = "oidc"
	// authTypeIntrospection is the JWT auth type for token introspection.
	// Used by newJWTVerifier to dispatch to newIntrospectionVerifier.
	authTypeIntrospection = "introspection"
	// jwtHTTPTimeout is the HTTP client timeout for OIDC discovery and
	// introspection requests made by JWT verifiers.
	jwtHTTPTimeout = 10 * time.Second
)

// JWTConfig holds configuration for JWT-based MCP auth.
type JWTConfig struct {
	Type             string
	JWKSURL          string
	Issuer           string
	Audience         string
	IntrospectionURL string
	ClientID         string
	ClientSecret     string
}

type introspectionResponse struct {
	Active    bool   `json:"active"`
	Sub       string `json:"sub,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Scope     string `json:"scope,omitempty"`
	TokenType string `json:"token_type,omitempty"`
}

type oidcConfig struct {
	Issuer  string `json:"issuer"`
	JWKSURI string `json:"jwks_uri"`
}

// newJWTVerifier creates a TokenVerifier from JWTConfig.
func newJWTVerifier(config *JWTConfig) auth.TokenVerifier {
	switch config.Type {
	case authTypeJWKS:
		return newJWKSVerifier(config)
	case authTypeOIDC:
		return newOIDCVerifier(config)
	case authTypeIntrospection:
		return newIntrospectionVerifier(config)
	default:
		return func(_ context.Context, _ string, _ *http.Request) (*auth.TokenInfo, error) {
			return nil, fmt.Errorf("unsupported JWT auth type: %s", config.Type)
		}
	}
}

// newJWKSVerifier creates a JWKS-based token verifier. It launches a background
// goroutine that periodically refreshes the JWK Set from the remote URL.
// The goroutine lives until the process exits.
func newJWKSVerifier(config *JWTConfig) auth.TokenVerifier {
	k, err := keyfunc.NewDefault([]string{config.JWKSURL})
	if err != nil {
		return func(_ context.Context, _ string, _ *http.Request) (*auth.TokenInfo, error) {
			return nil, fmt.Errorf("failed to create JWKS keyfunc: %w", err)
		}
	}

	return func(ctx context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
		parsed, err := jwt.Parse(token, k.KeyfuncCtx(ctx))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", auth.ErrInvalidToken, err)
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("%w: invalid claims format", auth.ErrInvalidToken)
		}

		if err := validateJWTClaims(claims, config); err != nil {
			return nil, err
		}

		sub, _ := claims.GetSubject()
		exp, err := claims.GetExpirationTime()
		if err != nil {
			exp = nil
		}

		info := &auth.TokenInfo{
			UserID: sub,
		}
		if exp != nil {
			info.Expiration = exp.Time
		}

		return info, nil
	}
}

func validateJWTClaims(claims jwt.MapClaims, config *JWTConfig) error {
	if config.Issuer != "" {
		iss, ok := claims["iss"].(string)
		if !ok || iss != config.Issuer {
			return fmt.Errorf("%w: invalid issuer", auth.ErrInvalidToken)
		}
	}

	if config.Audience != "" {
		aud, err := claims.GetAudience()
		if err != nil || !slices.Contains(aud, config.Audience) {
			return fmt.Errorf("%w: invalid audience", auth.ErrInvalidToken)
		}
	}

	return nil
}

// newOIDCVerifier creates an OIDC Discovery-based token verifier.
// It fetches the OIDC configuration once at startup (using [context.Background]
// since this runs during server initialization, not per-request), extracts the
// JWKS URI, and delegates to newJWKSVerifier.
func newOIDCVerifier(config *JWTConfig) auth.TokenVerifier {
	httpClient := &http.Client{Timeout: jwtHTTPTimeout}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		config.Issuer+"/.well-known/openid-configuration",
		nil,
	)
	if err != nil {
		return func(_ context.Context, _ string, _ *http.Request) (*auth.TokenInfo, error) {
			return nil, fmt.Errorf("failed to create OIDC discovery request: %w", err)
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return func(_ context.Context, _ string, _ *http.Request) (*auth.TokenInfo, error) {
			return nil, fmt.Errorf("failed to fetch OIDC configuration: %w", err)
		}
	}
	defer resp.Body.Close()

	var oidc oidcConfig
	if err := json.NewDecoder(resp.Body).Decode(&oidc); err != nil {
		return func(_ context.Context, _ string, _ *http.Request) (*auth.TokenInfo, error) {
			return nil, fmt.Errorf("failed to decode OIDC configuration: %w", err)
		}
	}

	jwksConfig := *config
	jwksConfig.JWKSURL = oidc.JWKSURI
	return newJWKSVerifier(&jwksConfig)
}

func newIntrospectionVerifier(config *JWTConfig) auth.TokenVerifier {
	httpClient := &http.Client{Timeout: jwtHTTPTimeout}

	return func(ctx context.Context, token string, _ *http.Request) (*auth.TokenInfo, error) {
		body := []byte("token=" + token + "&token_type_hint=access_token")
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.IntrospectionURL, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create introspection request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(config.ClientID, config.ClientSecret)

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%w: introspection request failed: %w", auth.ErrInvalidToken, err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to read introspection response: %w", auth.ErrInvalidToken, err)
		}

		var introResp introspectionResponse
		if err := json.Unmarshal(data, &introResp); err != nil {
			return nil, fmt.Errorf("%w: failed to decode introspection response: %w", auth.ErrInvalidToken, err)
		}

		if !introResp.Active {
			return nil, fmt.Errorf("%w: token is not active", auth.ErrInvalidToken)
		}

		info := &auth.TokenInfo{
			UserID: introResp.Sub,
		}
		if introResp.Exp > 0 {
			info.Expiration = time.Unix(introResp.Exp, 0)
		}

		return info, nil
	}
}
