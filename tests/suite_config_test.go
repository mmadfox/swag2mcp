package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	BaseSuite
}

func (s *ConfigSuite) TestValidateValidConfig() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	stdout, _, code := s.RunCommandInWS("validate", ".")
	s.Equal(0, code)
	s.Contains(stdout, "valid")
}

func (s *ConfigSuite) TestValidateDuplicateDomain() {
	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: petstore
    llm_title: Duplicate
    base_url: https://api.example.com
    collections:
      - title: Store
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	_, stderr, code := s.RunCommandInWS("validate", ".")
	s.NotEqual(0, code)
	s.Contains(stderr, "duplicate")
}

func (s *ConfigSuite) TestValidateInvalidDomainFormat() {
	configContent := `specs:
  - domain: "UPPERCASE INVALID"
    llm_title: Bad API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	_, stderr, code := s.RunCommandInWS("validate", ".")
	s.NotEqual(0, code)
	s.Contains(stderr, "Domain")
}

func (s *ConfigSuite) TestValidateUnreachableLocation() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Missing
        location: ./nonexistent.yaml
`
	s.WriteConfig(configContent)
	_, _, code := s.RunCommandInWS("validate", ".")
	s.NotEqual(0, code)
}

func (s *ConfigSuite) TestValidateTagFilter() {
	configContent := `specs:
  - domain: public-api
    llm_title: Public API
    base_url: https://api.example.com
    tags: ["public"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: internal-api
    llm_title: Internal API
    base_url: https://api.example.com
    tags: ["internal"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	_, _, code := s.RunCommandInWS("validate", "-t", "public", ".")
	s.Equal(0, code)
}

func (s *ConfigSuite) TestAddSpecFromYAML() {
	yamlData := `domain: added-spec
llm_title: Added Spec
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	stdout, stderr, code := s.RunCommandInWS("add", "spec", "--yaml", yamlData, ".")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "added")

	stdout2, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout2, "added-spec")
}

func (s *ConfigSuite) TestAddSpecFromStdin() {
	yamlData := `domain: stdin-spec
llm_title: Stdin Spec
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	stdout, stderr, code := s.RunCommandWithStdinInWS(yamlData, "add", "spec", "--yaml", "-", ".")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "added")

	stdout2, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout2, "stdin-spec")
}

func (s *ConfigSuite) TestAddSpecInvalidYAML() {
	_, _, code := s.RunCommandInWS("add", "spec", "--yaml", "invalid: [yaml: broken", ".")
	s.NotEqual(0, code)
}

func (s *ConfigSuite) TestAddCollectionFromYAML() {
	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - title: Existing
    location: ./testdata/petstore.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	collectionYAML := `spec_domain: test-api
llm_title: Added Collection
location: ./testdata/petstore.yaml
`
	stdout, stderr, code := s.RunCommandInWS("add", "collection", "--yaml", collectionYAML, ".")
	s.Equal(0, code)
	s.Contains(stdout+stderr, "added")

	stdout2, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout2, "Added Collection")
}

func (s *ConfigSuite) TestDeleteSpec() {
	specYAML := `domain: to-delete
llm_title: To Delete
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	specYAML2 := `domain: keep-me
llm_title: Keep Me
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML2, ".")

	stdout, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout, "to-delete")
	s.Contains(stdout, "keep-me")

	_, _, code := s.RunCommandWithStdinInWS("1\ny\n", "delete", "spec", ".")
	s.Equal(0, code)

	stdout2, _, _ := s.RunCommandInWS("ls", ".")
	s.NotContains(stdout2, "to-delete")
	s.Contains(stdout2, "keep-me")
}

func (s *ConfigSuite) TestDeleteSpecCancel() {
	specYAML := `domain: keep-me
llm_title: Keep Me
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")
	s.RunCommandWithStdinInWS("n\n", "delete", "spec", ".")

	stdout, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout, "keep-me")
}

func (s *ConfigSuite) TestListSpecs() {
	specYAML := `domain: list-test
llm_title: List Test
base_url: https://api.example.com
collections:
  - title: Pets
    location: ./testdata/petstore.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	stdout, _, code := s.RunCommandInWS("ls", ".")
	s.Equal(0, code)
	s.Contains(stdout, "list-test")
	s.Contains(stdout, "List Test")
	s.Contains(stdout, "petstore.yaml")
}

func (s *ConfigSuite) TestListSpecsEmpty() {
	_, _, code := s.RunCommandInWS("ls", ".")
	s.Equal(0, code)
}

func (s *ConfigSuite) TestListSpecsTagFilter() {
	configContent := `specs:
  - domain: public-api
    llm_title: Public API
    base_url: https://api.example.com
    tags: ["public"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
  - domain: internal-api
    llm_title: Internal API
    base_url: https://api.example.com
    tags: ["internal"]
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	stdout, _, _ := s.RunCommandInWS("ls", "-t", "public", ".")
	s.Contains(stdout, "public-api")
	s.NotContains(stdout, "internal-api")
}

func (s *ConfigSuite) TestUpdateReCachesSpecs() {
	configContent := `specs:
  - domain: test-api
    llm_title: Test API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	_, stderr, code := s.RunCommandInWS("update", ".")
	s.Equal(0, code)
	s.Contains(stderr, "processed")
}

func (s *ConfigSuite) TestUpdateInvalidConfig() {
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
	s.NotEqual(0, code)
}

func (s *ConfigSuite) TestCleanRemovesCache() {
	root := s.Workspace
	cacheDir := filepath.Join(root, "cache")
	s.Require().NoError(os.MkdirAll(cacheDir, 0755))
	dummyFile := filepath.Join(cacheDir, "test.cache")
	s.Require().NoError(os.WriteFile(dummyFile, []byte("data"), 0644))

	responsesDir := filepath.Join(root, "responses")
	s.Require().NoError(os.MkdirAll(responsesDir, 0755))
	dummyResp := filepath.Join(responsesDir, "test.json")
	s.Require().NoError(os.WriteFile(dummyResp, []byte("{}"), 0644))

	stdout, _, code := s.RunCommandInWS("clean", ".")
	s.Equal(0, code)
	s.Contains(stdout, "Removed")

	_, err := os.Stat(dummyFile)
	s.True(os.IsNotExist(err), "cache file was not removed")
	_, err = os.Stat(dummyResp)
	s.True(os.IsNotExist(err), "response file was not removed")
}

func (s *ConfigSuite) TestCleanPreservesSpecs() {
	root := s.Workspace
	specsDir := filepath.Join(root, "specs")
	s.Require().NoError(os.MkdirAll(specsDir, 0755))
	specFile := filepath.Join(specsDir, "test.yaml")
	s.Require().NoError(os.WriteFile(specFile, []byte("spec: test"), 0644))

	s.RunCommandInWS("clean", ".")

	_, err := os.Stat(specFile)
	s.Require().NoError(err, "spec file was removed but should be preserved")
}

func (s *ConfigSuite) TestEnvVarResolution() {
	s.T().Setenv("TEST_BASE_URL", "https://env-test.example.com")

	configContent := `specs:
  - domain: env-test
    llm_title: Env Test
    base_url: https://env-test.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	s.WriteConfig(configContent)
	_, _, code := s.RunCommandInWS("validate", ".")
	s.Equal(0, code)
}

func (s *ConfigSuite) TestConfigCascade() {
	configContent := `http_client:
  timeout: 10s
  headers:
    X-Global: "true"
specs:
  - domain: cascade-test
    llm_title: Cascade Test
    base_url: https://spec.example.com
    http_client:
      timeout: 30s
      headers:
        X-Spec: "spec-only"
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
        http_client:
          headers:
            X-Collection: "collection-only"
`
	s.WriteConfig(configContent)
	_, _, code := s.RunCommandInWS("validate", ".")
	s.Equal(0, code)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
