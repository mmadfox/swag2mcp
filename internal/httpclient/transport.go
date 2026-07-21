package httpclient

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"net/http"
)

// randomizingTransport adds browser-like headers to requests.
type randomizingTransport struct {
	Base      http.RoundTripper
	UserAgent string
	Headers   map[string]string
	Cookies   []Cookie
}

// RoundTrip decorates the request with browser-like headers and cookies before dispatching it.
func (t *randomizingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" && t.UserAgent != "" {
		req.Header.Set("User-Agent", t.UserAgent)
	}
	for k, v := range t.Headers {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}
	for _, c := range t.Cookies {
		req.AddCookie(&http.Cookie{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Secure:   c.Secure,
			HttpOnly: c.HTTPOnly,
		})
	}
	return t.Base.RoundTrip(req)
}
