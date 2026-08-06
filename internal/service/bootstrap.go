/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

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
		s.log.ErrorContext(ctx, "bootstrap failed: load config", "config", request.ConfFilePath, "error", err)
		return fmt.Errorf("failed to load config: %w", err)
	}
	init.storeConfig(cfg)

	if err := init.initWorkspace(filepath.Dir(request.ConfFilePath)); err != nil {
		s.log.ErrorContext(ctx, "bootstrap failed: init workspace", "config", request.ConfFilePath, "error", err)
		return err
	}

	httpCfg := BuildGlobalHTTPConfig(cfg.HTTPClient)
	if err := init.setupHTTPClient(httpCfg); err != nil {
		s.log.ErrorContext(ctx, "bootstrap failed: setup http client", "error", err)
		return err
	}
	s.cache.SetHTTPClient(s.ctx.loadHTTPClient())

	if err := init.processSpecs(ctx, cfg, request.Tags); err != nil {
		s.log.ErrorContext(ctx, "bootstrap failed: process specs", "error", err)
		return err
	}

	init.buildSnapshot()

	return nil
}
