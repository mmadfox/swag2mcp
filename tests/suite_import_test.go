package tests

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ImportSuite struct {
	BaseSuite
}

func (s *ImportSuite) specServer() (*httptest.Server, string) {
	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(specContent))
	}))
	s.T().Cleanup(srv.Close)
	return srv, specContent
}

func (s *ImportSuite) TestSingleFile() {
	s.InitWorkspace()
	srv, _ := s.specServer()

	stdout, stderr, code := s.RunCommand("import", s.Workspace, srv.URL, "meteo.yaml")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "meteo.yaml")

	specPath := filepath.Join(s.Workspace, "specs", "meteo.yaml")
	_, err := os.Stat(specPath)
	s.Require().NoError(err, "spec file not created at %s", specPath)
}

func (s *ImportSuite) TestSingleFileDuplicate() {
	s.InitWorkspace()

	specDir := filepath.Join(s.Workspace, "specs")
	s.Require().NoError(os.WriteFile(filepath.Join(specDir, "meteo.yaml"), []byte("existing"), 0600))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("openapi: 3.0.0"))
	}))
	s.T().Cleanup(srv.Close)

	_, _, code := s.RunCommand("import", s.Workspace, srv.URL, "meteo.yaml")
	s.NotEqual(0, code)
}

func (s *ImportSuite) TestWithPath() {
	s.InitWorkspace()
	srv, _ := s.specServer()

	stdout, stderr, code := s.RunCommand("import", s.Workspace, srv.URL, "meteo.yaml")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "meteo.yaml")
}

func (s *ImportSuite) TestLocalFile() {
	s.InitWorkspace()

	specContent := `openapi: 3.0.0
info:
  title: Local Spec
  version: 1.0.0
paths:
  /items:
    get:
      operationId: listItems
      summary: List all items
`
	localPath := filepath.Join(s.Workspace, "local-spec.yaml")
	s.Require().NoError(os.WriteFile(localPath, []byte(specContent), 0600))

	stdout, stderr, code := s.RunCommand("import", s.Workspace, localPath, "local-spec.yaml")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "local-spec.yaml")

	specPath := filepath.Join(s.Workspace, "specs", "local-spec.yaml")
	_, err := os.Stat(specPath)
	s.Require().NoError(err, "spec file not created at %s", specPath)
}

func (s *ImportSuite) TestWithSpec() {
	s.InitWorkspace()
	srv, _ := s.specServer()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	s.WriteConfig(configContent)

	stdout, stderr, code := s.RunCommand("import", s.Workspace, "--spec", "meteo")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Imported")

	specPath := filepath.Join(s.Workspace, "specs")
	entries, err := os.ReadDir(specPath)
	s.Require().NoError(err)
	s.NotEmpty(entries, "no spec files were imported")
}

func (s *ImportSuite) TestWithPathAndSpec() {
	s.InitWorkspace()
	srv, _ := s.specServer()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
`
	s.WriteConfig(configContent)

	stdout, stderr, code := s.RunCommand("import", s.Workspace, "--spec", "meteo")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Imported")
}

func (s *ImportSuite) TestWithMultipleSpecs() {
	s.InitWorkspace()
	srv, _ := s.specServer()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: ` + srv.URL + `
  - domain: store
    llm_title: Store API
    base_url: https://api.store.com
    collections:
      - title: Products
        location: ` + srv.URL + `
`
	s.WriteConfig(configContent)

	stdout, stderr, code := s.RunCommand("import", s.Workspace, "--spec", "meteo,store")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Imported")
}

func (s *ImportSuite) TestWithSpecNoMatch() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: meteo
    llm_title: Open-Meteo API
    base_url: https://api.meteo.com
    collections:
      - title: Pets
        location: https://example.com/spec.yaml
`
	s.WriteConfig(configContent)

	_, _, code := s.RunCommand("import", s.Workspace, "--spec", "nonexistent")
	s.NotEqual(0, code)
}

func (s *ImportSuite) TestWithSpecNoConfig() {
	_, _, code := s.RunCommand("import", s.Workspace, "--spec", "meteo")
	s.NotEqual(0, code)
}

func (s *ImportSuite) TestMissingArgs() {
	_, _, code := s.RunCommand("import", s.Workspace)
	s.NotEqual(0, code)
}

func TestImportSuite(t *testing.T) {
	suite.Run(t, new(ImportSuite))
}
