package auth

import "net/http"

// ScriptAuthClient holds a script source used to generate authentication dynamically.
type ScriptAuthClient struct {
	Source string `yaml:"source" validate:"required"`
}

func NewScriptAuthClient(source string) *ScriptAuthClient {
	return &ScriptAuthClient{
		Source: source,
	}
}

func (c *ScriptAuthClient) Type() Type {
	return ScriptAuth
}

func (c *ScriptAuthClient) Apply(_ *http.Request) error {
	return nil
}

func (c *ScriptAuthClient) Validate() error {
	return authValidator.Struct(c)
}
