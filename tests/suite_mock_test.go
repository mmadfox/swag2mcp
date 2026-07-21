package tests

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bytes"
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

type MockSuite struct {
	BaseSuite
}

func (s *MockSuite) TestServerStartAndQuery() {
	s.InitWorkspace()

	specContent := `openapi: 3.0.0
info:
  title: Mock Test API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
      responses:
        '200':
          description: A list of pets
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                    name:
                      type: string
`
	specPath := filepath.Join(s.Workspace, "specs", "mock-test.yaml")
	s.Require().NoError(os.WriteFile(specPath, []byte(specContent), 0600))

	port := s.NextPort()
	configContent := fmt.Sprintf(`mock_enabled: true
specs:
  - domain: mock-test
    llm_title: Mock Test API
    base_url: http://localhost:9999
    collections:
      - title: Pets
        location: ./specs/mock-test.yaml
        base_mock_url: localhost:%d
`, port)
	configPath := filepath.Join(s.Workspace, "swag2mcp.yaml")
	s.Require().NoError(os.WriteFile(configPath, []byte(configContent), 0600))

	cmd := exec.Command(MockBinPath, s.Workspace)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/pets", port)
	resp, err := http.Get(url)
	s.Require().NoError(err, "mock server request failed\nstderr: %s", stderrBuf.String())
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var result interface{}
	s.Require().NoError(json.Unmarshal(body, &result), "response should be valid JSON")
}

func (s *MockSuite) TestMultipleSpecs() {
	s.InitWorkspace()

	specContent := `openapi: 3.0.0
info:
  title: Multi Spec API
  version: 1.0.0
paths:
  /items:
    get:
      operationId: listItems
      summary: List all items
      responses:
        '200':
          description: A list of items
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
`
	specPath := filepath.Join(s.Workspace, "specs", "multi.yaml")
	s.Require().NoError(os.WriteFile(specPath, []byte(specContent), 0600))

	port := s.NextPort()
	configContent := fmt.Sprintf(`mock_enabled: true
specs:
  - domain: multi-api
    llm_title: Multi Spec API
    base_url: http://localhost:9998
    collections:
      - title: Items
        location: ./specs/multi.yaml
        base_mock_url: localhost:%d
`, port)
	configPath := filepath.Join(s.Workspace, "swag2mcp.yaml")
	s.Require().NoError(os.WriteFile(configPath, []byte(configContent), 0600))

	cmd := exec.Command(MockBinPath, s.Workspace)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	s.Require().NoError(cmd.Start())
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/items", port)
	resp, err := http.Get(url)
	s.Require().NoError(err, "mock server request failed\nstderr: %s", stderrBuf.String())
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	var result interface{}
	s.Require().NoError(json.Unmarshal(body, &result), "response should be valid JSON")
}

func TestMockSuite(t *testing.T) {
	suite.Run(t, new(MockSuite))
}
