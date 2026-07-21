package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

var (
	binPath      string
	MockBinPath  string
	projectRoot  string
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

	MockBinPath = filepath.Join(tmpDir, "swag2mcp-mock")
	cmd = exec.Command("go", "build", "-o", MockBinPath, filepath.Join(projectRoot, "./cmd/swag2mcp-mock/"))
	out, err = cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "mock build failed: %v\n%s", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

type mcpClient struct {
	cmd    *exec.Cmd
	stdin  *os.File
	stdout *os.File
	stderr *bytes.Buffer
	mu     sync.Mutex
}

const mcpStartTimeout = 10 * time.Second

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
	if err := decodeWithTimeout(dec, &resp, mcpStartTimeout); err != nil {
		t.Fatalf("decode initialize response: %v\nstderr: %s", err, c.stderr.String())
	}

	if resp.Error != nil {
		t.Fatalf("initialize error: code=%d message=%s data=%s", resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

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

	reqData, err := c.writeCall(name, params)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	_, err = fmt.Fprintf(c.stdin, "%s\n", reqData)
	if err != nil {
		t.Fatalf("write request: %v", err)
	}

	var resp jsonRPCResponse
	dec := json.NewDecoder(c.stdout)
	if err := decodeWithTimeout(dec, &resp, mcpStartTimeout); err != nil {
		t.Fatalf("decode response: %v\nstderr: %s", err, c.stderr.String())
	}

	if resp.Error != nil {
		t.Fatalf("tool %s error: code=%d message=%s data=%s", name, resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

	var toolResult struct {
		StructuredContent json.RawMessage `json:"structuredContent"`
	}
	if err := json.Unmarshal(resp.Result, &toolResult); err != nil {
		t.Fatalf("parse tool result: %v", err)
	}
	return toolResult.StructuredContent
}

func (c *mcpClient) callToolRaw(t *testing.T, name string, params interface{}) (json.RawMessage, error) {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()

	reqData, err := c.writeCall(name, params)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Fprintf(c.stdin, "%s\n", reqData)
	if err != nil {
		return nil, err
	}

	var resp jsonRPCResponse
	dec := json.NewDecoder(c.stdout)
	if err := decodeWithTimeout(dec, &resp, mcpStartTimeout); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("tool %s error: code=%d message=%s", name, resp.Error.Code, resp.Error.Message)
	}

	var toolResult struct {
		StructuredContent json.RawMessage `json:"structuredContent"`
	}
	if err := json.Unmarshal(resp.Result, &toolResult); err != nil {
		return nil, err
	}
	return toolResult.StructuredContent, nil
}

func (c *mcpClient) writeCall(name string, params interface{}) ([]byte, error) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      name,
			"arguments": params,
		},
	}
	return json.Marshal(req)
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
	if err := decodeWithTimeout(dec, &resp, mcpStartTimeout); err != nil {
		t.Fatalf("decode response: %v\nstderr: %s", err, c.stderr.String())
	}

	if resp.Error != nil {
		t.Fatalf("list tools error: code=%d message=%s data=%s", resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

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

func decodeWithTimeout(dec *json.Decoder, v any, timeout time.Duration) error {
	type result struct {
		err error
	}
	ch := make(chan result, 1)
	go func() {
		ch <- result{err: dec.Decode(v)}
	}()
	select {
	case r := <-ch:
		return r.err
	case <-time.After(timeout):
		return fmt.Errorf("timed out after %v waiting for JSON response", timeout)
	}
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
