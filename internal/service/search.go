/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/index"
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
		return SearchResponse{}, NewSearchQueryError(err)
	}

	eps, err := ss.index.Search(ctx, strings.ToLower(rq.Query), rq.Limit)
	if err != nil {
		if errors.Is(err, index.ErrInvalidQuery) {
			return SearchResponse{}, NewSearchQueryError(err)
		}
		return SearchResponse{}, NewSearchNoResultsError()
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
			return nil, NewSpecNotFoundError(e.SpecID, err)
		}
		coll, err := index.CollectionByID(e.CollectionID)
		if err != nil {
			return nil, NewCollectionNotFoundError(e.CollectionID, err)
		}
		tag, err := index.TagByID(e.TagID)
		if err != nil {
			return nil, NewTagNotFoundError(e.TagID, err)
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
