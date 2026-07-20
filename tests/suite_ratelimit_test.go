package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RateLimitSuite struct {
	BaseSuite
}

func (s *RateLimitSuite) TestBlocksSecondCall() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1}]`))
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

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/v1/forecast")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})
}

func (s *RateLimitSuite) TestRecoversAfterWait() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1}]`))
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

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/v1/forecast")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})

	time.Sleep(11 * time.Second)

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})
}

func (s *RateLimitSuite) TestDifferentEndpoints() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1}]`))
	})
	mux.HandleFunc("/store/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total":100}`))
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

	specID := s.GetSpecID(client)

	petsID := s.GetEndpointID(client, specID, "GET", "/v1/forecast")
	inventoryID := s.GetEndpointID(client, specID, "GET", "/store/inventory")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": petsID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": inventoryID,
	})
}

func TestRateLimitSuite(t *testing.T) {
	suite.Run(t, new(RateLimitSuite))
}
