package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

// InvokeRequest represents a request to invoke an API endpoint.
type InvokeRequest struct {
	EndpointID  string         `json:"endpointId"            validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to invoke"`
	Parameters  map[string]any `json:"parameters,omitempty"                          jsonschema:"optional,Path, query, and header parameters as key-value pairs"`
	RequestBody map[string]any `json:"requestBody,omitempty"                         jsonschema:"optional,Request body for POST/PUT/PATCH requests"`
}

// InvokeResponse represents the response from invoking an API endpoint.
type InvokeResponse struct {
	StatusCode int               `json:"statusCode" jsonschema:"required,HTTP response status code"`
	Headers    map[string]string `json:"headers"    jsonschema:"required,HTTP response headers"`
	Body       any               `json:"body"       jsonschema:"required,Response body data"`
}

// Invoke validates the request, builds an HTTP request, sends it, and returns the response.
func (s *Service) Invoke(ctx context.Context, request InvokeRequest) (InvokeResponse, error) {
	if err := s.validateRequest(request); err != nil {
		return InvokeResponse{}, NewValidationError(
			"endpointId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	endpoint, err := s.index.EndpointByID(request.EndpointID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf("endpoint %q not found", request.EndpointID), err,
		)
	}

	if endpoint.Operation == nil {
		return InvokeResponse{}, NewValidationError("endpoint has no operation definition", nil)
	}

	specification, err := s.index.SpecByID(endpoint.SpecID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf("spec %q not found", endpoint.SpecID), err,
		)
	}

	collection, err := s.index.CollectionByID(endpoint.CollectionID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(
			fmt.Sprintf("collection %q not found", endpoint.CollectionID), err,
		)
	}

	if validationError := validateParameters(endpoint.Operation, request.Parameters); validationError != nil {
		return InvokeResponse{}, NewValidationError("parameter validation failed", validationError)
	}

	if validationError := validateRequestBody(endpoint.Operation, request.RequestBody); validationError != nil {
		return InvokeResponse{}, NewValidationError("request body validation failed", validationError)
	}

	httpRequest, buildError := newRequestBuilder(
		withContext(ctx),
		withSpec(specification),
		withCollection(collection),
		withEndpoint(endpoint),
		withParameters(request.Parameters),
		withBody(request.RequestBody),
		withHTTPConfig(mergeHTTPClientConfigs(specification.HTTPClient, collection.HTTPClient)),
	).build()
	if buildError != nil {
		return InvokeResponse{}, fmt.Errorf("failed to build request: %w", buildError)
	}

	s.dumpRequest(httpRequest, specification.Domain)

	httpClient := newAuthHTTPClient(specification, mergeHTTPClientConfigs(specification.HTTPClient, collection.HTTPClient))
	response, doError := httpClient.Do(httpRequest)
	if doError != nil {
		return InvokeResponse{}, fmt.Errorf("request failed: %w", doError)
	}
	defer response.Body.Close()

	return newInvokeResponse(response)
}

// requestBuilder builds an [http.Request] from spec, collection, endpoint, and parameters.
type requestBuilder struct {
	context    context.Context
	spec       *types.Spec
	collection *types.Collection
	endpoint   *types.Endpoint
	parameters map[string]any
	body       map[string]any
	httpConfig *types.HTTPClientConfig
}

// requestOption is a functional option for configuring a requestBuilder.
type requestOption func(*requestBuilder)

// newRequestBuilder creates a new requestBuilder with the given options.
func newRequestBuilder(options ...requestOption) *requestBuilder {
	builder := &requestBuilder{
		context: context.Background(),
	}
	for _, option := range options {
		option(builder)
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
func withSpec(specification *types.Spec) requestOption {
	return func(builder *requestBuilder) {
		builder.spec = specification
	}
}

// withCollection sets the collection.
func withCollection(collection *types.Collection) requestOption {
	return func(builder *requestBuilder) {
		builder.collection = collection
	}
}

// withEndpoint sets the endpoint.
func withEndpoint(endpoint *types.Endpoint) requestOption {
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
func withHTTPConfig(config *types.HTTPClientConfig) requestOption {
	return func(builder *requestBuilder) {
		builder.httpConfig = config
	}
}

// build constructs the [http.Request] from the configured options.
func (builder *requestBuilder) build() (*http.Request, error) {
	targetURL := builder.resolveBaseURL()
	targetURL = strings.TrimRight(targetURL, "/")
	requestURL := targetURL + "/" + strings.TrimLeft(builder.endpoint.Path, "/")

	pathParameters := builder.filterParametersByLocation("path")
	for parameterName, parameterValue := range pathParameters {
		requestURL = strings.ReplaceAll(
			requestURL,
			"{"+parameterName+"}",
			url.PathEscape(parameterValue),
		)
	}

	parsedURL, parseError := url.Parse(requestURL)
	if parseError != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", requestURL, parseError)
	}

	queryParameters := builder.filterParametersByLocation("query")
	queryValues := parsedURL.Query()
	for parameterName, parameterValue := range queryParameters {
		queryValues.Set(parameterName, parameterValue)
	}
	parsedURL.RawQuery = queryValues.Encode()

	var bodyReader io.Reader
	if builder.body != nil {
		bodyBytes, marshalError := json.Marshal(builder.body)
		if marshalError != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", marshalError)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpRequest, requestError := http.NewRequestWithContext(
		builder.context,
		builder.endpoint.Name,
		parsedURL.String(),
		bodyReader,
	)
	if requestError != nil {
		return nil, fmt.Errorf("failed to create request: %w", requestError)
	}

	builder.applyHeaders(httpRequest)
	builder.applyHTTPClientConfig(httpRequest)

	return httpRequest, nil
}

// resolveBaseURL returns the base URL, preferring the collection's over the spec's.
func (builder *requestBuilder) resolveBaseURL() string {
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
func (builder *requestBuilder) applyHeaders(httpRequest *http.Request) {
	headerParameters := builder.filterParametersByLocation("header")
	for parameterName, parameterValue := range headerParameters {
		httpRequest.Header.Set(parameterName, parameterValue)
	}

	if builder.body != nil && httpRequest.Header.Get("Content-Type") == "" {
		httpRequest.Header.Set("Content-Type", "application/json")
	}

	if httpRequest.Header.Get("Accept") == "" {
		httpRequest.Header.Set("Accept", "application/json")
	}
}

// applyHTTPClientConfig applies HTTP client config (headers, cookies, user-agent) to the request.
func (builder *requestBuilder) applyHTTPClientConfig(httpRequest *http.Request) {
	if builder.httpConfig == nil {
		return
	}

	for headerName, headerValue := range builder.httpConfig.Headers {
		httpRequest.Header.Set(headerName, headerValue)
	}

	if builder.httpConfig.UserAgent != "" {
		httpRequest.Header.Set("User-Agent", builder.httpConfig.UserAgent)
	}

	if len(builder.httpConfig.Cookies) > 0 {
		for _, cookie := range builder.httpConfig.Cookies {
			//nolint:gosec // cookies are user-configured, not secrets
			httpRequest.AddCookie(&http.Cookie{
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
func newInvokeResponse(response *http.Response) (InvokeResponse, error) {
	body, readError := io.ReadAll(response.Body)
	if readError != nil {
		return InvokeResponse{}, fmt.Errorf("failed to read response: %w", readError)
	}

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
			}, nil
		}
	}

	return InvokeResponse{
		StatusCode: response.StatusCode,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

// newAuthHTTPClient returns an [http.Client] that applies authentication and
// HTTP client config (timeout, redirects) to every request.
func newAuthHTTPClient(specification *types.Spec, httpConfig *types.HTTPClientConfig) *http.Client {
	client := &http.Client{}

	if specification.Auth != nil {
		client.Transport = &auth.Transport{
			Base: http.DefaultTransport,
			Auth: specification.Auth,
		}
	}

	if httpConfig != nil {
		applyHTTPClientTimeout(client, httpConfig)
		applyHTTPClientRedirects(client, httpConfig)
	}

	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}

	return client
}

// applyHTTPClientTimeout sets the timeout on the client if configured.
func applyHTTPClientTimeout(client *http.Client, config *types.HTTPClientConfig) {
	if config.Timeout > 0 {
		client.Timeout = config.Timeout
	}
}

// applyHTTPClientRedirects sets redirect behavior on the client if configured.
func applyHTTPClientRedirects(client *http.Client, config *types.HTTPClientConfig) {
	if config.FollowRedirects != nil && !*config.FollowRedirects {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if config.MaxRedirects != nil {
		maxRedirects := *config.MaxRedirects
		client.CheckRedirect = func(_ *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("too many redirects (max %d)", maxRedirects)
			}
			return nil
		}
	}
}

// validateParameters checks that all required parameters are present and that no
// unknown parameters are passed. Every parameter must be declared in the operation spec.
func validateParameters(operation *spec.Operation, parameters map[string]any) error {
	declaredParameterNames := make(map[string]struct{}, len(operation.Parameters))
	for _, parameter := range operation.Parameters {
		declaredParameterNames[parameter.Name] = struct{}{}
	}

	for parameterName := range parameters {
		if _, exists := declaredParameterNames[parameterName]; !exists {
			return fmt.Errorf(
				"unknown parameter %q, all parameters must match the operation schema",
				parameterName,
			)
		}
	}

	var missingRequiredParameters []string
	for _, parameter := range operation.Parameters {
		if !parameter.Required {
			continue
		}
		if _, exists := parameters[parameter.Name]; !exists {
			missingRequiredParameters = append(missingRequiredParameters, parameter.Name)
		}
	}

	if len(missingRequiredParameters) > 0 {
		return fmt.Errorf(
			"missing required parameters: %s",
			strings.Join(missingRequiredParameters, ", "),
		)
	}

	return nil
}

// validateRequestBody validates a request body against the operation's request body schema.
// It checks that all required properties are present and that no unknown keys are passed.
// Type validation is not performed.
func validateRequestBody(operation *spec.Operation, body map[string]any) error {
	if operation.RequestBody == nil {
		return nil
	}

	if operation.RequestBody.Required && body == nil {
		return errors.New("request body is required for this endpoint")
	}

	if body == nil {
		return nil
	}

	schema := schemaForContentType(operation.RequestBody.Content)
	if schema == nil {
		return nil
	}

	return validateSchemaValue(schema, body, "$")
}

// schemaForContentType extracts the JSON schema from a content map, preferring application/json.
func schemaForContentType(content map[string]*spec.MediaType) *spec.Schema {
	if content == nil {
		return nil
	}
	mediaType, exists := content["application/json"]
	if !exists || mediaType == nil {
		return nil
	}
	return mediaType.Schema
}

// validateSchemaValue recursively validates a value against a schema path.
// It is used for request body validation.
func validateSchemaValue(schema *spec.Schema, value any, path string) error {
	if schema == nil {
		return nil
	}

	switch schema.Type {
	case "object":
		return validateObjectSchema(schema, value, path)
	case "array":
		return validateArraySchema(schema, value, path)
	}

	return nil
}

// validateObjectSchema validates a map value against an object schema.
func validateObjectSchema(schema *spec.Schema, value any, path string) error {
	objectValue, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	for _, requiredField := range schema.Required {
		if _, exists := objectValue[requiredField]; !exists {
			return fmt.Errorf("missing required field %q at %s", requiredField, path)
		}
	}

	for key := range objectValue {
		if _, defined := schema.Properties[key]; !defined {
			return fmt.Errorf(
				"unknown field %q at %s, all fields must match the schema",
				key, path,
			)
		}
	}

	for key, propertySchema := range schema.Properties {
		childValue, exists := objectValue[key]
		if !exists {
			continue
		}
		childPath := path + "." + key
		if validationError := validateSchemaValue(propertySchema, childValue, childPath); validationError != nil {
			return validationError
		}
	}

	return nil
}

// validateArraySchema validates a slice value against an array schema.
func validateArraySchema(schema *spec.Schema, value any, path string) error {
	arrayValue, ok := value.([]any)
	if !ok {
		return nil
	}

	for index, item := range arrayValue {
		childPath := fmt.Sprintf("%s[%d]", path, index)
		if validationError := validateSchemaValue(schema.Items, item, childPath); validationError != nil {
			return validationError
		}
	}

	return nil
}

// dumpRequest writes the HTTP request to a file for debugging if dumpDir is configured.
func (s *Service) dumpRequest(request *http.Request, domain string) {
	if len(s.dumpDir) == 0 {
		return
	}

	dump, dumpError := httputil.DumpRequestOut(request, true)
	if dumpError != nil {
		return
	}

	timestamp := time.Now().UnixMilli()
	filename := fmt.Sprintf("invoke-%s-%d.txt", domain, timestamp)
	filePath := filepath.Join(s.dumpDir, filename)

	_ = os.MkdirAll(s.dumpDir, 0750)
	_ = os.WriteFile(filePath, dump, 0600)
}

// mergeHTTPClientConfigs merges two HTTP client configs. Collection overrides spec.
func mergeHTTPClientConfigs(spec, collection *types.HTTPClientConfig) *types.HTTPClientConfig {
	if spec == nil {
		return collection
	}
	if collection == nil {
		return spec
	}

	result := &types.HTTPClientConfig{
		Headers:         make(map[string]string),
		Cookies:         collection.Cookies,
		UserAgent:       collection.UserAgent,
		Timeout:         collection.Timeout,
		FollowRedirects: collection.FollowRedirects,
		MaxRedirects:    collection.MaxRedirects,
	}

	maps.Copy(result.Headers, spec.Headers)
	maps.Copy(result.Headers, collection.Headers)

	if result.UserAgent == "" {
		result.UserAgent = spec.UserAgent
	}
	if result.Timeout == 0 {
		result.Timeout = spec.Timeout
	}
	if result.FollowRedirects == nil {
		result.FollowRedirects = spec.FollowRedirects
	}
	if result.MaxRedirects == nil {
		result.MaxRedirects = spec.MaxRedirects
	}
	if len(result.Cookies) == 0 {
		result.Cookies = spec.Cookies
	}

	return result
}
