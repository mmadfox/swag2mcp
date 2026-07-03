package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/workspace"
)

type Service struct {
	index *index.Index
	cache *cache.Cache
	ws    *workspace.Workspace
	v     *validator.Validate
}

func New() (*Service, error) {
	idx, err := index.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	ws, err := workspace.New("")
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}
	return &Service{
		index: idx,
		cache: cache.New(ws.Root()),
		ws:    ws,
		v: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}, nil
}

func (s *Service) validateRequest(typ any) error {
	return s.v.Struct(typ)
}
