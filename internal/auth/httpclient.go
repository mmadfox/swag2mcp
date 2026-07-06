package auth

import (
	"net/http"
	"time"
)

// defaultHTTPTimeout and defaultHTTPClient are used for OAuth2 token exchange
// and similar internal auth flows.
var defaultHTTPTimeout = 30 * time.Second //nolint:gochecknoglobals // shared timeout

var defaultHTTPClient = &http.Client{Timeout: defaultHTTPTimeout} //nolint:gochecknoglobals // shared client
