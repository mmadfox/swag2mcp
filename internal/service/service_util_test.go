package service

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/types"
)

func TestWorkspace(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	ws := svc.Workspace()
	if ws == nil {
		t.Fatal("Workspace() returned nil")
	}
}

func TestNewInvokeResponse_JSON(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	body := []byte(`{"key": "value", "number": 42}`)

	resp := newInvokeResponse(response, body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type = %q, want %q", resp.Headers["Content-Type"], "application/json")
	}

	parsed, ok := resp.Body.(map[string]any)
	if !ok {
		t.Fatalf("Body type = %T, want map[string]any", resp.Body)
	}
	if parsed["key"] != "value" {
		t.Errorf("Body.key = %v, want %v", parsed["key"], "value")
	}
	if parsed["number"] != float64(42) {
		t.Errorf("Body.number = %v, want %v", parsed["number"], float64(42))
	}
}

func TestNewInvokeResponse_RawString(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
	}
	body := []byte("plain text response")

	resp := newInvokeResponse(response, body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	raw, ok := resp.Body.(string)
	if !ok {
		t.Fatalf("Body type = %T, want string", resp.Body)
	}
	if raw != "plain text response" {
		t.Errorf("Body = %q, want %q", raw, "plain text response")
	}
}

func TestNewInvokeResponse_EmptyBody(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusNoContent,
		Header:     http.Header{},
	}
	body := []byte{}

	resp := newInvokeResponse(response, body)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
	if resp.Body != "" {
		t.Errorf("Body = %v, want empty string", resp.Body)
	}
}

func TestNewInvokeResponse_Headers(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": {"application/json"},
			"X-Custom":     {"value1", "value2"},
		},
	}
	body := []byte(`{}`)

	resp := newInvokeResponse(response, body)
	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type = %q, want %q", resp.Headers["Content-Type"], "application/json")
	}
	if resp.Headers["X-Custom"] != "value1, value2" {
		t.Errorf("X-Custom = %q, want %q", resp.Headers["X-Custom"], "value1, value2")
	}
}

func TestMergeHTTPClientConfigs_BothNil(t *testing.T) {
	t.Parallel()

	result := mergeHTTPClientConfigs(nil, nil)
	if result != nil {
		t.Errorf("got %v, want nil", result)
	}
}

func TestMergeHTTPClientConfigs_SpecOnly(t *testing.T) {
	t.Parallel()

	spec := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Spec": "spec-val"},
	}
	result := mergeHTTPClientConfigs(spec, nil)
	if result == nil {
		t.Fatal("result is nil")
	}
	if result.Headers["X-Spec"] != "spec-val" {
		t.Errorf("X-Spec = %q, want %q", result.Headers["X-Spec"], "spec-val")
	}
}

func TestMergeHTTPClientConfigs_CollectionOnly(t *testing.T) {
	t.Parallel()

	collection := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Coll": "coll-val"},
	}
	result := mergeHTTPClientConfigs(nil, collection)
	if result == nil {
		t.Fatal("result is nil")
	}
	if result.Headers["X-Coll"] != "coll-val" {
		t.Errorf("X-Coll = %q, want %q", result.Headers["X-Coll"], "coll-val")
	}
}

func TestMergeHTTPClientConfigs_CollectionOverridesSpec(t *testing.T) {
	t.Parallel()

	spec := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "spec-val"},
	}
	collection := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "coll-val"},
	}
	result := mergeHTTPClientConfigs(spec, collection)
	// Collection overrides spec
	if result.Headers["X-Header"] != "coll-val" {
		t.Errorf("X-Header = %q, want %q", result.Headers["X-Header"], "coll-val")
	}
}

func TestMergeHTTPClientConfigs_Headers(t *testing.T) {
	t.Parallel()

	spec := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Spec": "spec-val"},
	}
	collection := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Coll": "coll-val"},
	}
	result := mergeHTTPClientConfigs(spec, collection)
	if result.Headers["X-Spec"] != "spec-val" {
		t.Errorf("X-Spec = %q, want %q", result.Headers["X-Spec"], "spec-val")
	}
	if result.Headers["X-Coll"] != "coll-val" {
		t.Errorf("X-Coll = %q, want %q", result.Headers["X-Coll"], "coll-val")
	}
}

func TestMergeHTTPClientConfigs_CollectionHeadersOverrideSpec(t *testing.T) {
	t.Parallel()

	spec := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "spec-val"},
	}
	collection := &types.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "coll-val"},
	}
	result := mergeHTTPClientConfigs(spec, collection)
	if result.Headers["X-Header"] != "coll-val" {
		t.Errorf("X-Header = %q, want %q", result.Headers["X-Header"], "coll-val")
	}
}

func TestBootstrap_Success(t *testing.T) {
	tmpDir := t.TempDir()
	specPath, _ := filepath.Abs("./testdata/valid_v311_openapi.yaml")

	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: test-api\n    llm_title: Test API v1\n    base_url: https://api.example.com\n    collections:\n      - llm_title: Main\n        location: " + specPath + "\n")
	if err := os.WriteFile(configPath, configContent, 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	svc, svcErr := New()
	if svcErr != nil {
		t.Fatalf("New() = %v", svcErr)
	}

	if bootErr := svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilepath: configPath,
	}); bootErr != nil {
		t.Fatalf("Bootstrap() = %v", bootErr)
	}

	specs, specErr := svc.Specs(context.Background())
	if specErr != nil {
		t.Fatalf("Specs() = %v", specErr)
	}
	if len(specs.Specs) == 0 {
		t.Fatal("no specs after bootstrap")
	}
}

func TestBootstrap_ConfigNotFound(t *testing.T) {
	t.Parallel()

	svc, svcErr := New()
	if svcErr != nil {
		t.Fatalf("New() = %v", svcErr)
	}

	bootErr := svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilepath: "/nonexistent/config.yaml",
	})
	if bootErr == nil {
		t.Fatal("expected error for nonexistent config")
	}
}

func TestBootstrap_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	if wrErr := os.WriteFile(configPath, []byte("invalid: yaml: ["), 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}

	svc, svcErr := New()
	if svcErr != nil {
		t.Fatalf("New() = %v", svcErr)
	}

	bootErr := svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilepath: configPath,
	})
	if bootErr == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestBootstrap_WithTags(t *testing.T) {
	tmpDir := t.TempDir()
	specPath, _ := filepath.Abs("./testdata/valid_v311_openapi.yaml")

	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: public-api\n    llm_title: Public API v1\n    base_url: https://api.example.com\n    tags: [public]\n    collections:\n      - llm_title: Main\n        location: " + specPath + "\n  - domain: internal-api\n    llm_title: Internal API v1\n    base_url: https://internal.example.com\n    tags: [internal]\n    collections:\n      - llm_title: Internal\n        location: " + specPath + "\n")
	if wrErr := os.WriteFile(configPath, configContent, 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}

	svc, svcErr := New()
	if svcErr != nil {
		t.Fatalf("New() = %v", svcErr)
	}

	if bootErr := svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilepath: configPath,
		Tags:         []string{"public"},
	}); bootErr != nil {
		t.Fatalf("Bootstrap() = %v", bootErr)
	}

	specs, specErr := svc.Specs(context.Background())
	if specErr != nil {
		t.Fatalf("Specs() = %v", specErr)
	}
	if len(specs.Specs) != 1 {
		t.Fatalf("Specs = %d, want 1", len(specs.Specs))
	}
	if specs.Specs[0].Domain != "public-api" {
		t.Errorf("Domain = %q, want %q", specs.Specs[0].Domain, "public-api")
	}
}

func TestBootstrap_InvalidSpecFile(t *testing.T) {
	tmpDir := t.TempDir()
	invalidSpecPath := filepath.Join(tmpDir, "invalid.yaml")
	if wrErr := os.WriteFile(invalidSpecPath, []byte("invalid: yaml: ["), 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}

	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: test-api\n    llm_title: Test API v1\n    base_url: https://api.example.com\n    collections:\n      - llm_title: Main\n        location: " + invalidSpecPath + "\n")
	if wrErr := os.WriteFile(configPath, configContent, 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}

	svc, svcErr := New()
	if svcErr != nil {
		t.Fatalf("New() = %v", svcErr)
	}

	bootErr := svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilepath: configPath,
	})
	if bootErr == nil {
		t.Fatal("expected error for invalid spec file")
	}
}

func TestBootstrap_ConfigValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	// Valid YAML but missing required fields (no base_url)
	configContent := []byte("specs:\n  - domain: test-api\n    llm_title: Test\n    collections:\n      - llm_title: Main\n        location: /tmp/test.yaml\n")
	if wrErr := os.WriteFile(configPath, configContent, 0600); wrErr != nil {
		t.Fatalf("WriteFile() = %v", wrErr)
	}

	svc, svcErr := New()
	if svcErr != nil {
		t.Fatalf("New() = %v", svcErr)
	}

	bootErr := svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilepath: configPath,
	})
	if bootErr == nil {
		t.Fatal("expected error for config validation failure")
	}
}
