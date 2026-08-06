/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"log/slog"
	"sort"
)

type specService struct {
	index IndexReader
	v     RequestValidator
	log   *slog.Logger
}

func newSpecService(index IndexReader, v RequestValidator, log *slog.Logger) *specService {
	return &specService{index: index, v: v, log: log}
}

// Specs returns a list of all available specifications.
func (ss *specService) Specs(_ context.Context) (SpecsResponse, error) {
	allSpecs := ss.index.AllSpecs()
	r := SpecsResponse{
		Specs: make([]SpecItem, len(allSpecs)),
	}

	for i, sp := range allSpecs {
		r.Specs[i] = SpecItem{
			ID:     sp.ID,
			Domain: sp.Domain,
		}
	}

	sort.Slice(r.Specs, func(i, j int) bool {
		return r.Specs[i].ID < r.Specs[j].ID
	})

	return r, nil
}

// SpecByID returns the specification identified by the given spec ID,
// along with its associated collections.
func (ss *specService) SpecByID(
	ctx context.Context,
	rq SpecByIDRequest,
) (SpecByIDResponse, error) {
	if err := ss.v.Struct(rq); err != nil {
		return SpecByIDResponse{}, NewInvalidSpecIDError(err)
	}

	var r SpecByIDResponse
	sp, err := ss.index.SpecByID(rq.ID)
	if err != nil {
		ss.log.ErrorContext(ctx, "spec_by_id failed: spec not found", "spec_id", rq.ID, "error", err)
		return SpecByIDResponse{}, NewSpecNotFoundError(rq.ID, err)
	}
	r.Spec = Spec{
		ID:     sp.ID,
		Domain: sp.Domain,
	}

	colls, err := ss.index.CollectionsBySpec(rq.ID)
	if err == nil {
		r.Collections = make([]CollectionItem, 0, len(colls))
		for _, c := range colls {
			r.Collections = append(r.Collections, CollectionItem{
				ID:           c.ID,
				Title:        c.Title,
				CountTags:    c.Stats.Tags,
				CountMethods: c.Stats.Methods,
			})
		}
		sort.Slice(r.Collections, func(i, j int) bool {
			return r.Collections[i].ID < r.Collections[j].ID
		})
	}

	return r, nil
}
