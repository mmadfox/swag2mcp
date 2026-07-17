package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParsingSuite struct {
	BaseSuite
}

func (s *ParsingSuite) TestOpenAPI300() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: oas300
    llm_title: OAS 3.0.0
    base_url: ` + srv.URL + `
    collections:
      - title: Users
        location: ./internal/service/testdata/valid_v300_openapi.yaml
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
	s.Require().NotEmpty(specsResp.Specs, "no specs found")
	s.Equal("oas300", specsResp.Specs[0].Domain)
}

func (s *ParsingSuite) TestOpenAPI311() {
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: oas311
    llm_title: OAS 3.1.1
    base_url: ` + srv.URL + `
    collections:
      - title: Orders
        location: ./internal/service/testdata/valid_v311_openapi.yaml
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
	s.Require().NotEmpty(specsResp.Specs, "no specs found")
	s.Equal("oas311", specsResp.Specs[0].Domain)
}

func (s *ParsingSuite) TestSwagger20() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `specs:
  - domain: swagger20
    llm_title: Swagger 2.0
    base_url: ` + srv.URL + `
    collections:
      - title: Users
        location: ./internal/service/testdata/valid_v20_swagger.yaml
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
	s.Require().NotEmpty(specsResp.Specs, "no specs found")
	s.Equal("swagger20", specsResp.Specs[0].Domain)
}

func (s *ParsingSuite) TestInvalidSpec() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: invalid-spec
    llm_title: Invalid Spec
    base_url: https://api.example.com
    collections:
      - title: Bad
        location: ./tests/testdata/invalid.yaml
`
	s.WriteConfig(configContent)

	_, _, code := s.RunCommandInWS("validate", ".")
	s.NotEqual(0, code, "expected validation to fail with invalid spec")
}

func (s *ParsingSuite) TestEmptySpec() {
	emptySpec := `openapi: 3.0.0
info:
  title: Empty
  version: 1.0.0
paths: {}
`
	s.WriteSpec("empty.yaml", emptySpec)

	configContent := `specs:
  - domain: empty-spec
    llm_title: Empty Spec
    base_url: https://api.example.com
    collections:
      - title: Empty
        location: ./empty.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	result := client.callTool(s.T(), "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"specs"`
	}
	s.Require().NoError(json.Unmarshal(result, &specsResp))
	s.Require().NotEmpty(specsResp.Specs, "no specs found")

	epResult := client.callTool(s.T(), "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	if len(epResult) == 0 {
		return
	}
	var epResp struct {
		Endpoints []interface{} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(epResult, &epResp))
	s.Empty(epResp.Endpoints, "expected 0 endpoints for empty spec")
}

func TestParsingSuite(t *testing.T) {
	suite.Run(t, new(ParsingSuite))
}
