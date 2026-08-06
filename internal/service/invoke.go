/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
)

// globalRateLimitKey is the fixed key used for the global rate limiter.
// It caps the total number of invoke requests per second across all endpoints.
const globalRateLimitKey = "global"

// InvokeRequest represents a request to invoke an API endpoint.
type InvokeRequest struct {
	EndpointID  string         `json:"endpointId"            validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to invoke"`
	Parameters  map[string]any `json:"parameters,omitempty"                          jsonschema:"optional,Path, query, and header parameters as key-value pairs"`
	RequestBody map[string]any `json:"requestBody,omitempty"                         jsonschema:"optional,Request body for POST/PUT/PATCH requests"`
}

// FileReference holds information about a response saved to disk.
type FileReference struct {
	Path        string `json:"path"`
	Size        int    `json:"size"`
	SizeHint    string `json:"sizeHint"`
	MaxSizeHint string `json:"maxSizeHint"`
	Message     string `json:"message"`
	OpenCmd     string `json:"openCmd"`
}

// InvokeResponse represents the response from invoking an API endpoint.
type InvokeResponse struct {
	StatusCode int               `json:"statusCode"        jsonschema:"required,HTTP response status code"`
	Headers    map[string]string `json:"headers"           jsonschema:"required,HTTP response headers"`
	Body       any               `json:"body"              jsonschema:"required,Response body data"`
	FileRef    *FileReference    `json:"fileRef,omitempty"`
}

type invokeService struct {
	ctx     *serviceContext
	index   IndexReader
	ws      WorkspaceOps
	v       RequestValidator
	dumpDir string
	log     *slog.Logger
}

func newInvokeService(
	ctx *serviceContext,
	index IndexReader,
	ws WorkspaceOps,
	v RequestValidator,
	dumpDir string,
	log *slog.Logger,
) *invokeService {
	return &invokeService{
		ctx:     ctx,
		index:   index,
		ws:      ws,
		v:       v,
		dumpDir: dumpDir,
		log:     log,
	}
}

// Invoke validates the request, builds an HTTP request, sends it, and returns the response.
func (is *invokeService) Invoke(ctx context.Context, rq InvokeRequest) (InvokeResponse, error) {
	if err := is.v.Struct(rq); err != nil {
		return InvokeResponse{}, NewInvalidEndpointIDError(err)
	}

	if !is.ctx.disableRateLimiter.Load() {
		if gl := is.ctx.loadGlobalRateLimiter(); gl != nil {
			if err := gl.Allow(globalRateLimitKey); err != nil {
				is.log.ErrorContext(ctx, "invoke failed: global rate limit", "endpoint_id", rq.EndpointID, "error", err)
				return InvokeResponse{}, NewGlobalRateLimitError(err)
			}
		}

		if err := is.ctx.loadRateLimiter().Allow(rq.EndpointID); err != nil {
			is.log.ErrorContext(ctx, "invoke failed: rate limit", "endpoint_id", rq.EndpointID, "error", err)
			return InvokeResponse{}, NewRateLimitError(err)
		}
	}

	ep, err := is.index.EndpointByID(rq.EndpointID)
	if err != nil {
		is.log.ErrorContext(ctx, "invoke failed: endpoint not found", "endpoint_id", rq.EndpointID, "error", err)
		return InvokeResponse{}, NewEndpointNotFoundError(rq.EndpointID, err)
	}

	if ep.Operation == nil {
		is.log.ErrorContext(ctx, "invoke failed: no operation", "endpoint_id", rq.EndpointID)
		return InvokeResponse{}, NewInvokeError(
			"This endpoint has no operation definition. It may be malformed or incomplete.",
			nil,
		)
	}

	sp, err := is.index.SpecByID(ep.SpecID)
	if err != nil {
		is.log.ErrorContext(ctx, "invoke failed: spec not found", "spec_id", ep.SpecID, "error", err)
		return InvokeResponse{}, NewSpecNotFoundError(ep.SpecID, err)
	}

	coll, err := is.index.CollectionByID(ep.CollectionID)
	if err != nil {
		is.log.ErrorContext(ctx, "invoke failed: collection not found", "collection_id", ep.CollectionID, "error", err)
		return InvokeResponse{}, NewCollectionNotFoundError(ep.CollectionID, err)
	}

	if err := validateParameters(ep.Operation, rq.Parameters); err != nil {
		is.log.ErrorContext(ctx, "invoke failed: parameter validation", "endpoint_id", rq.EndpointID, "error", err)
		return InvokeResponse{}, NewParameterValidationError(err)
	}

	if err := validateRequestBody(ep.Operation, rq.RequestBody); err != nil {
		is.log.ErrorContext(ctx, "invoke failed: request body validation", "endpoint_id", rq.EndpointID, "error", err)
		return InvokeResponse{}, NewRequestBodyValidationError(err)
	}

	req, err := is.buildRequest(ctx, sp, coll, ep, rq)
	if err != nil {
		is.log.ErrorContext(ctx, "invoke failed: build request", "endpoint_id", rq.EndpointID, "spec", sp.Domain, "error", err)
		return InvokeResponse{}, NewBuildRequestError(err)
	}

	return is.executeRequest(ctx, req, sp, ep)
}

func (is *invokeService) buildRequest(
	ctx context.Context,
	sp *model.Spec,
	coll *model.Collection,
	ep *model.Endpoint,
	rq InvokeRequest,
) (*http.Request, error) {
	cfg := is.ctx.loadConfig()
	mockEnabled := cfg != nil && cfg.MockEnabled
	return newRequestBuilder(
		withSpec(sp),
		withCollection(coll),
		withEndpoint(ep),
		withParameters(rq.Parameters),
		withBody(rq.RequestBody),
		withHTTPConfig(mergeHTTPClientConfigs(sp.HTTPClient, coll.HTTPClient)),
		withGlobalHeaders(is.ctx.loadGlobalHeaders()),
		withGlobalUserAgent(is.ctx.loadGlobalUserAgent()),
		withGlobalCookies(is.ctx.loadGlobalCookies()),
		withMockEnabled(mockEnabled),
	).build(ctx)
}

func (is *invokeService) executeRequest(
	_ context.Context,
	req *http.Request,
	sp *model.Spec,
	ep *model.Endpoint,
) (InvokeResponse, error) {
	client := is.ctx.loadHTTPClient()
	needsAuth := sp.Auth != nil && (!sp.HasSecureEndpoints || ep.RequiresAuth)
	if needsAuth {
		base := client.Transport
		if base == nil {
			base = http.DefaultTransport
		}
		client = &http.Client{
			Transport: &auth.Transport{
				Base: base,
				Auth: sp.Auth,
			},
			Timeout:       client.Timeout,
			CheckRedirect: client.CheckRedirect,
		}
		if err := sp.Auth.Apply(req, nil); err != nil {
			is.log.ErrorContext(req.Context(), "invoke failed: apply auth", "spec", sp.Domain, "error", err)
			return InvokeResponse{}, NewAuthError("failed to apply auth", err)
		}
	}

	is.dumpRequest(req, sp.Domain)

	response, err := client.Do(req)
	if err != nil {
		is.log.ErrorContext(req.Context(), "invoke failed: http request", "spec", sp.Domain, "error", err)
		return InvokeResponse{}, NewHTTPRequestError(err)
	}
	defer response.Body.Close()

	maxSize := is.ctx.MaxResponseSize()

	return is.streamOrBuffer(response, sp.Domain, ep, maxSize)
}

// streamOrBuffer reads the response body into a buffer up to maxSize+1 bytes.
// If the response fits within maxSize, it returns the body inline.
// If it exceeds maxSize, it writes the buffer to a file and streams the rest from response.Body.
func (is *invokeService) streamOrBuffer(
	r *http.Response,
	domain string,
	ep *model.Endpoint,
	maxSize int,
) (InvokeResponse, error) {
	headers := make(map[string]string, len(r.Header))
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ", ")
	}

	reqCtx := context.Background()
	if r.Request != nil {
		reqCtx = r.Request.Context()
	}

	buf := make([]byte, maxSize+1)
	n, readErr := io.ReadFull(r.Body, buf)
	buf = buf[:n]

	if readErr != nil && readErr != io.ErrUnexpectedEOF && readErr != io.EOF {
		is.log.ErrorContext(reqCtx, "invoke failed: read response", "spec", domain, "error", readErr)
		return InvokeResponse{}, NewResponseReadError(readErr)
	}

	if n <= maxSize {
		return newInvokeResponse(r, buf), nil
	}

	fp := is.responseFilePath(domain, ep)
	if err := os.MkdirAll(is.ws.ResponsesDir(), 0750); err != nil {
		is.log.ErrorContext(reqCtx, "invoke failed: create responses dir", "error", err)
		return InvokeResponse{}, NewDirCreateError(err)
	}

	f, err := os.Create(fp)
	if err != nil {
		is.log.ErrorContext(reqCtx, "invoke failed: create response file", "path", fp, "error", err)
		return InvokeResponse{}, NewFileCreateError(err)
	}
	defer f.Close()

	if _, err := f.Write(buf); err != nil {
		is.log.ErrorContext(reqCtx, "invoke failed: write response file", "path", fp, "error", err)
		return InvokeResponse{}, NewFileWriteError(err)
	}

	written, err := io.Copy(f, r.Body)
	if err != nil {
		is.log.ErrorContext(reqCtx, "invoke failed: stream response", "path", fp, "error", err)
		return InvokeResponse{}, NewStreamError(err)
	}

	totalSize := n + int(written)
	size := formatSize(totalSize)
	maxSizeStr := formatSize(maxSize)
	msg := fmt.Sprintf(
		"Response body (%s) exceeds the maximum size limit (%s). The full response has been saved to disk.",
		size, maxSizeStr,
	)

	return InvokeResponse{
		StatusCode: r.StatusCode,
		Headers:    headers,
		Body: map[string]string{
			"message": msg,
		},
		FileRef: &FileReference{
			Path:        fp,
			Size:        totalSize,
			SizeHint:    size,
			MaxSizeHint: maxSizeStr,
			Message:     msg,
			OpenCmd:     openCommand(fp),
		},
	}, nil
}

// responseFilePath generates a unique file path for a response file.
func (is *invokeService) responseFilePath(domain string, ep *model.Endpoint) string {
	m := strings.ToLower(ep.Name)
	p := strings.TrimPrefix(ep.Path, "/")
	p = strings.ReplaceAll(p, "/", "_")
	p = strings.ReplaceAll(p, "{", "")
	p = strings.ReplaceAll(p, "}", "")
	suf := randomSuffix(config.RandSuffixLen)
	fname := fmt.Sprintf("%s-%s-%s-%s.json", domain, m, p, suf)
	return filepath.Join(is.ws.ResponsesDir(), fname)
}

// dumpRequest writes the HTTP request to a file for debugging if dumpDir is configured.
func (is *invokeService) dumpRequest(req *http.Request, domain string) {
	if len(is.dumpDir) == 0 {
		return
	}

	d, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return
	}

	ts := time.Now().UnixMilli()
	fname := fmt.Sprintf("invoke-%s-%d.txt", domain, ts)
	fp := filepath.Join(is.dumpDir, fname)

	if err := os.MkdirAll(is.dumpDir, 0750); err != nil {
		is.log.WarnContext(req.Context(), "failed to create dump dir", "error", err)
		return
	}
	if err := os.WriteFile(fp, d, 0600); err != nil {
		is.log.WarnContext(req.Context(), "failed to write dump file", "error", err)
	}
}
