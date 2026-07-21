package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import "net/http"

// Transport is an [http.RoundTripper] that applies authentication
// to every request before delegating to the underlying transport.
type Transport struct {
	Base http.RoundTripper
	Auth Authenticator
}

// RoundTrip implements [http.RoundTripper].
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.Auth.Apply(req, nil); err != nil {
		return nil, err
	}
	return t.Base.RoundTrip(req)
}

// newHTTPClient returns an [http.Client] that applies the given authentication
// to every outgoing request.
func newHTTPClient(auth Authenticator) *http.Client {
	return &http.Client{
		Transport: &Transport{
			Base: http.DefaultTransport,
			Auth: auth,
		},
	}
}
