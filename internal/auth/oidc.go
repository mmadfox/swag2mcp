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
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

// oidcDiscoveryPath is the well-known OIDC discovery document path.
const oidcDiscoveryPath = "/.well-known/openid-configuration"

// oidcConfig is the subset of the OIDC discovery document needed to obtain a token.
type oidcConfig struct {
	TokenEndpoint string `json:"token_endpoint"`
}

// OIDCAuthClient holds configuration for OIDC Discovery-based authentication.
// It discovers the token endpoint from the issuer, obtains a Bearer token via
// the client credentials grant, and caches it until expiry.
type OIDCAuthClient struct {
	Issuer       string   `yaml:"issuer"           validate:"required,url"`
	ClientID     string   `yaml:"client_id"        validate:"required"`
	ClientSecret string   `yaml:"client_secret"    validate:"required"`
	Scopes       []string `yaml:"scopes,omitempty"`

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

// New resolves environment variable references in the issuer and credentials.
func (c *OIDCAuthClient) New() error {
	c.Issuer = resolveEnv(c.Issuer)
	c.ClientID = resolveEnv(c.ClientID)
	c.ClientSecret = resolveEnv(c.ClientSecret)
	return nil
}

// Type returns the authentication type for OIDC Discovery auth.
func (c *OIDCAuthClient) Type() Type {
	return OpenIDConnect
}

// Apply obtains a Bearer token via OIDC Discovery and sets it on the request,
// caching the token until expiry.
func (c *OIDCAuthClient) Apply(req *http.Request, out *Info) error {
	if token, ok := c.readCachedToken(); ok {
		setAuthHeader(req, out, headerAuthorization, bearerToken(token))
		return nil
	}

	token, expiresIn, err := c.fetchToken(req.Context())
	if err != nil {
		return fmt.Errorf("openid-connect: %w", err)
	}

	c.writeToken(token, expiresIn)
	setAuthHeader(req, out, headerAuthorization, bearerToken(token))
	return nil
}

func (c *OIDCAuthClient) readCachedToken() (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		return c.token, true
	}
	return "", false
}

func (c *OIDCAuthClient) writeToken(token string, expiresIn int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
	c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// fetchToken discovers the token endpoint from the issuer and exchanges the
// client credentials for a Bearer token.
func (c *OIDCAuthClient) fetchToken(ctx context.Context) (string, int, error) {
	tokenURL, err := c.discoverTokenEndpoint(ctx)
	if err != nil {
		return "", 0, err
	}

	body, contentType := c.buildTokenBody()
	resp, err := doTokenRequest(ctx, tokenURL, contentType, body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return decodeTokenResponse(resp.Body)
}

// discoverTokenEndpoint fetches the OIDC discovery document and returns the
// token endpoint URL.
func (c *OIDCAuthClient) discoverTokenEndpoint(ctx context.Context) (string, error) {
	issuer := strings.TrimRight(c.Issuer, "/")
	discoveryURL := issuer + oidcDiscoveryPath

	ctx, cancel := context.WithTimeout(ctx, tokenRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return "", fmt.Errorf("create discovery request: %w", err)
	}

	client, err := httpclient.NewDefault()
	if err != nil {
		return "", fmt.Errorf("create http client: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("discovery request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("discovery endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var cfg oidcConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return "", fmt.Errorf("decode discovery response: %w", err)
	}
	if cfg.TokenEndpoint == "" {
		return "", errors.New("discovery response missing token_endpoint")
	}
	return cfg.TokenEndpoint, nil
}

func (c *OIDCAuthClient) buildTokenBody() (io.Reader, string) {
	params := url.Values{}
	params.Set("grant_type", grantTypeClientCredentials)
	params.Set("client_id", c.ClientID)
	params.Set("client_secret", c.ClientSecret)
	if len(c.Scopes) > 0 {
		params.Set("scope", strings.Join(c.Scopes, " "))
	}
	return strings.NewReader(params.Encode()), contentTypeForm
}

// Validate checks that the Issuer, ClientID, and ClientSecret fields are present and valid.
func (c *OIDCAuthClient) Validate() error {
	return authValidator.Struct(c)
}

// QueryParamNames returns nil since OIDC auth only sets a header.
func (c *OIDCAuthClient) QueryParamNames() []string {
	return nil
}
