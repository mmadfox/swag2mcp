package auth

import (
	"net/http"
	"time"
)

// defaultHTTPTimeout and defaultHTTPClient are reserved for future OAuth2
// token exchange and similar internal auth flows.
//
//nolint:unused,gochecknoglobals // reserved for future OAuth2 token exchange.
var defaultHTTPTimeout = 30 * time.Second

//nolint:unused,gochecknoglobals // reserved for future OAuth2 token exchange.
var defaultHTTPClient = &http.Client{Timeout: defaultHTTPTimeout}
