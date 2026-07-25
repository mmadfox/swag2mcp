package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"

	"github.com/mmadfox/swag2mcp/internal/model"
)

func (s *Service) indexSpec(
	sp *model.Spec,
	colls map[string]*model.Collection,
	tags map[string]*model.Tag,
	eps map[string]*model.Endpoint,
) error {
	collections := make([]*model.Collection, 0, len(colls))
	for _, c := range colls {
		collections = append(collections, c)
	}

	ts := make([]*model.Tag, 0, len(tags))
	for _, t := range tags {
		ts = append(ts, t)
	}

	endpoints := make([]*model.Endpoint, 0, len(eps))
	for _, e := range eps {
		endpoints = append(endpoints, e)
	}

	if err := s.index.EnsureIndex(sp, collections, ts, endpoints); err != nil {
		return fmt.Errorf("failed to ensure index: %w", err)
	}

	return nil
}
