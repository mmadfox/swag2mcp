package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"strings"
)

type inspectService struct {
	index IndexReader
	v     RequestValidator
}

func newInspectService(index IndexReader, v RequestValidator) *inspectService {
	return &inspectService{index: index, v: v}
}

// Inspect returns the full endpoint details for the given endpoint ID,
// including the HTTP method, path, base URL, full URL, and the complete
// OpenAPI operation specification.
func (is *inspectService) Inspect(
	_ context.Context,
	rq InspectRequest,
) (InspectResponse, error) {
	if err := is.v.Struct(rq); err != nil {
		return InspectResponse{}, NewValidationError(
			"The endpoint ID is invalid. It must be a 32-character hex string. "+
				"Use the search tool to find the correct endpoint ID.",
			err,
		)
	}

	e, err := is.index.EndpointByID(rq.EndpointID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(
			fmt.Sprintf("Endpoint %q was not found. Use the search tool to find the correct endpoint ID.", rq.EndpointID),
			err,
		)
	}

	sp, err := is.index.SpecByID(e.SpecID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(
			fmt.Sprintf("Spec %q was not found. The endpoint references a spec that no longer exists.", e.SpecID),
			err,
		)
	}
	coll, err := is.index.CollectionByID(e.CollectionID)
	if err != nil {
		return InspectResponse{}, NewNotFoundError(
			fmt.Sprintf("Collection %q was not found. The endpoint references a collection that no longer exists.", e.CollectionID),
			err,
		)
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
