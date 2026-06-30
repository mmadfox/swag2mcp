package service

import (
	"context"
	"fmt"
)

type (
	InspectRequest struct {
		EndpointID string `json:"endpointId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to inspect"`
	}

	InspectResponse struct {
		ID      string `json:"id" jsonschema:"required,Endpoint ID"`
		Method  string `json:"method" jsonschema:"required,HTTP method"`
		Path    string `json:"path" jsonschema:"required,API path"`
		Summary string `json:"summary" jsonschema:"required,Endpoint summary"`
	}
)

func (s *Service) Inspect(_ context.Context, req InspectRequest) (InspectResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return InspectResponse{}, NewValidationError("endpointId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	ep, err := s.index.EndpointByID(req.EndpointID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.EndpointID), err)
	}

	return InspectResponse{
		ID:      ep.ID,
		Method:  ep.Name,
		Path:    ep.Path,
		Summary: ep.SummaryOrFallback(),
	}, nil
}
