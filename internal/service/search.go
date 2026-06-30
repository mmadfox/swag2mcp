package service

import (
	"context"
)

type (
	// SearchRequest represents a request to search endpoints.
	SearchRequest struct {
		Query string `json:"query" jsonschema:"required," validate:"required"`
		Limit int    `json:"limit" jsonschema:"required,Maximum number of results to return" validate:"required,min=1,max=50"`
	}

	// SearchResponse represents a response to search endpoints.
	SearchResponse struct {
		Endpoints []EndpointItem `json:"endpoints" jsonschema:"required,List of endpoints matching the search query"`
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
		Endpoints: make([]EndpointItem, 0, len(endpoints)),
	}
	for _, ep := range endpoints {
		resp.Endpoints = append(resp.Endpoints, EndpointItem{
			ID:           ep.ID,
			TagID:        ep.TagID,
			CollectionID: ep.CollectionID,
			SpecID:       ep.SpecID,
			Method:       ep.Name,
			Path:         ep.Path,
			Summary:      ep.SummaryOrFallback(),
			Deprecated:   ep.Operation != nil && ep.Operation.Deprecated,
		})
	}

	return resp, nil
}
