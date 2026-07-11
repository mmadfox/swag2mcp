package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestScript_Transport_Stdio(t *testing.T) {
	ws := newTestWorkspace(t)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false")
	client.initialize(t)
	result := client.listTools(t)

	var toolsResp struct {
		Tools []interface{} `json:"tools"`
	}
	if err := json.Unmarshal(result, &toolsResp); err != nil {
		t.Fatalf("parse tools: %v", err)
	}
	if len(toolsResp.Tools) == 0 {
		t.Errorf("expected tools from stdio transport")
	}
}

func TestScript_Transport_SSE(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	port := nextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", ws,
		"--transport", "sse",
		"--http-addr", addr,
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	if err := cmd.Start(); err != nil {
		t.Fatalf("mcp sse start: %v", err)
	}
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("SSE request failed: %v\nstderr: %s", err, stderrBuf.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("SSE expected 200, got %d: %s\nstderr: %s", resp.StatusCode, body, stderrBuf.String())
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "event:") {
		t.Errorf("expected SSE event stream, got: %s", string(body))
	}
}

func TestScript_Transport_StreamableHTTP(t *testing.T) {
	t.Skip("needs HTTP server")
}

func TestScript_Transport_AuthToken(t *testing.T) {
	ws := newTestWorkspace(t)
	initWorkspace(t, ws)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: https://api.example.com
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	writeConfig(t, ws, configContent)

	port := nextPort()
	addr := fmt.Sprintf(":%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath, "mcp", ws,
		"--transport", "sse",
		"--http-addr", addr,
		"--auth-token", "my-secret-token",
	)
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	if err := cmd.Start(); err != nil {
		t.Fatalf("mcp auth start: %v", err)
	}
	defer func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_ = cmd.Wait()
	}()

	time.Sleep(2 * time.Second)

	url := fmt.Sprintf("http://localhost:%d/mcp", port)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer my-secret-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("auth request failed: %v\nstderr: %s", err, stderrBuf.String())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 with valid token, got %d", resp.StatusCode)
	}

	resp2, err := http.Get(url)
	if err != nil {
		t.Fatalf("request without token failed: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", resp2.StatusCode)
	}
}

func TestScript_Transport_DumpDir(t *testing.T) {
	ws := newTestWorkspace(t)

	dumpDir := filepath.Join(ws, "dumps")
	_ = os.MkdirAll(dumpDir, 0755)

	mux := http.NewServeMux()
	mux.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	})
	srv := startHTTPServer(t, mux)

	configContent := `specs:
  - domain: petstore
    llm_title: Petstore API
    base_url: ` + srv.URL + `
    collections:
      - title: Pets
        location: ./testdata/petstore.yaml
`
	client := startMCPStdio(t, ws, configContent, "--disable-llm-auth=false", "--dump-dir", dumpDir)
	client.initialize(t)

	specsResult := client.callTool(t, "spec_list", map[string]interface{}{})
	var specsResp struct {
		Specs []struct {
			ID string `json:"id"`
		} `json:"specs"`
	}
	if err := json.Unmarshal(specsResult, &specsResp); err != nil {
		t.Fatalf("parse spec_list: %v", err)
	}
	if len(specsResp.Specs) == 0 {
		t.Fatal("no specs found")
	}

	epResult := client.callTool(t, "endpoint_by_spec", map[string]interface{}{
		"specId": specsResp.Specs[0].ID,
	})
	var epResp struct {
		Endpoints []struct {
			ID     string `json:"id"`
			Method string `json:"method"`
			Path   string `json:"path"`
		} `json:"endpoints"`
	}
	if err := json.Unmarshal(epResult, &epResp); err != nil {
		t.Fatalf("parse endpoints: %v", err)
	}

	var getEndpointID string
	for _, ep := range epResp.Endpoints {
		if ep.Method == "GET" && ep.Path == "/pets" {
			getEndpointID = ep.ID
			break
		}
	}
	if getEndpointID == "" {
		t.Fatal("GET /pets endpoint not found")
	}

	client.callTool(t, "invoke", map[string]interface{}{
		"endpointId": getEndpointID,
	})

	entries, _ := os.ReadDir(dumpDir)
	if len(entries) == 0 {
		t.Errorf("expected dump files in %s", dumpDir)
	}
}
