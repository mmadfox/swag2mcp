/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"sort"
)

type endpointService struct {
	index IndexReader
	v     RequestValidator
}

func newEndpointService(index IndexReader, v RequestValidator) *endpointService {
	return &endpointService{index: index, v: v}
}

// EndpointsByTag returns all endpoints associated with the given tag,
// along with the parent spec, collection, and tag metadata.
func (es *endpointService) EndpointsByTag(
	_ context.Context,
	rq EndpointsByTagRequest,
) (EndpointsByTagResponse, error) {
	if err := es.v.Struct(rq); err != nil {
		return EndpointsByTagResponse{}, NewInvalidTagIDError(err)
	}

	tag, err := es.index.TagByID(rq.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewTagNotFoundError(rq.TagID, err)
	}

	coll, err := es.index.CollectionByID(tag.CollectionID)
	if err != nil {
		return EndpointsByTagResponse{}, NewCollectionNotFoundError(tag.CollectionID, err)
	}

	sp, err := es.index.SpecByID(coll.SpecID)
	if err != nil {
		return EndpointsByTagResponse{}, NewSpecNotFoundError(coll.SpecID, err)
	}

	eps, err := es.index.EndpointsByTag(rq.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewTagNotFoundError(rq.TagID, err)
	}

	r := EndpointsByTagResponse{
		Spec: Spec{
			ID:     sp.ID,
			Domain: sp.Domain,
		},
		Collection: Collection{
			ID:           coll.ID,
			Title:        coll.Title,
			CountMethods: coll.Stats.Methods,
		},
		Tag: TagListItem{
			ID:           tag.ID,
			Title:        tag.Name,
			CountMethods: tag.Stats.Methods,
		},
		Endpoints: make([]EndpointTagItem, 0, len(eps)),
	}
	for _, e := range eps {
		r.Endpoints = append(r.Endpoints, EndpointTagItem{
			ID:      e.ID,
			Method:  e.Name,
			Path:    e.Path,
			Summary: e.SummaryOrFallback(),
		})
	}

	sort.Slice(r.Endpoints, func(i, j int) bool {
		return r.Endpoints[i].ID < r.Endpoints[j].ID
	})

	return r, nil
}

// EndpointsByCollection returns all endpoints within the given collection,
// along with the parent spec and collection metadata.
func (es *endpointService) EndpointsByCollection(
	_ context.Context,
	rq EndpointsByCollectionRequest,
) (EndpointsByCollectionResponse, error) {
	if err := es.v.Struct(rq); err != nil {
		return EndpointsByCollectionResponse{}, NewInvalidCollectionIDError(err)
	}

	coll, err := es.index.CollectionByID(rq.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewCollectionNotFoundError(rq.CollectionID, err)
	}

	sp, err := es.index.SpecByID(coll.SpecID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewSpecNotFoundError(coll.SpecID, err)
	}

	eps, err := es.index.EndpointByCollection(rq.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewCollectionNotFoundError(rq.CollectionID, err)
	}

	r := EndpointsByCollectionResponse{
		Spec: Spec{
			ID:     sp.ID,
			Domain: sp.Domain,
		},
		Collection: Collection{
			ID:           coll.ID,
			Title:        coll.Title,
			CountMethods: coll.Stats.Methods,
		},
		Endpoints: make([]EndpointCollectionItem, 0, len(eps)),
	}
	for _, e := range eps {
		tg, err := es.index.TagByID(e.TagID)
		if err != nil {
			return EndpointsByCollectionResponse{}, NewTagNotFoundError(e.TagID, err)
		}
		r.Endpoints = append(r.Endpoints, EndpointCollectionItem{
			ID:      e.ID,
			TagID:   e.TagID,
			TagName: tg.Name,
			Method:  e.Name,
			Path:    e.Path,
			Summary: e.SummaryOrFallback(),
		})
	}

	sort.Slice(r.Endpoints, func(i, j int) bool {
		return r.Endpoints[i].ID < r.Endpoints[j].ID
	})

	return r, nil
}

// EndpointsBySpec returns all endpoints belonging to the given spec.
func (es *endpointService) EndpointsBySpec(
	_ context.Context,
	rq EndpointsBySpecRequest,
) (EndpointsBySpecResponse, error) {
	if err := es.v.Struct(rq); err != nil {
		return EndpointsBySpecResponse{}, NewInvalidSpecIDError(err)
	}

	eps, err := es.index.EndpointsBySpec(rq.SpecID)
	if err != nil {
		return EndpointsBySpecResponse{}, NewSpecNotFoundError(rq.SpecID, err)
	}

	is, err := mapEndpointsToSearchItems(es.index, eps)
	if err != nil {
		return EndpointsBySpecResponse{}, err
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

	return EndpointsBySpecResponse{Endpoints: is}, nil
}

// EndpointByID returns the full details for a single endpoint identified by
// its unique endpoint ID, including the parent spec, collection, and tag.
func (es *endpointService) EndpointByID(
	_ context.Context,
	rq EndpointByIDRequest,
) (EndpointByIDResponse, error) {
	if err := es.v.Struct(rq); err != nil {
		return EndpointByIDResponse{}, NewInvalidEndpointIDError(err)
	}

	e, err := es.index.EndpointByID(rq.ID)
	if err != nil {
		return EndpointByIDResponse{}, NewEndpointNotFoundError(rq.ID, err)
	}

	sp, err := es.index.SpecByID(e.SpecID)
	if err != nil {
		return EndpointByIDResponse{}, NewSpecNotFoundError(e.SpecID, err)
	}
	coll, err := es.index.CollectionByID(e.CollectionID)
	if err != nil {
		return EndpointByIDResponse{}, NewCollectionNotFoundError(e.CollectionID, err)
	}
	tag, err := es.index.TagByID(e.TagID)
	if err != nil {
		return EndpointByIDResponse{}, NewTagNotFoundError(e.TagID, err)
	}

	r := EndpointByIDResponse{
		Spec: Spec{
			ID:     sp.ID,
			Domain: sp.Domain,
		},
		Collection: Collection{
			ID:           coll.ID,
			Title:        coll.Title,
			CountMethods: coll.Stats.Methods,
		},
		Tag: TagListItem{
			ID:           tag.ID,
			Title:        tag.Name,
			CountMethods: tag.Stats.Methods,
		},
		Endpoint: Endpoint{
			ID:      e.ID,
			Method:  e.Name,
			Path:    e.Path,
			Summary: e.SummaryOrFallback(),
		},
	}

	return r, nil
}
