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

const defaultExpiresIn = 3600

// OAuth2ClientCredentialsAuthClient holds configuration for the OAuth2 client credentials flow.
type OAuth2ClientCredentialsAuthClient struct {
	ClientID     string   `yaml:"client_id"        validate:"required"`
	ClientSecret string   `yaml:"client_secret"    validate:"required"`
	TokenURL     string   `yaml:"token_url"        validate:"required,url"`
	Scopes       []string `yaml:"scopes,omitempty"`

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

type oauth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *OAuth2ClientCredentialsAuthClient) New() error {
	c.ClientID = resolveEnv(c.ClientID)
	c.ClientSecret = resolveEnv(c.ClientSecret)
	return nil
}

func (c *OAuth2ClientCredentialsAuthClient) Type() Type {
	return OAuth2ClientCredentials
}

func (c *OAuth2ClientCredentialsAuthClient) Apply(req *http.Request, out *Info) error {
	c.mu.Lock()
	if c.token != "" && time.Now().Before(c.expiresAt) {
		setAuthHeader(req, out, "Authorization", "Bearer "+c.token)
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	token, expiresIn, err := c.fetchToken()
	if err != nil {
		return fmt.Errorf("oauth2-cc: %w", err)
	}

	c.mu.Lock()
	c.token = token
	c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	setAuthHeader(req, out, "Authorization", "Bearer "+c.token)
	c.mu.Unlock()
	return nil
}

func (c *OAuth2ClientCredentialsAuthClient) fetchToken() (string, int, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}
	if len(c.Scopes) > 0 {
		form.Set("scope", strings.Join(c.Scopes, " "))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) //nolint:mnd // token request timeout
	defer cancel()

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("create token request: %w", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cli, cliErr := httpclient.NewDefault()
	if cliErr != nil {
		return "", 0, fmt.Errorf("create http client: %w", cliErr)
	}
	resp, err := cli.Do(tokenReq)
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

func (c *OAuth2ClientCredentialsAuthClient) SetTokenURL(url string) {
	c.TokenURL = url
}

func (c *OAuth2ClientCredentialsAuthClient) Validate() error {
	return authValidator.Struct(c)
}
