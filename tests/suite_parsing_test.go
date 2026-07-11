package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestScript_Parsing_OpenAPI300(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: oas300
    llm_title: OAS 3.0.0
    base_url: ` + srv.URL + `
    collections:
      - title: Users
        location: ./internal/service/testdata/valid_v300_openapi.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	result := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(result, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}
	assertEqual(t, "domain", specsResp.Specs[0].Domain, "oas300")
}

func TestScript_Parsing_OpenAPI311(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: oas311
    llm_title: OAS 3.1.1
    base_url: ` + srv.URL + `
    collections:
      - title: Orders
        location: ./internal/service/testdata/valid_v311_openapi.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	result := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(result, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}
	assertEqual(t, "domain", specsResp.Specs[0].Domain, "oas311")
}

func TestScript_Parsing_Swagger20(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: swagger20
    llm_title: Swagger 2.0
    base_url: ` + srv.URL + `
    collections:
      - title: Users
        location: ./internal/service/testdata/valid_v20_swagger.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	result := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(result, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}
	assertEqual(t, "domain", specsResp.Specs[0].Domain, "swagger20")
}

func TestScript_Parsing_InvalidSpec(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	configContent := `specs:
  - domain: invalid-spec
    llm_title: Invalid Spec
    base_url: https://api.example.com
    collections:
      - title: Bad
        location: ./tests/testdata/invalid.yaml
`
	writeConfig(t, ws, configContent)

	_, _, code := runCommandInWS(t, ws, "validate", ".")
	if code == 0 {
		t.Errorf("expected validation to fail with invalid spec")
	}
}

func TestScript_Parsing_EmptySpec(t *testing.T) {
	ws := newTestWorkspace(t)

	emptySpec := `openapi: 3.0.0
info:
  title: Empty
  version: 1.0.0
paths: {}
`
	writeSpec(t, ws, "empty.yaml", emptySpec)

	configContent := `specs:
  - domain: empty-spec
    llm_title: Empty Spec
    base_url: https://api.example.com
    collections:
      - title: Empty
        location: ./empty.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	result := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(result, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	epResult := client.callTool(t, "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	if len(epResult) == 0 {
		return
	}
	var epResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	if err := json.Unmarshal(epResult, &epResp); err != nil {
		t.Fatalf("parse endpoints: %v", err)
	}
	if len(epResp.Endpoints) != 0 {
		t.Errorf("expected 0 endpoints for empty spec, got %d", len(epResp.Endpoints))
	}
}
