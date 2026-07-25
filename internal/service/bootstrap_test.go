package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/stretchr/testify/require"
)

func TestBootstrap_ConfigError(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: "/nonexistent/config.yaml",
	})
	require.Error(t, err)
}

func TestInitializer_LoadConfig(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	cfg, err := init.loadConfig("testdata/swag2mcp.yaml", nil)
	require.Error(t, err) // config is missing base_url, validation fails
	require.Nil(t, cfg)
}

func TestInitializer_LoadConfig_InvalidPath(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	_, err = init.loadConfig("/nonexistent/path.yaml", nil)
	require.Error(t, err)
}

func TestInitializer_StoreConfig(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	cfg := &config.Config{DisableRateLimiter: true}
	init.storeConfig(cfg)

	loaded := svc.ctx.loadConfig()
	require.NotNil(t, loaded)
	require.True(t, loaded.DisableRateLimiter)
	require.True(t, svc.ctx.disableRateLimiter.Load())
}

func TestInitializer_SetupHTTPClient(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	httpCfg := httpclient.Config{
		UserAgent: "test-agent",
		Timeout:   30,
	}
	err = init.setupHTTPClient(httpCfg)
	require.NoError(t, err)

	client := svc.ctx.loadHTTPClient()
	require.NotNil(t, client)
	require.Equal(t, httpCfg.Timeout, svc.ctx.loadHTTPClientConfig().Timeout)
}

func TestInitializer_SetStartedAt(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	init := newInitializer(svc)
	init.setStartedAt()

	started := svc.ctx.startedAt.Load()
	require.Positive(t, started)
}
