package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type VersionSuite struct {
	BaseSuite
}

func (s *VersionSuite) TestPrintsVersion() {
	stdout, _, code := s.RunCommand("version")
	s.Equal(0, code)
	s.NotEmpty(stdout, "expected version output")
}

func TestVersionSuite(t *testing.T) {
	suite.Run(t, new(VersionSuite))
}
