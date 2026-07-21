package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/httpclient"
	"github.com/mmadfox/swag2mcp/internal/model"
)

// requestBuilder builds an [http.Request] from spec, collection, endpoint, and parameters.
type requestBuilder struct {
	context         context.Context
	spec            *model.Spec
	collection      *model.Collection
	endpoint        *model.Endpoint
	parameters      map[string]any
	body            map[string]any
	httpConfig      *model.HTTPClientConfig
	invokeHeaders   map[string]string
	invokeCookies   map[string]string
	globalHeaders   map[string]string
	globalUserAgent string
	globalCookies   []httpclient.Cookie
}

// requestOption is a functional option for configuring a requestBuilder.
type requestOption func(*requestBuilder)

// newRequestBuilder creates a new requestBuilder with the given options.
func newRequestBuilder(options ...requestOption) *requestBuilder {
	builder := &requestBuilder{}
	for _, option := range options {
		option(builder)
	}
	if builder.context == nil {
		builder.context = context.Background()
	}
	return builder
}

func withContext(ctx context.Context) requestOption {
	return func(builder *requestBuilder) {
		builder.context = ctx
	}
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

func withInvokeHeaders(headers map[string]string) requestOption {
	return func(builder *requestBuilder) {
		builder.invokeHeaders = headers
	}
}

func withInvokeCookies(cookies map[string]string) requestOption {
	return func(builder *requestBuilder) {
		builder.invokeCookies = cookies
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

// build constructs the [http.Request] from the configured options.
func (builder *requestBuilder) build() (*http.Request, error) {
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
		builder.context,
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
	if builder.collection.BaseMockURL != "" {
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
		result[parameter.Name] = fmt.Sprintf("%v", value)
	}
	return result
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
	builder.applyInvokeOverrides(req)
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
		req.AddCookie(&http.Cookie{Name: c.Name, Value: c.Value, Domain: c.Domain, Path: c.Path, Secure: c.Secure, HttpOnly: c.HTTPOnly})
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

func (builder *requestBuilder) applyInvokeOverrides(req *http.Request) {
	for name, val := range builder.invokeHeaders {
		req.Header.Set(name, val)
	}

	for name, val := range builder.invokeCookies {
		req.AddCookie(&http.Cookie{Name: name, Value: val})
	}
}
