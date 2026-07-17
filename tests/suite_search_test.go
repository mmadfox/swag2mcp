package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SearchSuite struct {
	BaseSuite
}

func (s *SearchSuite) TestBasic() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "pet",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []struct {
			ID string `json:"id"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for 'pet'")
}

func (s *SearchSuite) TestByMethod() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "method:GET",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Method string `json:"method"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	for _, ep := range searchResp.Endpoints {
		s.Equal("GET", ep.Method, "expected only GET endpoints")
	}
}

func (s *SearchSuite) TestByTag() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "tag:pets",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			TagName string `json:"tagName"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for tag:pets")
	for _, ep := range searchResp.Endpoints {
		s.Equal("pets", ep.TagName, "expected only 'pets' tag")
	}
}

func (s *SearchSuite) TestByPath() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "/pets",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Path string `json:"path"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for /pets")
}

func (s *SearchSuite) TestBooleanAND() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "+method:GET +summary:pet",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Method string `json:"method"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	for _, ep := range searchResp.Endpoints {
		s.Equal("GET", ep.Method, "expected only GET endpoints with boolean AND")
	}
}

func (s *SearchSuite) TestEmptyResults() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "zzzzzznonexistent",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.Empty(searchResp.Endpoints, "expected 0 results for nonexistent query")
}

func (s *SearchSuite) TestWildcard() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "list*",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for wildcard 'list*'")
}

func (s *SearchSuite) TestAllEndpoints() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "*",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for '*' query")
}

func (s *SearchSuite) TestLimitBounds() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
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

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "pet",
		"limit": 1,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.LessOrEqual(len(searchResp.Endpoints), 1, "expected at most 1 result with limit=1")
}

func TestSearchSuite(t *testing.T) {
	suite.Run(t, new(SearchSuite))
}
