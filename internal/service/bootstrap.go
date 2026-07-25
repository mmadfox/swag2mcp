package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"path/filepath"
)

// BootstrapRequest is the request for the Bootstrap method.
type BootstrapRequest struct {
	ConfFilePath string
	Tags         []string
}

// Bootstrap loads the configuration, initializes the workspace, creates the
// global HTTP client, and indexes all specs, collections, tags, and endpoints.
func (s *Service) Bootstrap(ctx context.Context, request BootstrapRequest) error {
	init := newInitializer(s)
	init.setStartedAt()

	cfg, err := init.loadConfig(request.ConfFilePath, request.Tags)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	init.storeConfig(cfg)

	if err := init.initWorkspace(filepath.Dir(request.ConfFilePath)); err != nil {
		return err
	}

	httpCfg := BuildGlobalHTTPConfig(cfg.HTTPClient)
	if err := init.setupHTTPClient(httpCfg); err != nil {
		return err
	}
	s.cache.SetHTTPClient(s.ctx.loadHTTPClient())

	if err := init.processSpecs(ctx, cfg, request.Tags); err != nil {
		return err
	}

	init.buildSnapshot()

	return nil
}
