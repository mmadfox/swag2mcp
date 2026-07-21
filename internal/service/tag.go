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
		return TagsByCollectionResponse{}, NewValidationError(
			"The collection ID is invalid. It must be a 32-character hex string.",
			err,
		)
	}

	coll, err := ts.index.CollectionByID(rq.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Collection %q was not found. Use collection_by_spec to list collections.",
				rq.CollectionID,
			),
			err,
		)
	}

	sp, err := ts.index.SpecByID(coll.SpecID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Spec %q was not found. The collection references a spec that no longer exists.",
				coll.SpecID,
			),
			err,
		)
	}

	tgs, err := ts.index.TagsByCollection(rq.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Collection %q was not found. Use collection_by_spec to list collections.",
				rq.CollectionID,
			),
			err,
		)
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
		return TagByIDResponse{}, NewValidationError(
			"The tag ID is invalid. It must be a 32-character hex string.",
			err,
		)
	}

	tag, err := ts.index.TagByID(rq.ID)
	if err != nil {
		return TagByIDResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Tag %q was not found. Use tag_by_collection or tag_by_spec to list tags.",
				rq.ID,
			),
			err,
		)
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
		return TagsBySpecResponse{}, NewValidationError(
			"The spec ID is invalid. It must be a 32-character hex string. "+
				"Use spec_list to find the correct spec ID.",
			err,
		)
	}

	tgs, err := ts.index.TagsBySpec(rq.SpecID)
	if err != nil {
		return TagsBySpecResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Spec %q was not found. Use spec_list to see all available specs.",
				rq.SpecID,
			),
			err,
		)
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
