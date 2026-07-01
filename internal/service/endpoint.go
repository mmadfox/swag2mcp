package service

import (
	"context"
	"fmt"
	"sort"
)

type (
	// EndpointsByTagRequest represents a request to list all endpoints for a given tag.
	EndpointsByTagRequest struct {
		TagID string `json:"tagId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsByTagResponse represents a response to list all endpoints for a given tag.
	EndpointsByTagResponse struct {
		Spec       Spec              `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection        `json:"collection" jsonschema:"required,Collection"`
		Tag        TagListItem       `json:"tag"        jsonschema:"required,Tag"`
		Endpoints  []EndpointTagItem `json:"endpoints"  jsonschema:"required,List of endpoints associated with the tag"`
	}

	// EndpointsByCollectionRequest represents a request to list all endpoints for a given collection.
	EndpointsByCollectionRequest struct {
		CollectionID string `json:"collectionId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsByCollectionResponse represents a response to list all endpoints for a given collection.
	EndpointsByCollectionResponse struct {
		Spec       Spec                     `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection               `json:"collection" jsonschema:"required,Collection"`
		Endpoints  []EndpointCollectionItem `json:"endpoints"  jsonschema:"required,List of endpoints associated with the collection"`
	}

	// EndpointsBySpecRequest represents a request to list all endpoints for a given spec.
	EndpointsBySpecRequest struct {
		SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsBySpecResponse represents a response to list all endpoints for a given spec.
	EndpointsBySpecResponse struct {
		Endpoints []EndpointSearchItem `json:"endpoints" jsonschema:"required,List of endpoints associated with the spec"`
	}

	// EndpointByIDRequest represents a request to get an endpoint by its ID.
	EndpointByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"required,Unique identifier for the endpoint"`
	}

	// EndpointByIDResponse represents a response to get an endpoint by its ID.
	EndpointByIDResponse struct {
		Spec       Spec        `json:"spec"       jsonschema:"required,Specification"`
		Collection Collection  `json:"collection" jsonschema:"required,Collection"`
		Tag        TagListItem `json:"tag"        jsonschema:"required,Tag"`
		Endpoint   Endpoint    `json:"endpoint"   jsonschema:"required,"`
	}
)

// EndpointsByTag returns a list of all available endpoints for a given tag.
func (s *Service) EndpointsByTag(_ context.Context, req EndpointsByTagRequest) (EndpointsByTagResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointsByTagResponse{}, NewValidationError(
			"tagId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	tag, err := s.index.TagByID(req.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", req.TagID), err)
	}

	collection, err := s.index.CollectionByID(tag.CollectionID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", tag.CollectionID), err)
	}

	spec, err := s.index.SpecByID(collection.SpecID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", collection.SpecID), err)
	}

	endpoints, err := s.index.EndpointsByTag(req.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", req.TagID), err)
	}

	resp := EndpointsByTagResponse{
		Spec: Spec{
			ID:     spec.ID,
			Domain: spec.Domain,
		},
		Collection: Collection{
			ID:           collection.ID,
			Title:        collection.Title,
			CountMethods: collection.Stats.Methods,
		},
		Tag: TagListItem{
			ID:           tag.ID,
			Title:        tag.Name,
			CountMethods: tag.Stats.Methods,
		},
		Endpoints: make([]EndpointTagItem, 0, len(endpoints)),
	}
	for _, ep := range endpoints {
		resp.Endpoints = append(resp.Endpoints, EndpointTagItem{
			ID:      ep.ID,
			Method:  ep.Name,
			Path:    ep.Path,
			Summary: ep.SummaryOrFallback(),
		})
	}

	sort.Slice(resp.Endpoints, func(i, j int) bool {
		return resp.Endpoints[i].Method < resp.Endpoints[j].Method
	})

	return resp, nil
}

// EndpointsByCollection returns a list of all available endpoints for a given collection.
func (s *Service) EndpointsByCollection(
	_ context.Context,
	req EndpointsByCollectionRequest,
) (EndpointsByCollectionResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointsByCollectionResponse{}, NewValidationError(
			"collectionId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	collection, cErr := s.index.CollectionByID(req.CollectionID)
	if cErr != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("collection %q not found", req.CollectionID),
			cErr,
		)
	}

	spec, sErr := s.index.SpecByID(collection.SpecID)
	if sErr != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("spec %q not found", collection.SpecID),
			sErr,
		)
	}

	endpoints, eErr := s.index.EndpointByCollection(req.CollectionID)
	if eErr != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(
			fmt.Sprintf("collection %q not found", req.CollectionID),
			eErr,
		)
	}

	resp := EndpointsByCollectionResponse{
		Spec: Spec{
			ID:     spec.ID,
			Domain: spec.Domain,
		},
		Collection: Collection{
			ID:           collection.ID,
			Title:        collection.Title,
			CountMethods: collection.Stats.Methods,
		},
		Endpoints: make([]EndpointCollectionItem, 0, len(endpoints)),
	}
	for _, ep := range endpoints {
		epTag, tErr := s.index.TagByID(ep.TagID)
		if tErr != nil {
			return EndpointsByCollectionResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", ep.TagID), tErr)
		}
		resp.Endpoints = append(resp.Endpoints, EndpointCollectionItem{
			ID:      ep.ID,
			TagID:   ep.TagID,
			TagName: epTag.Name,
			Method:  ep.Name,
			Path:    ep.Path,
			Summary: ep.SummaryOrFallback(),
		})
	}

	sort.Slice(resp.Endpoints, func(i, j int) bool {
		return resp.Endpoints[i].TagID < resp.Endpoints[j].TagID
	})

	return resp, nil
}

// EndpointsBySpec returns a list of all available endpoints for a given spec.
func (s *Service) EndpointsBySpec(_ context.Context, req EndpointsBySpecRequest) (EndpointsBySpecResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointsBySpecResponse{}, NewValidationError(
			"specId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	endpoints, err := s.index.EndpointsBySpec(req.SpecID)
	if err != nil {
		return EndpointsBySpecResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", req.SpecID), err)
	}

	items, itemsErr := s.mapEndpointsToSearchItems(endpoints)
	if itemsErr != nil {
		return EndpointsBySpecResponse{}, itemsErr
	}

	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.SpecID != b.SpecID {
			return a.SpecID < b.SpecID
		}
		if a.CollectionID != b.CollectionID {
			return a.CollectionID < b.CollectionID
		}
		return a.TagID < b.TagID
	})

	return EndpointsBySpecResponse{Endpoints: items}, nil
}

// EndpointByID returns an endpoint by its ID.
func (s *Service) EndpointByID(_ context.Context, req EndpointByIDRequest) (EndpointByIDResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointByIDResponse{}, NewValidationError(
			"endpointId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	ep, err := s.index.EndpointByID(req.ID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.ID), err)
	}

	spec, err := s.index.SpecByID(ep.SpecID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", ep.SpecID), err)
	}
	collection, err := s.index.CollectionByID(ep.CollectionID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", ep.CollectionID), err)
	}
	tag, err := s.index.TagByID(ep.TagID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", ep.TagID), err)
	}

	resp := EndpointByIDResponse{
		Spec: Spec{
			ID:     spec.ID,
			Domain: spec.Domain,
		},
		Collection: Collection{
			ID:           collection.ID,
			Title:        collection.Title,
			CountMethods: collection.Stats.Methods,
		},
		Tag: TagListItem{
			ID:           tag.ID,
			Title:        tag.Name,
			CountMethods: tag.Stats.Methods,
		},
		Endpoint: Endpoint{
			ID:      ep.ID,
			Method:  ep.Name,
			Path:    ep.Path,
			Summary: ep.SummaryOrFallback(),
		},
	}

	return resp, nil
}
