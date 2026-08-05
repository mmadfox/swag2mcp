/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package config

// MockAuthConfig holds port configuration for mock auth servers.
type MockAuthConfig struct {
	OAuth2Port int `yaml:"oauth2_port,omitempty" validate:"omitempty,min=1024,max=65535"`
	DigestPort int `yaml:"digest_port,omitempty" validate:"omitempty,min=1024,max=65535"`
	HMACPort   int `yaml:"hmac_port,omitempty"   validate:"omitempty,min=1024,max=65535"`
	JWTPort    int `yaml:"jwt_port,omitempty"    validate:"omitempty,min=1024,max=65535"`
}
