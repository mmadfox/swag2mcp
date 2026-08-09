/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// OAuth2ClientCredentialsAuthClient holds configuration for the OAuth2 client credentials flow.
type OAuth2ClientCredentialsAuthClient struct {
	ClientID     string   `yaml:"client_id"        validate:"required"`
	ClientSecret string   `yaml:"client_secret"    validate:"required"`
	TokenURL     string   `yaml:"token_url"        validate:"required,url"`
	Scopes       []string `yaml:"scopes,omitempty"`

	// RequestFormat specifies the token request body format: "form" (default) or "json".
	RequestFormat string `yaml:"request_format,omitempty"`

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

// New resolves environment variables in ClientID and ClientSecret and returns nil.
func (c *OAuth2ClientCredentialsAuthClient) New() error {
	c.ClientID = resolveEnv(c.ClientID)
	c.ClientSecret = resolveEnv(c.ClientSecret)
	return nil
}

// Type returns the authentication type for OAuth2 client credentials flow.
func (c *OAuth2ClientCredentialsAuthClient) Type() Type {
	return OAuth2ClientCredentials
}

// Apply obtains a Bearer token via the client credentials grant and sets it on the request, caching the token until expiry.
func (c *OAuth2ClientCredentialsAuthClient) Apply(req *http.Request, out *Info) error {
	if token, ok := c.readCachedToken(); ok {
		setAuthHeader(req, out, headerAuthorization, bearerToken(token))
		return nil
	}

	token, expiresIn, err := c.fetchToken(req.Context())
	if err != nil {
		return fmt.Errorf("oauth2-cc: %w", err)
	}

	c.writeToken(token, expiresIn)
	setAuthHeader(req, out, headerAuthorization, bearerToken(token))
	return nil
}

func (c *OAuth2ClientCredentialsAuthClient) readCachedToken() (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		return c.token, true
	}
	return "", false
}

func (c *OAuth2ClientCredentialsAuthClient) writeToken(token string, expiresIn int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
	c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (c *OAuth2ClientCredentialsAuthClient) fetchToken(ctx context.Context) (string, int, error) {
	body, contentType := c.buildTokenBody()
	resp, err := doTokenRequest(ctx, c.TokenURL, contentType, body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return decodeTokenResponse(resp.Body)
}

func (c *OAuth2ClientCredentialsAuthClient) buildTokenBody() (io.Reader, string) {
	params := map[string]string{
		"grant_type":    grantTypeClientCredentials,
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
	}
	if len(c.Scopes) > 0 {
		params["scope"] = strings.Join(c.Scopes, " ")
	}

	if c.RequestFormat == RequestFormatJSON {
		data, _ := json.Marshal(params)
		return bytes.NewReader(data), contentTypeJSON
	}

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}
	return strings.NewReader(form.Encode()), contentTypeForm
}

// SetTokenURL sets the token endpoint URL for the client credentials flow.
func (c *OAuth2ClientCredentialsAuthClient) SetTokenURL(url string) {
	c.TokenURL = url
}

// Validate checks that the ClientID, ClientSecret, and TokenURL fields are present and valid.
func (c *OAuth2ClientCredentialsAuthClient) Validate() error {
	return authValidator.Struct(c)
}

// QueryParamNames returns nil since OAuth2 client credentials auth only sets a header.
func (c *OAuth2ClientCredentialsAuthClient) QueryParamNames() []string {
	return nil
}
