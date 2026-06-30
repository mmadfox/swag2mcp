package service

import (
	"context"
	"fmt"
)

type (
	// EndpointsByTagRequest represents a request to list all endpoints for a given tag.
	EndpointsByTagRequest struct {
		TagID string `json:"tagId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsByTagResponse represents a response to list all endpoints for a given tag.
	EndpointsByTagResponse struct {
		Endpoints []EndpointItem `json:"endpoints" jsonschema:"required,List of endpoints associated with the tag"`
	}

	// EndpointsByCollectionRequest represents a request to list all endpoints for a given collection.
	EndpointsByCollectionRequest struct {
		CollectionID string `json:"collectionId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsByCollectionResponse represents a response to list all endpoints for a given collection.
	EndpointsByCollectionResponse struct {
		Endpoints []EndpointItem `json:"endpoints" jsonschema:"required,List of endpoints associated with the collection"`
	}

	// EndpointsBySpecRequest represents a request to list all endpoints for a given spec.
	EndpointsBySpecRequest struct {
		SpecID string `json:"specId" jsonschema:"required," validate:"required,md5"`
	}

	// EndpointsBySpecResponse represents a response to list all endpoints for a given spec.
	EndpointsBySpecResponse struct {
		Endpoints []EndpointItem `json:"endpoints" jsonschema:"required,List of endpoints associated with the spec"`
	}

	// EndpointByIDRequest represents a request to get an endpoint by its ID.
	EndpointByIDRequest struct {
		ID string `json:"id" validate:"required,md5" jsonschema:"required,Unique identifier for the endpoint"`
	}

	// EndpointByIDResponse represents a response to get an endpoint by its ID.
	EndpointByIDResponse struct {
		Endpoint EndpointItem `json:"endpoint" jsonschema:"required,"`
	}
)

// EndpointsByTag returns a list of all available endpoints for a given tag.
func (s *Service) EndpointsByTag(_ context.Context, req EndpointsByTagRequest) (EndpointsByTagResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointsByTagResponse{}, NewValidationError("tagId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	endpoints, err := s.index.EndpointsByTag(req.TagID)
	if err != nil {
		return EndpointsByTagResponse{}, NewNotFoundError(fmt.Sprintf("tag %q not found", req.TagID), err)
	}

	resp := EndpointsByTagResponse{
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

// EndpointsByCollection returns a list of all available endpoints for a given collection.
func (s *Service) EndpointsByCollection(_ context.Context, req EndpointsByCollectionRequest) (EndpointsByCollectionResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointsByCollectionResponse{}, NewValidationError("collectionId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	endpoints, err := s.index.EndpointByCollection(req.CollectionID)
	if err != nil {
		return EndpointsByCollectionResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", req.CollectionID), err)
	}

	resp := EndpointsByCollectionResponse{
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

// EndpointsBySpec returns a list of all available endpoints for a given spec.
func (s *Service) EndpointsBySpec(_ context.Context, req EndpointsBySpecRequest) (EndpointsBySpecResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointsBySpecResponse{}, NewValidationError("specId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	endpoints, err := s.index.EndpointsBySpec(req.SpecID)
	if err != nil {
		return EndpointsBySpecResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", req.SpecID), err)
	}

	resp := EndpointsBySpecResponse{
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

// EndpointByID returns an endpoint by its ID.
func (s *Service) EndpointByID(_ context.Context, req EndpointByIDRequest) (EndpointByIDResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return EndpointByIDResponse{}, NewValidationError("endpointId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	ep, err := s.index.EndpointByID(req.ID)
	if err != nil {
		return EndpointByIDResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.ID), err)
	}

	resp := EndpointByIDResponse{
		Endpoint: EndpointItem{
			ID:           ep.ID,
			TagID:        ep.TagID,
			CollectionID: ep.CollectionID,
			SpecID:       ep.SpecID,
			Method:       ep.Name,
			Path:         ep.Path,
			Summary:      ep.SummaryOrFallback(),
			Deprecated:   ep.Operation != nil && ep.Operation.Deprecated,
		},
	}

	return resp, nil
}
