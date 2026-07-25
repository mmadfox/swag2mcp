package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// Service is the core business logic layer for swag2mcp.
type Service struct {
	ctx             *serviceContext
	index           *index.Index
	cache           *cache.Cache
	ws              *workspace.Workspace
	v               *validator.Validate
	dumpDir         string
	version         string
	indexNoFullText bool
	snapshot        atomic.Value // stores *InfoSnapshot

	specSvc       *specService
	collectionSvc *collectionService
	tagSvc        *tagService
	endpointSvc   *endpointService
	searchSvc     *searchService
	inspectSvc    *inspectService
	authSvc       *authService
	toolsSvc      *toolsService
	invokeSvc     *invokeService
	infoSvc       *infoService
	exportSvc     *exportService
	importSvc     *importService
	responseSvc   *responseService
}

// Option is a functional option for configuring a Service.
type Option func(*Service)

// WithDisableLLMAuth configures whether the auth tool is disabled.
func WithDisableLLMAuth(disable bool) Option {
	return func(s *Service) {
		s.ctx.disableLLMAuth.Store(disable)
	}
}

// WithDumpDir configures the directory for dumping HTTP request traces.
func WithDumpDir(dir string) Option {
	return func(s *Service) {
		s.dumpDir = dir
	}
}

// WithVersion configures the version string for the service.
func WithVersion(version string) Option {
	return func(s *Service) {
		s.version = version
	}
}

// WithIndexNoFullText disables full-text search indexing.
func WithIndexNoFullText() Option {
	return func(s *Service) {
		s.indexNoFullText = true
	}
}

// WithWorkspace sets a custom workspace for the service.
func WithWorkspace(ws *workspace.Workspace) Option {
	return func(s *Service) {
		s.ws = ws
		s.cache = cache.New(ws.Root())
	}
}

// New creates a new Service with the given options.
func New(opts ...Option) (*Service, error) {
	s := &Service{
		ctx: newServiceContext(),
		v:   validator.New(validator.WithRequiredStructEnabled()),
	}

	for _, opt := range opts {
		opt(s)
	}

	idxOpts := []index.Option{}
	if s.indexNoFullText {
		idxOpts = append(idxOpts, index.WithNoFullText())
	}
	idx, err := index.New(idxOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	s.index = idx

	if s.ws == nil {
		ws, wsErr := workspace.New("")
		if wsErr != nil {
			return nil, fmt.Errorf("failed to create workspace: %w", wsErr)
		}
		s.ws = ws
	}
	if s.cache == nil {
		s.cache = cache.New(s.ws.Root())
	}

	llmAuthDisabled := func() bool { return s.ctx.disableLLMAuth.Load() }

	s.specSvc = newSpecService(s.index, s.v)
	s.collectionSvc = newCollectionService(s.index, s.v)
	s.tagSvc = newTagService(s.index, s.v)
	s.endpointSvc = newEndpointService(s.index, s.v)
	s.searchSvc = newSearchService(s.index, s.v)
	s.inspectSvc = newInspectService(s.index, s.v)
	s.authSvc = newAuthService(s.index, llmAuthDisabled)
	s.toolsSvc = newToolsService(s.index, llmAuthDisabled)
	s.invokeSvc = newInvokeService(s.ctx, s.index, s.ws, s.v, s.dumpDir)
	s.infoSvc = newInfoService(s.ctx, realClock{}, s.index, s.ws, s.version, &s.snapshot, s.ctx.startedAt.Load())
	s.exportSvc = newExportService(s.ws, s.version)
	s.importSvc = newImportService(s.ws)
	s.responseSvc = newResponseService(s.ctx, s.ws, s.v)

	return s, nil
}

// Workspace returns the workspace associated with the service.
func (s *Service) Workspace() *workspace.Workspace {
	return s.ws
}

// Specs returns a list of all available specifications.
func (s *Service) Specs(ctx context.Context) (SpecsResponse, error) {
	return s.specSvc.Specs(ctx)
}

// SpecByID returns the specification identified by the given spec ID,
// along with its associated collections.
func (s *Service) SpecByID(ctx context.Context, rq SpecByIDRequest) (SpecByIDResponse, error) {
	return s.specSvc.SpecByID(ctx, rq)
}

// CollectionsBySpec returns a list of all available collections for a given spec.
func (s *Service) CollectionsBySpec(ctx context.Context, rq CollectionsRequest) (CollectionsResponse, error) {
	return s.collectionSvc.CollectionsBySpec(ctx, rq)
}

// CollectionByID returns a collection by its ID, including its spec and tags.
func (s *Service) CollectionByID(ctx context.Context, rq CollectionByIDRequest) (CollectionByIDResponse, error) {
	return s.collectionSvc.CollectionByID(ctx, rq)
}

// TagsByCollection returns a list of all available tags for a given collection.
func (s *Service) TagsByCollection(ctx context.Context, rq TagsByCollectionRequest) (TagsByCollectionResponse, error) {
	return s.tagSvc.TagsByCollection(ctx, rq)
}

// TagByID returns a tag by its ID.
func (s *Service) TagByID(ctx context.Context, rq TagByIDRequest) (TagByIDResponse, error) {
	return s.tagSvc.TagByID(ctx, rq)
}

// TagsBySpec returns a list of all available tags for a given spec.
func (s *Service) TagsBySpec(ctx context.Context, rq TagsBySpecRequest) (TagsBySpecResponse, error) {
	return s.tagSvc.TagsBySpec(ctx, rq)
}

// EndpointsByTag returns all endpoints associated with the given tag,
// along with the parent spec, collection, and tag metadata.
func (s *Service) EndpointsByTag(ctx context.Context, rq EndpointsByTagRequest) (EndpointsByTagResponse, error) {
	return s.endpointSvc.EndpointsByTag(ctx, rq)
}

// EndpointsByCollection returns all endpoints within the given collection,
// along with the parent spec and collection metadata.
func (s *Service) EndpointsByCollection(ctx context.Context, rq EndpointsByCollectionRequest) (EndpointsByCollectionResponse, error) {
	return s.endpointSvc.EndpointsByCollection(ctx, rq)
}

// EndpointsBySpec returns all endpoints belonging to the given spec.
func (s *Service) EndpointsBySpec(ctx context.Context, rq EndpointsBySpecRequest) (EndpointsBySpecResponse, error) {
	return s.endpointSvc.EndpointsBySpec(ctx, rq)
}

// EndpointByID returns the full details for a single endpoint identified by
// its unique endpoint ID, including the parent spec, collection, and tag.
func (s *Service) EndpointByID(ctx context.Context, rq EndpointByIDRequest) (EndpointByIDResponse, error) {
	return s.endpointSvc.EndpointByID(ctx, rq)
}

// Search performs a full-text search across all endpoints using the given
// query string and returns up to the specified limit of matching results.
func (s *Service) Search(ctx context.Context, rq SearchRequest) (SearchResponse, error) {
	return s.searchSvc.Search(ctx, rq)
}

// Inspect returns the full endpoint details for the given endpoint ID,
// including the HTTP method, path, base URL, full URL, and the complete
// OpenAPI operation specification.
func (s *Service) Inspect(ctx context.Context, rq InspectRequest) (InspectResponse, error) {
	return s.inspectSvc.Inspect(ctx, rq)
}

// Auth retrieves authentication information for the spec identified by the
// request SpecID. It applies the spec's auth configuration and returns the
// resulting token, headers, and query parameters.
func (s *Service) Auth(ctx context.Context, rq AuthRequest) (AuthResponse, error) {
	return s.authSvc.Auth(ctx, rq)
}

// MakeToolDefinitions loads tool descriptions from embedded markdown files
// and returns the complete set of MCP tools with their descriptions.
func (s *Service) MakeToolDefinitions() (ToolDefinitions, error) {
	return s.toolsSvc.MakeToolDefinitions()
}

// Invoke validates the request, builds an HTTP request, sends it, and returns the response.
func (s *Service) Invoke(ctx context.Context, rq InvokeRequest) (InvokeResponse, error) {
	return s.invokeSvc.Invoke(ctx, rq)
}

// Info returns a comprehensive summary of the current service state from the snapshot.
func (s *Service) Info(ctx context.Context) (InfoResponse, error) {
	return s.infoSvc.Info(ctx)
}

// Export creates a portable ZIP backup of the workspace.
func (s *Service) Export(ctx context.Context, req ExportRequest) (ExportResponse, error) {
	return s.exportSvc.Export(ctx, req)
}

// Import imports spec files into the workspace specs/ directory.
func (s *Service) Import(ctx context.Context, req ImportRequest) (ImportResponse, error) {
	return s.importSvc.Import(ctx, req)
}

// ResponseOutline returns a high-level structural summary of a saved response.
func (s *Service) ResponseOutline(ctx context.Context, req ResponseOutlineRequest) (ResponseOutlineResponse, error) {
	return s.responseSvc.ResponseOutline(ctx, req)
}

// ResponseCompress reduces a JSON value in a saved response file.
func (s *Service) ResponseCompress(ctx context.Context, req ResponseCompressRequest) (ResponseCompressResponse, error) {
	return s.responseSvc.ResponseCompress(ctx, req)
}

// ResponseSlice extracts a fragment of a saved response file.
func (s *Service) ResponseSlice(ctx context.Context, req ResponseSliceRequest) (ResponseSliceResponse, error) {
	return s.responseSvc.ResponseSlice(ctx, req)
}
