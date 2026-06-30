package service

import (
	"context"
	"fmt"
)

type (
	InvokeRequest struct {
		EndpointID  string         `json:"endpointId" validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to invoke"`
		Parameters  map[string]any `json:"parameters,omitempty" jsonschema:"optional,Path, query, and header parameters as key-value pairs"`
		RequestBody map[string]any `json:"requestBody,omitempty" jsonschema:"optional,Request body for POST/PUT/PATCH requests"`
	}

	InvokeResponse struct {
		StatusCode int              `json:"statusCode" jsonschema:"required,HTTP response status code"`
		Headers    map[string]string `json:"headers" jsonschema:"required,HTTP response headers"`
		Body       any              `json:"body" jsonschema:"required,Response body data"`
	}
)

func (s *Service) Invoke(_ context.Context, req InvokeRequest) (InvokeResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return InvokeResponse{}, NewValidationError("endpointId must be a 32-character lowercase hex string (MD5 format)", err)
	}

	if _, err := s.index.EndpointByID(req.EndpointID); err != nil {
		return InvokeResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.EndpointID), err)
	}

	return InvokeResponse{}, fmt.Errorf("invoke not yet implemented")
}
