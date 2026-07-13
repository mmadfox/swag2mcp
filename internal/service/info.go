package service

import (
	"context"
	"maps"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
)

// InfoSnapshot is a point-in-time snapshot of the service state,
// computed once after Bootstrap and served on every Info call.
type InfoSnapshot struct {
	Version    string         `json:"version"`
	Workspace  string         `json:"workspace"`
	Uptime     time.Duration  `json:"-"`
	Specs      SpecsSummary   `json:"specs"`
	HTTPClient HTTPClientInfo `json:"http_client"`
	MCP        MCPInfo        `json:"mcp"`
	Auth       AuthInfo       `json:"auth"`
	Mock       MockInfo       `json:"mock"`
}

// InfoResponse holds the complete runtime and configuration summary.
type InfoResponse struct {
	Version    string         `json:"version"`
	Workspace  string         `json:"workspace"`
	Uptime     string         `json:"uptime"`
	Specs      SpecsSummary   `json:"specs"`
	HTTPClient HTTPClientInfo `json:"http_client"`
	MCP        MCPInfo        `json:"mcp"`
	Auth       AuthInfo       `json:"auth"`
	Mock       MockInfo       `json:"mock"`
}

// SpecsSummary holds aggregate spec statistics.
type SpecsSummary struct {
	Total       int `json:"total"`
	Active      int `json:"active"`
	Disabled    int `json:"disabled"`
	Collections int `json:"collections"`
	Endpoints   int `json:"endpoints"`
}

// HTTPClientInfo holds the effective HTTP client configuration.
type HTTPClientInfo struct {
	Randomize       bool              `json:"randomize"`
	UserAgent       string            `json:"user_agent,omitempty"`
	Timeout         string            `json:"timeout,omitempty"`
	FollowRedirects *bool             `json:"follow_redirects,omitempty"`
	MaxRedirects    *int              `json:"max_redirects,omitempty"`
	MaxResponseSize int               `json:"max_response_size"`
	Proxy           *ProxyInfo        `json:"proxy,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Cookies         []CookieInfo      `json:"cookies,omitempty"`
}

// ProxyInfo holds proxy configuration details.
type ProxyInfo struct {
	URL      string   `json:"url"`
	Username string   `json:"username,omitempty"`
	Bypass   []string `json:"bypass,omitempty"`
}

// CookieInfo holds a single cookie configuration.
type CookieInfo struct {
	Name     string `json:"name"`
	Domain   string `json:"domain,omitempty"`
	Path     string `json:"path,omitempty"`
	Secure   bool   `json:"secure"`
	HTTPOnly bool   `json:"http_only"`
}

// MCPInfo holds the MCP server configuration.
type MCPInfo struct {
	Transport   string `json:"transport"`
	Addr        string `json:"addr,omitempty"`
	Path        string `json:"path,omitempty"`
	AuthEnabled bool   `json:"auth_enabled"`
}

// AuthInfo holds authentication method information.
type AuthInfo struct {
	Methods []string `json:"methods,omitempty"`
}

// MockInfo holds mock server configuration.
type MockInfo struct {
	Enabled bool `json:"enabled"`
}

// Info returns a comprehensive summary of the current service state from the snapshot.
func (s *Service) Info(_ context.Context) (InfoResponse, error) {
	snap, ok := s.snapshot.Load().(*InfoSnapshot)
	if ok && snap != nil {
		return InfoResponse{
			Version:    snap.Version,
			Workspace:  snap.Workspace,
			Uptime:     time.Since(s.startedAt).Round(time.Second).String(),
			Specs:      snap.Specs,
			HTTPClient: snap.HTTPClient,
			MCP:        snap.MCP,
			Auth:       snap.Auth,
			Mock:       snap.Mock,
		}, nil
	}

	return InfoResponse{
		Version:   s.version,
		Workspace: s.ws.Root(),
		Uptime:    time.Since(s.startedAt).Round(time.Second).String(),
	}, nil
}

// buildSnapshot computes a point-in-time snapshot of the service state.
func (s *Service) buildSnapshot() {
	snap := &InfoSnapshot{
		Version:   s.version,
		Workspace: s.ws.Root(),
	}

	if s.config != nil {
		snap.Specs = s.buildSpecsSummary(s.config)
		snap.HTTPClient = s.buildHTTPClientInfo(s.config)
		snap.MCP = s.buildMCPInfo(s.config)
		snap.Auth = s.buildAuthInfo(s.config)
		snap.Mock = MockInfo{Enabled: s.config.MockEnabled}
	}

	s.snapshot.Store(snap)
}

func (s *Service) buildSpecsSummary(cfg *config.Config) SpecsSummary {
	summary := SpecsSummary{}

	if cfg != nil {
		for i := range cfg.Specs {
			spec := &cfg.Specs[i]
			summary.Total++
			if spec.Disable {
				summary.Disabled++
			} else {
				summary.Active++
			}
		}
	}

	specs := s.index.AllSpecs()
	for _, spec := range specs {
		summary.Collections += spec.Stats.Collections
		summary.Endpoints += spec.Stats.Methods
	}

	return summary
}

func (s *Service) buildHTTPClientInfo(cfg *config.Config) HTTPClientInfo {
	info := HTTPClientInfo{
		MaxResponseSize: s.maxResponseSize,
	}

	if cfg == nil || cfg.HTTPClient == nil {
		return info
	}

	g := cfg.HTTPClient
	info.Randomize = s.httpClientConfig.Randomize
	info.UserAgent = s.httpClientConfig.UserAgent
	if info.UserAgent == "" {
		info.UserAgent = g.UserAgent
	}
	if g.Timeout > 0 {
		info.Timeout = g.Timeout.String()
	}
	info.FollowRedirects = g.FollowRedirects
	info.MaxRedirects = g.MaxRedirects
	if g.MaxResponseSize != nil {
		info.MaxResponseSize = *g.MaxResponseSize
	}

	if len(g.Headers) > 0 {
		info.Headers = make(map[string]string, len(g.Headers))
		maps.Copy(info.Headers, g.Headers)
	}

	if len(g.Cookies) > 0 {
		info.Cookies = make([]CookieInfo, 0, len(g.Cookies))
		for _, c := range g.Cookies {
			info.Cookies = append(info.Cookies, CookieInfo{
				Name:     c.Name,
				Domain:   c.Domain,
				Path:     c.Path,
				Secure:   c.Secure,
				HTTPOnly: c.HTTPOnly,
			})
		}
	}

	if g.Proxy != nil {
		info.Proxy = &ProxyInfo{
			URL:      g.Proxy.URL,
			Username: g.Proxy.Username,
			Bypass:   append([]string{}, g.Proxy.Bypass...),
		}
	}

	return info
}

func (s *Service) buildMCPInfo(cfg *config.Config) MCPInfo {
	info := MCPInfo{Transport: defaultMCPTransport}

	if cfg == nil || cfg.MCP == nil {
		return info
	}

	m := cfg.MCP
	if m.Transport != "" {
		info.Transport = m.Transport
	}
	info.Addr = m.Addr
	info.Path = m.Path
	info.AuthEnabled = m.Auth != nil && m.Auth.Token != ""

	return info
}

func (s *Service) buildAuthInfo(cfg *config.Config) AuthInfo {
	seen := make(map[string]struct{})

	if cfg != nil {
		for i := range cfg.Specs {
			spec := &cfg.Specs[i]
			if spec.Disable || spec.Auth.Client == nil {
				continue
			}
			t := string(spec.Auth.Client.Type())
			if t == string(auth.NoAuth) {
				continue
			}
			if _, ok := seen[t]; !ok {
				seen[t] = struct{}{}
			}
		}
	}

	if len(seen) == 0 {
		return AuthInfo{}
	}

	methods := make([]string, 0, len(seen))
	for m := range seen {
		methods = append(methods, m)
	}
	return AuthInfo{Methods: methods}
}

const (
	defaultMCPTransport = "stdio"
)
