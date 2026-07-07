package service

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/spec"
)

type (
	// InspectRequest represents a request to inspect an endpoint.
	InspectRequest struct {
		EndpointID string `json:"endpointId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to inspect"`
	}

	// InspectResponse represents a response to inspect an endpoint.
	InspectResponse struct {
		ID           string            `json:"id"                jsonschema:"required,Unique identifier for the endpoint"`
		TagID        string            `json:"tagId"             jsonschema:"required,Unique identifier for the tag"`
		CollectionID string            `json:"collectionId"      jsonschema:"required,Unique identifier for the collection"`
		SpecID       string            `json:"specId"            jsonschema:"required,Unique identifier for the spec"`
		SpecDomain   string            `json:"specDomain"        jsonschema:"required,Domain of the spec"`
		Method       string            `json:"method"            jsonschema:"required,HTTP method (GET, POST, etc.)"`
		Path         string            `json:"path"              jsonschema:"required,API path"`
		BaseURL      string            `json:"baseUrl"           jsonschema:"required,Base URL of the API"`
		FullURL      string            `json:"fullUrl"           jsonschema:"required,Full URL of the endpoint"`
		Operation    *spec.Operation   `json:"operation"         jsonschema:"required,Operation details"`
		Headers      map[string]string `json:"headers,omitempty" jsonschema:"optional,Headers to be sent with the request"`
	}
)

// Inspect returns the details of an endpoint.
func (s *Service) Inspect(_ context.Context, req InspectRequest) (InspectResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return InspectResponse{}, NewValidationError(
			"endpointId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	ep, eErr := s.index.EndpointByID(req.EndpointID)
	if eErr != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.EndpointID), eErr)
	}

	spec, sErr := s.index.SpecByID(ep.SpecID)
	if sErr != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", ep.SpecID), sErr)
	}
	collection, cErr := s.index.CollectionByID(ep.CollectionID)
	if cErr != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", ep.CollectionID), cErr)
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
		SpecDomain:   spec.Domain,
		Method:       ep.Name,
		Path:         ep.Path,
		Operation:    ep.Operation,
		BaseURL:      baseURL,
		FullURL:      baseURL + "/" + strings.TrimLeft(ep.Path, "/"),
	}

	if collection.HTTPClient != nil && len(collection.HTTPClient.Headers) > 0 {
		resp.Headers = make(map[string]string, len(collection.HTTPClient.Headers))
		maps.Copy(resp.Headers, collection.HTTPClient.Headers)
	}

	return resp, nil
}
