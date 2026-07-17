package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/model"
	"github.com/mmadfox/swag2mcp/internal/spec"
)

const (
	defaultMaxResponseSize = 1048    // 1 KB
	maxMaxResponseSize     = 1048576 // 1 MB
	randSuffixLen          = 6

	schemaTypeObject = "object"
	schemaTypeArray  = "array"
)

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

// Invoke validates the request, builds an HTTP request, sends it, and returns the response.
func (s *Service) Invoke(ctx context.Context, rq InvokeRequest) (InvokeResponse, error) {
	if err := s.validateRequest(rq); err != nil {
		return InvokeResponse{}, NewValidationError(
			"The endpoint ID is invalid — it must be a 32-character hex string. Use the search tool to find the correct endpoint ID.",
			err,
		)
	}

	if !s.disableRateLimiter.Load() {
		if err := s.rateLimiter.allow(rq.EndpointID); err != nil {
			return InvokeResponse{}, NewRateLimitError(err)
		}
	}

	ep, err := s.index.EndpointByID(rq.EndpointID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Endpoint %q not found — use the search tool to find the correct endpoint ID.",
				rq.EndpointID,
			),
			err,
		)
	}

	if ep.Operation == nil {
		return InvokeResponse{}, NewValidationError(
			"This endpoint has no operation definition — it may be malformed or incomplete.",
			nil,
		)
	}

	sp, err := s.index.SpecByID(ep.SpecID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Spec %q not found — the endpoint references a specification that no longer exists.",
				ep.SpecID,
			),
			err,
		)
	}

	coll, err := s.index.CollectionByID(ep.CollectionID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf(
				"Collection %q not found — the endpoint references a collection that no longer exists.",
				ep.CollectionID,
			),
			err,
		)
	}

	if err := validateParameters(ep.Operation, rq.Parameters); err != nil {
		return InvokeResponse{}, NewValidationError(
			"Parameter validation failed — check that all required parameters are provided and match the expected names.",
			err,
		)
	}

	if err := validateRequestBody(ep.Operation, rq.RequestBody); err != nil {
		return InvokeResponse{}, NewValidationError(
			"Request body validation failed — check that all required fields are present and no unknown fields are included.",
			err,
		)
	}

	req, err := newRequestBuilder(
		withContext(ctx),
		withSpec(sp),
		withCollection(coll),
		withEndpoint(ep),
		withParameters(rq.Parameters),
		withBody(rq.RequestBody),
		withHTTPConfig(mergeHTTPClientConfigs(sp.HTTPClient, coll.HTTPClient)),
	).build()
	if err != nil {
		return InvokeResponse{}, NewInvokeError(
			"Failed to build the HTTP request — check the endpoint parameters and try again.",
			err,
		)
	}

	s.dumpRequest(req, sp.Domain)

	client := s.httpClient
	if sp.Auth != nil {
		base := s.httpClient.Transport
		if base == nil {
			base = http.DefaultTransport
		}
		client = &http.Client{
			Transport: &auth.Transport{
				Base: base,
				Auth: sp.Auth,
			},
			Timeout:       s.httpClient.Timeout,
			CheckRedirect: s.httpClient.CheckRedirect,
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return InvokeResponse{}, NewInvokeError(
			"The API request failed — the server may be unreachable or returned an error.",
			err,
		)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return InvokeResponse{}, NewInvokeError(
			"Failed to read the API response — the connection may have been interrupted.",
			err,
		)
	}

	maxSize := s.maxResponseSize
	if len(body) > maxSize {
		return s.saveLargeResponse(response, body, sp.Domain, ep, maxSize)
	}

	return newInvokeResponse(response, body), nil
}

// requestBuilder builds an [http.Request] from spec, collection, endpoint, and parameters.
type requestBuilder struct {
	context    context.Context
	spec       *model.Spec
	collection *model.Collection
	endpoint   *model.Endpoint
	parameters map[string]any
	body       map[string]any
	httpConfig *model.HTTPClientConfig
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

// withContext sets the context for the request.
func withContext(ctx context.Context) requestOption {
	return func(builder *requestBuilder) {
		builder.context = ctx
	}
}

// withSpec sets the API specification.
func withSpec(specification *model.Spec) requestOption {
	return func(builder *requestBuilder) {
		builder.spec = specification
	}
}

// withCollection sets the collection.
func withCollection(collection *model.Collection) requestOption {
	return func(builder *requestBuilder) {
		builder.collection = collection
	}
}

// withEndpoint sets the endpoint.
func withEndpoint(endpoint *model.Endpoint) requestOption {
	return func(builder *requestBuilder) {
		builder.endpoint = endpoint
	}
}

// withParameters sets the request parameters.
func withParameters(parameters map[string]any) requestOption {
	return func(builder *requestBuilder) {
		builder.parameters = parameters
	}
}

// withBody sets the request body.
func withBody(body map[string]any) requestOption {
	return func(builder *requestBuilder) {
		builder.body = body
	}
}

// withHTTPConfig sets the HTTP client configuration.
func withHTTPConfig(config *model.HTTPClientConfig) requestOption {
	return func(builder *requestBuilder) {
		builder.httpConfig = config
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

// resolveBaseURL returns the base URL, preferring the collection's over the spec's.
// If BaseMockURL is set on the collection, it is used with http:// prefix.
func (builder *requestBuilder) resolveBaseURL() string {
	if builder.collection.BaseMockURL != "" {
		return "http://" + builder.collection.BaseMockURL
	}
	if builder.collection.BaseURL != "" {
		return builder.collection.BaseURL
	}
	return builder.spec.BaseURL
}

// filterParametersByLocation returns parameters that match the given location (path, query, header).
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

// applyHeaders sets operation-level headers and defaults on the request.
func (builder *requestBuilder) applyHeaders(req *http.Request) {
	headers := builder.filterParametersByLocation("header")
	for name, val := range headers {
		req.Header.Set(name, val)
	}

	if builder.body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if req.Header.Get("Accept") == "" {
		isJSON := builder.body != nil ||
			req.Header.Get("Content-Type") == "application/json"
		if isJSON {
			req.Header.Set("Accept", "application/json, text/plain, */*")
		} else {
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		}
	}
}

// applyHTTPClientConfig applies per-request HTTP config (headers, cookies) to the request.
func (builder *requestBuilder) applyHTTPClientConfig(req *http.Request) {
	if builder.httpConfig == nil {
		return
	}

	for name, val := range builder.httpConfig.Headers {
		req.Header.Set(name, val)
	}

	if len(builder.httpConfig.Cookies) > 0 {
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
}

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

// saveLargeResponse saves a response body that exceeds the max size to a file
// and returns an InvokeResponse with a FileReference instead of the full body.
func (s *Service) saveLargeResponse(
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
	suf := randomSuffix(randSuffixLen)
	fname := fmt.Sprintf("%s-%s-%s-%s.json", domain, m, p, suf)
	fp := filepath.Join(s.ws.ResponsesDir(), fname)

	if err := os.MkdirAll(s.ws.ResponsesDir(), 0750); err != nil {
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

// resolveMaxResponseSize returns the effective max response size.
// Default is 2 KB, maximum is 1 MB.
func resolveMaxResponseSize(maxResponseSize *int) int {
	if maxResponseSize == nil {
		return defaultMaxResponseSize
	}
	if *maxResponseSize > maxMaxResponseSize {
		return maxMaxResponseSize
	}
	if *maxResponseSize <= 0 {
		return defaultMaxResponseSize
	}
	return *maxResponseSize
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
	byteLen := (n + 1) / 2 //nolint:mnd // Hex encoding uses 2 characters per byte.
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%0*x", n, 0)
	}
	return hex.EncodeToString(b)[:n]
}

// validateParameters checks that all required parameters are present and that no
// unknown parameters are passed. Every parameter must be declared in the operation spec.
func validateParameters(op *spec.Operation, params map[string]any) error {
	paramNames := make(map[string]struct{}, len(op.Parameters))
	for _, p := range op.Parameters {
		paramNames[p.Name] = struct{}{}
	}

	for name := range params {
		if _, exists := paramNames[name]; !exists {
			return fmt.Errorf(
				"unknown parameter %q, all parameters must match the operation schema",
				name,
			)
		}
	}

	var missing []string
	for _, p := range op.Parameters {
		if !p.Required {
			continue
		}
		if _, exists := params[p.Name]; !exists {
			missing = append(missing, p.Name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"missing required parameters: %s",
			strings.Join(missing, ", "),
		)
	}

	return nil
}

// validateRequestBody validates a request body against the operation's request body schema.
// It checks that all required properties are present and that no unknown keys are passed.
// Type validation is not performed.
func validateRequestBody(op *spec.Operation, body map[string]any) error {
	if op.RequestBody == nil {
		return nil
	}

	if op.RequestBody.Required && body == nil {
		return errors.New("request body is required for this endpoint")
	}

	if body == nil {
		return nil
	}

	sc := schemaForContentType(op.RequestBody.Content)
	if sc == nil {
		return nil
	}

	return validateSchemaValue(sc, body, "$")
}

// schemaForContentType extracts the JSON schema from a content map, preferring application/json.
func schemaForContentType(ct map[string]*spec.MediaType) *spec.Schema {
	if ct == nil {
		return nil
	}
	mt, exists := ct["application/json"]
	if !exists || mt == nil {
		return nil
	}
	return mt.Schema
}

// validateSchemaValue recursively validates a value against a schema path.
// It is used for request body validation.
func validateSchemaValue(sc *spec.Schema, value any, path string) error {
	if sc == nil {
		return nil
	}

	switch sc.Type {
	case schemaTypeObject:
		return validateObjectSchema(sc, value, path)
	case schemaTypeArray:
		return validateArraySchema(sc, value, path)
	}

	return nil
}

// validateObjectSchema validates a map value against an object schema.
func validateObjectSchema(sc *spec.Schema, value any, path string) error {
	obj, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	for _, requiredField := range sc.Required {
		if _, exists := obj[requiredField]; !exists {
			return fmt.Errorf("missing required field %q at %s", requiredField, path)
		}
	}

	for key := range obj {
		if _, defined := sc.Properties[key]; !defined {
			return fmt.Errorf(
				"unknown field %q at %s, all fields must match the schema",
				key, path,
			)
		}
	}

	for key, ps := range sc.Properties {
		cv, exists := obj[key]
		if !exists {
			continue
		}
		cp := path + "." + key
		if err := validateSchemaValue(ps, cv, cp); err != nil {
			return err
		}
	}

	return nil
}

// validateArraySchema validates a slice value against an array schema.
func validateArraySchema(sc *spec.Schema, value any, path string) error {
	arr, ok := value.([]any)
	if !ok {
		return nil
	}

	for i, item := range arr {
		cp := fmt.Sprintf("%s[%d]", path, i)
		if err := validateSchemaValue(sc.Items, item, cp); err != nil {
			return err
		}
	}

	return nil
}

// dumpRequest writes the HTTP request to a file for debugging if dumpDir is configured.
func (s *Service) dumpRequest(req *http.Request, domain string) {
	if len(s.dumpDir) == 0 {
		return
	}

	d, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return
	}

	ts := time.Now().UnixMilli()
	fname := fmt.Sprintf("invoke-%s-%d.txt", domain, ts)
	fp := filepath.Join(s.dumpDir, fname)

	if err := os.MkdirAll(s.dumpDir, 0750); err != nil {
		slog.Default().WarnContext(req.Context(), "failed to create dump dir", "error", err)
		return
	}
	if err := os.WriteFile(fp, d, 0600); err != nil {
		slog.Default().WarnContext(req.Context(), "failed to write dump file", "error", err)
	}
}

// mergeHTTPClientConfigs merges two per-request HTTP configs. Collection overrides spec.
func mergeHTTPClientConfigs(sp, coll *model.HTTPClientConfig) *model.HTTPClientConfig {
	if sp == nil {
		return coll
	}
	if coll == nil {
		return sp
	}

	result := &model.HTTPClientConfig{
		Headers: make(map[string]string),
		Cookies: coll.Cookies,
	}

	maps.Copy(result.Headers, sp.Headers)
	maps.Copy(result.Headers, coll.Headers)

	if len(result.Cookies) == 0 {
		result.Cookies = sp.Cookies
	}

	return result
}
