package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"sort"
)

type specService struct {
	index IndexReader
	v     RequestValidator
}

func newSpecService(index IndexReader, v RequestValidator) *specService {
	return &specService{index: index, v: v}
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
	_ context.Context,
	rq SpecByIDRequest,
) (SpecByIDResponse, error) {
	if err := ss.v.Struct(rq); err != nil {
		return SpecByIDResponse{}, NewValidationError(
			"The spec ID is invalid. It must be a 32-character hex string. "+
				"Use spec_list to find the correct spec ID.",
			err,
		)
	}

	var r SpecByIDResponse
	sp, err := ss.index.SpecByID(rq.ID)
	if err != nil {
		return SpecByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. Use spec_list to see all available specs.", rq.ID),
			err,
		)
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
