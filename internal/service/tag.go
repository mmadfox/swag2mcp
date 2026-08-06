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

type tagService struct {
	index IndexReader
	v     RequestValidator
	log   *slog.Logger
}

func newTagService(index IndexReader, v RequestValidator, log *slog.Logger) *tagService {
	return &tagService{index: index, v: v, log: log}
}

// TagsByCollection returns a list of all available tags for a given collection.
func (ts *tagService) TagsByCollection(
	ctx context.Context,
	rq TagsByCollectionRequest,
) (TagsByCollectionResponse, error) {
	if err := ts.v.Struct(rq); err != nil {
		return TagsByCollectionResponse{}, NewInvalidCollectionIDError(err)
	}

	coll, err := ts.index.CollectionByID(rq.CollectionID)
	if err != nil {
		ts.log.ErrorContext(ctx, "tags_by_collection failed: collection not found", "collection_id", rq.CollectionID, "error", err)
		return TagsByCollectionResponse{}, NewCollectionNotFoundError(rq.CollectionID, err)
	}

	sp, err := ts.index.SpecByID(coll.SpecID)
	if err != nil {
		ts.log.ErrorContext(ctx, "tags_by_collection failed: spec not found", "spec_id", coll.SpecID, "error", err)
		return TagsByCollectionResponse{}, NewSpecNotFoundError(coll.SpecID, err)
	}

	tgs, err := ts.index.TagsByCollection(rq.CollectionID)
	if err != nil {
		ts.log.ErrorContext(ctx, "tags_by_collection failed: tags not found", "collection_id", rq.CollectionID, "error", err)
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
	ctx context.Context,
	rq TagByIDRequest,
) (TagByIDResponse, error) {
	if err := ts.v.Struct(rq); err != nil {
		return TagByIDResponse{}, NewInvalidTagIDError(err)
	}

	tag, err := ts.index.TagByID(rq.ID)
	if err != nil {
		ts.log.ErrorContext(ctx, "tag_by_id failed: tag not found", "tag_id", rq.ID, "error", err)
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
	ctx context.Context,
	rq TagsBySpecRequest,
) (TagsBySpecResponse, error) {
	if err := ts.v.Struct(rq); err != nil {
		return TagsBySpecResponse{}, NewInvalidSpecIDError(err)
	}

	tgs, err := ts.index.TagsBySpec(rq.SpecID)
	if err != nil {
		ts.log.ErrorContext(ctx, "tags_by_spec failed: spec not found", "spec_id", rq.SpecID, "error", err)
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
