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
func (s *Service) TagsByCollection(_ context.Context, req TagsByCollectionRequest) (TagsByCollectionResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return TagsByCollectionResponse{}, NewValidationError(
			"collectionId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	collection, err := s.index.CollectionByID(req.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("collection %q not found", req.CollectionID),
			err,
		)
	}

	spec, err := s.index.SpecByID(collection.SpecID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", collection.SpecID), err)
	}

	tags, err := s.index.TagsByCollection(req.CollectionID)
	if err != nil {
		return TagsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("collection %q not found", req.CollectionID),
			err,
		)
	}

	resp := TagsByCollectionResponse{
		Spec: Spec{
			ID:     spec.ID,
			Domain: spec.Domain,
		},
		Collection: Collection{
			ID:           collection.ID,
			Title:        collection.Title,
			CountMethods: collection.Stats.Methods,
		},
		Tags: make([]TagListItem, 0, len(tags)),
	}
	for _, t := range tags {
		resp.Tags = append(resp.Tags, TagListItem{
			ID:           t.ID,
			Title:        t.Name,
			CountMethods: t.Stats.Methods,
		})
	}

	sort.Slice(resp.Tags, func(i, j int) bool {
		return resp.Tags[i].ID < resp.Tags[j].ID
	})

	return resp, nil
}

// TagByID returns a tag by its ID.
func (s *Service) TagByID(_ context.Context, req TagByIDRequest) (TagByIDResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return TagByIDResponse{}, NewValidationError(
			"tagId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	tag, err := s.index.TagByID(req.ID)
	if err != nil {
		return TagByIDResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", req.ID), err)
	}

	resp := TagByIDResponse{
		Tag: TagListItem{
			ID:           tag.ID,
			Title:        tag.Name,
			CountMethods: tag.Stats.Methods,
		},
	}

	return resp, nil
}

// TagsBySpec returns a list of all available tags for a given spec.
func (s *Service) TagsBySpec(_ context.Context, req TagsBySpecRequest) (TagsBySpecResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return TagsBySpecResponse{}, NewValidationError(
			"specId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	tags, err := s.index.TagsBySpec(req.SpecID)
	if err != nil {
		return TagsBySpecResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", req.SpecID), err)
	}

	resp := TagsBySpecResponse{
		Tags: make([]TagListItem, 0, len(tags)),
	}
	for _, t := range tags {
		resp.Tags = append(resp.Tags, TagListItem{
			ID:           t.ID,
			Title:        t.Name,
			CountMethods: t.Stats.Methods,
		})
	}

	sort.Slice(resp.Tags, func(i, j int) bool {
		return resp.Tags[i].ID < resp.Tags[j].ID
	})

	return resp, nil
}
