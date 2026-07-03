package config

import (
	"errors"
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

// Auth holds the parsed authentication configuration for a spec.
// It is populated during YAML unmarshalling by reading the "type" and "config" keys.
type Auth struct {
	Client auth.Authenticator
}

// authRaw is the intermediate YAML structure used to extract the type
// and the raw config node before dispatching to the concrete client.
type authRaw struct {
	Type   string    `yaml:"type"`
	Config yaml.Node `yaml:"config"`
}

// UnmarshalYAML implements yaml.Unmarshaler. It reads the "type" field,
// normalises underscores to hyphens, and decodes "config" into the
// corresponding auth client from the auth package.
func (a *Auth) UnmarshalYAML(value *yaml.Node) error {
	var raw authRaw
	if err := value.Decode(&raw); err != nil {
		return err
	}

	raw.Type = strings.ReplaceAll(raw.Type, "_", "-")

	switch raw.Type {
	case auth.NoAuth.String(), "":
		a.Client = auth.NewNoAuthClient()

	case auth.BasicAuth.String():
		var client auth.BasicAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	case auth.BearerTokenAuth.String():
		var client auth.BearerTokenAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	case auth.DigestAuth.String():
		var client auth.DigestAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	case auth.OAuth2ClientCredentials.String():
		var client auth.OAuth2ClientCredentialsAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	case auth.OAuth2Password.String():
		var client auth.OAuth2PasswordAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	case auth.APIKeyAuth.String():
		var client auth.APIKeyAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	case auth.ScriptAuth.String():
		var client auth.ScriptAuthClient
		if err := decodeConfig(&raw.Config, &client); err != nil {
			return err
		}
		a.Client = &client

	default:
		return fmt.Errorf("unsupported auth type %q", raw.Type)
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (a Auth) MarshalYAML() (any, error) {
	if a.Client == nil || a.Client.Type() == auth.NoAuth {
		return nil, nil
	}

	configData, err := yaml.Marshal(a.Client)
	if err != nil {
		return nil, err
	}

	var config any
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return map[string]any{
		"type":   string(a.Client.Type()),
		"config": config,
	}, nil
}

// decodeConfig decodes a YAML node into the provided config struct.
// It returns an error if the node is missing (Kind == 0) for types that require a config block.
func decodeConfig(node *yaml.Node, cfg any) error {
	if node.Kind == 0 {
		return errors.New("auth config is required for this type")
	}
	if err := node.Decode(cfg); err != nil {
		return fmt.Errorf("invalid auth config: %w", err)
	}
	return nil
}
