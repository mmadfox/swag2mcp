package tests

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

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
