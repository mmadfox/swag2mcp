/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"sort"
)

type tagService struct {
	index IndexReader
	v     RequestValidator
}

func newTagService(index IndexReader, v RequestValidator) *tagService {
	return &tagService{index: index, v: v}
}

// TagsByCollection returns a list of all available tags for a given collection.
func (ts *tagService) TagsByCollection(
	_ context.Context,
	rq TagsByCollectionRequest,
) (TagsByCollectionResponse, error) {
	if err := ts.v.Struct(rq); err != nil {
		return TagsByCollectionResponse{}, NewInvalidCollectionIDError(err)
	}

	coll, err := ts.index.CollectionByID(rq.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewCollectionNotFoundError(rq.CollectionID, err)
	}

	sp, err := ts.index.SpecByID(coll.SpecID)
	if err != nil {
		return TagsByCollectionResponse{}, NewSpecNotFoundError(coll.SpecID, err)
	}

	tgs, err := ts.index.TagsByCollection(rq.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewCollectionNotFoundError(rq.CollectionID, err)
	}

	r := TagsByCollectionResponse{
		Spec: Spec{
			ID:     sp.ID,
			Domain: sp.Domain,
		},
		Collection: Collection{
			ID:           coll.ID,
			Title:        coll.Title,
			CountMethods: coll.Stats.Methods,
		},
		Tags: make([]TagListItem, 0, len(tgs)),
	}
	for _, tg := range tgs {
		r.Tags = append(r.Tags, TagListItem{
			ID:           tg.ID,
			Title:        tg.Name,
			CountMethods: tg.Stats.Methods,
		})
	}

	sort.Slice(r.Tags, func(i, j int) bool {
		return r.Tags[i].ID < r.Tags[j].ID
	})

	return r, nil
}

// TagByID returns a tag by its ID.
func (ts *tagService) TagByID(
	_ context.Context,
	rq TagByIDRequest,
) (TagByIDResponse, error) {
	if err := ts.v.Struct(rq); err != nil {
		return TagByIDResponse{}, NewInvalidTagIDError(err)
	}

	tag, err := ts.index.TagByID(rq.ID)
	if err != nil {
		return TagByIDResponse{}, NewTagNotFoundError(rq.ID, err)
	}

	r := TagByIDResponse{
		Tag: TagListItem{
			ID:           tag.ID,
			Title:        tag.Name,
			CountMethods: tag.Stats.Methods,
		},
	}

	return r, nil
}

// TagsBySpec returns a list of all available tags for a given spec.
func (ts *tagService) TagsBySpec(
	_ context.Context,
	rq TagsBySpecRequest,
) (TagsBySpecResponse, error) {
	if err := ts.v.Struct(rq); err != nil {
		return TagsBySpecResponse{}, NewInvalidSpecIDError(err)
	}

	tgs, err := ts.index.TagsBySpec(rq.SpecID)
	if err != nil {
		return TagsBySpecResponse{}, NewSpecNotFoundError(rq.SpecID, err)
	}

	r := TagsBySpecResponse{
		Tags: make([]TagListItem, 0, len(tgs)),
	}
	for _, tg := range tgs {
		r.Tags = append(r.Tags, TagListItem{
			ID:           tg.ID,
			Title:        tg.Name,
			CountMethods: tg.Stats.Methods,
		})
	}

	sort.Slice(r.Tags, func(i, j int) bool {
		return r.Tags[i].ID < r.Tags[j].ID
	})

	return r, nil
}
