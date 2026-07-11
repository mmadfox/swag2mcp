package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestScript_RateLimit_BlocksSecondCall(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
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

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})
}

func TestScript_RateLimit_RecoversAfterWait(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
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

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	time.Sleep(11 * time.Second)

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})
}

func TestScript_RateLimit_DifferentEndpoints(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})
	mux.HandleFunc("/store/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total":100}`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
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

	var petsID, inventoryID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "GET" && ep.Path == "/pets" {
			petsID = ep.ID
		}
		if ep.Method == "GET" && ep.Path == "/store/inventory" {
			inventoryID = ep.ID
		}
	}
	if petsID == "" || inventoryID == "" {
		t.Fatal("could not find both endpoints")
	}

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": petsID,
	})

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": inventoryID,
	})
}
