package tests

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	basePort  = 9000
	portRange = 1000
)

const specIDKey = "specId"

var portCounter int32

type BaseSuite struct {
	suite.Suite

	Workspace string
}

func (s *BaseSuite) SetupTest() {
	s.Workspace = s.newTestWorkspace()
}

func (s *BaseSuite) TearDownTest() {
	matches, _ := filepath.Glob(filepath.Join(projectRoot, "swag2mcp-backup-*.zip"))
	for _, m := range matches {
		_ = os.Remove(m)
	}
}

func (s *BaseSuite) newTestWorkspace() string {
	t := s.T()
	t.Helper()
	ws, err := os.MkdirTemp("", "swag2mcp-ws-*")
	s.Require().NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(ws) })
	return ws
}

func (s *BaseSuite) WriteConfig(content string) {
	t := s.T()
	t.Helper()

	src := filepath.Join(projectRoot, "tests", "testdata")
	dst := filepath.Join(s.Workspace, "testdata")
	s.Require().NoError(copyDir(src, dst))

	src = filepath.Join(projectRoot, "internal", "service", "testdata")
	dst = filepath.Join(s.Workspace, "internal", "service", "testdata")
	s.Require().NoError(copyDir(src, dst))

	path := filepath.Join(s.Workspace, "swag2mcp.yaml")
	s.Require().NoError(os.WriteFile(path, []byte(content), 0600))
}

func (s *BaseSuite) WriteSpec(name, content string) {
	t := s.T()
	t.Helper()
	path := filepath.Join(s.Workspace, name)
	s.Require().NoError(os.WriteFile(path, []byte(content), 0600))
}

func (s *BaseSuite) InitWorkspace() {
	t := s.T()
	t.Helper()
	stdout, stderr, code := s.RunCommand("init", s.Workspace)
	s.Require().Equalf(0, code, "init failed:\nstdout: %s\nstderr: %s", stdout, stderr)
}

func (s *BaseSuite) runCommand(cmd *exec.Cmd) (stdout, stderr string, exitCode int) {
	t := s.T()
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			s.Require().Fail("failed to run command", err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func (s *BaseSuite) RunCommand(args ...string) (stdout, stderr string, exitCode int) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = projectRoot
	return s.runCommand(cmd)
}

func (s *BaseSuite) RunCommandInWS(args ...string) (stdout, stderr string, exitCode int) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = s.Workspace
	return s.runCommand(cmd)
}

func (s *BaseSuite) RunCommandWithStdin(stdin string, args ...string) (stdout, stderr string, exitCode int) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = projectRoot
	cmd.Stdin = strings.NewReader(stdin)
	return s.runCommand(cmd)
}

func (s *BaseSuite) RunCommandWithStdinInWS(stdin string, args ...string) (stdout, stderr string, exitCode int) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = s.Workspace
	cmd.Stdin = strings.NewReader(stdin)
	return s.runCommand(cmd)
}

func (s *BaseSuite) RunCommandWithEnv(env []string, args ...string) (stdout, stderr string, exitCode int) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), env...)
	return s.runCommand(cmd)
}

func (s *BaseSuite) StartMCPStdio(configContent string, extraArgs ...string) *mcpClient {
	t := s.T()
	t.Helper()
	s.WriteConfig(configContent)

	args := []string{"mcp", "."}
	args = append(args, extraArgs...)
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = s.Workspace

	stdinRead, stdinWrite, err := os.Pipe()
	s.Require().NoError(err)
	stdoutRead, stdoutWrite, err := os.Pipe()
	s.Require().NoError(err)
	cmd.Stdin = stdinRead
	cmd.Stdout = stdoutWrite

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	s.Require().NoError(cmd.Start())

	client := &mcpClient{
		cmd:    cmd,
		stdin:  stdinWrite,
		stdout: stdoutRead,
		stderr: &stderrBuf,
	}

	t.Cleanup(func() {
		_ = stdinWrite.Close()
		_ = stdoutRead.Close()
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	})

	return client
}

func (s *BaseSuite) StartHTTPServer(handler http.Handler) *httptest.Server {
	t := s.T()
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func (s *BaseSuite) StartHTTPSServer(handler http.Handler) *httptest.Server {
	t := s.T()
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func (s *BaseSuite) WithTimeout(d time.Duration, fn func(ctx context.Context)) {
	t := s.T()
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	done := make(chan struct{})
	go func() {
		fn(ctx)
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		s.Require().Fail("timed out", d)
	}
}

func (s *BaseSuite) NextPort() int {
	return int(basePort + (atomic.AddInt32(&portCounter, 1) % portRange))
}

func (s *BaseSuite) GetSpecID(client *mcpClient) string {
	t := s.T()
	t.Helper()
	result := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	s.Require().NoError(json.Unmarshal(result, &specsResp))
	s.Require().NotEmpty(specsResp.Specs, "no specs found")
	return specsResp.Specs[0].ID
}

func (s *BaseSuite) GetEndpointID(client *mcpClient, specID, method, path string) string {
	t := s.T()
	t.Helper()
	epResult := client.callTool(t, "endpoint_by_spec", map[string]interface{}{
		specIDKey: specID,
	})
	var epResp struct {
		Endpoints []struct {
			ID     string `json:"id"`
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"endpoints"`
	}
	s.Require().NoError(json.Unmarshal(epResult, &epResp))
	for _, ep := range epResp.Endpoints {
		if ep.Method == method && ep.Path == path {
			return ep.ID
		}
	}
	s.Require().FailNowf("endpoint %s %s not found", method, path)
	return ""
}
