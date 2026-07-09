package cache

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

const defaultHTTPTimeout = 30 * time.Second

type httpClient struct {
	cli *http.Client
}

func defaultHTTPClient() *httpClient {
	cli, err := httpclient.New(httpclient.Config{
		Timeout: defaultHTTPTimeout,
	})
	if err != nil {
		cli = &http.Client{Timeout: defaultHTTPTimeout}
	}
	return &httpClient{cli: cli}
}

func (h *httpClient) Get(specURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, specURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := h.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.New("empty response body")
	}

	return data, nil
}

func fileURIToPath(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("expected file scheme, got %s", u.Scheme)
	}
	p, err := url.PathUnescape(u.Path)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" && len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	return filepath.FromSlash(p), nil
}
