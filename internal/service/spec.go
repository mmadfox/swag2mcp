package service

import (
	"context"
	"fmt"
	"sort"
)

// SpecByIDRequest contains the spec ID used to look up a specific
// specification and its collections.
type (
	SpecByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
	}

	// SpecByIDResponse contains the requested spec and its associated collections.
	SpecByIDResponse struct {
		Spec        Spec             `json:"spec"        jsonschema:"required,Specification"`
		Collections []CollectionItem `json:"collections" jsonschema:"required,List of collections associated with the spec"`
	}

	// SpecsResponse contains the list of all available specifications.
	SpecsResponse struct {
		Specs []SpecItem `json:"specs" jsonschema:"required,List of specifications"`
	}
)

// Specs returns a list of all available specifications.
func (s *Service) Specs(_ context.Context) (SpecsResponse, error) {
	allSpecs := s.index.AllSpecs()
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
func (s *Service) SpecByID(_ context.Context, rq SpecByIDRequest) (SpecByIDResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return SpecByIDResponse{}, NewValidationError(
			"The spec ID is invalid — it must be a 32-character hex string. Use spec_list to find available specs.",
			err,
		)
	}

	var r SpecByIDResponse
	sp, err := s.index.SpecByID(rq.ID)
	if err != nil {
		return SpecByIDResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — use spec_list to see all available specs.", rq.ID), err)
	}
	r.Spec = Spec{
		ID:     sp.ID,
		Domain: sp.Domain,
	}

	colls, err := s.index.CollectionsBySpec(rq.ID)
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
