package service

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// BootstrapRequest is the request for the Bootstrap method.
type BootstrapRequest struct {
	ConfFilepath string
	Tags         []string
}

// Bootstrap loads the configuration, initializes the workspace, and indexes
// all specs, collections, tags, and endpoints into the service.
func (s *Service) Bootstrap(_ context.Context, request BootstrapRequest) error {
	configuration, loadError := s.loadConfiguration(request.ConfFilepath, request.Tags)
	if loadError != nil {
		return fmt.Errorf("failed to load config: %w", loadError)
	}

	if initError := s.initializeWorkspace(filepath.Dir(request.ConfFilepath)); initError != nil {
		return initError
	}

	filter := config.NewFilter(request.Tags)

	for specConfig := range configuration.Iterate(filter) {
		if specError := s.processSpec(configuration, specConfig); specError != nil {
			return specError
		}
	}

	return nil
}

func (s *Service) loadConfiguration(configFilepath string, tags []string) (*config.Config, error) {
	configuration, loadError := config.Load(configFilepath)
	if loadError != nil {
		return nil, loadError
	}

	filter := config.NewFilter(tags)
	if err := configuration.Validate(filter); err != nil {
		return nil, err
	}

	return configuration, nil
}

func (s *Service) initializeWorkspace(workspaceDirectory string) error {
	if workspaceDirectory != "" && workspaceDirectory != s.ws.Root() {
		newWorkspace, workspaceError := workspace.New(workspaceDirectory)
		if workspaceError != nil {
			return fmt.Errorf("failed to create workspace: %w", workspaceError)
		}
		s.ws = newWorkspace
	}

	if initError := s.ws.Init(); initError != nil {
		return fmt.Errorf("failed to init workspace: %w", initError)
	}

	s.cache.SetWorkspaceDir(s.ws.Root())
	return nil
}

func (s *Service) processSpec(configuration *config.Config, specConfig *config.Spec) error {
	specification, specError := s.buildSpecInfo(specConfig)
	if specError != nil {
		return specError
	}

	allTags := make(map[string]*types.Tag)
	allCollections := make(map[string]*types.Collection)
	allEndpoints := make(map[string]*types.Endpoint)

	for index := range specConfig.Collections {
		collectionConfig := &specConfig.Collections[index]
		if collectionConfig.Disable {
			continue
		}

		collectionInfo, processError := s.processCollection(
			configuration, specification, specConfig, collectionConfig,
			allTags, allEndpoints,
		)
		if processError != nil {
			return processError
		}

		allCollections[collectionInfo.ID] = collectionInfo
		specification.Stats.Collections++
	}

	return s.indexSpec(specification, allCollections, allTags, allEndpoints)
}

func (s *Service) processCollection(
	configuration *config.Config,
	specification *types.Spec,
	specConfig *config.Spec,
	collectionConfig *config.Collection,
	allTags map[string]*types.Tag,
	allEndpoints map[string]*types.Endpoint,
) (*types.Collection, error) {
	collectionInfo := &types.Collection{
		ID:             id.Collection(specification.ID, collectionConfig.Location),
		SpecID:         specification.ID,
		LLMTitle:       collectionConfig.LLMTitle,
		LLMInstruction: collectionConfig.LLMInstruction,
		BaseURL:        collectionConfig.BaseURL,
		HTTPClient: mergeHTTPClientConfig(
			configuration.HTTPClient,
			specConfig.HTTPClient,
			collectionConfig.HTTPClient,
		),
	}

	specDocument, parseError := s.parseSpecDocument(collectionConfig.Location)
	if parseError != nil {
		return nil, parseError
	}

	applySpecMetadata(collectionInfo, specDocument)

	for _, pathItem := range specDocument.PathItems {
		operation := pathItem.Operation
		if operation == nil {
			continue
		}

		tagName := resolveTagName(operation.Tags)
		tagID := id.Tag(specification.ID, collectionInfo.ID, tagName)

		tagInfo, tagExists := allTags[tagID]
		if !tagExists {
			collectionInfo.Stats.Tags++
			tagInfo = &types.Tag{
				ID:           tagID,
				SpecID:       specification.ID,
				CollectionID: collectionInfo.ID,
				Name:         tagName,
			}
			allTags[tagID] = tagInfo
		}

		collectionInfo.Stats.Methods++
		tagInfo.Stats.Methods++

		endpoint := types.Endpoint{
			ID: id.Method(
				specification.ID,
				collectionInfo.ID,
				tagID,
				pathItem.Method,
				pathItem.Path,
				operation.ID,
			),
			SpecID:       specification.ID,
			CollectionID: collectionInfo.ID,
			TagID:        tagID,
			Tag:          tagName,
			Name:         pathItem.Method,
			Path:         pathItem.Path,
			Operation:    operation,
		}
		allEndpoints[endpoint.ID] = &endpoint
	}

	return collectionInfo, nil
}

func (s *Service) parseSpecDocument(location string) (*spec.Doc, error) {
	localPath := location
	var resolveError error

	if s.cache != nil {
		localPath, resolveError = s.cache.Resolve(location)
		if resolveError != nil {
			return nil, fmt.Errorf("collection %q: %w", location, resolveError)
		}
	}

	data, readError := os.ReadFile(localPath)
	if readError != nil {
		return nil, fmt.Errorf("collection %q: read file: %w", location, readError)
	}

	specDocument, parseError := spec.Parse(data)
	if parseError != nil {
		return nil, fmt.Errorf("collection %q: parse spec: %w", location, parseError)
	}

	return specDocument, nil
}

func (s *Service) buildSpecInfo(specConfig *config.Spec) (*types.Spec, error) {
	specification := &types.Spec{
		ID:             id.Domain(specConfig.Domain),
		Domain:         specConfig.Domain,
		LLMTitle:       specConfig.LLMTitle,
		LLMInstruction: specConfig.LLMInstruction,
		BaseURL:        specConfig.BaseURL,
		Auth:           specConfig.Auth.Client,
	}

	if specConfig.HTTPClient != nil {
		specification.HTTPClient = &types.HTTPClientConfig{
			Headers:         specConfig.HTTPClient.Headers,
			Cookies:         convertCookies(specConfig.HTTPClient.Cookies),
			UserAgent:       specConfig.HTTPClient.UserAgent,
			Timeout:         specConfig.HTTPClient.Timeout,
			FollowRedirects: specConfig.HTTPClient.FollowRedirects,
			MaxRedirects:    specConfig.HTTPClient.MaxRedirects,
			MaxResponseSize: specConfig.HTTPClient.MaxResponseSize,
		}
	}

	if specConfig.Auth.Client != nil {
		if initError := specification.InitAuthenticator(); initError != nil {
			return nil, fmt.Errorf(
				"spec %s, failed to initialize authenticator: %w",
				specConfig.Domain, initError,
			)
		}

		if scriptClient, isScript := specConfig.Auth.Client.(*auth.ScriptAuthClient); isScript {
			scriptClient.SetWorkspaceDir(s.ws.Root())
		}
	}

	return specification, nil
}

func (s *Service) indexSpec(
	specification *types.Spec,
	allCollections map[string]*types.Collection,
	allTags map[string]*types.Tag,
	allEndpoints map[string]*types.Endpoint,
) error {
	collections := make([]*types.Collection, 0, len(allCollections))
	for _, collection := range allCollections {
		collections = append(collections, collection)
	}

	tags := make([]*types.Tag, 0, len(allTags))
	for _, tag := range allTags {
		tags = append(tags, tag)
	}

	endpoints := make([]*types.Endpoint, 0, len(allEndpoints))
	for _, endpoint := range allEndpoints {
		endpoints = append(endpoints, endpoint)
	}

	if indexError := s.index.EnsureIndex(specification, collections, tags, endpoints); indexError != nil {
		return fmt.Errorf("failed to ensure index: %w", indexError)
	}

	return nil
}

func resolveTagName(tags []string) string {
	if len(tags) > 0 {
		return strings.Join(tags, ",")
	}
	return "default"
}

func applySpecMetadata(collection *types.Collection, specDocument *spec.Doc) {
	if len(collection.LLMTitle) == 0 && len(specDocument.Title) > 0 {
		collection.LLMTitle = specDocument.Title
	}
	if len(collection.LLMInstruction) == 0 && len(specDocument.Description) > 0 {
		collection.LLMInstruction = specDocument.Description
	}
	collection.Title = specDocument.Title
}

func convertCookies(cookies []config.Cookie) []types.Cookie {
	if len(cookies) == 0 {
		return nil
	}

	result := make([]types.Cookie, len(cookies))
	for index, cookie := range cookies {
		result[index] = types.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
		}
	}

	return result
}

// mergeHTTPClientConfig merges HTTP client configs with cascade:
// global → spec → collection. Fields set at a lower level override higher levels.
func mergeHTTPClientConfig(
	global, spec, collection *config.HTTPClientConfig,
) *types.HTTPClientConfig {
	result := &types.HTTPClientConfig{}

	levels := []*config.HTTPClientConfig{global, spec, collection}

	for _, level := range levels {
		if level == nil {
			continue
		}
		mergeHeaders(result, level)
		mergeScalars(result, level)
		mergeCookies(result, level)
	}

	return result
}

func mergeHeaders(result *types.HTTPClientConfig, level *config.HTTPClientConfig) {
	if result.Headers == nil && len(level.Headers) > 0 {
		result.Headers = make(map[string]string, len(level.Headers))
		maps.Copy(result.Headers, level.Headers)
	}
}

func mergeScalars(result *types.HTTPClientConfig, level *config.HTTPClientConfig) {
	if result.UserAgent == "" && level.UserAgent != "" {
		result.UserAgent = level.UserAgent
	}
	if result.Timeout == 0 && level.Timeout != 0 {
		result.Timeout = level.Timeout
	}
	if result.FollowRedirects == nil && level.FollowRedirects != nil {
		result.FollowRedirects = level.FollowRedirects
	}
	if result.MaxRedirects == nil && level.MaxRedirects != nil {
		result.MaxRedirects = level.MaxRedirects
	}
	if result.MaxResponseSize == nil && level.MaxResponseSize != nil {
		result.MaxResponseSize = level.MaxResponseSize
	}
}

func mergeCookies(result *types.HTTPClientConfig, level *config.HTTPClientConfig) {
	if len(result.Cookies) > 0 || len(level.Cookies) == 0 {
		return
	}

	result.Cookies = make([]types.Cookie, len(level.Cookies))
	for index, cookie := range level.Cookies {
		result.Cookies[index] = types.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
		}
	}
}
