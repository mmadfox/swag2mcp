package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
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

	token, expiresIn, err := c.fetchToken()
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

func (c *OAuth2PasswordAuthClient) fetchToken() (string, int, error) {
	form := c.buildTokenForm()
	resp, err := doTokenRequest(c.TokenURL, form)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return decodeTokenResponse(resp.Body)
}

func (c *OAuth2PasswordAuthClient) buildTokenForm() url.Values {
	const grantTypePassword = "password"
	form := url.Values{
		"grant_type": {grantTypePassword},
		"username":   {c.Username},
		"password":   {c.Password}, //nolint:goconst // Form field name, not the grant type value.
		"client_id":  {c.ClientID},
	}
	if c.ClientSecret != "" {
		form.Set("client_secret", c.ClientSecret)
	}
	if len(c.Scopes) > 0 {
		form.Set("scope", strings.Join(c.Scopes, " "))
	}
	return form
}

// SetTokenURL sets the token endpoint URL for the password grant flow.
func (c *OAuth2PasswordAuthClient) SetTokenURL(url string) {
	c.TokenURL = url
}

// Validate checks that the Username, Password, ClientID, and TokenURL fields are present and valid.
func (c *OAuth2PasswordAuthClient) Validate() error {
	return authValidator.Struct(c)
}
