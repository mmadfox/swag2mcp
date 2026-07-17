package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/model"
)

// SearchRequest contains the search query and result limit for searching
// endpoints across all loaded specs.
type (
	SearchRequest struct {
		Query string `json:"query" jsonschema:"required,"                                    validate:"required"`
		Limit int    `json:"limit" jsonschema:"required,Maximum number of results to return" validate:"required,min=1,max=50"`
	}

	// SearchResponse contains the list of endpoints that matched the search query.
	SearchResponse struct {
		Endpoints []EndpointSearchItem `json:"endpoints" jsonschema:"required,List of endpoints matching the search query"`
	}
)

// Search performs a full-text search across all endpoints using the given
// query string and returns up to the specified limit of matching results.
func (s *Service) Search(ctx context.Context, rq SearchRequest) (SearchResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return SearchResponse{}, NewValidationError("A search query is required and the limit must be between 1 and 50.", err)
	}

	eps, err := s.index.Search(ctx, strings.ToLower(rq.Query), rq.Limit)
	if err != nil {
		return SearchResponse{}, NewNotFoundError("search failed", err)
	}

	is, err := s.mapEndpointsToSearchItems(eps)
	if err != nil {
		return SearchResponse{}, err
	}

	sort.Slice(is, func(i, j int) bool {
		a, b := is[i], is[j]
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

	return SearchResponse{Endpoints: is}, nil
}

func (s *Service) mapEndpointsToSearchItems(eps []*model.Endpoint) ([]EndpointSearchItem, error) {
	items := make([]EndpointSearchItem, 0, len(eps))
	for _, e := range eps {
		sp, err := s.index.SpecByID(e.SpecID)
		if err != nil {
			return nil, NewNotFoundError(fmt.Sprintf("spec %q not found", e.SpecID), err)
		}
		coll, err := s.index.CollectionByID(e.CollectionID)
		if err != nil {
			return nil, NewNotFoundError(fmt.Sprintf("collection %q not found", e.CollectionID), err)
		}
		tag, err := s.index.TagByID(e.TagID)
		if err != nil {
			return nil, NewNotFoundError(fmt.Sprintf("tag %q not found", e.TagID), err)
		}
		items = append(items, EndpointSearchItem{
			ID:              e.ID,
			TagID:           e.TagID,
			TagName:         tag.Name,
			CollectionID:    e.CollectionID,
			CollectionTitle: coll.Title,
			SpecID:          e.SpecID,
			SpecDomain:      sp.Domain,
			Method:          e.Name,
			Path:            e.Path,
			Summary:         e.SummaryOrFallback(),
		})
	}
	return items, nil
}
