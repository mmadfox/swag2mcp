package service

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
)

// InfoResponse holds the complete runtime and configuration summary.
type InfoResponse struct {
	Version       string         `json:"version"`
	LatestVersion string         `json:"latest_version,omitempty"`
	Workspace     string         `json:"workspace"`
	Uptime        string         `json:"uptime"`
	Specs         SpecsSummary   `json:"specs"`
	HTTPClient    HTTPClientInfo `json:"http_client"`
	MCP           MCPInfo        `json:"mcp"`
	Auth          AuthInfo       `json:"auth"`
	Mock          MockInfo       `json:"mock"`
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

// Info returns a comprehensive summary of the current service state.
// If cfg is nil, the method uses the config stored during Bootstrap.
func (s *Service) Info(ctx context.Context, cfg *config.Config) (InfoResponse, error) {
	resp := InfoResponse{
		Version:   s.version,
		Workspace: s.ws.Root(),
		Uptime:    time.Since(s.startedAt).Round(time.Second).String(),
	}

	if cfg == nil {
		cfg = s.config
	}

	if cfg != nil {
		resp.Specs = s.buildSpecsSummary(cfg)
		resp.HTTPClient = s.buildHTTPClientInfo(cfg)
		resp.MCP = s.buildMCPInfo(cfg)
		resp.Auth = s.buildAuthInfo(cfg)
		resp.Mock = MockInfo{Enabled: cfg.MockEnabled}
	}

	resp.LatestVersion = s.fetchLatestVersion(ctx)

	return resp, nil
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
	info.Randomize = g.Randomize
	info.UserAgent = g.UserAgent
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
	githubFetchTimeout  = 3
	defaultMCPTransport = "stdio"
)

func (s *Service) fetchLatestVersion(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, githubFetchTimeout*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.github.com/repos/mmadfox/swag2mcp/releases/latest", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/json")

	resp, doErr := http.DefaultClient.Do(req)
	if doErr != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&release); decodeErr != nil {
		return ""
	}

	return release.TagName
}

// InfoRequest is an empty request for the info tool.
type InfoRequest struct{}

// InfoResponseWrapper wraps InfoResponse for the MCP tool.
type InfoResponseWrapper struct {
	InfoResponse
}
