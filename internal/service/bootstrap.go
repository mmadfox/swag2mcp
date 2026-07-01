package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/id"
	specparser "github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

// BootstrapRequest is the request for the Bootstrap method.
type BootstrapRequest struct {
	ConfFilepath string
	Tags         []string
}

// Bootstrap bootstraps the service.
//
//nolint:gocognit,funlen
func (s *Service) Bootstrap(_ context.Context, r BootstrapRequest) error {
	filter := config.NewFilter(r.Tags)
	conf, err := s.loadAndValidateConfig(r.ConfFilepath, filter)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	s.makeWorkspaceDir(conf.WorkspaceDir)

	for spec := range conf.Iterate(filter) {
		specInfo, serr := s.specFromConfig(spec)
		if serr != nil {
			return serr
		}

		allTags := make(map[string]*types.Tag)               // per spec
		allCollections := make(map[string]*types.Collection) // per spec
		allEndpoints := make(map[string]*types.Endpoint)     // per spec

		for _, col := range spec.Collections {
			if col.Disable {
				continue
			}

			colInfo := &types.Collection{
				ID:             id.Collection(specInfo.ID, col.Location),
				SpecID:         specInfo.ID,
				LLMTitle:       col.LLMTitle,
				LLMInstruction: col.LLMInstruction,
				BaseURL:        col.BaseURL,
				Headers:        col.Headers,
			}
			allCollections[colInfo.ID] = colInfo

			specInfo.Stats.Collections++

			localPath := col.Location
			if s.cache != nil {
				if localPath, err = s.cache.Resolve(col.Location); err != nil {
					return fmt.Errorf("collection %q: %w", col.Location, err)
				}
			}

			data, readErr := os.ReadFile(localPath)
			if readErr != nil {
				return fmt.Errorf("collection %q: read file: %w", col.Location, readErr)
			}

			specDoc, parseErr := specparser.Parse(data)
			if err != nil {
				return fmt.Errorf("collection %q: parse spec: %w", col.Location, parseErr)
			}
			if len(colInfo.LLMTitle) == 0 && len(specDoc.Title) > 0 {
				colInfo.LLMTitle = specDoc.Title
			}
			if len(colInfo.LLMInstruction) == 0 && len(specDoc.Description) > 0 {
				colInfo.LLMInstruction = specDoc.Description
			}
			colInfo.Title = specDoc.Title

			for _, pi := range specDoc.PathItems {
				op := pi.Operation
				if op == nil {
					continue
				}

				endpointTags := op.Tags
				if len(endpointTags) == 0 {
					endpointTags = []string{"default"}
				}

				tagName := strings.Join(endpointTags, ",")
				tagID := id.Tag(specInfo.ID, colInfo.ID, tagName)
				tagInfo, ok := allTags[tagID]
				if !ok {
					colInfo.Stats.Tags++
					tagInfo = &types.Tag{
						ID:           tagID,
						SpecID:       specInfo.ID,
						CollectionID: colInfo.ID,
						Name:         tagName,
					}
					allTags[tagID] = tagInfo
				}
				colInfo.Stats.Methods++
				tagInfo.Stats.Methods++

				endpoint := types.Endpoint{
					ID: id.Method(
						specInfo.ID,
						colInfo.ID,
						tagID,
						pi.Method,
						pi.Path,
						op.ID,
					),
					SpecID:       specInfo.ID,
					CollectionID: colInfo.ID,
					TagID:        tagID,
					Name:         pi.Method,
					Path:         pi.Path,
					Operation:    op,
				}
				allEndpoints[endpoint.ID] = &endpoint
			}
		}

		if indexErr := s.indexSpec(specInfo, allCollections, allTags, allEndpoints); indexErr != nil {
			return indexErr
		}
	}

	return nil
}

func (s *Service) specFromConfig(spec *config.Spec) (*types.Spec, error) {
	specInfo := &types.Spec{
		ID:             id.Domain(spec.Domain),
		Domain:         spec.Domain,
		LLMTitle:       spec.LLMTitle,
		LLMInstruction: spec.LLMInstruction,
		BaseURL:        spec.BaseURL,
		Headers:        spec.Headers,
		Auth:           spec.Auth.Client,
	}

	if err := specInfo.InitAuthenticator(); err != nil {
		return nil, fmt.Errorf("spec %s, failed to initialize authenticator: %w", spec.Domain, err)
	}

	return specInfo, nil
}

func (s *Service) loadAndValidateConfig(filepath string, filter *config.Filter) (*config.Config, error) {
	conf, loadErr := config.Load(filepath)
	if loadErr != nil {
		return nil, loadErr
	}

	if err := conf.Validate(filter); err != nil {
		return nil, err
	}

	return conf, nil
}

func (s *Service) makeWorkspaceDir(workspaceDir string) {
	s.cache.SetWorkspaceDir(workspaceDir)
	// TODO: make all dirs
}

func (s *Service) indexSpec(
	specInfo *types.Spec,
	allCollections map[string]*types.Collection,
	allTags map[string]*types.Tag,
	allEndpoints map[string]*types.Endpoint,
) error {
	colls := make([]*types.Collection, 0, len(allCollections))
	for _, c := range allCollections {
		colls = append(colls, c)
	}

	tags := make([]*types.Tag, 0, len(allTags))
	for _, t := range allTags {
		tags = append(tags, t)
	}

	ends := make([]*types.Endpoint, 0, len(allEndpoints))
	for _, e := range allEndpoints {
		ends = append(ends, e)
	}

	if err := s.index.EnsureIndex(specInfo, colls, tags, ends); err != nil {
		return fmt.Errorf("failed to ensure index: %w", err)
	}
	return nil
}
