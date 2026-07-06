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
)

// OAuth2PasswordAuthClient holds configuration for the OAuth2 resource owner password flow.
type OAuth2PasswordAuthClient struct {
	Username     string   `yaml:"username"         validate:"required"`
	Password     string   `yaml:"password"         validate:"required"`
	ClientID     string   `yaml:"client_id"        validate:"required"`
	ClientSecret string   `yaml:"client_secret"    validate:"required"`
	TokenURL     string   `yaml:"token_url"        validate:"required,url"`
	Scopes       []string `yaml:"scopes,omitempty"`

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

func (c *OAuth2PasswordAuthClient) New() error {
	c.Username = resolveEnv(c.Username)
	c.Password = resolveEnv(c.Password)
	c.ClientID = resolveEnv(c.ClientID)
	c.ClientSecret = resolveEnv(c.ClientSecret)
	return nil
}

func (c *OAuth2PasswordAuthClient) Type() Type {
	return OAuth2Password
}

func (c *OAuth2PasswordAuthClient) Apply(req *http.Request, out *Info) error {
	c.mu.Lock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		setAuthHeader(req, out, "Authorization", "Bearer "+c.token)
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	token, expiresIn, err := c.fetchToken()
	if err != nil {
		return fmt.Errorf("oauth2-pwd: %w", err)
	}

	c.mu.Lock()
	c.token = token
	c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	setAuthHeader(req, out, "Authorization", "Bearer "+c.token)
	c.mu.Unlock()
	return nil
}

func (c *OAuth2PasswordAuthClient) fetchToken() (string, int, error) {
	form := url.Values{
		"grant_type":    {"password"},
		"username":      {c.Username},
		"password":      {c.Password},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}
	if len(c.Scopes) > 0 {
		form.Set("scope", strings.Join(c.Scopes, " "))
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("create token request: %w", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := defaultHTTPClient.Do(tokenReq)
	if err != nil {
		return "", 0, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tr oauth2TokenResponse
	if decodeErr := json.NewDecoder(resp.Body).Decode(&tr); decodeErr != nil {
		return "", 0, fmt.Errorf("decode token response: %w", decodeErr)
	}

	if tr.AccessToken == "" {
		return "", 0, errors.New("empty access_token in response")
	}

	expiresIn := tr.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = defaultExpiresIn
	}

	return tr.AccessToken, expiresIn, nil
}

func (c *OAuth2PasswordAuthClient) Validate() error {
	return authValidator.Struct(c)
}
