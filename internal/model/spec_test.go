package model

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
}

func TestInitAuthenticator_Success(t *testing.T) {
	t.Parallel()

	s := &Spec{
		Auth: &mockAuth{newFunc: func() error { return nil }},
	}
	err := s.InitAuthenticator()
	require.NoError(t, err)
}

func TestInitAuthenticator_Error(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("auth init failed")
	s := &Spec{
		Auth: &mockAuth{newFunc: func() error { return expectedErr }},
	}
	err := s.InitAuthenticator()
	require.Error(t, err)
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
	assert.Equal(t, "List all users", result)
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
	assert.Equal(t, "Returns a list of all users", result)
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
	assert.Equal(t, "POST /orders", result)
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
	assert.Equal(t, "Delete a user by ID", result)
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
	assert.Equal(t, "session", c.Name)
	assert.Equal(t, "abc123", c.Value)
	assert.Equal(t, "example.com", c.Domain)
	assert.Equal(t, "/", c.Path)
	assert.True(t, c.Secure)
	assert.True(t, c.HTTPOnly)
}

func TestHTTPClientConfig_Fields(t *testing.T) {
	t.Parallel()

	cfg := &HTTPClientConfig{
		Headers: map[string]string{"X-Custom": "val"},
		Cookies: []httpclient.Cookie{{Name: "c", Value: "v"}},
	}
	assert.Equal(t, "val", cfg.Headers["X-Custom"])
	assert.Len(t, cfg.Cookies, 1)
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
	assert.Equal(t, "spec-1", s.ID)
	assert.Equal(t, "test-api", s.Domain)
	assert.Equal(t, "Test API", s.LLMTitle)
	assert.Equal(t, "Test instruction", s.LLMInstruction)
	assert.Equal(t, "https://api.example.com", s.BaseURL)
}

func TestSpec_Stats(t *testing.T) {
	t.Parallel()

	s := &Spec{}
	s.Stats.Collections = 3
	s.Stats.Tags = 10
	s.Stats.Methods = 25
	assert.Equal(t, 3, s.Stats.Collections)
	assert.Equal(t, 10, s.Stats.Tags)
	assert.Equal(t, 25, s.Stats.Methods)
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
	assert.Equal(t, "coll-1", c.ID)
	assert.Equal(t, "spec-1", c.SpecID)
	assert.Equal(t, "My Collection", c.LLMTitle)
	assert.Equal(t, "Collection instruction", c.LLMInstruction)
	assert.Equal(t, "My Collection Title", c.Title)
	assert.Equal(t, "https://coll.example.com", c.BaseURL)
}

func TestCollection_Stats(t *testing.T) {
	t.Parallel()

	c := &Collection{}
	c.Stats.Tags = 5
	c.Stats.Methods = 15
	assert.Equal(t, 5, c.Stats.Tags)
	assert.Equal(t, 15, c.Stats.Methods)
}

func TestTag_Fields(t *testing.T) {
	t.Parallel()

	tg := &Tag{
		ID:           "tag-1",
		CollectionID: "coll-1",
		SpecID:       "spec-1",
		Name:         "pets",
	}
	assert.Equal(t, "tag-1", tg.ID)
	assert.Equal(t, "coll-1", tg.CollectionID)
	assert.Equal(t, "spec-1", tg.SpecID)
	assert.Equal(t, "pets", tg.Name)
}

func TestTag_Stats(t *testing.T) {
	t.Parallel()

	tg := &Tag{}
	tg.Stats.Methods = 7
	assert.Equal(t, 7, tg.Stats.Methods)
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
	assert.Equal(t, "ep-1", e.ID)
	assert.Equal(t, "tag-1", e.TagID)
	assert.Equal(t, "coll-1", e.CollectionID)
	assert.Equal(t, "spec-1", e.SpecID)
	assert.Equal(t, "GET", e.Name)
	assert.Equal(t, "/users", e.Path)
	assert.Equal(t, "users-tag", e.Tag)
	require.NotNil(t, e.Operation)
	assert.Equal(t, "getUsers", e.Operation.ID)
}
