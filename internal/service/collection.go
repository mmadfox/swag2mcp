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

type collectionService struct {
	index IndexReader
	v     RequestValidator
}

func newCollectionService(index IndexReader, v RequestValidator) *collectionService {
	return &collectionService{index: index, v: v}
}

// CollectionsBySpec returns a list of all available collections for a given spec.
func (cs *collectionService) CollectionsBySpec(
	_ context.Context,
	rq CollectionsRequest,
) (CollectionsResponse, error) {
	if err := cs.v.Struct(rq); err != nil {
		return CollectionsResponse{}, NewValidationError(
			"The spec ID is invalid. It must be a 32-character hex string. "+
				"Use spec_list to find the correct spec ID.",
			err,
		)
	}

	sp, err := cs.index.SpecByID(rq.SpecID)
	if err != nil {
		return CollectionsResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. Use spec_list to see all available specs.", rq.SpecID),
			err,
		)
	}

	colls, err := cs.index.CollectionsBySpec(rq.SpecID)
	if err != nil {
		return CollectionsResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. Use spec_list to see all available specs.", rq.SpecID),
			err,
		)
	}

	r := CollectionsResponse{
		Spec: Spec{
			ID:     sp.ID,
			Domain: sp.Domain,
		},
		Collections: make([]CollectionItem, 0, len(colls)),
	}
	for _, c := range colls {
		r.Collections = append(r.Collections, CollectionItem{
			ID:           c.ID,
			Title:        c.Title,
			LLMTitle:     c.LLMTitle,
			CountTags:    c.Stats.Tags,
			CountMethods: c.Stats.Methods,
		})
	}

	sort.Slice(r.Collections, func(i, j int) bool {
		return r.Collections[i].ID < r.Collections[j].ID
	})

	return r, nil
}

// CollectionByID returns a collection by its ID, including its spec and tags.
func (cs *collectionService) CollectionByID(
	_ context.Context,
	rq CollectionByIDRequest,
) (CollectionByIDResponse, error) {
	if err := cs.v.Struct(rq); err != nil {
		return CollectionByIDResponse{}, NewValidationError(
			"The collection ID is invalid. It must be a 32-character hex string.",
			err,
		)
	}

	var r CollectionByIDResponse
	coll, err := cs.index.CollectionByID(rq.ID)
	if err != nil {
		return CollectionByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q was not found. Use collection_by_spec to list collections.", rq.ID),
			err,
		)
	}
	r.Collection = Collection{
		ID:           coll.ID,
		Title:        coll.Title,
		CountMethods: coll.Stats.Methods,
	}

	sp, err := cs.index.SpecByID(coll.SpecID)
	if err != nil {
		return CollectionByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. The collection references a spec that no longer exists.", coll.SpecID),
			err,
		)
	}
	r.Spec = Spec{
		ID:     sp.ID,
		Domain: sp.Domain,
	}

	ts, err := cs.index.TagsByCollection(rq.ID)
	if err == nil {
		r.Tags = make([]TagListItem, 0, len(ts))
		for _, tg := range ts {
			r.Tags = append(r.Tags, TagListItem{
				ID:           tg.ID,
				Title:        tg.Name,
				CountMethods: tg.Stats.Methods,
			})
		}
		sort.Slice(r.Tags, func(i, j int) bool {
			return r.Tags[i].ID < r.Tags[j].ID
		})
	}

	return r, nil
}
