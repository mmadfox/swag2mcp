/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
)

// requestBuilder builds an [http.Request] from spec, collection, endpoint, and parameters.
type requestBuilder struct {
	spec            *model.Spec
	collection      *model.Collection
	endpoint        *model.Endpoint
	parameters      map[string]any
	body            map[string]any
	httpConfig      *model.HTTPClientConfig
	globalHeaders   map[string]string
	globalUserAgent string
	globalCookies   []httpclient.Cookie
	mockEnabled     bool
}

// requestOption is a functional option for configuring a requestBuilder.
type requestOption func(*requestBuilder)

// newRequestBuilder creates a new requestBuilder with the given options.
func newRequestBuilder(options ...requestOption) *requestBuilder {
	builder := &requestBuilder{}
	for _, option := range options {
		option(builder)
	}
	return builder
}

func withSpec(specification *model.Spec) requestOption {
	return func(builder *requestBuilder) {
		builder.spec = specification
	}
}

func withCollection(collection *model.Collection) requestOption {
	return func(builder *requestBuilder) {
		builder.collection = collection
	}
}

func withEndpoint(endpoint *model.Endpoint) requestOption {
	return func(builder *requestBuilder) {
		builder.endpoint = endpoint
	}
}

func withParameters(parameters map[string]any) requestOption {
	return func(builder *requestBuilder) {
		builder.parameters = parameters
	}
}

func withBody(body map[string]any) requestOption {
	return func(builder *requestBuilder) {
		builder.body = body
	}
}

func withHTTPConfig(config *model.HTTPClientConfig) requestOption {
	return func(builder *requestBuilder) {
		builder.httpConfig = config
	}
}

func withGlobalHeaders(headers map[string]string) requestOption {
	return func(builder *requestBuilder) {
		builder.globalHeaders = headers
	}
}

func withGlobalUserAgent(ua string) requestOption {
	return func(builder *requestBuilder) {
		builder.globalUserAgent = ua
	}
}

func withGlobalCookies(cookies []httpclient.Cookie) requestOption {
	return func(builder *requestBuilder) {
		builder.globalCookies = cookies
	}
}

func withMockEnabled(enabled bool) requestOption {
	return func(builder *requestBuilder) {
		builder.mockEnabled = enabled
	}
}

// build constructs the [http.Request] from the configured options.
func (builder *requestBuilder) build(ctx context.Context) (*http.Request, error) {
	baseURL := builder.resolveBaseURL()
	baseURL = strings.TrimRight(baseURL, "/")
	reqURL := baseURL + "/" + strings.TrimLeft(builder.endpoint.Path, "/")

	pathParams := builder.filterParametersByLocation("path")
	for name, val := range pathParams {
		reqURL = strings.ReplaceAll(
			reqURL,
			"{"+name+"}",
			url.PathEscape(val),
		)
	}

	u, err := url.Parse(reqURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", reqURL, err)
	}

	queryParams := builder.filterParametersByLocation("query")
	params := u.Query()
	for name, val := range queryParams {
		params.Set(name, val)
	}
	u.RawQuery = params.Encode()

	var body io.Reader
	if builder.body != nil {
		data, err := json.Marshal(builder.body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		builder.endpoint.Name,
		u.String(),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	builder.applyHeaders(req)
	builder.applyHTTPClientConfig(req)

	return req, nil
}

func (builder *requestBuilder) resolveBaseURL() string {
	if builder.mockEnabled && builder.collection.BaseMockURL != "" {
		return "http://" + builder.collection.BaseMockURL
	}
	if builder.collection.BaseURL != "" {
		return builder.collection.BaseURL
	}
	return builder.spec.BaseURL
}

func (builder *requestBuilder) filterParametersByLocation(location string) map[string]string {
	result := make(map[string]string, len(builder.parameters))
	for _, parameter := range builder.endpoint.Operation.Parameters {
		if parameter.In != location {
			continue
		}
		value, exists := builder.parameters[parameter.Name]
		if !exists {
			continue
		}
		if isEmptyParam(value) {
			continue
		}
		result[parameter.Name] = formatParamValue(value)
	}
	return result
}

// isEmptyParam returns true for empty/zero parameter values that should not
// be included in the request URL. This prevents LLM-provided empty defaults
// (e.g. signature: "", timestamp: 0) from being sent to the API.
func isEmptyParam(v any) bool {
	switch val := v.(type) {
	case string:
		return val == ""
	case int:
		return val == 0
	case int8:
		return val == 0
	case int16:
		return val == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	default:
		return false
	}
}

// formatParamValue converts a parameter value to its string representation,
// handling numeric types correctly to avoid scientific notation (e.g. 7.31434e+06).
func formatParamValue(value any) string {
	switch v := value.(type) {
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		if v == float32(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (builder *requestBuilder) applyHeaders(req *http.Request) {
	headers := builder.filterParametersByLocation("header")
	for name, val := range headers {
		req.Header.Set(name, val)
	}

	if builder.body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}
}

func (builder *requestBuilder) applyHTTPClientConfig(req *http.Request) {
	builder.applyGlobalConfig(req)
	builder.applySpecConfig(req)
	builder.applyDefaultAccept(req)
}

func (builder *requestBuilder) applyGlobalConfig(req *http.Request) {
	for name, val := range builder.globalHeaders {
		if req.Header.Get(name) == "" {
			req.Header.Set(name, val)
		}
	}

	if req.Header.Get("User-Agent") == "" && builder.globalUserAgent != "" {
		req.Header.Set("User-Agent", builder.globalUserAgent)
	}

	for _, c := range builder.globalCookies {
		req.AddCookie(&http.Cookie{
			Name: c.Name, Value: c.Value, Domain: c.Domain,
			Path: c.Path, Secure: c.Secure, HttpOnly: c.HTTPOnly,
		})
	}
}

func (builder *requestBuilder) applySpecConfig(req *http.Request) {
	if builder.httpConfig == nil {
		return
	}

	for name, val := range builder.httpConfig.Headers {
		req.Header.Set(name, val)
	}

	for _, cookie := range builder.httpConfig.Cookies {
		req.AddCookie(&http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HttpOnly: cookie.HTTPOnly,
		})
	}
}

func (builder *requestBuilder) applyDefaultAccept(req *http.Request) {
	if req.Header.Get("Accept") != "" {
		return
	}

	isJSON := builder.body != nil ||
		req.Header.Get("Content-Type") == contentTypeJSON
	if isJSON {
		req.Header.Set("Accept", acceptHeaderJSON)
	} else {
		req.Header.Set("Accept", acceptHeaderOther)
	}
}
