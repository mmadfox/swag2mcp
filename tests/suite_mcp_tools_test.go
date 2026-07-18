package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MCPToolsSuite struct {
	BaseSuite
}

func (s *MCPToolsSuite) TestSpecList() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	result := client.listTools(s.T())

	var toolsResp struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	s.Require().NoError(json.Unmarshal(result, &toolsResp))

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
		s.True(toolNames[name], "missing tool: %s", name)
	}
}

func (s *MCPToolsSuite) TestSpecListNoAuthTool() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth")
	client.initialize(s.T())
	result := client.listTools(s.T())

	var toolsResp struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	s.Require().NoError(json.Unmarshal(result, &toolsResp))

	for _, tool := range toolsResp.Tools {
		s.NotEqual("auth", tool.Name, "auth tool should be disabled with --disable-llm-auth")
	}
}

func (s *MCPToolsSuite) TestSpecListEmpty() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	result := client.callTool(s.T(), "spec_list", map[string]interface{}{})

	var specsResp struct {
		Specs []struct {
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	s.Require().NoError(json.Unmarshal(result, &specsResp))
	s.Require().NotEmpty(specsResp.Specs, "expected at least 1 spec")
	s.Equal("test-api", specsResp.Specs[0].Domain)
}

func (s *MCPToolsSuite) TestSpecByID() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "spec_by_id", map[string]interface{}{
		"id": specID,
	})

	var specResp struct {
		Spec struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"spec"`
	}
	s.Require().NoError(json.Unmarshal(result, &specResp))
	s.Equal(specID, specResp.Spec.ID)
	s.Equal("petstore", specResp.Spec.Domain)
}

func (s *MCPToolsSuite) TestSpecByIDNotFound() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	_ = client.callTool(s.T(), "spec_by_id", map[string]interface{}{
		"id": "00000000000000000000000000000000",
	})
}

func (s *MCPToolsSuite) TestCollectionBySpec() {
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
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "collection_by_spec", map[string]interface{}{
		"specId": specID,
	})

	var collResp struct {
		Collections []struct {
			Title string `json:"title"`
		} `json:"collections"`
	}
	s.Require().NoError(json.Unmarshal(result, &collResp))
	s.Len(collResp.Collections, 1)
}

func (s *MCPToolsSuite) TestCollectionByID() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	collsResult := client.callTool(s.T(), "collection_by_spec", map[string]interface{}{
		"specId": specID,
	})
	var collsResp struct {
		Collections []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"collections"`
	}
	s.Require().NoError(json.Unmarshal(collsResult, &collsResp))
	s.Require().NotEmpty(collsResp.Collections, "no collections found")

	result := client.callTool(s.T(), "collection_by_id", map[string]interface{}{
		"id": collsResp.Collections[0].ID,
	})

	var collResp struct {
		Collection struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"collection"`
	}
	s.Require().NoError(json.Unmarshal(result, &collResp))
	s.Equal("Petstore API", collResp.Collection.Title)
}

func (s *MCPToolsSuite) TestTagBySpec() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "tag_by_spec", map[string]interface{}{
		"specId": specID,
	})

	var tagsResp struct {
		Tags []struct {
			Name string `json:"name"`
		} `json:"tags"`
	}
	s.Require().NoError(json.Unmarshal(result, &tagsResp))
	s.NotEmpty(tagsResp.Tags, "expected at least 1 tag")
}

func (s *MCPToolsSuite) TestEndpointBySpec() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "endpoint_by_spec", map[string]interface{}{
		"specId": specID,
	})

	var endpointsResp struct {
		Endpoints []struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &endpointsResp))
	s.NotEmpty(endpointsResp.Endpoints, "expected at least 1 endpoint")
}

func (s *MCPToolsSuite) TestEndpointByID() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	result := client.callTool(s.T(), "endpoint_by_id", map[string]interface{}{
		"id": endpointID,
	})

	var epByIDResp struct {
		Endpoint struct {
			ID     string `json:"id"`
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"endpoint"`
	}
	s.Require().NoError(json.Unmarshal(result, &epByIDResp))
	s.Equal(endpointID, epByIDResp.Endpoint.ID)
}

func (s *MCPToolsSuite) TestInspect() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	result := client.callTool(s.T(), "inspect", map[string]interface{}{
		"endpointId": endpointID,
	})

	var inspectResp struct {
		Operation interface{} `json:"operation"`
	}
	s.Require().NoError(json.Unmarshal(result, &inspectResp))
	s.NotNil(inspectResp.Operation, "expected operation to be present")
}

func (s *MCPToolsSuite) TestInvoke() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1,"name":"Fluffy","status":"available"}]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	result := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})

	var invokeResp struct {
		StatusCode int `json:"statusCode"`
	}
	s.Require().NoError(json.Unmarshal(result, &invokeResp))
	s.Equal(200, invokeResp.StatusCode)
}

func (s *MCPToolsSuite) TestInvokeWithQueryParams() {
	var capturedQuery string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"limit":  "5",
			"status": "available",
		},
	})

	s.Contains(capturedQuery, "limit=5")
	s.Contains(capturedQuery, "status=available")
}

func (s *MCPToolsSuite) TestInvokeWithPathParams() {
	var capturedPath string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets/42", func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":42,"name":"Fluffy"}`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets/{petId}")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"petId": "42",
		},
	})

	s.Equal("/pets/42", capturedPath)
}

func (s *MCPToolsSuite) TestInvokeWithRequestBody() {
	var capturedBody string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(buf)
		capturedBody = string(buf)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":1}`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "POST", "/pets")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"requestBody": map[string]interface{}{
			"name":   "Fluffy",
			"status": "available",
		},
	})

	s.Contains(capturedBody, "Fluffy")
}

func (s *MCPToolsSuite) TestInvokeServerError() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal error"}`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	result := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})

	var invokeResp struct {
		StatusCode int `json:"statusCode"`
	}
	s.Require().NoError(json.Unmarshal(result, &invokeResp))
	s.Equal(500, invokeResp.StatusCode)
}

func (s *MCPToolsSuite) TestAuth() {
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
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "auth", map[string]interface{}{
		"specId": specID,
	})

	var authResp struct {
		Token string `json:"token"`
	}
	s.Require().NoError(json.Unmarshal(result, &authResp))
	s.Equal("Bearer test-token-123", authResp.Token)
}

func (s *MCPToolsSuite) TestTagFilter() {
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
	client := s.StartMCPStdio(configContent, "-t", "public")
	client.initialize(s.T())

	result := client.callTool(s.T(), "spec_list", map[string]interface{}{})

	var specsResp struct {
		Specs []struct {
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	s.Require().NoError(json.Unmarshal(result, &specsResp))
	s.Require().Len(specsResp.Specs, 1, "expected 1 spec")
	s.Equal("public-api", specsResp.Specs[0].Domain)
}

func (s *MCPToolsSuite) TestInvokeWithCustomHeaders() {
	var capturedHeaders map[string]string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = map[string]string{
			"X-Custom-Header": r.Header.Get("X-Custom-Header"),
			"X-Another":       r.Header.Get("X-Another"),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"headers": map[string]interface{}{
			"X-Custom-Header": "custom-value",
			"X-Another":       "another-value",
		},
	})

	s.Equal("custom-value", capturedHeaders["X-Custom-Header"], "X-Custom-Header should be passed through")
	s.Equal("another-value", capturedHeaders["X-Another"], "X-Another should be passed through")
}

func (s *MCPToolsSuite) TestInvokeWithCustomCookies() {
	var capturedCookies string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		capturedCookies = r.Header.Get("Cookie")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"cookies": map[string]interface{}{
			"session": "abc123",
			"theme":   "dark",
		},
	})

	s.Contains(capturedCookies, "session=abc123", "session cookie should be passed through")
	s.Contains(capturedCookies, "theme=dark", "theme cookie should be passed through")
}

func (s *MCPToolsSuite) TestInvokeWithGlobalHeaders() {
	var capturedAccept, capturedUA string
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		capturedAccept = r.Header.Get("Accept")
		capturedUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `http_client:
  headers:
    Accept: application/json
  user_agent: swag2mcp-test/1.0
specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})

	s.Equal("application/json", capturedAccept, "Accept should come from global http_client.headers")
	s.Equal("swag2mcp-test/1.0", capturedUA, "User-Agent should come from global http_client.user_agent")
}

func TestMCPToolsSuite(t *testing.T) {
	suite.Run(t, new(MCPToolsSuite))
}
