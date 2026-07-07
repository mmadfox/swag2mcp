package service

import (
	"context"
	"fmt"
	"sort"
)

type (
	// CollectionsRequest represents a request to list all collections for a given spec.
	CollectionsRequest struct {
		SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
	}

	// CollectionsResponse represents a response to list all collections for a given spec.
	CollectionsResponse struct {
		Spec        Spec             `json:"spec"        jsonschema:"required,Specification"`
		Collections []CollectionItem `json:"collections" jsonschema:"List of collections associated with the spec,required"`
	}

	// CollectionByIDRequest represents a request to get a collection by its ID.
	CollectionByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"Unique identifier for the collection,required"`
	}

	// CollectionByIDResponse represents a response to get a collection by its ID.
	CollectionByIDResponse struct {
		Spec       Spec          `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection    `json:"collection" jsonschema:"required,Collection"`
		Tags       []TagListItem `json:"tags"       jsonschema:"List of tags associated with the collection,required"`
	}
)

// CollectionsBySpec returns a list of all available collections for a given spec.
func (s *Service) CollectionsBySpec(_ context.Context, req CollectionsRequest) (CollectionsResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return CollectionsResponse{}, NewValidationError(
			"specId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	spec, err := s.index.SpecByID(req.SpecID)
	if err != nil {
		return CollectionsResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", req.SpecID), err)
	}

	collections, err := s.index.CollectionsBySpec(req.SpecID)
	if err != nil {
		return CollectionsResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", req.SpecID), err)
	}

	resp := CollectionsResponse{
		Spec: Spec{
			ID:     spec.ID,
			Domain: spec.Domain,
		},
		Collections: make([]CollectionItem, 0, len(collections)),
	}
	for _, c := range collections {
		resp.Collections = append(resp.Collections, CollectionItem{
			ID:           c.ID,
			Title:        c.Title,
			LLMTitle:     c.LLMTitle,
			CountTags:    c.Stats.Tags,
			CountMethods: c.Stats.Methods,
		})
	}

	sort.Slice(resp.Collections, func(i, j int) bool {
		return resp.Collections[i].ID < resp.Collections[j].ID
	})

	return resp, nil
}

// CollectionByID returns a collection by its ID.
func (s *Service) CollectionByID(_ context.Context, req CollectionByIDRequest) (CollectionByIDResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return CollectionByIDResponse{}, NewValidationError(
			"specId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	var resp CollectionByIDResponse
	collection, err := s.index.CollectionByID(req.ID)
	if err != nil {
		return CollectionByIDResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", req.ID), err)
	}
	resp.Collection = Collection{
		ID:           collection.ID,
		Title:        collection.Title,
		CountMethods: collection.Stats.Methods,
	}

	spec, err := s.index.SpecByID(collection.SpecID)
	if err != nil {
		return CollectionByIDResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", collection.SpecID), err)
	}
	resp.Spec = Spec{
		ID:     spec.ID,
		Domain: spec.Domain,
	}

	tags, err := s.index.TagsByCollection(req.ID)
	if err == nil {
		resp.Tags = make([]TagListItem, 0, len(tags))
		for _, t := range tags {
			resp.Tags = append(resp.Tags, TagListItem{
				ID:           t.ID,
				Title:        t.Name,
				CountMethods: t.Stats.Methods,
			})
		}
	}

	return resp, nil
}
