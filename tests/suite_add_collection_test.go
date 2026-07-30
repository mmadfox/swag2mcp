/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AddCollectionSuite struct {
	BaseSuite
}

func (s *AddCollectionSuite) TestAddCollectionYAML_DuplicateLocation() {
	s.WriteConfig("specs: []")
	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - llm_title: Forecast
    location: ./testdata/meteo.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	dupYAML := `spec_domain: test-api
llm_title: Forecast Duplicate
location: ./testdata/meteo.yaml
`
	stdout, stderr, code := s.RunCommandInWS("add", "collection", "--yaml", dupYAML, ".")
	s.NotEqual(0, code, "expected non-zero exit code for duplicate location")
	s.Contains(stdout+stderr, "already exists")
}

func (s *AddCollectionSuite) TestAddCollectionYAML_SpecNotFound() {
	s.WriteConfig("specs: []")
	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - llm_title: Forecast
    location: ./testdata/meteo.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	badYAML := `spec_domain: nonexistent
llm_title: New Collection
location: ./testdata/meteo.yaml
`
	stdout, stderr, code := s.RunCommandInWS("add", "collection", "--yaml", badYAML, ".")
	s.NotEqual(0, code, "expected non-zero exit code for nonexistent spec")
	s.Contains(stdout+stderr, "not found")
}

func (s *AddCollectionSuite) TestAddCollectionYAML_Success() {
	s.WriteConfig("specs: []")
	specYAML := `domain: test-api
llm_title: Test API
base_url: https://api.example.com
collections:
  - llm_title: Forecast
    location: ./testdata/meteo.yaml
`
	s.RunCommandInWS("add", "spec", "--yaml", specYAML, ".")

	newYAML := `spec_domain: test-api
llm_title: New Collection
location: ./testdata/duplicate_domain.yaml
`
	stdout, stderr, code := s.RunCommandInWS("add", "collection", "--yaml", newYAML, ".")
	s.Equal(0, code, "stdout: %s, stderr: %s", stdout, stderr)
	s.Contains(stdout+stderr, "New Collection")
}

func TestAddCollectionSuite(t *testing.T) {
	suite.Run(t, new(AddCollectionSuite))
}
