package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mmadfox/swag2mcp/internal/cache"
	"github.com/mmadfox/swag2mcp/internal/index"
)

type Service struct {
	index *index.Index
	cache *cache.Cache
	v     *validator.Validate
}

func New() (*Service, error) {
	idx, err := index.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	return &Service{
		index: idx,
		cache: cache.New("./"),
		v: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}, nil
}

func (s *Service) validateRequest(typ any) error {
	return s.v.Struct(typ)
}
