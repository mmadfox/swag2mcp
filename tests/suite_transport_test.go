package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TransportSuite struct {
	BaseSuite
}

func (s *TransportSuite) TestStdio() {
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
		Tools []interface{} `json:"tools"`
	}
	s.Require().NoError(json.Unmarshal(result, &toolsResp))
	s.NotEmpty(toolsResp.Tools, "expected tools from stdio transport")
}

func (s *TransportSuite) TestSSE() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
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
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
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

func (s *TransportSuite) TestDumpDir() {
	dumpDir := filepath.Join(s.Workspace, "dumps")
	s.Require().NoError(os.MkdirAll(dumpDir, 0755))

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
	client := s.StartMCPStdio(configContent, "--disable-llm-auth=false", "--dump-dir", dumpDir)
	client.initialize(s.T())

	specID := s.GetSpecID(client)
	endpointID := s.GetEndpointID(client, specID, "GET", "/pets")

	client.callTool(s.T(), "invoke", map[string]interface{}{
		"endpointId": endpointID,
	})

	entries, _ := os.ReadDir(dumpDir)
	s.NotEmpty(entries, "expected dump files in %s", dumpDir)
}

func TestTransportSuite(t *testing.T) {
	suite.Run(t, new(TransportSuite))
}
