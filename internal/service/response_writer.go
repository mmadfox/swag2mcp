package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"runtime"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
)

// newInvokeResponse reads the HTTP response and creates an InvokeResponse.
func newInvokeResponse(response *http.Response, body []byte) InvokeResponse {
	headers := make(map[string]string, len(response.Header))
	for key, values := range response.Header {
		headers[key] = strings.Join(values, ", ")
	}

	var parsedBody any
	if len(body) > 0 {
		if jsonError := json.Unmarshal(body, &parsedBody); jsonError == nil {
			return InvokeResponse{
				StatusCode: response.StatusCode,
				Headers:    headers,
				Body:       parsedBody,
			}
		}
	}

	return InvokeResponse{
		StatusCode: response.StatusCode,
		Headers:    headers,
		Body:       string(body),
	}
}

// resolveMaxResponseSize returns the effective max response size.
// Default is 1 MB, maximum is 10 MB.
func resolveMaxResponseSize(size *int) int {
	if size == nil {
		return config.DefaultMaxResponseSize
	}
	if *size > config.MaxAllowedResponseSize {
		return config.MaxAllowedResponseSize
	}
	if *size <= 0 {
		return config.DefaultMaxResponseSize
	}
	return *size
}

// openCommand returns the OS-specific command to open a file.
func openCommand(path string) string {
	switch runtime.GOOS {
	case "darwin":
		return "open " + path
	case "windows":
		return "start " + path
	default:
		return "xdg-open " + path
	}
}

// formatSize returns a human-readable size string.
func formatSize(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// randomSuffix generates a random hex string of length n.
func randomSuffix(n int) string {
	byteLen := (n + 1) / 2 //nolint:mnd // hex encoding: 2 chars per byte
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%0*x", n, 0)
	}
	return hex.EncodeToString(b)[:n]
}

// mergeHTTPClientConfigs merges two per-request HTTP configs. Collection overrides spec (last-wins).
func mergeHTTPClientConfigs(sp, coll *model.HTTPClientConfig) *model.HTTPClientConfig {
	if sp == nil && coll == nil {
		return nil
	}

	result := &model.HTTPClientConfig{}

	levels := []*model.HTTPClientConfig{sp, coll}

	for _, level := range levels {
		if level == nil {
			continue
		}
		if len(level.Headers) > 0 {
			if result.Headers == nil {
				result.Headers = make(map[string]string, len(level.Headers))
			}
			maps.Copy(result.Headers, level.Headers)
		}
		if len(level.Cookies) > 0 {
			result.Cookies = level.Cookies
		}
	}

	return result
}
