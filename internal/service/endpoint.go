package service

import (
	"context"
	"fmt"
	"sort"
)

// EndpointsByTagRequest contains the tag ID used to look up endpoints
// associated with a specific tag.
type (
	EndpointsByTagRequest struct {
		TagID string `json:"tagId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsByTagResponse contains the spec, collection, tag, and the list
	// of endpoints associated with that tag.
	EndpointsByTagResponse struct {
		Spec       Spec              `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection        `json:"collection" jsonschema:"required,Collection"`
		Tag        TagListItem       `json:"tag"        jsonschema:"required,Tag"`
		Endpoints  []EndpointTagItem `json:"endpoints"  jsonschema:"required,List of endpoints associated with the tag"`
	}

	// EndpointsByCollectionRequest contains the collection ID used to look up
	// endpoints within a specific collection.
	EndpointsByCollectionRequest struct {
		CollectionID string `json:"collectionId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsByCollectionResponse contains the spec, collection, and the list
	// of endpoints associated with that collection.
	EndpointsByCollectionResponse struct {
		Spec       Spec                     `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection               `json:"collection" jsonschema:"required,Collection"`
		Endpoints  []EndpointCollectionItem `json:"endpoints"  jsonschema:"required,List of endpoints associated with the collection"`
	}

	// EndpointsBySpecRequest contains the spec ID used to look up all endpoints
	// belonging to a specific spec.
	EndpointsBySpecRequest struct {
		SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsBySpecResponse contains the list of endpoints associated with
	// the requested spec.
	EndpointsBySpecResponse struct {
		Endpoints []EndpointSearchItem `json:"endpoints" jsonschema:"required,List of endpoints associated with the spec"`
	}

	// EndpointByIDRequest contains the unique endpoint ID to look up a single
	// endpoint by its identifier.
	EndpointByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"required,Unique identifier for the endpoint"`
	}

	// EndpointByIDResponse contains the spec, collection, tag, and full endpoint
	// details for the requested endpoint ID.
	EndpointByIDResponse struct {
		Spec       Spec        `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection  `json:"collection" jsonschema:"required,Collection"`
		Tag        TagListItem `json:"tag"        jsonschema:"required,Tag"`
		Endpoint   Endpoint    `json:"endpoint"   jsonschema:"required,"`
	}
)

// EndpointsByTag returns all endpoints associated with the given tag,
// along with the parent spec, collection, and tag metadata.
func (s *Service) EndpointsByTag(_ context.Context, rq EndpointsByTagRequest) (EndpointsByTagResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return EndpointsByTagResponse{}, NewValidationError(
			"The tag ID is invalid — it must be a 32-character hex string.",
			err,
		)
	}

	tag, err := s.index.TagByID(rq.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("Tag %q not found — use tag_by_collection or tag_by_spec to list tags.", rq.TagID), err)
	}

	coll, err := s.index.CollectionByID(tag.CollectionID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("Collection %q not found — the tag references a collection that no longer exists.", tag.CollectionID), err)
	}

	sp, err := s.index.SpecByID(coll.SpecID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — the collection references a spec that no longer exists.", coll.SpecID), err)
	}

	eps, err := s.index.EndpointsByTag(rq.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("Tag %q not found — use tag_by_collection or tag_by_spec to list tags.", rq.TagID), err)
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
func (s *Service) EndpointsByCollection(
	_ context.Context,
	rq EndpointsByCollectionRequest,
) (EndpointsByCollectionResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return EndpointsByCollectionResponse{}, NewValidationError(
			"The collection ID is invalid — it must be a 32-character hex string.",
			err,
		)
	}

	coll, err := s.index.CollectionByID(rq.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q not found — use collection_by_spec to list collections.", rq.CollectionID),
			err,
		)
	}

	sp, err := s.index.SpecByID(coll.SpecID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q not found — the collection references a spec that no longer exists.", coll.SpecID),
			err,
		)
	}

	eps, err := s.index.EndpointByCollection(rq.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q not found — use collection_by_spec to list collections.", rq.CollectionID),
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
		tg, err := s.index.TagByID(e.TagID)
		if err != nil {
			return EndpointsByCollectionResponse{}, NewNotFoundError(fmt.Sprintf("Tag %q not found — the endpoint references a tag that no longer exists.", e.TagID), err)
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
func (s *Service) EndpointsBySpec(_ context.Context, rq EndpointsBySpecRequest) (EndpointsBySpecResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return EndpointsBySpecResponse{}, NewValidationError(
			"The spec ID is invalid — it must be a 32-character hex string. Use spec_list to find available specs.",
			err,
		)
	}

	eps, err := s.index.EndpointsBySpec(rq.SpecID)
	if err != nil {
		return EndpointsBySpecResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — use spec_list to see all available specs.", rq.SpecID), err)
	}

	is, err := s.mapEndpointsToSearchItems(eps)
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
func (s *Service) EndpointByID(_ context.Context, rq EndpointByIDRequest) (EndpointByIDResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return EndpointByIDResponse{}, NewValidationError(
			"The endpoint ID is invalid — it must be a 32-character hex string. Use the search tool to find the correct endpoint ID.",
			err,
		)
	}

	e, err := s.index.EndpointByID(rq.ID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("Endpoint %q not found — use the search tool to find the correct endpoint ID.", rq.ID), err)
	}

	sp, err := s.index.SpecByID(e.SpecID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — the endpoint references a spec that no longer exists.", e.SpecID), err)
	}
	coll, err := s.index.CollectionByID(e.CollectionID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("Collection %q not found — the endpoint references a collection that no longer exists.", e.CollectionID), err)
	}
	tag, err := s.index.TagByID(e.TagID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("Tag %q not found — the endpoint references a tag that no longer exists.", e.TagID), err)
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
