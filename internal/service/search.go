/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/model"
)

type searchService struct {
	index IndexReader
	v     RequestValidator
	log   *slog.Logger
}

func newSearchService(index IndexReader, v RequestValidator, log *slog.Logger) *searchService {
	return &searchService{index: index, v: v, log: log}
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
			ss.log.ErrorContext(ctx, "search failed: invalid query", "query", rq.Query, "error", err)
			return SearchResponse{}, NewSearchQueryError(err)
		}
		ss.log.ErrorContext(ctx, "search failed: no results", "query", rq.Query, "error", err)
		return SearchResponse{}, NewSearchNoResultsError()
	}

	is, err := mapEndpointsToSearchItems(ss.index, eps, ss.log)
	if err != nil {
		ss.log.ErrorContext(ctx, "search failed: map results", "query", rq.Query, "error", err)
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
	log *slog.Logger,
) ([]EndpointSearchItem, error) {
	items := make([]EndpointSearchItem, 0, len(eps))
	for _, e := range eps {
		sp, err := index.SpecByID(e.SpecID)
		if err != nil {
			log.Error("search failed: spec not found", "spec_id", e.SpecID, "error", err)
			return nil, NewSpecNotFoundError(e.SpecID, err)
		}
		coll, err := index.CollectionByID(e.CollectionID)
		if err != nil {
			log.Error("search failed: collection not found", "collection_id", e.CollectionID, "error", err)
			return nil, NewCollectionNotFoundError(e.CollectionID, err)
		}
		tag, err := index.TagByID(e.TagID)
		if err != nil {
			log.Error("search failed: tag not found", "tag_id", e.TagID, "error", err)
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
