package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct {
	BaseSuite
}

func (s *ErrorsSuite) TestNotFound() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	_ = client.callTool(s.T(), "spec_by_id", map[string]interface{}{
		"id": "00000000000000000000000000000000",
	})
}

func (s *ErrorsSuite) TestInvalidID() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	_ = client.callTool(s.T(), "spec_by_id", map[string]interface{}{
		"id": "not-a-valid-id",
	})
}

func (s *ErrorsSuite) TestEmptyID() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	_ = client.callTool(s.T(), "spec_by_id", map[string]interface{}{
		"id": "",
	})
}

func (s *ErrorsSuite) TestInvokeConnectionRefused() {
	configContent := `specs:
  - domain: nowhere
    llm_title: Nowhere API
    base_url: http://localhost:1
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specsResult := client.callTool(s.T(), "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	s.Require().NoError(json.Unmarshal(specsResult, &specsResp))
	s.Require().NotEmpty(specsResp.Specs, "no specs found")

	epResult := client.callTool(s.T(), "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	var epResp struct {
		Endpoints []struct {
			ID     string `json:"id"`
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(epResult, &epResp))

	var getEndpointID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "GET" && ep.Path == "/v1/forecast" {
			getEndpointID = ep.ID
			break
		}
	}
	s.Require().NotEmpty(getEndpointID, "GET /v1/forecast endpoint not found")

	_ = client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})
}

func (s *ErrorsSuite) TestTimeout() {
	configContent := `http_client:
  timeout: 1s
specs:
  - domain: slow-api
    llm_title: Slow API
    base_url: http://localhost:1
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specsResult := client.callTool(s.T(), "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	s.Require().NoError(json.Unmarshal(specsResult, &specsResp))
	s.Require().NotEmpty(specsResp.Specs, "no specs found")

	epResult := client.callTool(s.T(), "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	var epResp struct {
		Endpoints []struct {
			ID     string `json:"id"`
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(epResult, &epResp))

	var getEndpointID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "GET" && ep.Path == "/v1/forecast" {
			getEndpointID = ep.ID
			break
		}
	}
	s.Require().NotEmpty(getEndpointID, "GET /v1/forecast endpoint not found")

	_ = client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}
