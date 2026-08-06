/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"log/slog"
	"strings"
)

type inspectService struct {
	index IndexReader
	v     RequestValidator
	log   *slog.Logger
}

func newInspectService(index IndexReader, v RequestValidator, log *slog.Logger) *inspectService {
	return &inspectService{index: index, v: v, log: log}
}

// Inspect returns the full endpoint details for the given endpoint ID,
// including the HTTP method, path, base URL, full URL, and the complete
// OpenAPI operation specification.
func (is *inspectService) Inspect(
	ctx context.Context,
	rq InspectRequest,
) (InspectResponse, error) {
	if err := is.v.Struct(rq); err != nil {
		return InspectResponse{}, NewInvalidEndpointIDError(err)
	}

	e, err := is.index.EndpointByID(rq.EndpointID)
	if err != nil {
		is.log.ErrorContext(ctx, "inspect failed: endpoint not found", "endpoint_id", rq.EndpointID, "error", err)
		return InspectResponse{}, NewEndpointNotFoundError(rq.EndpointID, err)
	}

	sp, err := is.index.SpecByID(e.SpecID)
	if err != nil {
		is.log.ErrorContext(ctx, "inspect failed: spec not found", "spec_id", e.SpecID, "error", err)
		return InspectResponse{}, NewSpecNotFoundError(e.SpecID, err)
	}
	coll, err := is.index.CollectionByID(e.CollectionID)
	if err != nil {
		is.log.ErrorContext(ctx, "inspect failed: collection not found", "collection_id", e.CollectionID, "error", err)
		return InspectResponse{}, NewCollectionNotFoundError(e.CollectionID, err)
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
