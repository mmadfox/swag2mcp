package service

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
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/model"
)

// InvokeRequest represents a request to invoke an API endpoint.
type InvokeRequest struct {
	EndpointID  string            `json:"endpointId"            validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to invoke"`
	Parameters  map[string]any    `json:"parameters,omitempty"                          jsonschema:"optional,Path, query, and header parameters as key-value pairs"`
	RequestBody map[string]any    `json:"requestBody,omitempty"                         jsonschema:"optional,Request body for POST/PUT/PATCH requests"`
	Headers     map[string]string `json:"headers,omitempty"                             jsonschema:"optional,Additional HTTP headers to send with the request"`
	Cookies     map[string]string `json:"cookies,omitempty"                             jsonschema:"optional,Additional HTTP cookies to send with the request"`
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
}

func newInvokeService(
	ctx *serviceContext,
	index IndexReader,
	ws WorkspaceOps,
	v RequestValidator,
	dumpDir string,
) *invokeService {
	return &invokeService{
		ctx:     ctx,
		index:   index,
		ws:      ws,
		v:       v,
		dumpDir: dumpDir,
	}
}

// Invoke validates the request, builds an HTTP request, sends it, and returns the response.
func (is *invokeService) Invoke(ctx context.Context, rq InvokeRequest) (InvokeResponse, error) {
	if err := is.v.Struct(rq); err != nil {
		return InvokeResponse{}, NewValidationError(
			"The endpoint ID is invalid. It must be a 32-character hex string. "+
				"Use the search tool to find the correct endpoint ID.",
			err,
		)
	}

	if !is.ctx.disableRateLimiter.Load() {
		if err := is.ctx.loadRateLimiter().Allow(rq.EndpointID); err != nil {
			return InvokeResponse{}, NewRateLimitError(err)
		}
	}

	ep, err := is.index.EndpointByID(rq.EndpointID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Endpoint %q was not found. Use the search tool to find the correct endpoint ID.",
				rq.EndpointID,
			),
			err,
		)
	}

	if ep.Operation == nil {
		return InvokeResponse{}, NewValidationError(
			"This endpoint has no operation definition. It may be malformed or incomplete.",
			nil,
		)
	}

	sp, err := is.index.SpecByID(ep.SpecID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Spec %q was not found. The endpoint references a spec that no longer exists.",
				ep.SpecID,
			),
			err,
		)
	}

	coll, err := is.index.CollectionByID(ep.CollectionID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Collection %q was not found. The endpoint references a collection that no longer exists.",
				ep.CollectionID,
			),
			err,
		)
	}

	if err := validateParameters(ep.Operation, rq.Parameters); err != nil {
		return InvokeResponse{}, NewValidationError(
			"Parameter validation failed. Check that all required parameters are provided "+
				"and match the expected names.",
			err,
		)
	}

	if err := validateRequestBody(ep.Operation, rq.RequestBody); err != nil {
		return InvokeResponse{}, NewValidationError(
			"Request body validation failed. Check that all required fields are present "+
				"and no unknown fields are included.",
			err,
		)
	}

	req, err := is.buildRequest(ctx, sp, coll, ep, rq)
	if err != nil {
		return InvokeResponse{}, NewInvokeError(
			"Failed to build the HTTP request. Check the endpoint parameters and try again.",
			err,
		)
	}

	is.dumpRequest(req, sp.Domain)

	return is.executeRequest(ctx, req, sp, ep)
}

func (is *invokeService) buildRequest(
	ctx context.Context,
	sp *model.Spec,
	coll *model.Collection,
	ep *model.Endpoint,
	rq InvokeRequest,
) (*http.Request, error) {
	return newRequestBuilder(
		withContext(ctx),
		withSpec(sp),
		withCollection(coll),
		withEndpoint(ep),
		withParameters(rq.Parameters),
		withBody(rq.RequestBody),
		withHTTPConfig(mergeHTTPClientConfigs(sp.HTTPClient, coll.HTTPClient)),
		withInvokeHeaders(rq.Headers),
		withInvokeCookies(rq.Cookies),
		withGlobalHeaders(is.ctx.loadGlobalHeaders()),
		withGlobalUserAgent(is.ctx.loadGlobalUserAgent()),
		withGlobalCookies(is.ctx.loadGlobalCookies()),
	).build()
}

func (is *invokeService) executeRequest(
	_ context.Context,
	req *http.Request,
	sp *model.Spec,
	ep *model.Endpoint,
) (InvokeResponse, error) {
	client := is.ctx.loadHTTPClient()
	if sp.Auth != nil {
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
	}

	response, err := client.Do(req)
	if err != nil {
		return InvokeResponse{}, NewInvokeError(
			"The API request failed.",
			err,
		)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return InvokeResponse{}, NewInvokeError(
			"Failed to read the API response.",
			err,
		)
	}

	maxSize := is.ctx.MaxResponseSize()
	if len(body) > maxSize {
		return is.saveLargeResponse(response, body, sp.Domain, ep, maxSize)
	}

	return newInvokeResponse(response, body), nil
}

// saveLargeResponse saves a response body that exceeds the max size to a file
// and returns an InvokeResponse with a FileReference instead of the full body.
func (is *invokeService) saveLargeResponse(
	r *http.Response,
	body []byte,
	domain string,
	ep *model.Endpoint,
	maxSize int,
) (InvokeResponse, error) {
	headers := make(map[string]string, len(r.Header))
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ", ")
	}

	m := strings.ToLower(ep.Name)
	p := strings.TrimPrefix(ep.Path, "/")
	p = strings.ReplaceAll(p, "/", "_")
	p = strings.ReplaceAll(p, "{", "")
	p = strings.ReplaceAll(p, "}", "")
	suf := randomSuffix(config.RandSuffixLen)
	fname := fmt.Sprintf("%s-%s-%s-%s.json", domain, m, p, suf)
	fp := filepath.Join(is.ws.ResponsesDir(), fname)

	if err := os.MkdirAll(is.ws.ResponsesDir(), 0750); err != nil {
		return InvokeResponse{}, fmt.Errorf("failed to create responses dir: %w", err)
	}

	if err := os.WriteFile(fp, body, 0600); err != nil {
		return InvokeResponse{}, fmt.Errorf("failed to write response file: %w", err)
	}

	size := formatSize(len(body))
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
			Size:        len(body),
			SizeHint:    size,
			MaxSizeHint: maxSizeStr,
			Message:     msg,
			OpenCmd:     openCommand(fp),
		},
	}, nil
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
		slog.Default().WarnContext(req.Context(), "failed to create dump dir", "error", err)
		return
	}
	if err := os.WriteFile(fp, d, 0600); err != nil {
		slog.Default().WarnContext(req.Context(), "failed to write dump file", "error", err)
	}
}
