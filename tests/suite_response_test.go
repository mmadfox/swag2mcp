/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseSuite struct {
	BaseSuite
}

func (s *ResponseSuite) TestDefaultLimit() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		largeBody := make([]byte, 2000)
		for i := range largeBody {
			largeBody[i] = 'A'
		}
		_, _ = w.Write([]byte(`[{"data":"` + string(largeBody) + `"}]`))
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

	result := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	var invokeResp struct {
		StatusCode int             `json:"statusCode"`
		Body       json.RawMessage `json:"body"`
		FileRef    *struct {
			Path    string `json:"path"`
			Size    int    `json:"size"`
			Message string `json:"message"`
		} `json:"fileReference,omitempty"`
	}
	s.Require().NoError(json.Unmarshal(result, &invokeResp))

	if invokeResp.FileRef != nil {
		_, err := os.Stat(invokeResp.FileRef.Path)
		s.Require().NoError(err, "file reference points to non-existent file: %s", invokeResp.FileRef.Path)
	}
}

func (s *ResponseSuite) TestConfigurable() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":1,"name":"Fluffy"}]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `http_client:
  max_response_size: 300
specs:
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

	result := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	var invokeResp struct {
		StatusCode int `json:"statusCode"`
	}
	s.Require().NoError(json.Unmarshal(result, &invokeResp))
	s.Equal(200, invokeResp.StatusCode)
}

func (s *ResponseSuite) TestFileReference() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		largeBody := make([]byte, 50000)
		for i := range largeBody {
			largeBody[i] = byte('A' + (i % 26))
		}
		_, _ = w.Write([]byte(`{"data":"` + string(largeBody) + `"}`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `http_client:
  max_response_size: 1000
specs:
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

	result := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	var invokeResp struct {
		StatusCode int             `json:"statusCode"`
		Body       json.RawMessage `json:"body"`
		FileRef    *struct {
			Path    string `json:"path"`
			Size    int    `json:"size"`
			Message string `json:"message"`
		} `json:"fileRef,omitempty"`
	}
	s.Require().NoError(json.Unmarshal(result, &invokeResp))
	s.Require().NotNil(invokeResp.FileRef, "expected fileReference for large response")

	_, err := os.Stat(invokeResp.FileRef.Path)
	s.Require().NoError(err, "file reference points to non-existent file: %s", invokeResp.FileRef.Path)
}

func (s *ResponseSuite) TestResponseReaderTools() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		items := ""
		for i := 1; i <= 200; i++ {
			if items != "" {
				items += ","
			}
			items += fmt.Sprintf(`{"id":%d,"name":"Pet-%d","description":"A friendly companion"}`, i, i)
		}
		_, _ = w.Write([]byte(`{"status":"ok","pets":[` + items + `]}`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `http_client:
  max_response_size: 10240
specs:
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

	invokeResult := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	var invokeResp struct {
		StatusCode int `json:"statusCode"`
		FileRef    *struct {
			Path string `json:"path"`
			Size int    `json:"size"`
		} `json:"fileRef,omitempty"`
	}
	s.Require().NoError(json.Unmarshal(invokeResult, &invokeResp))
	s.Require().NotNil(invokeResp.FileRef, "expected fileRef for response above max_response_size")

	outlineResult := client.callTool(s.T(), "response_outline", map[string]interface{}{
		"path": invokeResp.FileRef.Path,
	})
	var outlineResp struct {
		Outline struct {
			Type      string   `json:"type"`
			Keys      []string `json:"keys"`
			ItemCount int      `json:"itemCount"`
		} `json:"outline"`
	}
	s.Require().NoError(json.Unmarshal(outlineResult, &outlineResp))
	s.Equal("object", outlineResp.Outline.Type)
	s.Contains(outlineResp.Outline.Keys, "pets")

	compressResult := client.callTool(s.T(), "response_compress", map[string]interface{}{
		"path":     invokeResp.FileRef.Path,
		"jsonPath": "pets",
		"mode":     "first_of_array",
	})
	var compressResp struct {
		Body    map[string]any `json:"body,omitempty"`
		FileRef *struct {
			Path string `json:"path"`
		} `json:"fileRef,omitempty"`
	}
	s.Require().NoError(json.Unmarshal(compressResult, &compressResp))
	s.Require().NotNil(compressResp.Body, "expected compressed body inline")
	compressed, ok := compressResp.Body["compressed"].(map[string]any)
	s.Require().True(ok)
	s.Equal("array", compressed["type"])

	sliceResult := client.callTool(s.T(), "response_slice", map[string]interface{}{
		"path":     invokeResp.FileRef.Path,
		"jsonPath": "pets.0",
	})
	var sliceResp struct {
		Slice struct {
			Context string         `json:"context"`
			Value   map[string]any `json:"value"`
		} `json:"slice"`
	}
	s.Require().NoError(json.Unmarshal(sliceResult, &sliceResp))
	s.Equal("object", sliceResp.Slice.Context)
	s.Equal(float64(1), sliceResp.Slice.Value["id"])
}

func (s *ResponseSuite) TestResponseFilterSearch() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		items := ""
		for i := 1; i <= 500; i++ {
			if items != "" {
				items += ","
			}
			items += fmt.Sprintf(`{"id":%d,"name":"Pet-%d","status":"%s"}`, i, i, map[bool]string{true: "active", false: "inactive"}[i%2 == 0])
		}
		_, _ = w.Write([]byte(`{"status":"ok","pets":[` + items + `]}`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `http_client:
  max_response_size: 10240
specs:
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

	invokeResult := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	var invokeResp struct {
		StatusCode int `json:"statusCode"`
		FileRef    *struct {
			Path string `json:"path"`
		} `json:"fileRef,omitempty"`
	}
	s.Require().NoError(json.Unmarshal(invokeResult, &invokeResp))
	s.Require().NotNil(invokeResp.FileRef, "expected fileRef for response above max_response_size")

	// Search for "inactive" — matches all odd-numbered items
	filterResult := client.callTool(s.T(), "response_filter", map[string]interface{}{
		"path":     invokeResp.FileRef.Path,
		"jsonPath": "pets",
		"search":   "inactive",
	})
	var filterResp struct {
		Page       int           `json:"page"`
		PageSize   int           `json:"pageSize"`
		Total      int           `json:"total"`
		TotalPages int           `json:"totalPages"`
		Items      []interface{} `json:"items"`
		Strategy   string        `json:"strategy"`
	}
	s.Require().NoError(json.Unmarshal(filterResult, &filterResp))
	s.Equal(250, filterResp.Total)
	s.Len(filterResp.Items, 10)

	// Filter by exact status match
	filterResult2 := client.callTool(s.T(), "response_filter", map[string]interface{}{
		"path":     invokeResp.FileRef.Path,
		"jsonPath": "pets",
		"filter":   "status = active",
	})
	s.Require().NoError(json.Unmarshal(filterResult2, &filterResp))
	s.Equal(250, filterResp.Total)
	s.Len(filterResp.Items, 10)

	// Pagination
	filterResult3 := client.callTool(s.T(), "response_filter", map[string]interface{}{
		"path":     invokeResp.FileRef.Path,
		"jsonPath": "pets",
		"page":     2,
		"pageSize": 5,
	})
	s.Require().NoError(json.Unmarshal(filterResult3, &filterResp))
	s.Equal(2, filterResp.Page)
	s.Equal(5, filterResp.PageSize)
	s.Equal(500, filterResp.Total)
	s.Equal(100, filterResp.TotalPages)
	s.Len(filterResp.Items, 5)
}

func (s *ResponseSuite) TestResponseFilterRootArray() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		items := ""
		for i := 1; i <= 500; i++ {
			if items != "" {
				items += ","
			}
			items += fmt.Sprintf(`{"id":%d,"symbol":"coin-%d","price":%d}`, i, i, i*100)
		}
		_, _ = w.Write([]byte(`[` + items + `]`))
	})
	srv := s.StartHTTPServer(mux)

	configContent := `http_client:
  max_response_size: 10240
specs:
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

	invokeResult := client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
		"parameters": map[string]interface{}{
			"latitude":  0.0,
			"longitude": 0.0,
		},
	})

	var invokeResp struct {
		StatusCode int `json:"statusCode"`
		FileRef    *struct {
			Path string `json:"path"`
		} `json:"fileRef,omitempty"`
	}
	s.Require().NoError(json.Unmarshal(invokeResult, &invokeResp))
	s.Require().NotNil(invokeResp.FileRef, "expected fileRef for response above max_response_size")

	// Search root array with empty jsonPath — use exact filter to avoid substring matches
	filterResult := client.callTool(s.T(), "response_filter", map[string]interface{}{
		"path":   invokeResp.FileRef.Path,
		"filter": "symbol = coin-50",
	})
	var filterResp struct {
		Page       int           `json:"page"`
		PageSize   int           `json:"pageSize"`
		Total      int           `json:"total"`
		TotalPages int           `json:"totalPages"`
		Items      []interface{} `json:"items"`
		Strategy   string        `json:"strategy"`
	}
	s.Require().NoError(json.Unmarshal(filterResult, &filterResp))
	s.Equal(1, filterResp.Total)
	s.Len(filterResp.Items, 1)

	// Filter root array by price
	filterResult2 := client.callTool(s.T(), "response_filter", map[string]interface{}{
		"path":   invokeResp.FileRef.Path,
		"filter": "price > 5000",
	})
	s.Require().NoError(json.Unmarshal(filterResult2, &filterResp))
	s.Equal(450, filterResp.Total)
}

func TestResponseSuite(t *testing.T) {
	suite.Run(t, new(ResponseSuite))
}
