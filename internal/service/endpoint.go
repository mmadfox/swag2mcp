package service

import (
	"context"
	"fmt"
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
		return EndpointsByTagResponse{}, NewValidationError(
			"The tag ID is invalid. It must be a 32-character hex string. "+
				"Use tag_by_collection or tag_by_spec to find the correct tag ID.",
			err,
		)
	}

	tag, err := es.index.TagByID(rq.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(
			fmt.Sprintf("Tag %q was not found. Use tag_by_collection or tag_by_spec to list tags.", rq.TagID),
			err,
		)
	}

	coll, err := es.index.CollectionByID(tag.CollectionID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q was not found. The tag references a collection that no longer exists.", tag.CollectionID),
			err,
		)
	}

	sp, err := es.index.SpecByID(coll.SpecID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. The collection references a spec that no longer exists.", coll.SpecID),
			err,
		)
	}

	eps, err := es.index.EndpointsByTag(rq.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(
			fmt.Sprintf("Tag %q was not found. Use tag_by_collection or tag_by_spec to list tags.", rq.TagID),
			err,
		)
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
		return EndpointsByCollectionResponse{}, NewValidationError(
			"The collection ID is invalid. It must be a 32-character hex string.",
			err,
		)
	}

	coll, err := es.index.CollectionByID(rq.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q was not found. Use collection_by_spec to list collections.", rq.CollectionID),
			err,
		)
	}

	sp, err := es.index.SpecByID(coll.SpecID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. The collection references a spec that no longer exists.", coll.SpecID),
			err,
		)
	}

	eps, err := es.index.EndpointByCollection(rq.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q was not found. Use collection_by_spec to list collections.", rq.CollectionID),
			err,
		)
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
			return EndpointsByCollectionResponse{}, NewNotFoundError(
				fmt.Sprintf("Tag %q was not found. The endpoint references a tag that no longer exists.", e.TagID),
				err,
			)
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
		return EndpointsBySpecResponse{}, NewValidationError(
			"The spec ID is invalid. It must be a 32-character hex string. "+
				"Use spec_list to find the correct spec ID.",
			err,
		)
	}

	eps, err := es.index.EndpointsBySpec(rq.SpecID)
	if err != nil {
		return EndpointsBySpecResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. Use spec_list to see all available specs.", rq.SpecID),
			err,
		)
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
		return EndpointByIDResponse{}, NewValidationError(
			"The endpoint ID is invalid. It must be a 32-character hex string. "+
				"Use the search tool to find the correct endpoint ID.",
			err,
		)
	}

	e, err := es.index.EndpointByID(rq.ID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Endpoint %q was not found. Use the search tool to find the correct endpoint ID.", rq.ID),
			err,
		)
	}

	sp, err := es.index.SpecByID(e.SpecID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. The endpoint references a spec that no longer exists.", e.SpecID),
			err,
		)
	}
	coll, err := es.index.CollectionByID(e.CollectionID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q was not found. The endpoint references a collection that no longer exists.", e.CollectionID),
			err,
		)
	}
	tag, err := es.index.TagByID(e.TagID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(
			fmt.Sprintf("Tag %q was not found. The endpoint references a tag that no longer exists.", e.TagID),
			err,
		)
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
