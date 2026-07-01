package auth

import "net/http"

// Transport is an [http.RoundTripper] that applies authentication
// to every request before delegating to the underlying transport.
type Transport struct {
	Base http.RoundTripper
	Auth Authenticator
}

// RoundTrip implements [http.RoundTripper].
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.Auth.Apply(req); err != nil {
		return nil, err
	}
	return t.Base.RoundTrip(req)
}

// NewHTTPClient returns an [http.Client] that applies the given authentication
// to every outgoing request.
func NewHTTPClient(auth Authenticator) *http.Client {
	return &http.Client{
		Transport: &Transport{
			Base: http.DefaultTransport,
			Auth: auth,
		},
	}
}
