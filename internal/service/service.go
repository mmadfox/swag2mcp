package service

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

// Service is the core business logic layer for swag2mcp.
// It manages the search index, cache, workspace, HTTP client, rate limiter, and configuration.
type Service struct {
	index              *index.Index
	cache              *cache.Cache
	ws                 *workspace.Workspace
	v                  *validator.Validate
	disableLLMAuth     atomic.Bool
	dumpDir            string
	rateLimiter        *invokeRateLimiter
	httpClient         *http.Client
	httpClientConfig   httpclient.Config
	maxResponseSize    int
	version            string
	startedAt          time.Time
	config             *config.Config
	disableRateLimiter atomic.Bool
	indexNoFullText    bool
	snapshot           atomic.Value // stores *InfoSnapshot
}

// NewOption is a functional option for configuring a Service.
type NewOption func(*Service)

// WithDisableLLMAuth configures whether the auth tool is disabled.
func WithDisableLLMAuth(disable bool) NewOption {
	return func(s *Service) {
		s.disableLLMAuth.Store(disable)
	}
}

// WithDumpDir configures the directory for dumping HTTP request traces.
func WithDumpDir(dir string) NewOption {
	return func(s *Service) {
		s.dumpDir = dir
	}
}

// WithVersion configures the version string for the service.
func WithVersion(version string) NewOption {
	return func(s *Service) {
		s.version = version
	}
}

// WithIndexNoFullText disables full-text search indexing.
// Use this for CLI commands that only need in-memory lookups (e.g. info).
func WithIndexNoFullText() NewOption {
	return func(s *Service) {
		s.indexNoFullText = true
	}
}

// New creates a new Service with the given options.
func New(opts ...NewOption) (*Service, error) {
	s := &Service{
		v:               validator.New(validator.WithRequiredStructEnabled()),
		rateLimiter:     newInvokeRateLimiter(),
		httpClient:      &http.Client{Transport: http.DefaultTransport},
		maxResponseSize: defaultMaxResponseSize,
		startedAt:       time.Now(),
	}
	for _, opt := range opts {
		opt(s)
	}

	idxOpts := []index.NewOption{}
	if s.indexNoFullText {
		idxOpts = append(idxOpts, index.WithNoFullText())
	}
	idx, err := index.New(idxOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	s.index = idx

	ws, err := workspace.New("")
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}
	s.ws = ws
	s.cache = cache.New(ws.Root())

	return s, nil
}

func (s *Service) validateRequest(typ any) error {
	return s.v.Struct(typ)
}

// Workspace returns the workspace associated with the service.
func (s *Service) Workspace() *workspace.Workspace {
	return s.ws
}
