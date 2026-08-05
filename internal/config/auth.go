/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"errors"
	"fmt"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/mmadfox/swag2mcp/internal/auth"
)

// authDecoder decodes a YAML config node into an auth.Authenticator.
type authDecoder func(*yaml.Node) (auth.Authenticator, error)

// authDecoders maps auth type strings to their config decoders.
var authDecoders = map[string]authDecoder{ //nolint:gochecknoglobals // Static registry.
	"":                                    decodeNoAuth,
	auth.NoAuth.String():                  decodeNoAuth,
	auth.BasicAuth.String():               decodeBasicAuth,
	auth.BearerTokenAuth.String():         decodeBearerAuth,
	auth.DigestAuth.String():              decodeDigestAuth,
	auth.OAuth2ClientCredentials.String(): decodeOAuth2CC,
	auth.OAuth2Password.String():          decodeOAuth2Pwd,
	auth.APIKeyAuth.String():              decodeAPIKeyAuth,
	auth.ScriptAuth.String():              decodeScriptAuth,
	auth.HMACAuth.String():                decodeHMACAuth,
}

// Auth holds the parsed authentication configuration for a spec.
// It is populated during YAML unmarshalling by reading the "type" and "config" keys.
type Auth struct {
	Client auth.Authenticator `yaml:"-"`
}

// MarshalYAML implements yaml.Marshaler.
// It serialises the auth config as "type" and "config" keys matching the input format.
func (a Auth) MarshalYAML() (any, error) {
	if a.Client == nil {
		return nil, nil
	}
	return map[string]any{
		"type":   a.Client.Type().String(),
		"config": a.Client,
	}, nil
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

	fn, ok := authDecoders[raw.Type]
	if !ok {
		return fmt.Errorf("unsupported auth type %q", raw.Type)
	}

	client, err := fn(&raw.Config)
	if err != nil {
		return err
	}
	a.Client = client
	return nil
}

// decodeNoAuth returns a no-auth client. The config node is ignored.
func decodeNoAuth(_ *yaml.Node) (auth.Authenticator, error) {
	return auth.NewNoAuthClient(), nil
}

// decodeBasicAuth decodes a basic auth config block into a BasicAuthClient.
func decodeBasicAuth(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.BasicAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeBearerAuth decodes a bearer token auth config block into a BearerTokenAuthClient.
func decodeBearerAuth(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.BearerTokenAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeDigestAuth decodes a digest auth config block into a DigestAuthClient.
func decodeDigestAuth(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.DigestAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeOAuth2CC decodes an OAuth2 client credentials config block
// into an OAuth2ClientCredentialsAuthClient.
func decodeOAuth2CC(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.OAuth2ClientCredentialsAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeOAuth2Pwd decodes an OAuth2 password grant config block
// into an OAuth2PasswordAuthClient.
func decodeOAuth2Pwd(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.OAuth2PasswordAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeAPIKeyAuth decodes an API key auth config block into an APIKeyAuthClient.
func decodeAPIKeyAuth(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.APIKeyAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeScriptAuth decodes a script auth config block into a ScriptAuthClient.
func decodeScriptAuth(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.ScriptAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// decodeHMACAuth decodes an HMAC auth config block into an HMACAuthClient.
func decodeHMACAuth(node *yaml.Node) (auth.Authenticator, error) {
	var client auth.HMACAuthClient
	if err := decodeConfig(node, &client); err != nil {
		return nil, err
	}
	return &client, nil
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
