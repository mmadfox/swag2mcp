package auth

import "net/http"

// NoAuthClient is an authenticator that performs no authentication.
type NoAuthClient struct{}

// NewNoAuthClient creates a new NoAuthClient.
func NewNoAuthClient() NoAuthClient {
	return NoAuthClient{}
}

func (NoAuthClient) New() error {
	return nil
}

func (NoAuthClient) Type() Type {
	return NoAuth
}

func (NoAuthClient) Apply(_ *http.Request, _ *Info) error {
	return nil
}

func (NoAuthClient) Validate() error {
	return nil
}
