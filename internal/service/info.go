package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
)

const (
	defaultMCPTransport = "stdio"
	devVersion          = "dev"
	mbInBytes           = 1048576
	kbInBytes           = 1024
)

type infoService struct {
	settings  SettingsProvider
	Clock     Clock
	index     IndexReader
	ws        WorkspaceOps
	version   string
	snapshot  SnapshotStore
	startedAt int64 // UnixNano, set once during bootstrap
}

func newInfoService(
	settings SettingsProvider,
	clk Clock,
	index IndexReader,
	ws WorkspaceOps,
	version string,
	snapshot SnapshotStore,
	startedAt int64,
) *infoService {
	return &infoService{
		settings:  settings,
		Clock:     clk,
		index:     index,
		ws:        ws,
		version:   version,
		snapshot:  snapshot,
		startedAt: startedAt,
	}
}

// Info returns a comprehensive summary of the current service state from the snapshot.
func (is *infoService) Info(_ context.Context) (InfoResponse, error) {
	snap, ok := is.snapshot.Load().(*InfoSnapshot)
	if ok && snap != nil {
		started := is.startedAt
		uptime := time.Duration(0)
		if started > 0 {
			uptime = is.Clock.Now().Sub(time.Unix(0, started)).Round(time.Second)
		}
		uptimeStr := ""
		if uptime > 0 {
			uptimeStr = uptime.String()
		}
		version := snap.Version
		if version == devVersion {
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

	version := is.version
	if version == devVersion {
		version = ""
	}
	return InfoResponse{
		Version:   version,
		Workspace: is.ws.Root(),
	}, nil
}

// buildSnapshot computes a point-in-time snapshot of the service state.
func (is *infoService) buildSnapshot() {
	snap := &InfoSnapshot{
		Version:   is.version,
		Workspace: is.ws.Root(),
	}

	cfg := is.settings.Config()
	if cfg != nil {
		snap.Specs = is.buildSpecsSummary(cfg)
		snap.HTTPClient = is.buildHTTPClientInfo(cfg)
		snap.MCP = is.buildMCPInfo(cfg)
		snap.Auth = is.buildAuthInfo(cfg)
		snap.Mock = MockInfo{Enabled: cfg.MockEnabled}
	}

	is.snapshot.Store(snap)
}

func (is *infoService) buildSpecsSummary(c *config.Config) SpecsSummary {
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

	specs := is.index.AllSpecs()
	for _, sp := range specs {
		sum.Collections += sp.Stats.Collections
		sum.Endpoints += sp.Stats.Methods
	}

	return sum
}

func (is *infoService) buildHTTPClientInfo(c *config.Config) HTTPClientInfo {
	inf := HTTPClientInfo{
		MaxResponseSize: humanizeBytes(is.settings.MaxResponseSize()),
	}

	if c == nil || c.HTTPClient == nil {
		return inf
	}

	httpCfg := is.settings.HTTPClientConfig()
	gc := c.HTTPClient
	inf.Randomize = httpCfg.Randomize
	inf.UserAgent = httpCfg.UserAgent
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

func (is *infoService) buildMCPInfo(c *config.Config) MCPInfo {
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

func (is *infoService) buildAuthInfo(c *config.Config) AuthInfo {
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

// humanizeBytes converts a byte count to a human-readable string (e.g. 2048 becomes "2 KB").
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
