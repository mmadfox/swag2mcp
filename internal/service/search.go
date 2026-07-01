package service

import (
	"context"
	"fmt"
	"sort"
)

type (
	// SearchRequest represents a request to search endpoints.
	SearchRequest struct {
		Query string `json:"query" jsonschema:"required," validate:"required"`
		Limit int    `json:"limit" jsonschema:"required,Maximum number of results to return" validate:"required,min=1,max=50"`
	}

	// SearchResponse represents a response to search endpoints.
	SearchResponse struct {
		Endpoints []EndpointSearchItem `json:"endpoints" jsonschema:"required,List of endpoints matching the search query"`
	}
)

// Search returns endpoints matching the query.
func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return SearchResponse{}, NewValidationError("query is required and limit must be between 1 and 50", err)
	}

	endpoints, err := s.index.Search(ctx, req.Query, req.Limit)
	if err != nil {
		return SearchResponse{}, NewNotFoundError("search failed", err)
	}

	resp := SearchResponse{
		Endpoints: make([]EndpointSearchItem, 0, len(endpoints)),
	}
	for _, ep := range endpoints {
		spec, err := s.index.SpecByID(ep.SpecID)
		if err != nil {
			return SearchResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", ep.SpecID), err)
		}
		collection, err := s.index.CollectionByID(ep.CollectionID)
		if err != nil {
			return SearchResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", ep.CollectionID), err)
		}
		tag, err := s.index.TagByID(ep.TagID)
		if err != nil {
			return SearchResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", ep.TagID), err)
		}
		resp.Endpoints = append(resp.Endpoints, EndpointSearchItem{
			ID:              ep.ID,
			TagID:           ep.TagID,
			TagName:         tag.Name,
			CollectionID:    ep.CollectionID,
			CollectionTitle: collection.Title,
			SpecID:          ep.SpecID,
			SpecDomain:      spec.Domain,
			Method:          ep.Name,
			Path:            ep.Path,
			Summary:         ep.SummaryOrFallback(),
		})
	}

	sort.Slice(resp.Endpoints, func(i, j int) bool {
		a, b := resp.Endpoints[i], resp.Endpoints[j]
		if a.SpecID != b.SpecID {
			return a.SpecID < b.SpecID
		}
		if a.CollectionID != b.CollectionID {
			return a.CollectionID < b.CollectionID
		}
		return a.TagID < b.TagID
	})

	return resp, nil
}
