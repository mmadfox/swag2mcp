package auth

import "net/http"

// OAuth2PasswordAuthClient holds configuration for the OAuth2 resource owner password flow.
type OAuth2PasswordAuthClient struct {
	Username     string   `yaml:"username"         validate:"required"`
	Password     string   `yaml:"password"         validate:"required"`
	ClientID     string   `yaml:"client_id"        validate:"required"`
	ClientSecret string   `yaml:"client_secret"    validate:"required"`
	TokenURL     string   `yaml:"token_url"        validate:"required,url"`
	Scopes       []string `yaml:"scopes,omitempty"`
}

func NewOAuth2PasswordAuthClient(
	username, password, clientID, clientSecret, tokenURL string,
	scopes []string,
) *OAuth2PasswordAuthClient {
	return &OAuth2PasswordAuthClient{
		Username:     username,
		Password:     password,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}
}

func (c *OAuth2PasswordAuthClient) Type() Type {
	return OAuth2Password
}

func (c *OAuth2PasswordAuthClient) Apply(_ *http.Request) error {
	return nil
}

func (c *OAuth2PasswordAuthClient) Validate() error {
	return authValidator.Struct(c)
}
