package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DeleteCollectionSuite struct {
	BaseSuite
}

func (s *DeleteCollectionSuite) TestByIndex() {
	s.WriteConfig("specs: []")
	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - llm_title: Forecast
    location: ./testdata/meteo.yaml
  - llm_title: Store
    location: ./testdata/meteo.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	stdout, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout, "Forecast")

	_, _, code := s.RunCommandWithStdinInWS("1\n1\ny\n", "delete", "collection", ".")
	s.Equal(0, code)
}

func (s *DeleteCollectionSuite) TestCancel() {
	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - llm_title: Forecast
    location: ./testdata/meteo.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	s.RunCommandWithStdinInWS("n\n", "delete", "collection", ".")

	stdout, _, _ := s.RunCommandInWS("ls", ".")
	s.Contains(stdout, "Forecast")
}

func TestDeleteCollectionSuite(t *testing.T) {
	suite.Run(t, new(DeleteCollectionSuite))
}
