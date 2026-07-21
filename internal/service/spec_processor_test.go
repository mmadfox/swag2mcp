package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/stretchr/testify/require"
)

func TestResolveTagName_withTags(t *testing.T) {
	t.Parallel()

	result := resolveTagName([]string{"pets", "store"})
	require.Equal(t, "pets,store", result)
}

func TestResolveTagName_empty(t *testing.T) {
	t.Parallel()

	result := resolveTagName(nil)
	require.Equal(t, "default", result)
}

func TestApplySpecMetadata_setsTitle(t *testing.T) {
	t.Parallel()

	coll := &model.Collection{}
	doc := &spec.Doc{Title: "Pet Store API", Description: "A pet store"}
	applySpecMetadata(coll, doc)
	require.Equal(t, "Pet Store API", coll.Title)
	require.Equal(t, "Pet Store API", coll.LLMTitle)
	require.Equal(t, "A pet store", coll.LLMInstruction)
}

func TestApplySpecMetadata_preservesExisting(t *testing.T) {
	t.Parallel()

	coll := &model.Collection{
		LLMTitle:       "Custom Title",
		LLMInstruction: "Custom instruction",
	}
	doc := &spec.Doc{Title: "Doc Title", Description: "Doc description"}
	applySpecMetadata(coll, doc)
	require.Equal(t, "Custom Title", coll.LLMTitle)
	require.Equal(t, "Custom instruction", coll.LLMInstruction)
	require.Equal(t, "Doc Title", coll.Title)
}

func TestBuildSpecInfo_basic(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	sc := &config.Spec{
		Domain:   "test-api",
		BaseURL:  "https://api.example.com",
		LLMTitle: "Test API",
	}
	sp, err := svc.buildSpecInfo(sc, false, nil)
	require.NoError(t, err)
	require.NotNil(t, sp)
	require.Equal(t, "test-api", sp.Domain)
	require.Equal(t, "https://api.example.com", sp.BaseURL)
	require.Equal(t, "Test API", sp.LLMTitle)
}

func TestBuildSpecInfo_withAuth(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	sc := &config.Spec{
		Domain:  "auth-api",
		BaseURL: "https://api.example.com",
		Auth: config.Auth{
			Client: auth.NewNoAuthClient(),
		},
	}
	sp, err := svc.buildSpecInfo(sc, false, nil)
	require.NoError(t, err)
	require.NotNil(t, sp)
	require.NotNil(t, sp.Auth)
}

func TestBuildSpecInfo_withHTTPClient(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	sc := &config.Spec{
		Domain:  "http-api",
		BaseURL: "https://api.example.com",
		HTTPClient: &config.HTTPClientConfig{
			UserAgent: "spec-agent",
		},
	}
	sp, err := svc.buildSpecInfo(sc, false, nil)
	require.NoError(t, err)
	require.NotNil(t, sp)
	require.NotNil(t, sp.HTTPClient)
	require.Equal(t, "spec-agent", sp.HTTPClient.UserAgent)
}

func TestBuildSpecInfo_withMockAuth(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	sc := &config.Spec{
		Domain:  "mock-api",
		BaseURL: "https://api.example.com",
		Auth: config.Auth{
			Client: auth.NewNoAuthClient(),
		},
	}
	sp, err := svc.buildSpecInfo(sc, true, &config.MockAuthConfig{
		OAuth2Port: 9095,
		DigestPort: 9096,
	})
	require.NoError(t, err)
	require.NotNil(t, sp)
}
