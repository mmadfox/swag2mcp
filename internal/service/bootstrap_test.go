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

	result := mergeHTTPClientConfig(nil, nil, nil)
	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestMergeHTTPClientConfig_GlobalOnly(t *testing.T) {
	t.Parallel()

	ua := "global-agent"
	global := &config.HTTPClientConfig{UserAgent: ua}
	result := mergeHTTPClientConfig(global, nil, nil)
	if result.UserAgent != ua {
		t.Errorf("UserAgent = %q, want %q", result.UserAgent, ua)
	}
}

func TestMergeHTTPClientConfig_SpecOverridesGlobal(t *testing.T) {
	t.Parallel()

	global := &config.HTTPClientConfig{UserAgent: "global-agent", Timeout: 30}
	spec := &config.HTTPClientConfig{UserAgent: "spec-agent"}
	result := mergeHTTPClientConfig(global, spec, nil)
	// First-wins: global sets UserAgent first
	if result.UserAgent != "global-agent" {
		t.Errorf("UserAgent = %q, want %q", result.UserAgent, "global-agent")
	}
	if result.Timeout != 30 {
		t.Errorf("Timeout = %v, want %v", result.Timeout, 30)
	}
}

func TestMergeHTTPClientConfig_CollectionOverridesAll(t *testing.T) {
	t.Parallel()

	global := &config.HTTPClientConfig{UserAgent: "global-agent"}
	spec := &config.HTTPClientConfig{UserAgent: "spec-agent"}
	coll := &config.HTTPClientConfig{UserAgent: "coll-agent"}
	result := mergeHTTPClientConfig(global, spec, coll)
	// First-wins: global sets UserAgent first
	if result.UserAgent != "global-agent" {
		t.Errorf("UserAgent = %q, want %q", result.UserAgent, "global-agent")
	}
}

func TestMergeHTTPClientConfig_HeadersMerge(t *testing.T) {
	t.Parallel()

	global := &config.HTTPClientConfig{Headers: map[string]string{"X-Global": "g"}}
	spec := &config.HTTPClientConfig{Headers: map[string]string{"X-Spec": "s"}}
	coll := &config.HTTPClientConfig{Headers: map[string]string{"X-Coll": "c"}}
	result := mergeHTTPClientConfig(global, spec, coll)
	// First-wins: only global headers are set
	if result.Headers["X-Global"] != "g" {
		t.Errorf("X-Global = %q", result.Headers["X-Global"])
	}
	if result.Headers["X-Spec"] != "" {
		t.Errorf("X-Spec should be empty, got %q", result.Headers["X-Spec"])
	}
	if result.Headers["X-Coll"] != "" {
		t.Errorf("X-Coll should be empty, got %q", result.Headers["X-Coll"])
	}
}

func TestMergeHTTPClientConfig_Cookies(t *testing.T) {
	t.Parallel()

	global := &config.HTTPClientConfig{Cookies: []config.Cookie{{Name: "g", Value: "1"}}}
	spec := &config.HTTPClientConfig{Cookies: []config.Cookie{{Name: "s", Value: "2"}}}
	result := mergeHTTPClientConfig(global, spec, nil)
	// First-wins: global cookies are used
	if len(result.Cookies) != 1 {
		t.Fatalf("len = %d, want 1", len(result.Cookies))
	}
	if result.Cookies[0].Name != "g" {
		t.Errorf("Cookie name = %q, want %q", result.Cookies[0].Name, "g")
	}
}

func TestMergeHTTPClientConfig_FollowRedirects(t *testing.T) {
	t.Parallel()

	follow := false
	global := &config.HTTPClientConfig{FollowRedirects: &follow}
	result := mergeHTTPClientConfig(global, nil, nil)
	if result.FollowRedirects == nil || *result.FollowRedirects != false {
		t.Error("FollowRedirects not preserved")
	}
}

func TestMergeHTTPClientConfig_MaxRedirects(t *testing.T) {
	t.Parallel()

	maxRedirects := 5
	spec := &config.HTTPClientConfig{MaxRedirects: &maxRedirects}
	result := mergeHTTPClientConfig(nil, spec, nil)
	if result.MaxRedirects == nil || *result.MaxRedirects != 5 {
		t.Error("MaxRedirects not preserved")
	}
}

func TestMergeHTTPClientConfig_MaxResponseSize(t *testing.T) {
	t.Parallel()

	size := 4096
	coll := &config.HTTPClientConfig{MaxResponseSize: &size}
	result := mergeHTTPClientConfig(nil, nil, coll)
	if result.MaxResponseSize == nil || *result.MaxResponseSize != 4096 {
		t.Error("MaxResponseSize not preserved")
	}
}
