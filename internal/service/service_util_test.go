package service

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	t.Parallel()

	svc := newTestService(t)
	ws := svc.Workspace()
	require.NotNil(t, ws)
}

func TestNewInvokeResponse_JSON(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	body := []byte(`{"key": "value", "number": 42}`)

	resp := newInvokeResponse(response, body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.Headers["Content-Type"])

	parsed, ok := resp.Body.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "value", parsed["key"])
	require.Equal(t, float64(42), parsed["number"])
}

func TestNewInvokeResponse_RawString(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{},
	}
	body := []byte("plain text response")

	resp := newInvokeResponse(response, body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	raw, ok := resp.Body.(string)
	require.True(t, ok)
	require.Equal(t, "plain text response", raw)
}

func TestNewInvokeResponse_EmptyBody(t *testing.T) {
	t.Parallel()

	response := &http.Response{
		StatusCode: http.StatusNoContent,
		Header:     http.Header{},
	}
	body := []byte{}

	resp := newInvokeResponse(response, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Equal(t, "", resp.Body)
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
	require.Equal(t, "application/json", resp.Headers["Content-Type"])
	require.Equal(t, "value1, value2", resp.Headers["X-Custom"])
}

func TestMergeHTTPClientConfigs_BothNil(t *testing.T) {
	t.Parallel()

	result := mergeHTTPClientConfigs(nil, nil)
	require.Nil(t, result)
}

func TestMergeHTTPClientConfigs_SpecOnly(t *testing.T) {
	t.Parallel()

	spec := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Spec": "spec-val"},
	}
	result := mergeHTTPClientConfigs(spec, nil)
	require.NotNil(t, result)
	require.Equal(t, "spec-val", result.Headers["X-Spec"])
}

func TestMergeHTTPClientConfigs_CollectionOnly(t *testing.T) {
	t.Parallel()

	collection := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Coll": "coll-val"},
	}
	result := mergeHTTPClientConfigs(nil, collection)
	require.NotNil(t, result)
	require.Equal(t, "coll-val", result.Headers["X-Coll"])
}

func TestMergeHTTPClientConfigs_CollectionOverridesSpec(t *testing.T) {
	t.Parallel()

	spec := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "spec-val"},
	}
	collection := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "coll-val"},
	}
	result := mergeHTTPClientConfigs(spec, collection)
	require.Equal(t, "coll-val", result.Headers["X-Header"])
}

func TestMergeHTTPClientConfigs_Headers(t *testing.T) {
	t.Parallel()

	spec := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Spec": "spec-val"},
	}
	collection := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Coll": "coll-val"},
	}
	result := mergeHTTPClientConfigs(spec, collection)
	require.Equal(t, "spec-val", result.Headers["X-Spec"])
	require.Equal(t, "coll-val", result.Headers["X-Coll"])
}

func TestMergeHTTPClientConfigs_CollectionHeadersOverrideSpec(t *testing.T) {
	t.Parallel()

	spec := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "spec-val"},
	}
	collection := &model.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "coll-val"},
	}
	result := mergeHTTPClientConfigs(spec, collection)
	require.Equal(t, "coll-val", result.Headers["X-Header"])
}

func TestMergeHTTPClientConfig_ConfigLevel(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Spec": "spec-val"},
		Cookies: []config.Cookie{{Name: "spec-cookie", Value: "spec-val"}},
	}
	collection := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Coll": "coll-val"},
		Cookies: []config.Cookie{{Name: "coll-cookie", Value: "coll-val"}},
	}

	result := mergeHTTPClientConfig(spec, collection)
	require.NotNil(t, result)
	require.Equal(t, "spec-val", result.Headers["X-Spec"])
	require.Empty(t, result.Headers["X-Coll"])
	require.Len(t, result.Cookies, 1)
	require.Equal(t, "spec-cookie", result.Cookies[0].Name)
}

func TestMergeHTTPClientConfig_SpecNil(t *testing.T) {
	t.Parallel()

	collection := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Coll": "val"},
		Cookies: []config.Cookie{{Name: "c", Value: "v"}},
	}

	result := mergeHTTPClientConfig(nil, collection)
	require.NotNil(t, result)
	require.Equal(t, "val", result.Headers["X-Coll"])
	require.Len(t, result.Cookies, 1)
	require.Equal(t, "c", result.Cookies[0].Name)
}

func TestMergeHTTPClientConfig_BothNil(t *testing.T) {
	t.Parallel()

	result := mergeHTTPClientConfig(nil, nil)
	require.NotNil(t, result)
}

func TestMergeHTTPClientConfig_CollectionOverridesSpec(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "spec-val"},
	}
	collection := &config.HTTPClientConfig{
		Headers: map[string]string{"X-Header": "coll-val"},
	}

	result := mergeHTTPClientConfig(spec, collection)
	require.Equal(t, "spec-val", result.Headers["X-Header"])
}

func TestMergeHTTPClientConfig_NoHeaders(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{}
	collection := &config.HTTPClientConfig{}

	result := mergeHTTPClientConfig(spec, collection)
	require.Nil(t, result.Headers)
}

func TestMergeHTTPClientConfig_NoCookies(t *testing.T) {
	t.Parallel()

	spec := &config.HTTPClientConfig{}
	collection := &config.HTTPClientConfig{}

	result := mergeHTTPClientConfig(spec, collection)
	require.Nil(t, result.Cookies)
}

func TestBootstrap_Success(t *testing.T) {
	tmpDir := t.TempDir()
	specPath, _ := filepath.Abs("./testdata/valid_v311_openapi.yaml")

	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: test-api\n    llm_title: Test API v1\n    base_url: https://api.example.com\n    collections:\n      - llm_title: Main\n        location: " + specPath + "\n")
	err := os.WriteFile(configPath, configContent, 0600)
	require.NoError(t, err)

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: configPath,
	})
	require.NoError(t, err)

	specs, err := svc.Specs(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, specs.Specs)
}

func TestBootstrap_ConfigNotFound(t *testing.T) {
	t.Parallel()

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: "/nonexistent/config.yaml",
	})
	require.Error(t, err)
}

func TestBootstrap_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	err := os.WriteFile(configPath, []byte("invalid: yaml: ["), 0600)
	require.NoError(t, err)

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: configPath,
	})
	require.Error(t, err)
}

func TestBootstrap_WithTags(t *testing.T) {
	tmpDir := t.TempDir()
	specPath, _ := filepath.Abs("./testdata/valid_v311_openapi.yaml")

	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: public-api\n    llm_title: Public API v1\n    base_url: https://api.example.com\n    tags: [public]\n    collections:\n      - llm_title: Main\n        location: " + specPath + "\n  - domain: internal-api\n    llm_title: Internal API v1\n    base_url: https://internal.example.com\n    tags: [internal]\n    collections:\n      - llm_title: Internal\n        location: " + specPath + "\n")
	err := os.WriteFile(configPath, configContent, 0600)
	require.NoError(t, err)

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: configPath,
		Tags:         []string{"public"},
	})
	require.NoError(t, err)

	specs, err := svc.Specs(context.Background())
	require.NoError(t, err)
	require.Len(t, specs.Specs, 1)
	require.Equal(t, "public-api", specs.Specs[0].Domain)
}

func TestBootstrap_InvalidSpecFile(t *testing.T) {
	tmpDir := t.TempDir()
	invalidSpecPath := filepath.Join(tmpDir, "invalid.yaml")
	err := os.WriteFile(invalidSpecPath, []byte("invalid: yaml: ["), 0600)
	require.NoError(t, err)

	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: test-api\n    llm_title: Test API v1\n    base_url: https://api.example.com\n    collections:\n      - llm_title: Main\n        location: " + invalidSpecPath + "\n")
	err = os.WriteFile(configPath, configContent, 0600)
	require.NoError(t, err)

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: configPath,
	})
	require.Error(t, err)
}

func TestBootstrap_ConfigValidationError(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "swag2mcp.yaml")
	configContent := []byte("specs:\n  - domain: test-api\n    llm_title: Test\n    collections:\n      - llm_title: Main\n        location: /tmp/test.yaml\n")
	err := os.WriteFile(configPath, configContent, 0600)
	require.NoError(t, err)

	svc, err := New()
	require.NoError(t, err)

	err = svc.Bootstrap(context.Background(), BootstrapRequest{
		ConfFilePath: configPath,
	})
	require.Error(t, err)
}
