package cache

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
)

type httpClient struct {
	client *http.Client
}

// defaultHTTPClient creates an httpClient with the default timeout.
func defaultHTTPClient() *httpClient {
	cli, err := httpclient.New(httpclient.Config{
		Timeout: defaultHTTPTimeout,
	})
	if err != nil {
		slog.Default().Warn("failed to create default HTTP client, using fallback", "error", err)
		cli = &http.Client{Timeout: defaultHTTPTimeout}
	}
	return &httpClient{client: cli}
}

// SetClient replaces the underlying [http.Client].
func (h *httpClient) SetClient(cli *http.Client) {
	h.client = cli
}

// Do sends an HTTP request using the underlying client.
func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

// Get fetches a spec from the given URL and returns the response body.
func (h *httpClient) Get(ctx context.Context, specURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, specURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := h.Do(req)
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
		return nil, ErrEmptyBody
	}

	return data, nil
}

// fileURIToPath converts a file:// URI to a local filesystem path.
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
