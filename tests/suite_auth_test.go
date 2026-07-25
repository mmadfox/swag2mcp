package tests

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	BaseSuite
}

func (s *AuthSuite) TestNone() {
	configContent := `specs:
  - domain: noauth-api
    llm_title: No Auth API
    base_url: https://api.example.com
    auth:
      type: none
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
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
	s.Empty(authResp.Token)
}

func (s *AuthSuite) TestBasic() {
	configContent := `specs:
  - domain: basic-api
    llm_title: Basic Auth API
    base_url: https://api.example.com
    auth:
      type: basic
      config:
        username: testuser
        password: testpass
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "auth", map[string]interface{}{
		"specId": specID,
	})

	var authResp struct {
		Headers map[string]string `json:"headers"`
	}
	s.Require().NoError(json.Unmarshal(result, &authResp))

	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass"))
	s.Equal(expectedAuth, authResp.Headers["Authorization"])
}

func (s *AuthSuite) TestBearer() {
	configContent := `specs:
  - domain: bearer-api
    llm_title: Bearer Auth API
    base_url: https://api.example.com
    auth:
      type: bearer
      config:
        token: my-bearer-token
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "auth", map[string]interface{}{
		"specId": specID,
	})

	var authResp struct {
		Token   string            `json:"token"`
		Headers map[string]string `json:"headers"`
	}
	s.Require().NoError(json.Unmarshal(result, &authResp))
	s.Equal("Bearer my-bearer-token", authResp.Token)
	s.Equal("Bearer my-bearer-token", authResp.Headers["Authorization"])
}

func (s *AuthSuite) TestAPIKeyHeader() {
	configContent := `specs:
  - domain: apikey-api
    llm_title: API Key Auth API
    base_url: https://api.example.com
    auth:
      type: api-key
      config:
        key: X-API-Key
        value: my-api-key-value
        in: header
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "auth", map[string]interface{}{
		"specId": specID,
	})

	var authResp struct {
		Headers map[string]string `json:"headers"`
	}
	s.Require().NoError(json.Unmarshal(result, &authResp))
	s.Equal("my-api-key-value", authResp.Headers["X-API-Key"])
}

func (s *AuthSuite) TestAPIKeyQuery() {
	configContent := `specs:
  - domain: apikey-query-api
    llm_title: API Key Query Auth API
    base_url: https://api.example.com
    auth:
      type: api-key
      config:
        key: api_key
        value: query-key-value
        in: query
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "auth", map[string]interface{}{
		"specId": specID,
	})

	var authResp struct {
		QueryParams map[string]string `json:"queryParams"`
	}
	s.Require().NoError(json.Unmarshal(result, &authResp))
	s.Equal("query-key-value", authResp.QueryParams["api_key"])
}

func (s *AuthSuite) TestEnvVarResolution() {
	s.T().Setenv("AUTH_TOKEN", "resolved-token-value")

	configContent := `specs:
  - domain: env-auth-api
    llm_title: Env Auth API
    base_url: https://api.example.com
    auth:
      type: bearer
      config:
        token: $(AUTH_TOKEN)
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
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
	s.Equal("Bearer resolved-token-value", authResp.Token)
}

func (s *AuthSuite) TestInvokeWithBearer() {
	var authHeader string
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := fmt.Sprintf(`specs:
  - domain: invoke-auth-api
    llm_title: Invoke Auth API
    base_url: %s
    auth:
      type: bearer
      config:
        token: invoke-bearer-token
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`, srv.URL)
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/v1/forecast")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	s.Equal("Bearer invoke-bearer-token", authHeader)
}

func (s *AuthSuite) TestHMAC() {
	configContent := `specs:
  - domain: hmac-api
    llm_title: HMAC Auth API
    base_url: https://api.example.com
    auth:
      type: hmac
      config:
        api_key: test-api-key
        secret_key: test-secret-key
    collections:
      - title: Forecast
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	result := client.callTool(s.T(), "auth", map[string]interface{}{
		"specId": specID,
	})

	var authResp struct {
		Headers     map[string]string `json:"headers"`
		QueryParams map[string]string `json:"queryParams"`
	}
	s.Require().NoError(json.Unmarshal(result, &authResp))
	s.Equal("test-api-key", authResp.Headers["X-MBX-APIKEY"])
	s.NotEmpty(authResp.QueryParams["signature"])
	s.NotEmpty(authResp.QueryParams["timestamp"])
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
