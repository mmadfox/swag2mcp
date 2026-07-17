package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type InfoSuite struct {
	BaseSuite
}

func (s *InfoSuite) TestPrintsInfo() {
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

	stdout, _, code := s.RunCommandInWS("info", ".")
	s.Equal(0, code)
	s.NotEmpty(stdout, "expected info output")
}

func TestInfoSuite(t *testing.T) {
	suite.Run(t, new(InfoSuite))
}
