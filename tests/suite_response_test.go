package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func TestScript_ResponseSize_DefaultLimit(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		largeBody := make([]byte, 2000)
		for i := range largeBody {
			largeBody[i] = 'A'
		}
		_, _ = w.Write([]byte(`[{"data":"` + string(largeBody) + `"}]`))
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

	result := client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	var invokeResp struct {
		StatusCode  int             `json:"statusCode"`
		Body        json.RawMessage `json:"body"`
		FileRef     *struct {
			Path    string `json:"path"`
			Size    int    `json:"size"`
			Message string `json:"message"`
		} `json:"fileReference,omitempty"`
	}
	if err := json.Unmarshal(result, &invokeResp); err != nil {
		t.Fatalf("parse invoke: %v", err)
	}

	if invokeResp.FileRef != nil {
		if _, err := os.Stat(invokeResp.FileRef.Path); os.IsNotExist(err) {
			t.Errorf("file reference points to non-existent file: %s", invokeResp.FileRef.Path)
		}
	}
}

func TestScript_ResponseSize_Configurable(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1,"name":"Fluffy"}]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `http_client:
  max_response_size: 300
specs:
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

	result := client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	var invokeResp struct {
		StatusCode int             `json:"statusCode"`
		Body       json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(result, &invokeResp); err != nil {
		t.Fatalf("parse invoke: %v", err)
	}
	assertEqual(t, "statusCode", invokeResp.StatusCode, 200)
}

func TestScript_ResponseSize_FileReference(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		largeBody := make([]byte, 50000)
		for i := range largeBody {
			largeBody[i] = byte('A' + (i % 26))
		}
		_, _ = w.Write([]byte(`{"data":"` + string(largeBody) + `"}`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `http_client:
  max_response_size: 1000
specs:
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

	result := client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	var invokeResp struct {
		StatusCode int             `json:"statusCode"`
		Body        json.RawMessage `json:"body"`
		FileRef     *struct {
			Path    string `json:"path"`
			Size    int    `json:"size"`
			Message string `json:"message"`
		} `json:"fileRef,omitempty"`
	}
	if err := json.Unmarshal(result, &invokeResp); err != nil {
		t.Fatalf("parse invoke: %v", err)
	}

	if invokeResp.FileRef == nil {
		t.Fatal("expected fileReference for large response")
	}

	if _, err := os.Stat(invokeResp.FileRef.Path); os.IsNotExist(err) {
		t.Errorf("file reference points to non-existent file: %s", invokeResp.FileRef.Path)
	}
}
