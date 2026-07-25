package service

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/id"
	"github.com/mmadfox/swag2mcp/internal/index"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

// newTestService creates a Service with a fresh index and a single spec/collection/tag/endpoint.
// The domain is derived from t.Name() for parallel test isolation.
func newTestService(t *testing.T, opts ...NewOption) *Service {
	t.Helper()

	svc, err := New(opts...)
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	idx, err := index.New()
	if err != nil {
		t.Fatalf("index.New() = %v", err)
	}
	svc.index = idx

	return svc
}

// seedTestData populates the service index with a spec, collection, tag, and endpoint.
// Returns the spec, collection, tag, and endpoint for use in tests.
func seedTestData(t *testing.T, svc *Service, domain string) (*model.Spec, *model.Collection, *model.Tag, *model.Endpoint) {
	t.Helper()

	specID := id.Domain(domain)
	specInfo := &model.Spec{
		ID:      specID,
		Domain:  domain,
		BaseURL: "https://api.example.com",
	}

	collectionID := id.Collection(specID, domain+"/collection")
	collectionInfo := &model.Collection{
		ID:     collectionID,
		SpecID: specID,
		Title:  "Test Collection",
		Stats: struct {
			Tags    int `json:"tags"`
			Methods int `json:"methods"`
		}{Tags: 1, Methods: 1},
	}

	tagID := id.Tag(specID, collectionID, "test-tag")
	tagInfo := &model.Tag{
		ID:           tagID,
		SpecID:       specID,
		CollectionID: collectionID,
		Name:         "test-tag",
		Stats: struct {
			Methods int `json:"methods"`
		}{Methods: 1},
	}

	endpointID := id.Method(specID, collectionID, tagID, "GET", "/test", "testOp")
	endpointInfo := &model.Endpoint{
		ID:           endpointID,
		SpecID:       specID,
		CollectionID: collectionID,
		TagID:        tagID,
		Tag:          "test-tag",
		Name:         "GET",
		Path:         "/test",
		Operation: &spec.Operation{
			ID:          "testOp",
			Summary:     "Test endpoint",
			Description: "A test endpoint",
			Parameters: []*spec.Parameter{
				{Name: "id", In: "path", Required: true, Schema: &spec.Schema{Type: "string"}},
			},
		},
	}

	if err := svc.index.EnsureIndex(specInfo, []*model.Collection{collectionInfo}, []*model.Tag{tagInfo}, []*model.Endpoint{endpointInfo}); err != nil {
		t.Fatalf("EnsureIndex() = %v", err)
	}

	// Force refresh the bluge search reader so Search works immediately.
	svc.index.RefreshSearchReader()

	return specInfo, collectionInfo, tagInfo, endpointInfo
}

// seedTestDataWithAuth is like seedTestData but also sets an authenticator on the spec.
func seedTestDataWithAuth(t *testing.T, svc *Service, domain string, authenticator auth.Authenticator) (*model.Spec, *model.Collection, *model.Tag, *model.Endpoint) {
	t.Helper()

	specInfo, collectionInfo, tagInfo, endpointInfo := seedTestData(t, svc, domain)
	specInfo.Auth = authenticator
	return specInfo, collectionInfo, tagInfo, endpointInfo
}
