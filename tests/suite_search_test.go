package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestScript_Search_Basic(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "pet",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []struct {
			ID      string `json:"id"`
			Method  string `json:"method"`
			Path    string `json:"path"`
			Summary string `json:"summary"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) == 0 {
		t.Errorf("expected at least 1 result for 'pet'")
	}
}

func TestScript_Search_ByMethod(t *testing.T) {
	ws := newTestWorkspace(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "method:GET",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Method string `json:"method"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	for _, ep := range searchResp.Endpoints {
		if ep.Method != "GET" {
			t.Errorf("expected only GET endpoints, got %s", ep.Method)
		}
	}
}

func TestScript_Search_ByTag(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "tag:pets",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			TagName string `json:"tagName"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) == 0 {
		t.Errorf("expected at least 1 result for tag:pets")
	}
	for _, ep := range searchResp.Endpoints {
		if ep.TagName != "pets" {
			t.Errorf("expected only 'pets' tag, got %s", ep.TagName)
		}
	}
}

func TestScript_Search_ByPath(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "/pets",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Path string `json:"path"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) == 0 {
		t.Errorf("expected at least 1 result for /pets")
	}
}

func TestScript_Search_BooleanAND(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "+method:GET +summary:pet",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Method  string `json:"method"`
			Summary string `json:"summary"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	for _, ep := range searchResp.Endpoints {
		if ep.Method != "GET" {
			t.Errorf("expected only GET endpoints with boolean AND, got %s", ep.Method)
		}
	}
}

func TestScript_Search_EmptyResults(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "zzzzzznonexistent",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) != 0 {
		t.Errorf("expected 0 results for nonexistent query, got %d", len(searchResp.Endpoints))
	}
}

func TestScript_Search_Wildcard(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "list*",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) == 0 {
		t.Errorf("expected at least 1 result for wildcard 'list*'")
	}
}

func TestScript_Search_AllEndpoints(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "*",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) == 0 {
		t.Errorf("expected at least 1 result for '*' query")
	}
}

func TestScript_Search_LimitBounds(t *testing.T) {
	ws := newTestWorkspace(t)
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(t, "search", map[string]interface{}{
		"query": "pet",
		"limit": 1,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	if err := json.Unmarshal(result, &searchResp); err != nil {
		t.Fatalf("parse search: %v", err)
	}
	if len(searchResp.Endpoints) > 1 {
		t.Errorf("expected at most 1 result with limit=1, got %d", len(searchResp.Endpoints))
	}
}
