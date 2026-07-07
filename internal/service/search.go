package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/types"
)

type (
	// SearchRequest represents a request to search endpoints.
	SearchRequest struct {
		Query string `json:"query" jsonschema:"required,"                                    validate:"required"`
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

	endpoints, err := s.index.Search(ctx, strings.ToLower(req.Query), req.Limit)
	if err != nil {
		return SearchResponse{}, NewNotFoundError("search failed", err)
	}

	items, itemsErr := s.mapEndpointsToSearchItems(endpoints)
	if itemsErr != nil {
		return SearchResponse{}, itemsErr
	}

	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.SpecID != b.SpecID {
			return a.SpecID < b.SpecID
		}
		if a.CollectionID != b.CollectionID {
			return a.CollectionID < b.CollectionID
		}
		if a.TagID != b.TagID {
			return a.TagID < b.TagID
		}
		return a.ID < b.ID
	})

	return SearchResponse{Endpoints: items}, nil
}

// mapEndpointsToSearchItems maps raw endpoints to EndpointSearchItem with resolved spec/collection/tag.
func (s *Service) mapEndpointsToSearchItems(endpoints []*types.Endpoint) ([]EndpointSearchItem, error) {
	items := make([]EndpointSearchItem, 0, len(endpoints))
	for _, ep := range endpoints {
		epSpec, sErr := s.index.SpecByID(ep.SpecID)
		if sErr != nil {
			return nil, NewNotFoundError(fmt.Sprintf("spec %q not found", ep.SpecID), sErr)
		}
		epColl, cErr := s.index.CollectionByID(ep.CollectionID)
		if cErr != nil {
			return nil, NewNotFoundError(fmt.Sprintf("collection %q not found", ep.CollectionID), cErr)
		}
		epTag, tErr := s.index.TagByID(ep.TagID)
		if tErr != nil {
			return nil, NewNotFoundError(fmt.Sprintf("tag %q not found", ep.TagID), tErr)
		}
		items = append(items, EndpointSearchItem{
			ID:              ep.ID,
			TagID:           ep.TagID,
			TagName:         epTag.Name,
			CollectionID:    ep.CollectionID,
			CollectionTitle: epColl.Title,
			SpecID:          ep.SpecID,
			SpecDomain:      epSpec.Domain,
			Method:          ep.Name,
			Path:            ep.Path,
			Summary:         ep.SummaryOrFallback(),
		})
	}
	return items, nil
}
