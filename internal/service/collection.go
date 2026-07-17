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
func (s *Service) CollectionsBySpec(_ context.Context, rq CollectionsRequest) (CollectionsResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return CollectionsResponse{}, NewValidationError(
			"The spec ID is invalid — it must be a 32-character hex string. Use spec_list to find available specs.",
			err,
		)
	}

	sp, err := s.index.SpecByID(rq.SpecID)
	if err != nil {
		return CollectionsResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — use spec_list to see all available specs.", rq.SpecID), err)
	}

	colls, err := s.index.CollectionsBySpec(rq.SpecID)
	if err != nil {
		return CollectionsResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — use spec_list to see all available specs.", rq.SpecID), err)
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

// CollectionByID returns a collection by its ID.
func (s *Service) CollectionByID(_ context.Context, rq CollectionByIDRequest) (CollectionByIDResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return CollectionByIDResponse{}, NewValidationError(
			"The collection ID is invalid — it must be a 32-character hex string.",
			err,
		)
	}

	var r CollectionByIDResponse
	coll, err := s.index.CollectionByID(rq.ID)
	if err != nil {
		return CollectionByIDResponse{}, NewNotFoundError(fmt.Sprintf("Collection %q not found — use collection_by_spec to list collections.", rq.ID), err)
	}
	r.Collection = Collection{
		ID:           coll.ID,
		Title:        coll.Title,
		CountMethods: coll.Stats.Methods,
	}

	sp, err := s.index.SpecByID(coll.SpecID)
	if err != nil {
		return CollectionByIDResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — the collection references a spec that no longer exists.", coll.SpecID), err)
	}
	r.Spec = Spec{
		ID:     sp.ID,
		Domain: sp.Domain,
	}

	ts, err := s.index.TagsByCollection(rq.ID)
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
