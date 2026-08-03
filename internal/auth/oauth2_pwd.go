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

// OAuth2PasswordAuthClient holds configuration for the OAuth2 resource owner password flow.
type OAuth2PasswordAuthClient struct {
	Username string   `yaml:"username"         validate:"required"`
	Password string   `yaml:"password"         validate:"required"`
	ClientID string   `yaml:"client_id"        validate:"required"`
	TokenURL string   `yaml:"token_url"        validate:"required,url"`
	Scopes   []string `yaml:"scopes,omitempty"`

	ClientSecret string `yaml:"client_secret,omitempty"`

	// RequestFormat specifies the token request body format: "form" (default) or "json".
	RequestFormat string `yaml:"request_format,omitempty"`

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

// New resolves environment variables in Username, Password, ClientID, and ClientSecret and returns nil.
func (c *OAuth2PasswordAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	c.ClientID = resolveEnv(c.ClientID)
	c.ClientSecret = resolveEnv(c.ClientSecret)
	return nil
}

// Type returns the authentication type for OAuth2 password grant flow.
func (c *OAuth2PasswordAuthClient) Type() Type {
	return OAuth2Password
}

// Apply obtains a Bearer token via the resource owner password grant and sets it on the request, caching the token until expiry.
func (c *OAuth2PasswordAuthClient) Apply(req *http.Request, out *Info) error {
	if token, ok := c.readCachedToken(); ok {
		setAuthHeader(req, out, headerAuthorization, bearerToken(token))
		return nil
	}

	token, expiresIn, err := c.fetchToken(req.Context())
	if err != nil {
		return fmt.Errorf("oauth2-pwd: %w", err)
	}

	c.writeToken(token, expiresIn)
	setAuthHeader(req, out, headerAuthorization, bearerToken(token))
	return nil
}

func (c *OAuth2PasswordAuthClient) readCachedToken() (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		return c.token, true
	}
	return "", false
}

func (c *OAuth2PasswordAuthClient) writeToken(token string, expiresIn int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
	c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (c *OAuth2PasswordAuthClient) fetchToken(ctx context.Context) (string, int, error) {
	body, contentType := c.buildTokenBody()
	resp, err := doTokenRequest(ctx, c.TokenURL, contentType, body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return decodeTokenResponse(resp.Body)
}

func (c *OAuth2PasswordAuthClient) buildTokenBody() (io.Reader, string) {
	params := map[string]string{
		"grant_type": grantTypePassword,
		"username":   c.Username,
		"password":   c.Password, //nolint:goconst // Form field name, not the grant type value.
		"client_id":  c.ClientID,
	}
	if c.ClientSecret != "" {
		params["client_secret"] = c.ClientSecret
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

// SetTokenURL sets the token endpoint URL for the password grant flow.
func (c *OAuth2PasswordAuthClient) SetTokenURL(url string) {
	c.TokenURL = url
}

// Validate checks that the Username, Password, ClientID, and TokenURL fields are present and valid.
func (c *OAuth2PasswordAuthClient) Validate() error {
	return authValidator.Struct(c)
}
