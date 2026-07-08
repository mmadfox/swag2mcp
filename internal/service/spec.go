package service

import (
	"context"
	"fmt"
	"sort"
)

type (
	// SpecByIDRequest is the request for SpecByID.
	SpecByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"required,A unique 32-character MD5 hash identifier for the spec,pattern=^[0-9a-f]{32}$"`
	}

	// SpecByIDResponse is the response for SpecByID.
	SpecByIDResponse struct {
		Spec        Spec             `json:"spec"        jsonschema:"required,Specification"`
		Collections []CollectionItem `json:"collections" jsonschema:"required,List of collections associated with the spec"`
	}

	// SpecsResponse is the response for Specs.
	SpecsResponse struct {
		Specs []SpecItem `json:"specs" jsonschema:"required,List of specifications"`
	}
)

// Specs returns a list of all available openapi/swagger specifications.
func (s *Service) Specs(_ context.Context) (SpecsResponse, error) {
	allSpecs := s.index.AllSpecs()
	resp := SpecsResponse{
		Specs: make([]SpecItem, len(allSpecs)),
	}

	for i, spec := range allSpecs {
		resp.Specs[i] = SpecItem{
			ID:     spec.ID,
			Domain: spec.Domain,
		}
	}

	sort.Slice(resp.Specs, func(i, j int) bool {
		return resp.Specs[i].ID < resp.Specs[j].ID
	})

	return resp, nil
}

// SpecByID returns a specification by its ID.
func (s *Service) SpecByID(_ context.Context, req SpecByIDRequest) (SpecByIDResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return SpecByIDResponse{}, NewValidationError(
			"The spec ID is invalid — it must be a 32-character hex string. Use spec_list to find available specs.",
			err,
		)
	}

	var resp SpecByIDResponse
	spec, err := s.index.SpecByID(req.ID)
	if err != nil {
		return SpecByIDResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — use spec_list to see all available specs.", req.ID), err)
	}
	resp.Spec = Spec{
		ID:     spec.ID,
		Domain: spec.Domain,
	}

	collections, err := s.index.CollectionsBySpec(req.ID)
	if err == nil {
		resp.Collections = make([]CollectionItem, 0, len(collections))
		for _, c := range collections {
			resp.Collections = append(resp.Collections, CollectionItem{
				ID:           c.ID,
				Title:        c.Title,
				CountTags:    c.Stats.Tags,
				CountMethods: c.Stats.Methods,
			})
		}
		sort.Slice(resp.Collections, func(i, j int) bool {
			return resp.Collections[i].ID < resp.Collections[j].ID
		})
	}

	return resp, nil
}
