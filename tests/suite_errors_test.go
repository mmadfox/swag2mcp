package tests

import (
	"encoding/json"
	"testing"
)

func TestScript_Errors_NotFound(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	_ = client.callTool(t, "spec_by_id", map[string]interface{}{
		"id": "00000000000000000000000000000000",
	})
}

func TestScript_Errors_InvalidID(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	_ = client.callTool(t, "spec_by_id", map[string]interface{}{
		"id": "not-a-valid-id",
	})
}

func TestScript_Errors_EmptyID(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	_ = client.callTool(t, "spec_by_id", map[string]interface{}{
		"id": "",
	})
}

func TestScript_Errors_InvokeConnectionRefused(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: nowhere
    llm_title: Nowhere API
    base_url: http://localhost:1
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

	_ = client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})
}

func TestScript_Errors_Timeout(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `http_client:
  timeout: 1s
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: http://localhost:1
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

	_ = client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})
}
