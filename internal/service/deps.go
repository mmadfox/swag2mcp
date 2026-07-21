package service

import (
	"context"
	"time"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
)

// IndexReader provides read-only access to the search index.
type IndexReader interface {
	EndpointByID(id string) (*model.Endpoint, error)
	SpecByID(id string) (*model.Spec, error)
	CollectionByID(id string) (*model.Collection, error)
	AllSpecs() []*model.Spec
	EndpointsByTag(tagID string) ([]*model.Endpoint, error)
	EndpointByCollection(collectionID string) ([]*model.Endpoint, error)
	EndpointsBySpec(specID string) ([]*model.Endpoint, error)
	TagByID(id string) (*model.Tag, error)
	TagsByCollection(collectionID string) ([]*model.Tag, error)
	TagsBySpec(specID string) ([]*model.Tag, error)
	CollectionsBySpec(specID string) ([]*model.Collection, error)
	Search(ctx context.Context, query string, limit int) ([]*model.Endpoint, error)
}

// WorkspaceOps provides workspace operations needed by sub-services.
type WorkspaceOps interface {
	Root() string
	ResponsesDir() string
	ConfigPath() string
	ConfigNotExists() bool
	CreateExportDir() (string, error)
	DownloadSpec(ctx context.Context, location string) ([]byte, error)
	SaveSpec(name string, data []byte) (string, error)
	ListSpecs() ([]string, error)
	SpecPath(name string) string
	Init() error
	CopySpecsToWorkspace(src string) error
	CopyAuthScriptsToWorkspace(src string) error
	CopyAuthScriptsToExport(dst string) error
}

// RequestValidator validates request structs.
type RequestValidator interface {
	Struct(any) error
}

// SnapshotStore provides atomic load/store for the info snapshot.
type SnapshotStore interface {
	Load() any
	Store(any)
}

// Clock provides time-related operations for uptime calculation.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// SettingsProvider provides runtime settings needed by sub-services.
type SettingsProvider interface {
	MaxResponseSize() int
	HTTPClientConfig() httpclient.Config
	Config() *config.Config
}
