package service

import (
	"context"
	"fmt"

	"github.com/mmadfox/swag2mcp/internal/spec"
)

type (
	InspectRequest struct {
		EndpointID string `json:"endpointId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to inspect"`
	}

	InspectResponse struct {
		ID           string            `json:"id" jsonschema:"required,Unique identifier for the endpoint"`
		TagID        string            `json:"tagId" jsonschema:"required,Unique identifier for the tag"`
		CollectionID string            `json:"collectionId" jsonschema:"required,Unique identifier for the collection"`
		SpecID       string            `json:"specId" jsonschema:"required,Unique identifier for the spec"`
		Method       string            `json:"method" jsonschema:"required,HTTP method (GET, POST, etc.)"`
		Path         string            `json:"path" jsonschema:"required,API path"`
		BaseURL      string            `json:"baseUrl" jsonschema:"required,Base URL of the API"`
		Operation    *spec.Operation   `json:"operation" jsonschema:"required,Operation details"`
		Headers      map[string]string `json:"headers,omitempty" jsonschema:"optional,Headers to be sent with the request"`
	}
)

// Inspect returns the details of an endpoint.
func (s *Service) Inspect(_ context.Context, req InspectRequest) (InspectResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return InspectResponse{}, NewValidationError("endpointId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	ep, err := s.index.EndpointByID(req.EndpointID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.EndpointID), err)
	}

	spec, err := s.index.SpecByID(ep.SpecID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", ep.SpecID), err)
	}
	collection, err := s.index.CollectionByID(ep.CollectionID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", ep.CollectionID), err)
	}
	baseURL := spec.BaseURL
	if len(collection.BaseURL) > 0 {
		baseURL = collection.BaseURL
	}

	resp := InspectResponse{
		ID:           ep.ID,
		TagID:        ep.TagID,
		CollectionID: ep.CollectionID,
		SpecID:       ep.SpecID,
		Method:       ep.Name,
		Path:         ep.Path,
		Operation:    ep.Operation,
		BaseURL:      baseURL,
	}

	if len(spec.Headers) > 0 || len(collection.Headers) > 0 {
		resp.Headers = make(map[string]string, len(spec.Headers)+len(collection.Headers))
		for k, v := range spec.Headers {
			resp.Headers[k] = v
		}
		for k, v := range collection.Headers {
			resp.Headers[k] = v
		}
	}

	return resp, nil
}
