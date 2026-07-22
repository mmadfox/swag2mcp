package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"time"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/ratelimit"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// initializer encapsulates the bootstrap mutations on Service fields.
// This keeps bootstrap.go focused on orchestration, not field assignment.
type initializer struct {
	s *Service
}

func newInitializer(s *Service) *initializer {
	return &initializer{s: s}
}

func (init *initializer) setStartedAt() {
	init.s.ctx.startedAt.Store(time.Now().UnixNano())
}

func (init *initializer) storeConfig(cfg *config.Config) {
	init.s.ctx.storeConfig(cfg)
	init.s.ctx.disableRateLimiter.Store(cfg.DisableRateLimiter)
	init.s.ctx.storeRateLimiter(ratelimit.NewWithInterval(cfg.RateLimitInterval))
}

func (init *initializer) initWorkspace(wsDir string) error {
	if wsDir != "" && wsDir != init.s.ws.Root() {
		ws, err := workspace.New(wsDir)
		if err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}
		init.s.ws = ws
	}

	if err := init.s.ws.Init(); err != nil {
		return fmt.Errorf("failed to init workspace: %w", err)
	}

	init.s.cache.SetWorkspaceDir(init.s.ws.Root())
	return nil
}

func (init *initializer) setupHTTPClient(httpCfg httpclient.Config) error {
	if httpCfg.Randomize {
		httpclient.RandomizeConfig(&httpCfg)
	}

	client, err := httpclient.New(httpCfg)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	init.s.ctx.storeHTTPClient(client)
	init.s.ctx.storeHTTPClientConfig(httpCfg)
	init.s.ctx.maxResponseSize.Store(int64(resolveMaxResponseSize(httpCfg.MaxResponseSize)))
	init.s.ctx.storeGlobalHeaders(httpCfg.Headers)
	init.s.ctx.storeGlobalUserAgent(httpCfg.UserAgent)
	init.s.ctx.storeGlobalCookies(httpCfg.Cookies)

	httpclient.SetGlobalConfig(httpCfg)
	return nil
}

func (init *initializer) loadConfig(configFilepath string, tags []string) (*config.Config, error) {
	cfg, err := config.Load(configFilepath)
	if err != nil {
		return nil, err
	}

	filter := config.NewFilter(tags)
	if err := cfg.Validate(filter); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (init *initializer) processSpecs(ctx context.Context, cfg *config.Config, tags []string) error {
	filter := config.NewFilter(tags)
	for sc := range cfg.Iterate(filter) {
		if err := init.s.processSpec(ctx, sc, cfg.MockEnabled, cfg.MockAuth); err != nil {
			return err
		}
	}
	return nil
}

// buildSnapshot computes a point-in-time snapshot of the service state after bootstrap.
func (init *initializer) buildSnapshot() {
	init.s.infoSvc.buildSnapshot()
}
