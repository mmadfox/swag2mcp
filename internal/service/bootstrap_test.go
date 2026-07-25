package service

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

func TestResolveTagName_Empty(t *testing.T) {
	t.Parallel()

	name := resolveTagName(nil)
	if name != "default" {
		t.Errorf("got %q, want %q", name, "default")
	}
}

func TestResolveTagName_Single(t *testing.T) {
	t.Parallel()

	name := resolveTagName([]string{"pets"})
	if name != "pets" {
		t.Errorf("got %q, want %q", name, "pets")
	}
}

func TestResolveTagName_Multiple(t *testing.T) {
	t.Parallel()

	name := resolveTagName([]string{"pets", "store"})
	if name != "pets,store" {
		t.Errorf("got %q, want %q", name, "pets,store")
	}
}

func TestApplySpecMetadata_EmptyCollection(t *testing.T) {
	t.Parallel()

	coll := &types.Collection{}
	doc := &spec.Doc{
		Title:       "Pet Store API",
		Description: "A sample pet store API",
	}
	applySpecMetadata(coll, doc)
	if coll.LLMTitle != "Pet Store API" {
		t.Errorf("LLMTitle = %q, want %q", coll.LLMTitle, "Pet Store API")
	}
	if coll.LLMInstruction != "A sample pet store API" {
		t.Errorf("LLMInstruction = %q, want %q", coll.LLMInstruction, "A sample pet store API")
	}
	if coll.Title != "Pet Store API" {
		t.Errorf("Title = %q, want %q", coll.Title, "Pet Store API")
	}
}

func TestApplySpecMetadata_PreexistingValues(t *testing.T) {
	t.Parallel()

	coll := &types.Collection{
		LLMTitle:       "Custom Title",
		LLMInstruction: "Custom Instruction",
	}
	doc := &spec.Doc{
		Title:       "Doc Title",
		Description: "Doc Description",
	}
	applySpecMetadata(coll, doc)
	if coll.LLMTitle != "Custom Title" {
		t.Errorf("LLMTitle = %q, want %q", coll.LLMTitle, "Custom Title")
	}
	if coll.LLMInstruction != "Custom Instruction" {
		t.Errorf("LLMInstruction = %q, want %q", coll.LLMInstruction, "Custom Instruction")
	}
	if coll.Title != "Doc Title" {
		t.Errorf("Title = %q, want %q", coll.Title, "Doc Title")
	}
}

func TestApplySpecMetadata_EmptyDoc(t *testing.T) {
	t.Parallel()

	coll := &types.Collection{}
	doc := &spec.Doc{}
	applySpecMetadata(coll, doc)
	if coll.LLMTitle != "" {
		t.Errorf("LLMTitle = %q, want empty", coll.LLMTitle)
	}
	if coll.LLMInstruction != "" {
		t.Errorf("LLMInstruction = %q, want empty", coll.LLMInstruction)
	}
}

func TestConvertCookies_Nil(t *testing.T) {
	t.Parallel()

	result := convertCookies(nil)
	if result != nil {
		t.Fatal("expected nil")
	}
}

func TestConvertCookies_Empty(t *testing.T) {
	t.Parallel()

	result := convertCookies([]config.Cookie{})
	if result != nil {
		t.Fatal("expected nil")
	}
}

func TestConvertCookies_Populated(t *testing.T) {
	t.Parallel()

	input := []config.Cookie{
		{Name: "session", Value: "abc123", Domain: "example.com", Path: "/", Secure: true, HTTPOnly: true},
		{Name: "theme", Value: "dark"},
	}
	result := convertCookies(input)
	if len(result) != 2 {
		t.Fatalf("len = %d, want 2", len(result))
	}
	if result[0].Name != "session" || result[0].Value != "abc123" {
		t.Errorf("cookie[0] = %+v", result[0])
	}
	if !result[0].Secure || !result[0].HTTPOnly {
		t.Error("Secure/HTTPOnly not preserved")
	}
	if result[1].Name != "theme" || result[1].Value != "dark" {
		t.Errorf("cookie[1] = %+v", result[1])
	}
}

func TestMergeHTTPClientConfig_AllNil(t *testing.T) {
	t.Parallel()

	result := mergeHTTPClientConfig(nil, nil)
	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestMergeHTTPClientConfig_SpecOnly(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{Headers: map[string]string{"X-Spec": "s"}}
	result := mergeHTTPClientConfig(spec, nil)
	if result.Headers["X-Spec"] != "s" {
		t.Errorf("X-Spec = %q, want %q", result.Headers["X-Spec"], "s")
	}
}

func TestMergeHTTPClientConfig_CollectionOverridesSpec(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{Headers: map[string]string{"X-Header": "spec"}}
	coll := &config.HTTPClientConfig{Headers: map[string]string{"X-Header": "coll"}}
	result := mergeHTTPClientConfig(spec, coll)
	// First-wins: spec sets headers first
	if result.Headers["X-Header"] != "spec" {
		t.Errorf("X-Header = %q, want %q", result.Headers["X-Header"], "spec")
	}
}

func TestMergeHTTPClientConfig_HeadersMerge(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{Headers: map[string]string{"X-Spec": "s"}}
	coll := &config.HTTPClientConfig{Headers: map[string]string{"X-Coll": "c"}}
	result := mergeHTTPClientConfig(spec, coll)
	// First-wins: only spec headers are set
	if result.Headers["X-Spec"] != "s" {
		t.Errorf("X-Spec = %q", result.Headers["X-Spec"])
	}
	if result.Headers["X-Coll"] != "" {
		t.Errorf("X-Coll should be empty, got %q", result.Headers["X-Coll"])
	}
}

func TestMergeHTTPClientConfig_Cookies(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{Cookies: []config.Cookie{{Name: "s", Value: "2"}}}
	result := mergeHTTPClientConfig(spec, nil)
	if len(result.Cookies) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Cookies))
	}
	if result.Cookies[0].Name != "s" {
		t.Errorf("Cookie name = %q, want %q", result.Cookies[0].Name, "s")
	}
}

func TestMergeHTTPClientConfig_CookiesCollectionOverrides(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{Cookies: []config.Cookie{{Name: "s", Value: "2"}}}
	coll := &config.HTTPClientConfig{Cookies: []config.Cookie{{Name: "c", Value: "3"}}}
	result := mergeHTTPClientConfig(spec, coll)
	// First-wins: spec cookies are used
	if len(result.Cookies) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Cookies))
	}
	if result.Cookies[0].Name != "s" {
		t.Errorf("Cookie name = %q, want %q", result.Cookies[0].Name, "s")
	}
}
