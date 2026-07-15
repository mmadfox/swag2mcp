package tests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestScript_Auth_None(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: noauth-api
    llm_title: No Auth API
    base_url: https://api.example.com
    auth:
      type: none
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token   string            `json:"token"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}
	assertEqual(t, "token", authResp.Token, "")
}

func TestScript_Auth_Basic(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: basic-api
    llm_title: Basic Auth API
    base_url: https://api.example.com
    auth:
      type: basic
      config:
        username: testuser
        password: testpass
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token   string            `json:"token"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
	assertEqual(t, "Authorization header", authResp.Headers["Authorization"], expectedAuth)
}

func TestScript_Auth_Bearer(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: bearer-api
    llm_title: Bearer Auth API
    base_url: https://api.example.com
    auth:
      type: bearer
      config:
        token: my-bearer-token
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token   string            `json:"token"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}

	assertEqual(t, "token", authResp.Token, "Bearer my-bearer-token")
	assertEqual(t, "Authorization header", authResp.Headers["Authorization"], "Bearer my-bearer-token")
}

func TestScript_Auth_APIKey_Header(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: apikey-api
    llm_title: API Key Auth API
    base_url: https://api.example.com
    auth:
      type: api-key
      config:
        key: X-API-Key
        value: my-api-key-value
        in: header
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token       string            `json:"token"`
		Headers     map[string]string `json:"headers"`
		QueryParams map[string]string `json:"queryParams"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}

	assertEqual(t, "X-API-Key header", authResp.Headers["X-API-Key"], "my-api-key-value")
}

func TestScript_Auth_APIKey_Query(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: apikey-query-api
    llm_title: API Key Query Auth API
    base_url: https://api.example.com
    auth:
      type: api-key
      config:
        key: api_key
        value: query-key-value
        in: query
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token       string            `json:"token"`
		Headers     map[string]string `json:"headers"`
		QueryParams map[string]string `json:"queryParams"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}

	assertEqual(t, "api_key query param", authResp.QueryParams["api_key"], "query-key-value")
}

func TestScript_Auth_EnvVarResolution(t *testing.T) {
	ws := newTestWorkspace(t)
	t.Setenv("AUTH_TOKEN", "resolved-token-value")

	configContent := `specs:
  - domain: env-auth-api
    llm_title: Env Auth API
    base_url: https://api.example.com
    auth:
      type: bearer
      config:
        token: $(AUTH_TOKEN)
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token   string            `json:"token"`
		Headers map[string]string `json:"headers"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}

	assertEqual(t, "token", authResp.Token, "Bearer resolved-token-value")
}

func TestScript_Auth_InvokeWithBearer(t *testing.T) {
	ws := newTestWorkspace(t)

	var authHeader string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := fmt.Sprintf(`specs:
  - domain: invoke-auth-api
    llm_title: Invoke Auth API
    base_url: %s
    auth:
      type: bearer
      config:
        token: invoke-bearer-token
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`, srv.URL)
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	epResult := client.callTool(t, "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	var epResp struct {
		Endpoints []struct {
			ID   string `json:"id"`
			Method string `json:"method"`
			Path string `json:"path"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(epResult, &epResp); err != nil {
		t.Fatalf("parse endpoints: %v", err)
	}

	var getEndpointID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "GET" && ep.Path == "/pets" {
			getEndpointID = ep.ID
			break
		}
	}
	if getEndpointID == "" {
		t.Fatal("GET /pets endpoint not found")
	}

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	assertEqual(t, "Authorization header", authHeader, "Bearer invoke-bearer-token")
}

func TestScript_Auth_HMAC(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: hmac-api
    llm_title: HMAC Auth API
    base_url: https://api.example.com
    auth:
      type: hmac
      config:
        api_key: test-api-key
        secret_key: test-secret-key
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	result := client.callTool(t, "auth", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var authResp struct {
		Token       string            `json:"token"`
		Headers     map[string]string `json:"headers"`
		QueryParams map[string]string `json:"queryParams"`
	}
	if err := json.Unmarshal(result, &authResp); err != nil {
		t.Fatalf("parse auth: %v", err)
	}

	assertEqual(t, "X-MBX-APIKEY header", authResp.Headers["X-MBX-APIKEY"], "test-api-key")
	if authResp.QueryParams["signature"] == "" {
		t.Error("signature query param should not be empty")
	}
	if authResp.QueryParams["timestamp"] == "" {
		t.Error("timestamp query param should not be empty")
	}
}
