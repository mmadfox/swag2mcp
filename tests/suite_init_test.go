package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InitSuite struct {
	BaseSuite
}

func (s *InitSuite) TestCreatesWorkspace() {
	_, stderr, code := s.RunCommand("init", s.Workspace)
	s.Equal(0, code)
	s.Contains(stderr, "initialized")

	root := s.Workspace
	dirs := []string{"cache", "specs", "responses", "auth_scripts"}
	for _, d := range dirs {
		info, err := os.Stat(filepath.Join(root, d))
		if s.NoError(err, "missing directory %s", d) {
			s.True(info.IsDir(), "%s is not a directory", d)
		}
	}

	configPath := filepath.Join(root, "swag2mcp.yaml")
	_, err := os.Stat(configPath)
	s.Require().NoError(err, "missing swag2mcp.yaml")
}

func (s *InitSuite) TestForceOverwrite() {
	s.RunCommand("init", s.Workspace)

	_, stderr, code := s.RunCommand("init", s.Workspace)
	s.NotEqual(0, code, "expected failure without -f")
	s.Contains(stderr, "not empty")

	_, _, code = s.RunCommand("init", "-f", s.Workspace)
	s.Equal(0, code, "expected success with -f")
}

func (s *InitSuite) TestInteractive() {
	s.T().Skip("requires TTY")
}

func (s *InitSuite) TestCustomPath() {
	customPath := filepath.Join(s.Workspace, "custom", "nested", "workspace")
	stdout, stderr, code := s.RunCommand("init", customPath)
	s.Equal(0, code)
	s.Contains(stdout+stderr, "initialized")

	configPath := filepath.Join(customPath, "swag2mcp.yaml")
	_, err := os.Stat(configPath)
	s.Require().NoError(err, "config not created at custom path: %s", configPath)
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(InitSuite))
}
