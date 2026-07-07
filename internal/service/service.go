package service

import (
	"fmt"
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

type Service struct {
	index          *index.Index
	cache          *cache.Cache
	ws             *workspace.Workspace
	v              *validator.Validate
	disableLLMAuth atomic.Bool
	dumpDir        string
	rateLimiter    *invokeRateLimiter
}

type NewOption func(*Service)

func WithDisableLLMAuth(disable bool) NewOption {
	return func(s *Service) {
		s.disableLLMAuth.Store(disable)
	}
}

func WithDumpDir(dir string) NewOption {
	return func(s *Service) {
		s.dumpDir = dir
	}
}

func New(opts ...NewOption) (*Service, error) {
	idx, err := index.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	ws, err := workspace.New("")
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}
	s := &Service{
		index:       idx,
		cache:       cache.New(ws.Root()),
		ws:          ws,
		v:           validator.New(validator.WithRequiredStructEnabled()),
		rateLimiter: newInvokeRateLimiter(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) validateRequest(typ any) error {
	return s.v.Struct(typ)
}

// Workspace returns the workspace associated with the service.
func (s *Service) Workspace() *workspace.Workspace {
	return s.ws
}
