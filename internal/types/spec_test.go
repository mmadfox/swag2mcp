package types

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

type mockAuth struct {
	newFunc func() error
}

func (m *mockAuth) New() error {
	return m.newFunc()
}

func (m *mockAuth) Type() auth.Type {
	return auth.NoAuth
}

func (m *mockAuth) Apply(_ *http.Request, _ *auth.Info) error {
	return nil
}

func (m *mockAuth) Validate() error {
	return nil
}

func TestInitAuthenticator_Nil(t *testing.T) {
	t.Parallel()

	s := &Spec{}
	err := s.InitAuthenticator()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitAuthenticator_Success(t *testing.T) {
	t.Parallel()

	s := &Spec{
		Auth: &mockAuth{newFunc: func() error { return nil }},
	}
	err := s.InitAuthenticator()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitAuthenticator_Error(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("auth init failed")
	s := &Spec{
		Auth: &mockAuth{newFunc: func() error { return expectedErr }},
	}
	err := s.InitAuthenticator()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSummaryOrFallback_Summary(t *testing.T) {
	t.Parallel()

	e := &Endpoint{
		Name: "GET",
		Path: "/users",
		Operation: &spec.Operation{
			Summary:     "List all users",
			Description: "Returns a list of all users",
		},
	}
	result := e.SummaryOrFallback()
	if result != "List all users" {
		t.Errorf("got %q, want %q", result, "List all users")
	}
}

func TestSummaryOrFallback_Description(t *testing.T) {
	t.Parallel()

	e := &Endpoint{
		Name: "GET",
		Path: "/users",
		Operation: &spec.Operation{
			Summary:     "",
			Description: "Returns a list of all users",
		},
	}
	result := e.SummaryOrFallback()
	if result != "Returns a list of all users" {
		t.Errorf("got %q, want %q", result, "Returns a list of all users")
	}
}

func TestSummaryOrFallback_Fallback(t *testing.T) {
	t.Parallel()

	e := &Endpoint{
		Name: "POST",
		Path: "/orders",
		Operation: &spec.Operation{
			Summary:     "",
			Description: "",
		},
	}
	result := e.SummaryOrFallback()
	if result != "POST /orders" {
		t.Errorf("got %q, want %q", result, "POST /orders")
	}
}

func TestSummaryOrFallback_EmptySummary(t *testing.T) {
	t.Parallel()

	e := &Endpoint{
		Name: "DELETE",
		Path: "/users/{id}",
		Operation: &spec.Operation{
			Summary:     "",
			Description: "Delete a user by ID",
		},
	}
	result := e.SummaryOrFallback()
	if result != "Delete a user by ID" {
		t.Errorf("got %q, want %q", result, "Delete a user by ID")
	}
}

func TestCookie_Fields(t *testing.T) {
	t.Parallel()

	c := httpclient.Cookie{
		Name:     "session",
		Value:    "abc123",
		Domain:   "example.com",
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
	}
	if c.Name != "session" {
		t.Errorf("Name = %q", c.Name)
	}
	if c.Value != "abc123" {
		t.Errorf("Value = %q", c.Value)
	}
	if c.Domain != "example.com" {
		t.Errorf("Domain = %q", c.Domain)
	}
	if c.Path != "/" {
		t.Errorf("Path = %q", c.Path)
	}
	if !c.Secure {
		t.Error("Secure should be true")
	}
	if !c.HTTPOnly {
		t.Error("HTTPOnly should be true")
	}
}

func TestHTTPClientConfig_Fields(t *testing.T) {
	t.Parallel()

	cfg := &HTTPClientConfig{
		Headers: map[string]string{"X-Custom": "val"},
		Cookies: []httpclient.Cookie{{Name: "c", Value: "v"}},
	}
	if cfg.Headers["X-Custom"] != "val" {
		t.Error("Headers not set")
	}
	if len(cfg.Cookies) != 1 {
		t.Error("Cookies not set")
	}
}

func TestSpec_Fields(t *testing.T) {
	t.Parallel()

	s := &Spec{
		ID:             "spec-1",
		Domain:         "test-api",
		LLMTitle:       "Test API",
		LLMInstruction: "Test instruction",
		BaseURL:        "https://api.example.com",
	}
	if s.ID != "spec-1" {
		t.Errorf("ID = %q", s.ID)
	}
	if s.Domain != "test-api" {
		t.Errorf("Domain = %q", s.Domain)
	}
	if s.LLMTitle != "Test API" {
		t.Errorf("LLMTitle = %q", s.LLMTitle)
	}
	if s.LLMInstruction != "Test instruction" {
		t.Errorf("LLMInstruction = %q", s.LLMInstruction)
	}
	if s.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q", s.BaseURL)
	}
}

func TestSpec_Stats(t *testing.T) {
	t.Parallel()

	s := &Spec{}
	s.Stats.Collections = 3
	s.Stats.Tags = 10
	s.Stats.Methods = 25
	if s.Stats.Collections != 3 {
		t.Errorf("Collections = %d", s.Stats.Collections)
	}
	if s.Stats.Tags != 10 {
		t.Errorf("Tags = %d", s.Stats.Tags)
	}
	if s.Stats.Methods != 25 {
		t.Errorf("Methods = %d", s.Stats.Methods)
	}
}

func TestCollection_Fields(t *testing.T) {
	t.Parallel()

	c := &Collection{
		ID:             "coll-1",
		SpecID:         "spec-1",
		LLMTitle:       "My Collection",
		LLMInstruction: "Collection instruction",
		Title:          "My Collection Title",
		BaseURL:        "https://coll.example.com",
	}
	if c.ID != "coll-1" {
		t.Errorf("ID = %q", c.ID)
	}
	if c.SpecID != "spec-1" {
		t.Errorf("SpecID = %q", c.SpecID)
	}
	if c.LLMTitle != "My Collection" {
		t.Errorf("LLMTitle = %q", c.LLMTitle)
	}
	if c.LLMInstruction != "Collection instruction" {
		t.Errorf("LLMInstruction = %q", c.LLMInstruction)
	}
	if c.Title != "My Collection Title" {
		t.Errorf("Title = %q", c.Title)
	}
	if c.BaseURL != "https://coll.example.com" {
		t.Errorf("BaseURL = %q", c.BaseURL)
	}
}

func TestCollection_Stats(t *testing.T) {
	t.Parallel()

	c := &Collection{}
	c.Stats.Tags = 5
	c.Stats.Methods = 15
	if c.Stats.Tags != 5 {
		t.Errorf("Tags = %d", c.Stats.Tags)
	}
	if c.Stats.Methods != 15 {
		t.Errorf("Methods = %d", c.Stats.Methods)
	}
}

func TestTag_Fields(t *testing.T) {
	t.Parallel()

	tg := &Tag{
		ID:           "tag-1",
		CollectionID: "coll-1",
		SpecID:       "spec-1",
		Name:         "pets",
	}
	if tg.ID != "tag-1" {
		t.Errorf("ID = %q", tg.ID)
	}
	if tg.CollectionID != "coll-1" {
		t.Errorf("CollectionID = %q", tg.CollectionID)
	}
	if tg.SpecID != "spec-1" {
		t.Errorf("SpecID = %q", tg.SpecID)
	}
	if tg.Name != "pets" {
		t.Errorf("Name = %q", tg.Name)
	}
}

func TestTag_Stats(t *testing.T) {
	t.Parallel()

	tg := &Tag{}
	tg.Stats.Methods = 7
	if tg.Stats.Methods != 7 {
		t.Errorf("Methods = %d", tg.Stats.Methods)
	}
}

func TestEndpoint_Fields(t *testing.T) {
	t.Parallel()

	e := &Endpoint{
		ID:           "ep-1",
		TagID:        "tag-1",
		CollectionID: "coll-1",
		SpecID:       "spec-1",
		Name:         "GET",
		Path:         "/users",
		Tag:          "users-tag",
		Operation:    &spec.Operation{ID: "getUsers"},
	}
	if e.ID != "ep-1" {
		t.Errorf("ID = %q", e.ID)
	}
	if e.TagID != "tag-1" {
		t.Errorf("TagID = %q", e.TagID)
	}
	if e.CollectionID != "coll-1" {
		t.Errorf("CollectionID = %q", e.CollectionID)
	}
	if e.SpecID != "spec-1" {
		t.Errorf("SpecID = %q", e.SpecID)
	}
	if e.Name != "GET" {
		t.Errorf("Name = %q", e.Name)
	}
	if e.Path != "/users" {
		t.Errorf("Path = %q", e.Path)
	}
	if e.Tag != "users-tag" {
		t.Errorf("Tag = %q", e.Tag)
	}
	if e.Operation == nil {
		t.Fatal("Operation is nil")
	}
	if e.Operation.ID != "getUsers" {
		t.Errorf("Operation.ID = %q", e.Operation.ID)
	}
}
