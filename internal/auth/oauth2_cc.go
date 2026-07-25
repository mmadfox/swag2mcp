package auth

import (
	"fmt"
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

	token, expiresIn, err := c.fetchToken()
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

func (c *OAuth2ClientCredentialsAuthClient) fetchToken() (string, int, error) {
	form := c.buildTokenForm()
	resp, err := doTokenRequest(c.TokenURL, form)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	return decodeTokenResponse(resp.Body)
}

func (c *OAuth2ClientCredentialsAuthClient) buildTokenForm() url.Values {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}
	if len(c.Scopes) > 0 {
		form.Set("scope", strings.Join(c.Scopes, " "))
	}
	return form
}

// SetTokenURL sets the token endpoint URL for the client credentials flow.
func (c *OAuth2ClientCredentialsAuthClient) SetTokenURL(url string) {
	c.TokenURL = url
}

// Validate checks that the ClientID, ClientSecret, and TokenURL fields are present and valid.
func (c *OAuth2ClientCredentialsAuthClient) Validate() error {
	return authValidator.Struct(c)
}
