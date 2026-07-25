package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/env"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

// Type is the type of authentication used.
type Type string

// String returns the string representation of the auth type.
func (t Type) String() string {
	return string(t)
}

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
)

//nolint:gochecknoglobals // Validator is stateless and safe to reuse.
var authValidator = validator.New(validator.WithRequiredStructEnabled())

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

// doTokenRequest sends a form-encoded POST to the token URL and returns the response.
func doTokenRequest(tokenURL string, form url.Values) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tokenRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client, err := httpclient.NewDefault()
	if err != nil {
		return nil, fmt.Errorf("create http client: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
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
		return "", 0, errors.New("empty access_token in response")
	}

	expiresIn := tr.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = defaultExpiresIn
	}

	return token, expiresIn, nil
}
