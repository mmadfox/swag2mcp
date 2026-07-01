package auth

import "net/http"

// OAuth2ClientCredentialsAuthClient holds configuration for the OAuth2 client credentials flow.
type OAuth2ClientCredentialsAuthClient struct {
	ClientID     string   `yaml:"client_id"        validate:"required"`
	ClientSecret string   `yaml:"client_secret"    validate:"required"`
	TokenURL     string   `yaml:"token_url"        validate:"required,url"`
	Scopes       []string `yaml:"scopes,omitempty"`
}

func (c *OAuth2ClientCredentialsAuthClient) New() error {
	c.ClientID = resolveEnv(c.ClientID)
	c.ClientSecret = resolveEnv(c.ClientSecret)
	return nil
}

func (c *OAuth2ClientCredentialsAuthClient) Type() Type {
	return OAuth2ClientCredentials
}

func (c *OAuth2ClientCredentialsAuthClient) Apply(_ *http.Request) error {
	return nil
}

func (c *OAuth2ClientCredentialsAuthClient) Validate() error {
	return authValidator.Struct(c)
}
