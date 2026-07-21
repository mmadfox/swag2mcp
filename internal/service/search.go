package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/model"
)

type searchService struct {
	index IndexReader
	v     RequestValidator
}

func newSearchService(index IndexReader, v RequestValidator) *searchService {
	return &searchService{index: index, v: v}
}

// Search performs a full-text search across all endpoints using the given
// query string and returns up to the specified limit of matching results.
func (ss *searchService) Search(ctx context.Context, rq SearchRequest) (SearchResponse, error) {
	if err := ss.v.Struct(rq); err != nil {
		return SearchResponse{}, NewValidationError(
			"A search query is required and the limit must be between 1 and 50.",
			err,
		)
	}

	eps, err := ss.index.Search(ctx, strings.ToLower(rq.Query), rq.Limit)
	if err != nil {
		return SearchResponse{}, NewNotFoundError(
			"The search query did not match any endpoints. Try a different query.",
			err,
		)
	}

	is, err := mapEndpointsToSearchItems(ss.index, eps)
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

func mapEndpointsToSearchItems(
	index IndexReader,
	eps []*model.Endpoint,
) ([]EndpointSearchItem, error) {
	items := make([]EndpointSearchItem, 0, len(eps))
	for _, e := range eps {
		sp, err := index.SpecByID(e.SpecID)
		if err != nil {
			return nil, NewNotFoundError(
				fmt.Sprintf(
					"Spec %q was not found. The endpoint references a spec that no longer exists.",
					e.SpecID,
				),
				err,
			)
		}
		coll, err := index.CollectionByID(e.CollectionID)
		if err != nil {
			return nil, NewNotFoundError(
				fmt.Sprintf(
					"Collection %q was not found. The endpoint references a collection that no longer exists.",
					e.CollectionID,
				),
				err,
			)
		}
		tag, err := index.TagByID(e.TagID)
		if err != nil {
			return nil, NewNotFoundError(
				fmt.Sprintf(
					"Tag %q was not found. The endpoint references a tag that no longer exists.",
					e.TagID,
				),
				err,
			)
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
