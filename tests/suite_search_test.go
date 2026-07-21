package tests

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

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
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "forecast",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []struct {
			ID string `json:"id"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for 'forecast'")
}

func (s *SearchSuite) TestByMethod() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
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
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "tag:weather",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			TagName string `json:"tagName"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for tag:weather")
	for _, ep := range searchResp.Endpoints {
		s.Equal("weather", ep.TagName, "expected only 'weather' tag")
	}
}

func (s *SearchSuite) TestByPath() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "/v1/forecast",
		"limit": 50,
	})

	var searchResp struct {
		Endpoints []struct {
			Path string `json:"path"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for /v1/forecast")
}

func (s *SearchSuite) TestBooleanAND() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "+method:GET +summary:forecast",
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
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
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
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "open*",
		"limit": 10,
	})

	var searchResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(result, &searchResp))
	s.NotEmpty(searchResp.Endpoints, "expected at least 1 result for wildcard 'open*'")
}

func (s *SearchSuite) TestAllEndpoints() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
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
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: ` + srv.URL + `
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "search", map[string]interface{}{
		"query": "forecast",
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
