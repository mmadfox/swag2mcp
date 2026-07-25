package auth

import "net/http"

type NoAuthClient struct{}

func NewNoAuthClient() NoAuthClient {
	return NoAuthClient{}
}

func (NoAuthClient) New() error {
	return nil
}

func (NoAuthClient) Type() Type {
	return NoAuth
}

func (NoAuthClient) Apply(_ *http.Request) error {
	return nil
}

func (NoAuthClient) Validate() error {
	return nil
}
