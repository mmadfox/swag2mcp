package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExportSuite struct {
	BaseSuite
}

func (s *ExportSuite) specServer() (*httptest.Server, string) {
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

func (s *ExportSuite) TestSuccess() {
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

	outputPath := filepath.Join(s.Workspace, "backup.zip")
	stdout, stderr, code := s.RunCommand("export", s.Workspace, outputPath)
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Exported")

	_, err := os.Stat(outputPath)
	s.Require().NoError(err, "zip file not created at %s", outputPath)
}

func (s *ExportSuite) TestNoConfig() {
	_, _, code := s.RunCommand("export", s.Workspace)
	s.NotEqual(0, code)
}

func (s *ExportSuite) TestDefaultOutputPath() {
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

	stdout, stderr, code := s.RunCommandInWS("export", ".")
	s.Equal(0, code)
	s.Contains(stdout+stderr, ".zip")
}

func (s *ExportSuite) TestWithPathOnly() {
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

	stdout, stderr, code := s.RunCommand("export", s.Workspace)
	s.Equal(0, code)
	s.Contains(stdout+stderr, ".zip")
}

func (s *ExportSuite) TestWithSpecFilter() {
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

	outputPath := filepath.Join(s.Workspace, "backup.zip")
	stdout, stderr, code := s.RunCommand("export", s.Workspace, outputPath, "--spec", "meteo")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Exported")
}

func (s *ExportSuite) TestWithMultipleSpecFilter() {
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

	outputPath := filepath.Join(s.Workspace, "backup.zip")
	stdout, stderr, code := s.RunCommand("export", s.Workspace, outputPath, "--spec", "meteo,store")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Exported")
}

func (s *ExportSuite) TestWithPathOutputAndSpec() {
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

	outputPath := filepath.Join(s.Workspace, "backup.zip")
	stdout, stderr, code := s.RunCommand("export", s.Workspace, outputPath, "--spec", "meteo")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "Exported")
}

func (s *ExportSuite) TestWithSpecFilterNoMatch() {
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

	outputPath := filepath.Join(s.Workspace, "backup.zip")
	_, _, code := s.RunCommand("export", s.Workspace, outputPath, "--spec", "nonexistent")
	s.NotEqual(0, code)
}

func (s *ExportSuite) TestProducesValidSwag2mcpZip() {
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

	outputPath := filepath.Join(s.Workspace, "backup.zip")
	s.RunCommand("export", s.Workspace, outputPath)

	restoreWS := s.newTestWorkspace()
	stdout, stderr, code := s.RunCommand("import", restoreWS, "--from-zip", outputPath)
	s.Equal(0, code, "import exit code")
	s.Contains(stdout+stderr, "Restored")
}

func TestExportSuite(t *testing.T) {
	suite.Run(t, new(ExportSuite))
}
