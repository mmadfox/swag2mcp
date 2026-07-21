package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

func (s *Service) processSpec(ctx context.Context, sc *config.Spec, mockEnabled bool, ma *config.MockAuthConfig) error {
	sp, err := s.buildSpecInfo(sc, mockEnabled, ma)
	if err != nil {
		return err
	}

	tags := make(map[string]*model.Tag)
	colls := make(map[string]*model.Collection)
	eps := make(map[string]*model.Endpoint)

	for i := range sc.Collections {
		cc := &sc.Collections[i]
		if cc.Disable {
			continue
		}

		coll, err := s.processCollection(
			ctx, sp, sc, cc,
			tags, eps,
		)
		if err != nil {
			return err
		}

		colls[coll.ID] = coll
		sp.Stats.Collections++
	}

	sp.Stats.Tags = len(tags)
	sp.Stats.Methods = len(eps)

	return s.indexSpec(sp, colls, tags, eps)
}

func (s *Service) processCollection(
	ctx context.Context,
	sp *model.Spec,
	sc *config.Spec,
	cc *config.Collection,
	tags map[string]*model.Tag,
	eps map[string]*model.Endpoint,
) (*model.Collection, error) {
	coll := &model.Collection{
		ID:             id.Collection(sp.ID, cc.Location),
		SpecID:         sp.ID,
		LLMTitle:       cc.LLMTitle,
		LLMInstruction: cc.LLMInstruction,
		BaseURL:        cc.BaseURL,
		BaseMockURL:    cc.BaseMockURL,
		HTTPClient: mergeHTTPClientConfig(
			sc.HTTPClient,
			cc.HTTPClient,
		),
	}

	doc, err := s.parseSpecDocument(ctx, cc.Location)
	if err != nil {
		return nil, err
	}

	applySpecMetadata(coll, doc)

	for _, pi := range doc.PathItems {
		op := pi.Operation
		if op == nil {
			continue
		}

		tagName := resolveTagName(op.Tags)
		tagID := id.Tag(sp.ID, coll.ID, tagName)

		tag, exists := tags[tagID]
		if !exists {
			coll.Stats.Tags++
			tag = &model.Tag{
				ID:           tagID,
				SpecID:       sp.ID,
				CollectionID: coll.ID,
				Name:         tagName,
			}
			tags[tagID] = tag
		}

		coll.Stats.Methods++
		tag.Stats.Methods++

		ep := model.Endpoint{
			ID: id.Method(
				sp.ID,
				coll.ID,
				tagID,
				pi.Method,
				pi.Path,
				op.ID,
			),
			SpecID:       sp.ID,
			CollectionID: coll.ID,
			TagID:        tagID,
			Tag:          tagName,
			Name:         pi.Method,
			Path:         pi.Path,
			Operation:    op,
		}
		eps[ep.ID] = &ep
	}

	return coll, nil
}

func (s *Service) parseSpecDocument(ctx context.Context, location string) (*spec.Doc, error) {
	lp := location

	if s.cache != nil {
		var err error
		lp, err = s.cache.Resolve(ctx, location)
		if err != nil {
			return nil, fmt.Errorf("collection %q: %w", location, err)
		}
	}

	data, err := os.ReadFile(lp)
	if err != nil {
		return nil, fmt.Errorf("collection %q: read file: %w", location, err)
	}

	doc, err := spec.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("collection %q: parse spec: %w", location, err)
	}

	return doc, nil
}

func (s *Service) buildSpecInfo(sc *config.Spec, mockEnabled bool, ma *config.MockAuthConfig) (*model.Spec, error) {
	sp := &model.Spec{
		ID:             id.Domain(sc.Domain),
		Domain:         sc.Domain,
		LLMTitle:       sc.LLMTitle,
		LLMInstruction: sc.LLMInstruction,
		BaseURL:        sc.BaseURL,
		Auth:           sc.Auth.Client,
	}

	if sc.HTTPClient != nil {
		sc.HTTPClient.Resolve()
		sp.HTTPClient = configToModelHTTPClient(sc.HTTPClient)
	}

	if sc.Auth.Client != nil {
		if err := sp.InitAuthenticator(); err != nil {
			return nil, fmt.Errorf(
				"spec %s, failed to initialize authenticator: %w",
				sc.Domain, err,
			)
		}

		if scriptClient, isScript := sc.Auth.Client.(*auth.ScriptAuthClient); isScript {
			scriptClient.SetWorkspaceDir(s.ws.Root())
		}

		if mockEnabled {
			applyMockAuthURLs(sc.Auth.Client, ma)
		}
	}

	return sp, nil
}

func resolveTagName(tags []string) string {
	if len(tags) > 0 {
		return strings.Join(tags, ",")
	}
	return "default"
}

func applySpecMetadata(collection *model.Collection, specDocument *spec.Doc) {
	if len(collection.LLMTitle) == 0 && len(specDocument.Title) > 0 {
		collection.LLMTitle = specDocument.Title
	}
	if len(collection.LLMInstruction) == 0 && len(specDocument.Description) > 0 {
		collection.LLMInstruction = specDocument.Description
	}
	collection.Title = specDocument.Title
}
