package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import "net/http"

// NoAuthClient is an authenticator that performs no authentication.
type NoAuthClient struct{}

// NewNoAuthClient creates a new NoAuthClient.
func NewNoAuthClient() NoAuthClient {
	return NoAuthClient{}
}

// New initializes the NoAuthClient. It always returns nil.
func (NoAuthClient) New() error {
	return nil
}

// Type returns the authentication type for no authentication.
func (NoAuthClient) Type() Type {
	return NoAuth
}

// Apply performs no authentication. It always returns nil.
func (NoAuthClient) Apply(_ *http.Request, _ *Info) error {
	return nil
}

// Validate performs no validation. It always returns nil.
func (NoAuthClient) Validate() error {
	return nil
}
