/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type TransportSuite struct {
	BaseSuite
}

func (s *TransportSuite) TestStdio() {
	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/meteo.yaml
`
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false")
	client.initialize(s.T())
	result := client.listTools(s.T())

	var toolsResp struct {
		Tools []interface{} `json:"tools"`
	}
	s.Require().NoError(json.Unmarshal(result, &toolsResp))
	s.NotEmpty(toolsResp.Tools, "expected tools from stdio transport")
}

func (s *TransportSuite) TestSSE() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/meteo.yaml
`
	s.WriteConfig(configContent)

	port := s.NextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", s.Workspace,
		"--transport", "sse",
		"--http-addr", addr,
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)
	resp, err := http.Get(url)
	s.Require().NoError(err, "SSE request failed\nstderr: %s", stderrBuf.String())
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	s.Contains(string(body), "event:", "expected SSE event stream")
}

func (s *TransportSuite) TestStreamableHTTP() {
	s.T().Skip("needs HTTP server")
}

func (s *TransportSuite) TestAuthToken() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/meteo.yaml
`
	s.WriteConfig(configContent)

	port := s.NextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", s.Workspace,
		"--transport", "sse",
		"--http-addr", addr,
		"--auth-token", "my-secret-token",
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer my-secret-token")
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "auth request failed\nstderr: %s", stderrBuf.String())
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode, "expected 200 with valid token")

	resp2, err := http.Get(url)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusUnauthorized, resp2.StatusCode, "expected 401 without token")
}

func (s *TransportSuite) TestAuthJWT_JWKS() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/meteo.yaml
`
	s.WriteConfig(configContent)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)

	jwksMux := http.NewServeMux()
	jwksMux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwksMarshalFromKey(key))
	})
	jwksSrv := httptest.NewServer(jwksMux)
	defer jwksSrv.Close()

	port := s.NextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", s.Workspace,
		"--transport", "sse",
		"--http-addr", addr,
		"--auth-type", "jwks",
		"--auth-jwks-url", jwksSrv.URL+"/.well-known/jwks.json",
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString(key)
	s.Require().NoError(err)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "JWT auth request failed\nstderr: %s", stderrBuf.String())
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode, "expected 200 with valid JWT")

	resp2, err := http.Get(url)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusUnauthorized, resp2.StatusCode, "expected 401 without token")
}

func (s *TransportSuite) TestAuthJWT_Introspection() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/meteo.yaml
`
	s.WriteConfig(configContent)

	introMux := http.NewServeMux()
	introMux.HandleFunc("/introspect", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "my-client" || pass != "my-secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active": true,
			"sub":    "intro-user",
			"exp":    time.Now().Add(time.Hour).Unix(),
		})
	})
	introSrv := httptest.NewServer(introMux)
	defer introSrv.Close()

	port := s.NextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", s.Workspace,
		"--transport", "sse",
		"--http-addr", addr,
		"--auth-type", "introspection",
		"--auth-introspection-url", introSrv.URL+"/introspect",
		"--auth-client-id", "my-client",
		"--auth-client-secret", "my-secret",
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer any-valid-token")
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "introspection auth request failed\nstderr: %s", stderrBuf.String())
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode, "expected 200 with valid introspection token")

	resp2, err := http.Get(url)
	s.Require().NoError(err)
	defer resp2.Body.Close()
	s.Equal(http.StatusUnauthorized, resp2.StatusCode, "expected 401 without token")
}

func (s *TransportSuite) TestAuthJWT_InvalidToken() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/meteo.yaml
`
	s.WriteConfig(configContent)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)

	jwksMux := http.NewServeMux()
	jwksMux.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwksMarshalFromKey(key))
	})
	jwksSrv := httptest.NewServer(jwksMux)
	defer jwksSrv.Close()

	port := s.NextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", s.Workspace,
		"--transport", "sse",
		"--http-addr", addr,
		"--auth-type", "jwks",
		"--auth-jwks-url", jwksSrv.URL+"/.well-known/jwks.json",
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Equal(http.StatusUnauthorized, resp.StatusCode, "expected 401 with invalid JWT")
}

func jwksMarshalFromKey(key *rsa.PrivateKey) jwkset.JWKSMarshal {
	jwk, err := jwkset.NewJWKFromKey(&key.PublicKey, jwkset.JWKOptions{
		Metadata: jwkset.JWKMetadataOptions{
			KID: "test-kid",
			ALG: jwkset.AlgRS256,
		},
	})
	if err != nil {
		return jwkset.JWKSMarshal{}
	}

	storage := jwkset.NewMemoryStorage()
	if err := storage.KeyWrite(context.Background(), jwk); err != nil {
		return jwkset.JWKSMarshal{}
	}

	marshal, err := storage.Marshal(context.Background())
	if err != nil {
		return jwkset.JWKSMarshal{}
	}

	return marshal
}

func (s *TransportSuite) TestDumpDir() {
	dumpDir := filepath.Join(s.Workspace, "dumps")
	s.Require().NoError(os.MkdirAll(dumpDir, 0755))

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"latitude": 0.0}`))
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
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false", "--dump-dir", dumpDir)
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

	entries, _ := os.ReadDir(dumpDir)
	s.NotEmpty(entries, "expected dump files in %s", dumpDir)
}

func TestTransportSuite(t *testing.T) {
	suite.Run(t, new(TransportSuite))
}
