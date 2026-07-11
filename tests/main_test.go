package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	binPath     string
	mu          sync.Mutex
	portCounter int32
	projectRoot string
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot = filepath.Dir(filepath.Dir(filename))

	tmpDir, err := os.MkdirTemp("", "swag2mcp-test-*")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	binPath = filepath.Join(tmpDir, "swag2mcp")
	cmd := exec.Command("go", "build", "-o", binPath, filepath.Join(projectRoot, "./cmd/swag2mcp/"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "build failed: %v\n%s", err, out)
		os.Exit(1)
	}

	mockBinPath := filepath.Join(tmpDir, "swag2mcp-mock")
	cmd = exec.Command("go", "build", "-o", mockBinPath, filepath.Join(projectRoot, "./cmd/swag2mcp-mock/"))
	out, err = cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "mock build failed: %v\n%s", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func assertEqual[T any](t *testing.T, name string, actual, expected T) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("[%s] mismatch:\n  actual:   %v\n  expected: %v", name, actual, expected)
	}
}

func assertNotEqual[T any](t *testing.T, name string, actual, unexpected T) {
	t.Helper()
	if reflect.DeepEqual(actual, unexpected) {
		t.Errorf("[%s] should not equal %v", name, unexpected)
	}
}

func assertContains(t *testing.T, name, actual, substr string) {
	t.Helper()
	if !strings.Contains(actual, substr) {
		t.Errorf("[%s] expected to contain %q:\n%s", name, substr, actual)
	}
}

func assertNotContains(t *testing.T, name, actual, substr string) {
	t.Helper()
	if strings.Contains(actual, substr) {
		t.Errorf("[%s] should not contain %q:\n%s", name, substr, actual)
	}
}

func assertErrorCode(t *testing.T, name string, errBody []byte, expectedCode string) {
	t.Helper()
	var errResp struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(errBody, &errResp); err != nil {
		t.Fatalf("[%s] failed to parse error JSON: %v\nbody: %s", name, err, errBody)
	}
	assertEqual(t, name+".code", errResp.Code, expectedCode)
}

func runCommand(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = projectRoot
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run %v: %v", args, err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func runCommandInWS(t *testing.T, ws string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = ws
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run %v: %v", args, err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func runCommandWithStdin(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = projectRoot
	cmd.Stdin = strings.NewReader(stdin)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run %v: %v", args, err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func runCommandWithStdinInWS(t *testing.T, ws, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = ws
	cmd.Stdin = strings.NewReader(stdin)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run %v: %v", args, err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func runCommandWithEnv(t *testing.T, env []string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), env...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run %v: %v", args, err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func newTestWorkspace(t *testing.T) string {
	t.Helper()
	ws, err := os.MkdirTemp("", "swag2mcp-ws-*")
	if err != nil {
		t.Fatalf("failed to create workspace: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	return ws
}

func writeConfig(t *testing.T, ws, content string) string {
	t.Helper()
	wsDir := filepath.Join(ws, ".swag2mcp")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("failed to create .swag2mcp dir: %v", err)
	}

	// Copy testdata files into workspace root so relative paths work
	// when CWD is set to ws
	src := filepath.Join(projectRoot, "tests", "testdata")
	dst := filepath.Join(ws, "testdata")
	if err := copyDir(src, dst); err != nil {
		t.Fatalf("failed to copy testdata: %v", err)
	}

	// Copy internal service testdata for parsing tests
	src = filepath.Join(projectRoot, "internal", "service", "testdata")
	dst = filepath.Join(ws, "internal", "service", "testdata")
	if err := copyDir(src, dst); err != nil {
		t.Fatalf("failed to copy internal testdata: %v", err)
	}

	path := filepath.Join(wsDir, "swag2mcp.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	return path
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func writeSpec(t *testing.T, ws, name, content string) string {
	t.Helper()
	path := filepath.Join(ws, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write spec: %v", err)
	}
	return path
}

func initWorkspace(t *testing.T, ws string) {
	t.Helper()
	stdout, stderr, code := runCommand(t, "init", ws)
	if code != 0 {
		t.Fatalf("init failed (exit %d):\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
}

type mcpClient struct {
	cmd    *exec.Cmd
	stdin  *os.File
	stdout *os.File
	stderr *bytes.Buffer
	mu     sync.Mutex
}

func startMCPStdio(t *testing.T, ws, configContent string, extraArgs ...string) *mcpClient {
	t.Helper()
	writeConfig(t, ws, configContent)

	args := []string{"mcp", "."}
	args = append(args, extraArgs...)
	cmd := exec.Command(binPath, args...)
	cmd.Dir = ws

	stdinRead, stdinWrite, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdoutRead, stdoutWrite, err := os.Pipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	cmd.Stdin = stdinRead
	cmd.Stdout = stdoutWrite

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Start(); err != nil {
		t.Fatalf("mcp start: %v", err)
	}

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

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (c *mcpClient) initialize(t *testing.T) {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "swag2mcp-test",
				"version": "1.0.0",
			},
		},
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal initialize: %v", err)
	}

	_, err = fmt.Fprintf(c.stdin, "%s\n", reqData)
	if err != nil {
		t.Fatalf("write initialize: %v", err)
	}

	var resp jsonRPCResponse
	dec := json.NewDecoder(c.stdout)
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("decode initialize response: %v\nstderr: %s", err, c.stderr.String())
	}

	if resp.Error != nil {
		t.Fatalf("initialize error: code=%d message=%s data=%s", resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

	// Send initialized notification (no ID for notifications)
	notif := map[string]string{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}
	notifData, _ := json.Marshal(notif)
	_, _ = fmt.Fprintf(c.stdin, "%s\n", notifData)
}

func (c *mcpClient) callTool(t *testing.T, name string, params interface{}) json.RawMessage {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":   name,
			"arguments": params,
		},
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	_, err = fmt.Fprintf(c.stdin, "%s\n", reqData)
	if err != nil {
		t.Fatalf("write request: %v", err)
	}

	var resp jsonRPCResponse
	dec := json.NewDecoder(c.stdout)
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("decode response: %v\nstderr: %s", err, c.stderr.String())
	}

	if resp.Error != nil {
		t.Fatalf("tool %s error: code=%d message=%s data=%s", name, resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

	// Extract StructuredContent from CallToolResult
	var toolResult struct {
		StructuredContent json.RawMessage `json:"structuredContent"`
	}
	if err := json.Unmarshal(resp.Result, &toolResult); err != nil {
		t.Fatalf("parse tool result: %v", err)
	}
	return toolResult.StructuredContent
}

func (c *mcpClient) listTools(t *testing.T) json.RawMessage {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	_, err = fmt.Fprintf(c.stdin, "%s\n", reqData)
	if err != nil {
		t.Fatalf("write request: %v", err)
	}

	var resp jsonRPCResponse
	dec := json.NewDecoder(c.stdout)
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("decode response: %v\nstderr: %s", err, c.stderr.String())
	}

	if resp.Error != nil {
		t.Fatalf("list tools error: code=%d message=%s data=%s", resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

	// Extract tools from ListToolsResult
	var listResult struct {
		Tools []struct {
			Name string `json:"name"`
		} `json:"tools"`
	}
	if err := json.Unmarshal(resp.Result, &listResult); err != nil {
		t.Fatalf("parse list result: %v", err)
	}
	result, _ := json.Marshal(listResult)
	return result
}

func startHTTPServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func startHTTPSServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func withTimeout(t *testing.T, d time.Duration, fn func(ctx context.Context)) {
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
		t.Fatalf("timed out after %v", d)
	}
}

func nextPort() int {
	return int(9000 + (atomic.AddInt32(&portCounter, 1) % 1000))
}
