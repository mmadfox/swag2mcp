package service

import (
	"context"
	"fmt"
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
	Version    string         `json:"version,omitempty"`
	Workspace  string         `json:"workspace"`
	Uptime     string         `json:"uptime,omitempty"`
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
	MaxResponseSize string            `json:"max_response_size"`
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
		uptime := time.Since(s.startedAt).Round(time.Second)
		uptimeStr := ""
		if uptime > 0 {
			uptimeStr = uptime.String()
		}
		version := snap.Version
		if version == "dev" {
			version = ""
		}
		return InfoResponse{
			Version:    version,
			Workspace:  snap.Workspace,
			Uptime:     uptimeStr,
			Specs:      snap.Specs,
			HTTPClient: snap.HTTPClient,
			MCP:        snap.MCP,
			Auth:       snap.Auth,
			Mock:       snap.Mock,
		}, nil
	}

	version := s.version
	if version == "dev" {
		version = ""
	}
	return InfoResponse{
		Version:   version,
		Workspace: s.ws.Root(),
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

func (s *Service) buildSpecsSummary(c *config.Config) SpecsSummary {
	sum := SpecsSummary{}

	if c != nil {
		for i := range c.Specs {
			sp := &c.Specs[i]
			sum.Total++
			if sp.Disable {
				sum.Disabled++
			} else {
				sum.Active++
			}
		}
	}

	specs := s.index.AllSpecs()
	for _, sp := range specs {
		sum.Collections += sp.Stats.Collections
		sum.Endpoints += sp.Stats.Methods
	}

	return sum
}

func (s *Service) buildHTTPClientInfo(c *config.Config) HTTPClientInfo {
	inf := HTTPClientInfo{
		MaxResponseSize: humanizeBytes(s.maxResponseSize),
	}

	if c == nil || c.HTTPClient == nil {
		return inf
	}

	gc := c.HTTPClient
	inf.Randomize = s.httpClientConfig.Randomize
	inf.UserAgent = s.httpClientConfig.UserAgent
	if inf.UserAgent == "" {
		inf.UserAgent = gc.UserAgent
	}
	if gc.Timeout > 0 {
		inf.Timeout = gc.Timeout.String()
	}
	inf.FollowRedirects = gc.FollowRedirects
	inf.MaxRedirects = gc.MaxRedirects
	if gc.MaxResponseSize != nil {
		inf.MaxResponseSize = humanizeBytes(*gc.MaxResponseSize)
	}

	if len(gc.Headers) > 0 {
		inf.Headers = make(map[string]string, len(gc.Headers))
		maps.Copy(inf.Headers, gc.Headers)
	}

	if len(gc.Cookies) > 0 {
		inf.Cookies = make([]CookieInfo, 0, len(gc.Cookies))
		for _, ck := range gc.Cookies {
			inf.Cookies = append(inf.Cookies, CookieInfo{
				Name:     ck.Name,
				Domain:   ck.Domain,
				Path:     ck.Path,
				Secure:   ck.Secure,
				HTTPOnly: ck.HTTPOnly,
			})
		}
	}

	if gc.Proxy != nil {
		inf.Proxy = &ProxyInfo{
			URL:      gc.Proxy.URL,
			Username: gc.Proxy.Username,
			Bypass:   append([]string{}, gc.Proxy.Bypass...),
		}
	}

	return inf
}

func (s *Service) buildMCPInfo(c *config.Config) MCPInfo {
	inf := MCPInfo{Transport: defaultMCPTransport}

	if c == nil || c.MCP == nil {
		return inf
	}

	mt := c.MCP
	if mt.Transport != "" {
		inf.Transport = mt.Transport
	}
	inf.Addr = mt.Addr
	inf.Path = mt.Path
	inf.AuthEnabled = mt.Auth != nil && mt.Auth.Token != ""

	return inf
}

func (s *Service) buildAuthInfo(c *config.Config) AuthInfo {
	m := make(map[string]struct{})

	if c != nil {
		for i := range c.Specs {
			sp := &c.Specs[i]
			if sp.Disable || sp.Auth.Client == nil {
				continue
			}
			tp := string(sp.Auth.Client.Type())
			if tp == string(auth.NoAuth) {
				continue
			}
			if _, ok := m[tp]; !ok {
				m[tp] = struct{}{}
			}
		}
	}

	if len(m) == 0 {
		return AuthInfo{}
	}

	ms := make([]string, 0, len(m))
	for k := range m {
		ms = append(ms, k)
	}
	return AuthInfo{Methods: ms}
}

const (
	defaultMCPTransport = "stdio"
	mbInBytes           = 1048576
	kbInBytes           = 1024
)

// humanizeBytes converts a byte count to a human-readable string (e.g. 2048 → "2 KB").
func humanizeBytes(b int) string {
	switch {
	case b >= mbInBytes:
		return fmt.Sprintf("%d MB", b/mbInBytes)
	case b >= kbInBytes:
		return fmt.Sprintf("%d KB", b/kbInBytes)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
