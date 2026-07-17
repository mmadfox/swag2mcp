package auth

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
