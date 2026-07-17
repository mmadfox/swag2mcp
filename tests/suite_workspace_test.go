package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type WorkspaceSuite struct {
	BaseSuite
}

func (s *WorkspaceSuite) TestDirectoryStructure() {
	s.InitWorkspace()

	expectedDirs := []string{"cache", "specs", "responses", "auth_scripts"}
	root := s.Workspace
	for _, d := range expectedDirs {
		info, err := os.Stat(filepath.Join(root, d))
		if s.NoError(err, "missing directory %s", d) {
			s.True(info.IsDir(), "%s is not a directory", d)
		}
	}

	configPath := filepath.Join(root, "swag2mcp.yaml")
	_, err := os.Stat(configPath)
	s.Require().NoError(err, "missing swag2mcp.yaml")
}

func (s *WorkspaceSuite) TestCleanRemovesCacheAndResponses() {
	s.InitWorkspace()

	root := s.Workspace
	cacheFile := filepath.Join(root, "cache", "test.cache")
	s.Require().NoError(os.MkdirAll(filepath.Dir(cacheFile), 0755))
	s.Require().NoError(os.WriteFile(cacheFile, []byte("data"), 0644))
	respFile := filepath.Join(root, "responses", "test.json")
	s.Require().NoError(os.MkdirAll(filepath.Dir(respFile), 0755))
	s.Require().NoError(os.WriteFile(respFile, []byte("{}"), 0644))
	specFile := filepath.Join(root, "specs", "test.yaml")
	s.Require().NoError(os.MkdirAll(filepath.Dir(specFile), 0755))
	s.Require().NoError(os.WriteFile(specFile, []byte("spec: test"), 0644))

	s.RunCommandInWS("clean", ".")

	_, err := os.Stat(cacheFile)
	s.True(os.IsNotExist(err), "cache file was not removed")
	_, err = os.Stat(respFile)
	s.True(os.IsNotExist(err), "response file was not removed")
	_, err = os.Stat(specFile)
	s.Require().NoError(err, "spec file was removed but should be preserved")
}

func (s *WorkspaceSuite) TestUpdateReCachesSpecs() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)

	stdout, stderr, code := s.RunCommandInWS("update", ".")
	s.Equal(0, code, "update failed:\nstdout: %s\nstderr: %s", stdout, stderr)
}

func (s *WorkspaceSuite) TestUpdateWithInvalidConfig() {
	s.InitWorkspace()

	configContent := `specs:
  - domain: "INVALID DOMAIN"
    llm_title: Bad
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./nonexistent.yaml
`
	s.WriteConfig(configContent)

	_, _, code := s.RunCommandInWS("update", ".")
	s.NotEqual(0, code, "expected update to fail with invalid config")
}

func (s *WorkspaceSuite) TestOldResponsesCleaned() {
	s.InitWorkspace()

	root := s.Workspace
	oldResp := filepath.Join(root, "responses", "old.json")
	s.Require().NoError(os.MkdirAll(filepath.Dir(oldResp), 0755))
	s.Require().NoError(os.WriteFile(oldResp, []byte("old"), 0644))

	past := time.Now().Add(-49 * time.Hour)
	s.Require().NoError(os.Chtimes(oldResp, past, past))

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)

	s.RunCommandInWS("mcp", ".", "--auth-token", "test")

	_, err := os.Stat(oldResp)
	s.True(os.IsNotExist(err), "old response file (>48h) was not cleaned on mcp start")
}

func (s *WorkspaceSuite) TestRecentResponsesPreserved() {
	s.InitWorkspace()

	root := s.Workspace
	recentResp := filepath.Join(root, "responses", "recent.json")
	s.Require().NoError(os.MkdirAll(filepath.Dir(recentResp), 0755))
	s.Require().NoError(os.WriteFile(recentResp, []byte("recent"), 0644))

	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)

	s.RunCommandInWS("mcp", ".", "--auth-token", "test")

	_, err := os.Stat(recentResp)
	s.Require().NoError(err, "recent response file (<48h) was incorrectly cleaned")
}

func TestWorkspaceSuite(t *testing.T) {
	suite.Run(t, new(WorkspaceSuite))
}
