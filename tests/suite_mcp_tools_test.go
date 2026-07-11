package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestScript_MCP_SpecList(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)
	result := client.listTools(t)

	var toolsResp struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(result, &toolsResp); err != nil {
		t.Fatalf("parse tools: %v", err)
	}

	toolNames := make(map[string]bool)
	for _, t := range toolsResp.Tools {
		toolNames[t.Name] = true
	}

	expectedTools := []string{
		"spec_list", "spec_by_id", "collection_by_spec", "collection_by_id",
		"tag_by_spec", "tag_by_collection", "tag_by_id",
		"endpoint_by_spec", "endpoint_by_collection", "endpoint_by_tag", "endpoint_by_id",
		"search", "inspect", "invoke", "auth",
	}
	for _, name := range expectedTools {
		if !toolNames[name] {
			t.Errorf("missing tool: %s", name)
		}
	}
}

func TestScript_MCP_SpecList_NoAuthTool(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth")
	client.initialize(t)
	result := client.listTools(t)

	var toolsResp struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(result, &toolsResp); err != nil {
		t.Fatalf("parse tools: %v", err)
	}

	for _, tool := range toolsResp.Tools {
		if tool.Name == "auth" {
			t.Errorf("auth tool should be disabled with --disable-llm-auth")
		}
	}
}

func TestScript_MCP_SpecList_Empty(t *testing.T) {
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
	result := client.callTool(t, "spec_list", map[string]interface{}{})

	var specsResp struct {
		Specs []struct {
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(result, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v\nstderr: %s", err, client.stderr.String())
	}
	if len(specsResp.Specs) == 0 {
		t.Fatalf("expected at least 1 spec\nstderr: %s", client.stderr.String())
	}
	assertEqual(t, "domain", specsResp.Specs[0].Domain, "test-api")
}

func TestScript_MCP_SpecByID(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	specID := specsResp.Specs[0].ID
	result := client.callTool(t, "spec_by_id", map[string]interface{}{
		"id": specID,
	})

	var specResp struct {
		Spec struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"spec"`
	}
	if err := json.Unmarshal(result, &specResp); err != nil {
		t.Fatalf("parse spec_by_id: %v", err)
	}
	assertEqual(t, "spec.id", specResp.Spec.ID, specID)
	assertEqual(t, "spec.domain", specResp.Spec.Domain, "petstore")
}

func TestScript_MCP_SpecByID_NotFound(t *testing.T) {
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

func TestScript_MCP_CollectionBySpec(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
      - title: Store
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

	result := client.callTool(t, "collection_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var collResp struct {
		Collections []struct {
			Title string `json:"title"`
		} `json:"collections"`
	}
	if err := json.Unmarshal(result, &collResp); err != nil {
		t.Fatalf("parse collection_by_spec: %v", err)
	}
	if len(collResp.Collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(collResp.Collections))
	}
}

func TestScript_MCP_CollectionByID(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
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

	collsResult := client.callTool(t, "collection_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	var collsResp struct {
		Collections []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"collections"`
	}
	if err := json.Unmarshal(collsResult, &collsResp); err != nil {
		t.Fatalf("parse collections: %v", err)
	}
	if len(collsResp.Collections) == 0 {
		t.Fatal("no collections found")
	}

	result := client.callTool(t, "collection_by_id", map[string]interface{}{
		"id": collsResp.Collections[0].ID,
	})

	var collResp struct {
		Collection struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"collection"`
	}
	if err := json.Unmarshal(result, &collResp); err != nil {
		t.Fatalf("parse collection_by_id: %v", err)
	}
	assertEqual(t, "collection.title", collResp.Collection.Title, "Petstore API")
}

func TestScript_MCP_TagBySpec(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
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

	result := client.callTool(t, "tag_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var tagsResp struct {
		Tags []struct {
			Name string `json:"name"`
		} `json:"tags"`
	}
	if err := json.Unmarshal(result, &tagsResp); err != nil {
		t.Fatalf("parse tag_by_spec: %v", err)
	}
	if len(tagsResp.Tags) == 0 {
		t.Errorf("expected at least 1 tag")
	}
}

func TestScript_MCP_EndpointBySpec(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
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

	result := client.callTool(t, "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})

	var endpointsResp struct {
		Endpoints []struct {
			Method string `json:"method"`
			Path string `json:"path"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &endpointsResp); err != nil {
		t.Fatalf("parse endpoint_by_spec: %v", err)
	}
	if len(endpointsResp.Endpoints) == 0 {
		t.Errorf("expected at least 1 endpoint")
	}
}

func TestScript_MCP_EndpointByID(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
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
	if len(epResp.Endpoints) == 0 {
		t.Fatal("no endpoints found")
	}

	result := client.callTool(t, "endpoint_by_id", map[string]interface{}{
		"id": epResp.Endpoints[0].ID,
	})

	var epByIDResp struct {
		Endpoint struct {
			ID   string `json:"id"`
			Method string `json:"method"`
			Path string `json:"path"`
		} `json:"endpoint"`
	}
	if err := json.Unmarshal(result, &epByIDResp); err != nil {
		t.Fatalf("parse endpoint_by_id: %v", err)
	}
	assertEqual(t, "endpoint.id", epByIDResp.Endpoint.ID, epResp.Endpoints[0].ID)
}

func TestScript_MCP_Inspect(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
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
			ID string `json:"id"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(epResult, &epResp); err != nil {
		t.Fatalf("parse endpoints: %v", err)
	}
	if len(epResp.Endpoints) == 0 {
		t.Fatal("no endpoints found")
	}

	result := client.callTool(t, "inspect", map[string]interface{}{
		"endpointId": epResp.Endpoints[0].ID,
	})

	var inspectResp struct {
		Operation interface{} `json:"operation"`
	}
	if err := json.Unmarshal(result, &inspectResp); err != nil {
		t.Fatalf("parse inspect: %v", err)
	}
	if inspectResp.Operation == nil {
		t.Errorf("expected operation to be present")
	}
}

func TestScript_MCP_Invoke(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1,"name":"Fluffy","status":"available"}]`))
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
	if len(epResp.Endpoints) == 0 {
		t.Fatal("no endpoints found")
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
		StatusCode int               `json:"statusCode"`
		Headers    map[string]string `json:"headers"`
		Body       json.RawMessage  `json:"body"`
	}
	if err := json.Unmarshal(result, &invokeResp); err != nil {
		t.Fatalf("parse invoke: %v", err)
	}
	assertEqual(t, "statusCode", invokeResp.StatusCode, 200)
}

func TestScript_MCP_Invoke_WithQueryParams(t *testing.T) {
	ws := newTestWorkspace(t)

	var capturedQuery string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
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
		"parameters": map[string]interface{}{
			"limit":  "5",
			"status": "available",
		},
	})

	if capturedQuery != "limit=5&status=available" && capturedQuery != "status=available&limit=5" {
		t.Errorf("unexpected query: %s", capturedQuery)
	}
}

func TestScript_MCP_Invoke_WithPathParams(t *testing.T) {
	ws := newTestWorkspace(t)

	var capturedPath string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets/42", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":42,"name":"Fluffy"}`))
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

	var getByIDEndpointID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "GET" && ep.Path == "/pets/{petId}" {
			getByIDEndpointID = ep.ID
			break
		}
	}
	if getByIDEndpointID == "" {
		t.Fatal("GET /pets/{petId} endpoint not found")
	}

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getByIDEndpointID,
		"parameters": map[string]interface{}{
			"petId": "42",
		},
	})

	assertEqual(t, "captured path", capturedPath, "/pets/42")
}

func TestScript_MCP_Invoke_WithRequestBody(t *testing.T) {
	ws := newTestWorkspace(t)

	var capturedBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buf)
		capturedBody = string(buf)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":1}`))
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

	var postEndpointID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "POST" && ep.Path == "/pets" {
			postEndpointID = ep.ID
			break
		}
	}
	if postEndpointID == "" {
		t.Fatal("POST /pets endpoint not found")
	}

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": postEndpointID,
		"requestBody": map[string]interface{}{
			"name":   "Fluffy",
			"status": "available",
		},
	})

	assertContains(t, "request body", capturedBody, "Fluffy")
}

func TestScript_MCP_Invoke_ServerError(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal error"}`))
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
		StatusCode int         `json:"statusCode"`
		Body       interface{} `json:"body"`
	}
	if err := json.Unmarshal(result, &invokeResp); err != nil {
		t.Fatalf("parse invoke: %v", err)
	}
	assertEqual(t, "statusCode", invokeResp.StatusCode, 500)
}

func TestScript_MCP_Auth(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: secured-api
    llm_title: Secured API
    base_url: https://api.example.com
    auth:
      type: bearer
      config:
        token: test-token-123
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
	assertEqual(t, "token", authResp.Token, "Bearer test-token-123")
}

func TestScript_MCP_TagFilter(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: public-api
    llm_title: Public API
    base_url: https://api.example.com
    tags: ["public"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: internal-api
    llm_title: Internal API
    base_url: https://api.example.com
    tags: ["internal"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "-t", "public")
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
	if len(specsResp.Specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specsResp.Specs))
	}
	assertEqual(t, "domain", specsResp.Specs[0].Domain, "public-api")
}
