/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/env"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

const (
	// NoAuth represents no authentication.
	NoAuth Type = "none"
	// BasicAuth represents HTTP Basic authentication.
	BasicAuth Type = "basic"
	// BearerTokenAuth represents Bearer token authentication.
	BearerTokenAuth Type = "bearer"
	// DigestAuth represents HTTP Digest authentication.
	DigestAuth Type = "digest"
	// OAuth2ClientCredentials represents OAuth2 Client Credentials flow.
	OAuth2ClientCredentials Type = "oauth2-cc" //nolint:gosec // This is a type name, not a credential.
	// OAuth2Password represents OAuth2 Password grant flow.
	OAuth2Password Type = "oauth2-pwd" //nolint:gosec // This is a type name, not a credential.
	// APIKeyAuth represents API key authentication.
	APIKeyAuth Type = "api-key"
	// ScriptAuth represents authentication via an external script.
	ScriptAuth Type = "script"
	// HMACAuth represents HMAC-SHA256 signature authentication (Binance-style).
	HMACAuth Type = "hmac"
	// OpenIDConnect represents OIDC Discovery-based authentication.
	OpenIDConnect Type = "openid-connect"
)

const (
	headerAuthorization = "Authorization"
	headerValueBearer   = "Bearer "
	paramInQuery        = "query"
	// tokenRequestTimeout is the timeout for external HTTP requests
	// (token endpoints, digest challenges) and script execution.
	tokenRequestTimeout = 30 * time.Second
	// defaultExpiresIn is the fallback token expiry when the server omits expires_in.
	defaultExpiresIn = 3600
	// RequestFormatForm is the default form-encoded request format.
	RequestFormatForm = "form"
	// RequestFormatJSON is the JSON request format.
	RequestFormatJSON = "json"
	// contentTypeForm is the Content-Type for form-encoded requests.
	contentTypeForm = "application/x-www-form-urlencoded"
	// contentTypeJSON is the Content-Type for JSON requests.
	contentTypeJSON = "application/json"
	// grantTypePassword is the OAuth2 password grant type value.
	grantTypePassword = "password"
	// grantTypeClientCredentials is the OAuth2 client credentials grant type value.
	grantTypeClientCredentials = "client_credentials"
)

// Sentinel errors for common auth failure modes.
var (
	// ErrEmptyAccessToken is returned when a token endpoint returns an empty access_token.
	ErrEmptyAccessToken = errors.New("empty access_token in response")
)

//nolint:gochecknoglobals // Validator is stateless and safe to reuse.
var authValidator = validator.New(validator.WithRequiredStructEnabled())

// Type is the type of authentication used.
type Type string

// String returns the string representation of the auth type.
func (t Type) String() string {
	return string(t)
}

// Info holds the authentication details extracted during Apply.
type Info struct {
	Headers     map[string]string
	QueryParams map[string]string
}

// Authenticator is an interface for authenticating requests.
type Authenticator interface {
	New() error
	Type() Type
	Apply(req *http.Request, out *Info) error
	Validate() error
	// QueryParamNames returns the names of query parameters that this
	// authenticator injects into the request itself. These parameters are
	// considered satisfied by auth and must not be required from the caller.
	// It returns nil for authenticators that only set headers.
	QueryParamNames() []string
}

// TokenURLSetter is an optional interface for auth clients that have a configurable token URL.
type TokenURLSetter interface {
	SetTokenURL(url string)
}

// MockBaseURLSetter is an optional interface for auth clients that need a mock base URL.
type MockBaseURLSetter interface {
	SetMockBaseURL(url string)
}

func setAuthHeader(req *http.Request, out *Info, key, value string) {
	if value == "" {
		return
	}
	req.Header.Set(key, value)
	if out != nil {
		if out.Headers == nil {
			out.Headers = make(map[string]string)
		}
		out.Headers[key] = value
	}
}

func setAuthQuery(req *http.Request, out *Info, key, value string) {
	if value == "" {
		return
	}
	q := req.URL.Query()
	q.Set(key, value)
	req.URL.RawQuery = q.Encode()
	if out != nil {
		if out.QueryParams == nil {
			out.QueryParams = make(map[string]string)
		}
		out.QueryParams[key] = value
	}
}

// resolveEnv resolves $(VAR_NAME) patterns to environment variable values.
func resolveEnv(s string) string {
	return env.Parse(s)
}

// bearerToken returns the Authorization header value for a Bearer token.
func bearerToken(token string) string {
	return headerValueBearer + token
}

// oauth2TokenResponse is the JSON response from an OAuth2 token endpoint or auth script.
type oauth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	Token       string `json:"token,omitempty"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// doTokenRequest sends a POST to the token URL with the given body and content type.
func doTokenRequest(ctx context.Context, tokenURL, contentType string, body io.Reader) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, tokenRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, body)
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	client, err := httpclient.NewDefault()
	if err != nil {
		return nil, fmt.Errorf("create http client: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(respBody))
	}

	return resp, nil
}

// decodeTokenResponse parses an OAuth2 token response and returns the access token and expires_in.
func decodeTokenResponse(r io.Reader) (string, int, error) {
	var tr oauth2TokenResponse
	if err := json.NewDecoder(r).Decode(&tr); err != nil {
		return "", 0, fmt.Errorf("decode token response: %w", err)
	}

	token := tr.AccessToken
	if token == "" {
		token = tr.Token
	}
	if token == "" {
		return "", 0, ErrEmptyAccessToken
	}

	expiresIn := tr.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = defaultExpiresIn
	}

	return token, expiresIn, nil
}
