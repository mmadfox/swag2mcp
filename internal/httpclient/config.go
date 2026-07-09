package httpclient

import "time"

// ProxyConfig holds proxy connection settings.
type ProxyConfig struct {
	URL      string
	Username string
	Password string
	Bypass   []string
}

// Cookie represents an HTTP cookie.
type Cookie struct {
	Name     string
	Value    string
	Domain   string
	Path     string
	Secure   bool
	HTTPOnly bool
}

// Config holds all settings for creating an HTTP client.
type Config struct {
	Proxy           *ProxyConfig
	Timeout         time.Duration
	FollowRedirects *bool
	MaxRedirects    *int
	MaxResponseSize *int
	Randomize       bool
	UserAgent       string
	Headers         map[string]string
	Cookies         []Cookie
}
