package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestInfoService_Info(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		snapshot  *InfoSnapshot
		version   string
		wantEmpty bool
	}{
		{
			name:      "from snapshot",
			snapshot:  &InfoSnapshot{Version: "1.0.0", Workspace: "/tmp/ws", Specs: SpecsSummary{Total: 2}},
			version:   "1.0.0",
			wantEmpty: false,
		},
		{
			name:      "no snapshot",
			snapshot:  nil,
			version:   "dev",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			snap := NewMockSnapshotStore(ctrl)
			ws := NewMockWorkspaceOps(ctrl)
			if tt.snapshot != nil {
				snap.EXPECT().Load().Return(tt.snapshot)
			} else {
				snap.EXPECT().Load().Return(nil)
				ws.EXPECT().Root().Return("/tmp/ws")
			}

			svc := newInfoService(
				NewMockSettingsProvider(ctrl),
				NewMockClock(ctrl),
				NewMockIndexReader(ctrl),
				ws,
				tt.version,
				snap,
				0,
			)
			resp, err := svc.Info(context.Background())
			require.NoError(t, err)
			if tt.wantEmpty {
				require.Empty(t, resp.Version)
			} else {
				require.NotEmpty(t, resp.Workspace)
			}
		})
	}
}

func TestInfoService_Info_uptime(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	started := now.Add(-5 * time.Minute).UnixNano()

	ctrl := gomock.NewController(t)
	snap := NewMockSnapshotStore(ctrl)
	snap.EXPECT().Load().Return(&InfoSnapshot{Version: "1.0", Workspace: "/ws"})

	clk := NewMockClock(ctrl)
	clk.EXPECT().Now().Return(now)

	svc := newInfoService(
		NewMockSettingsProvider(ctrl),
		clk,
		NewMockIndexReader(ctrl),
		NewMockWorkspaceOps(ctrl),
		"1.0",
		snap,
		started,
	)
	resp, err := svc.Info(context.Background())
	require.NoError(t, err)
	require.Equal(t, "5m0s", resp.Uptime)
}

func TestHumanizeBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input int
		want  string
	}{
		{500, "500 B"},
		{1024, "1 KB"},
		{2048, "2 KB"},
		{1048576, "1 MB"},
		{2097152, "2 MB"},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, humanizeBytes(tt.input))
	}
}

func TestBuildSpecsSummary(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return([]*model.Spec{
		{Stats: struct {
			Collections int `json:"collections"`
			Tags        int `json:"tags"`
			Methods     int `json:"methods"`
		}{Collections: 2, Methods: 10}},
	})

	svc := newInfoService(
		NewMockSettingsProvider(ctrl),
		NewMockClock(ctrl),
		idx,
		NewMockWorkspaceOps(ctrl),
		"1.0",
		NewMockSnapshotStore(ctrl),
		0,
	)

	cfg := &config.Config{
		Specs: []config.Spec{
			{Domain: "active-spec"},
			{Domain: "disabled-spec", Disable: true},
		},
	}
	sum := svc.buildSpecsSummary(cfg)
	require.Equal(t, 2, sum.Total)
	require.Equal(t, 1, sum.Active)
	require.Equal(t, 1, sum.Disabled)
	require.Equal(t, 2, sum.Collections)
	require.Equal(t, 10, sum.Endpoints)
}

func TestBuildSpecsSummary_nilConfig(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return(nil)

	svc := newInfoService(
		NewMockSettingsProvider(ctrl),
		NewMockClock(ctrl),
		idx,
		NewMockWorkspaceOps(ctrl),
		"1.0",
		NewMockSnapshotStore(ctrl),
		0,
	)

	sum := svc.buildSpecsSummary(nil)
	require.Zero(t, sum.Total)
	require.Zero(t, sum.Active)
	require.Zero(t, sum.Disabled)
}

func TestBuildHTTPClientInfo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	settings := NewMockSettingsProvider(ctrl)
	settings.EXPECT().MaxResponseSize().Return(2048)
	settings.EXPECT().HTTPClientConfig().Return(httpclient.Config{
		UserAgent: "test-agent",
	})

	svc := newInfoService(
		settings,
		NewMockClock(ctrl),
		NewMockIndexReader(ctrl),
		NewMockWorkspaceOps(ctrl),
		"1.0",
		NewMockSnapshotStore(ctrl),
		0,
	)

	cfg := &config.Config{
		HTTPClient: &config.GlobalHTTPClientConfig{
			UserAgent:       "global-agent",
			Timeout:         30 * time.Second,
			FollowRedirects: func() *bool { v := true; return &v }(),
			MaxRedirects:    func() *int { v := 5; return &v }(),
			MaxResponseSize: func() *int { v := 4096; return &v }(),
			Headers:         map[string]string{"X-Global": "val"},
			Cookies:         []config.Cookie{{Name: "s", Value: "v"}},
			Proxy:           &config.ProxyConfig{URL: "http://proxy:8080"},
		},
	}
	inf := svc.buildHTTPClientInfo(cfg)
	require.Equal(t, "4 KB", inf.MaxResponseSize)
	require.Equal(t, "test-agent", inf.UserAgent)
	require.NotEmpty(t, inf.Timeout)
	require.True(t, *inf.FollowRedirects)
	require.Equal(t, 5, *inf.MaxRedirects)
	require.NotEmpty(t, inf.Headers)
	require.NotEmpty(t, inf.Cookies)
	require.NotNil(t, inf.Proxy)
}

func TestBuildMCPInfo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := newInfoService(
		NewMockSettingsProvider(ctrl),
		NewMockClock(ctrl),
		NewMockIndexReader(ctrl),
		NewMockWorkspaceOps(ctrl),
		"1.0",
		NewMockSnapshotStore(ctrl),
		0,
	)

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()
		inf := svc.buildMCPInfo(nil)
		require.Equal(t, "stdio", inf.Transport)
	})

	t.Run("with config", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			MCP: &config.MCPConfig{
				Transport: "sse",
				Addr:      ":8080",
				Path:      "/mcp",
				Auth:      &config.MCPAuthConfig{Token: "secret"},
			},
		}
		inf := svc.buildMCPInfo(cfg)
		require.Equal(t, "sse", inf.Transport)
		require.Equal(t, ":8080", inf.Addr)
		require.Equal(t, "/mcp", inf.Path)
		require.True(t, inf.AuthEnabled)
	})
}

func TestBuildAuthInfo(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	svc := newInfoService(
		NewMockSettingsProvider(ctrl),
		NewMockClock(ctrl),
		NewMockIndexReader(ctrl),
		NewMockWorkspaceOps(ctrl),
		"1.0",
		NewMockSnapshotStore(ctrl),
		0,
	)

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()
		inf := svc.buildAuthInfo(nil)
		require.Empty(t, inf.Methods)
	})

	t.Run("with auth methods", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Specs: []config.Spec{
				{
					Domain: "api1",
					Auth: config.Auth{
						Client: auth.NewNoAuthClient(),
					},
				},
				{
					Domain: "api2",
					Auth: config.Auth{
						Client: &auth.BasicAuthClient{Username: "u", Password: "p"},
					},
				},
			},
		}
		inf := svc.buildAuthInfo(cfg)
		require.Contains(t, inf.Methods, "basic")
	})
}

func TestBuildSnapshot(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	settings := NewMockSettingsProvider(ctrl)
	settings.EXPECT().Config().Return(&config.Config{
		DisableRateLimiter: true,
		Specs:              []config.Spec{{Domain: "api"}},
		HTTPClient: &config.GlobalHTTPClientConfig{
			UserAgent: "global-agent",
		},
	})
	settings.EXPECT().MaxResponseSize().Return(2048)
	settings.EXPECT().HTTPClientConfig().Return(httpclient.Config{UserAgent: "agent"})

	idx := NewMockIndexReader(ctrl)
	idx.EXPECT().AllSpecs().Return([]*model.Spec{
		{Stats: struct {
			Collections int `json:"collections"`
			Tags        int `json:"tags"`
			Methods     int `json:"methods"`
		}{Collections: 1, Methods: 5}},
	})

	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().Root().Return("/ws")

	snap := NewMockSnapshotStore(ctrl)
	snap.EXPECT().Store(gomock.Any())

	svc := newInfoService(
		settings,
		NewMockClock(ctrl),
		idx,
		ws,
		"1.0",
		snap,
		0,
	)

	svc.buildSnapshot()
}
