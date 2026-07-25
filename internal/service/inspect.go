package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/spec"
)

// InspectRequest contains the endpoint ID used to retrieve full endpoint
// details including the operation specification.
type (
	InspectRequest struct {
		EndpointID string `json:"endpointId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to inspect"`
	}

	// InspectResponse contains the full details of an endpoint, including its
	// HTTP method, path, URLs, and the complete OpenAPI operation specification.
	InspectResponse struct {
		ID           string          `json:"id"                jsonschema:"required,Unique identifier for the endpoint"`
		TagID        string          `json:"tagId"             jsonschema:"required,Unique identifier for the tag"`
		CollectionID string          `json:"collectionId"      jsonschema:"required,Unique identifier for the collection"`
		SpecID       string          `json:"specId"            jsonschema:"required,Unique identifier for the spec"`
		SpecDomain   string          `json:"specDomain"        jsonschema:"required,Domain of the spec"`
		Method       string          `json:"method"            jsonschema:"required,HTTP method (GET, POST, etc.)"`
		Path         string          `json:"path"              jsonschema:"required,API path"`
		BaseURL      string          `json:"baseUrl"           jsonschema:"required,Base URL of the API"`
		FullURL      string          `json:"fullUrl"           jsonschema:"required,Full URL of the endpoint"`
		Operation    *spec.Operation `json:"operation"         jsonschema:"required,Operation details"`
	}
)

// Inspect returns the full endpoint details for the given endpoint ID,
// including the HTTP method, path, base URL, full URL, and the complete
// OpenAPI operation specification.
func (s *Service) Inspect(_ context.Context, rq InspectRequest) (InspectResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return InspectResponse{}, NewValidationError(
			"The endpoint ID is invalid — it must be a 32-character hex string. Use the search tool to find the correct endpoint ID.",
			err,
		)
	}

	e, err := s.index.EndpointByID(rq.EndpointID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("Endpoint %q not found — use the search tool to find the correct endpoint ID.", rq.EndpointID), err)
	}

	sp, err := s.index.SpecByID(e.SpecID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("Spec %q not found — the endpoint references a spec that no longer exists.", e.SpecID), err)
	}
	coll, err := s.index.CollectionByID(e.CollectionID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("Collection %q not found — the endpoint references a collection that no longer exists.", e.CollectionID), err)
	}
	baseURL := sp.BaseURL
	if len(coll.BaseURL) > 0 {
		baseURL = coll.BaseURL
	}

	r := InspectResponse{
		ID:           e.ID,
		TagID:        e.TagID,
		CollectionID: e.CollectionID,
		SpecID:       e.SpecID,
		SpecDomain:   sp.Domain,
		Method:       e.Name,
		Path:         e.Path,
		Operation:    e.Operation,
		BaseURL:      baseURL,
		FullURL:      baseURL + "/" + strings.TrimLeft(e.Path, "/"),
	}

	return r, nil
}
