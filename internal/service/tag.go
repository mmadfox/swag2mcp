package service

import (
	"context"
	"fmt"
	"sort"
)

type (
	// TagsByCollectionRequest represents a request to list all tags for a given collection.
	TagsByCollectionRequest struct {
		CollectionID string `json:"collectionId" jsonschema:"required," validate:"required,md5"`
	}

	// TagsByCollectionResponse represents a response to list all tags for a given collection.
	TagsByCollectionResponse struct {
		Spec       Spec          `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection    `json:"collection" jsonschema:"required,Collection"`
		Tags       []TagListItem `json:"tags"       jsonschema:"required,List of tags associated with the collection"`
	}

	// TagByIDRequest represents a request to get a tag by its ID.
	TagByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"required,Unique identifier for the tag"`
	}

	// TagByIDResponse represents a response to get a tag by its ID.
	TagByIDResponse struct {
		Tag TagListItem `json:"tag" jsonschema:"required,"`
	}

	// TagsBySpecRequest represents a request to list all tags for a given spec.
	TagsBySpecRequest struct {
		SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
	}

	// TagsBySpecResponse represents a response to list all tags for a given spec.
	TagsBySpecResponse struct {
		Tags []TagListItem `json:"tags" jsonschema:"required,List of tags associated with the spec"`
	}
)

// TagsByCollection returns a list of all available tags for a given collection.
func (s *Service) TagsByCollection(_ context.Context, rq TagsByCollectionRequest) (TagsByCollectionResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return TagsByCollectionResponse{}, NewValidationError(
			"The collection ID is invalid — it must be a 32-character hex string.",
			err,
		)
	}

	coll, err := s.index.CollectionByID(rq.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q not found — use collection_by_spec to list collections.", rq.CollectionID),
			err,
		)
	}

	sp, err := s.index.SpecByID(coll.SpecID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — the collection references a spec that no longer exists.", coll.SpecID), err)
	}

	ts, err := s.index.TagsByCollection(rq.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q not found — use collection_by_spec to list collections.", rq.CollectionID),
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
		Tags: make([]TagListItem, 0, len(ts)),
	}
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

	return r, nil
}

// TagByID returns a tag by its ID.
func (s *Service) TagByID(_ context.Context, rq TagByIDRequest) (TagByIDResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return TagByIDResponse{}, NewValidationError(
			"The tag ID is invalid — it must be a 32-character hex string.",
			err,
		)
	}

	tag, err := s.index.TagByID(rq.ID)
	if err != nil {
		return TagByIDResponse{}, NewNotFoundError(fmt.Sprintf("Tag %q not found — use tag_by_collection or tag_by_spec to list tags.", rq.ID), err)
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
func (s *Service) TagsBySpec(_ context.Context, rq TagsBySpecRequest) (TagsBySpecResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return TagsBySpecResponse{}, NewValidationError(
			"The spec ID is invalid — it must be a 32-character hex string. Use spec_list to find available specs.",
			err,
		)
	}

	ts, err := s.index.TagsBySpec(rq.SpecID)
	if err != nil {
		return TagsBySpecResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — use spec_list to see all available specs.", rq.SpecID), err)
	}

	r := TagsBySpecResponse{
		Tags: make([]TagListItem, 0, len(ts)),
	}
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

	return r, nil
}
