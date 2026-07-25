package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/auth"
	"github.com/mmadfox/swag2mcp/internal/spec"
	"github.com/mmadfox/swag2mcp/internal/types"
)

type (
	// InvokeRequest represents a request to invoke an endpoint.
	InvokeRequest struct {
		EndpointID  string         `json:"endpointId"            validate:"required,md5" jsonschema:"required,The 32-character MD5 hash ID of the endpoint to invoke"`
		Parameters  map[string]any `json:"parameters,omitempty"                          jsonschema:"optional,Path, query, and header parameters as key-value pairs"`
		RequestBody map[string]any `json:"requestBody,omitempty"                         jsonschema:"optional,Request body for POST/PUT/PATCH requests"`
	}

	// InvokeResponse represents a response to invoke an endpoint.
	InvokeResponse struct {
		StatusCode int               `json:"statusCode" jsonschema:"required,HTTP response status code"`
		Headers    map[string]string `json:"headers"    jsonschema:"required,HTTP response headers"`
		Body       any               `json:"body"       jsonschema:"required,Response body data"`
	}
)

// Invoke invokes an endpoint.
// It validates the request, builds the HTTP request, and sends it.
func (s *Service) Invoke(ctx context.Context, req InvokeRequest) (InvokeResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return InvokeResponse{}, NewValidationError(
			"endpointId must be a 32-character lowercase hex string (MD5 format)",
			err,
		)
	}

	ep, err := s.index.EndpointByID(req.EndpointID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(fmt.Sprintf("endpoint %q not found", req.EndpointID), err)
	}

	if ep.Operation == nil {
		return InvokeResponse{}, NewValidationError("endpoint has no operation definition", nil)
	}

	spec, err := s.index.SpecByID(ep.SpecID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(fmt.Sprintf("spec %q not found", ep.SpecID), err)
	}

	collection, err := s.index.CollectionByID(ep.CollectionID)
	if err != nil {
		return InvokeResponse{}, NewNotFoundError(fmt.Sprintf("collection %q not found", ep.CollectionID), err)
	}

	if verr := validateParams(ep.Operation, req.Parameters); verr != nil {
		return InvokeResponse{}, NewValidationError("parameter validation failed", verr)
	}

	if verr := validateRequestBody(ep.Operation, req.RequestBody); verr != nil {
		return InvokeResponse{}, NewValidationError("request body validation failed", verr)
	}

	httpReq, buildErr := buildHTTPRequest(ctx, spec, collection, ep, req.Parameters, req.RequestBody)
	if buildErr != nil {
		return InvokeResponse{}, fmt.Errorf("failed to build request: %w", buildErr)
	}

	client := authHTTPClient(spec)
	resp, doErr := client.Do(httpReq)
	if doErr != nil {
		return InvokeResponse{}, fmt.Errorf("request failed: %w", doErr)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return InvokeResponse{}, fmt.Errorf("failed to read response: %w", readErr)
	}

	headers := make(map[string]string, len(resp.Header))
	for k, vals := range resp.Header {
		headers[k] = strings.Join(vals, ", ")
	}

	var bodyAny any
	if len(body) > 0 {
		var parsed any
		if jsonErr := json.Unmarshal(body, &parsed); jsonErr == nil {
			bodyAny = parsed
		} else {
			bodyAny = string(body)
		}
	}

	return InvokeResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       bodyAny,
	}, nil
}

// validateParams checks that all required parameters are present and that no
// unknown parameters are passed. Every parameter must be declared in the operation spec.
func validateParams(op *spec.Operation, params map[string]any) error {
	schemaParamNames := make(map[string]struct{}, len(op.Parameters))
	for _, p := range op.Parameters {
		schemaParamNames[p.Name] = struct{}{}
	}

	for name := range params {
		if _, ok := schemaParamNames[name]; !ok {
			return fmt.Errorf("unknown parameter %q, all parameters must match the operation schema", name)
		}
	}

	var missing []string
	for _, p := range op.Parameters {
		if !p.Required {
			continue
		}
		if _, ok := params[p.Name]; !ok {
			missing = append(missing, p.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required parameters: %s", strings.Join(missing, ", "))
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

	schema := schemaForContent(op.RequestBody.Content)
	if schema == nil {
		return nil
	}

	return validateSchemaValue(schema, body, "$")
}

// schemaForContent extracts the JSON schema from a content map, preferring application/json.
func schemaForContent(content map[string]*spec.MediaType) *spec.Schema {
	if content == nil {
		return nil
	}
	mt, ok := content["application/json"]
	if !ok || mt == nil {
		return nil
	}
	return mt.Schema
}

// validateSchemaValue recursively validates a value against a schema path.
// It is used for request body validation.
//
//nolint:gocognit // recursive schema validation is inherently complex
func validateSchemaValue(schema *spec.Schema, val any, path string) error {
	if schema == nil {
		return nil
	}

	switch schema.Type {
	case "object":
		m, ok := val.(map[string]any)
		if !ok {
			// skip type validation per user request
			return nil
		}

		for _, reqField := range schema.Required {
			if _, exists := m[reqField]; !exists {
				return fmt.Errorf("missing required field %q at %s", reqField, path)
			}
		}

		for key := range m {
			if _, defined := schema.Properties[key]; !defined {
				return fmt.Errorf("unknown field %q at %s, all fields must match the schema", key, path)
			}
		}

		for key, propSchema := range schema.Properties {
			if childVal, exists := m[key]; exists {
				childPath := path + "." + key
				if err := validateSchemaValue(propSchema, childVal, childPath); err != nil {
					return err
				}
			}
		}

	case "array":
		arr, ok := val.([]any)
		if !ok {
			return nil
		}
		for i, item := range arr {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			if err := validateSchemaValue(schema.Items, item, childPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func authHTTPClient(spec *types.Spec) *http.Client {
	if spec.Auth != nil {
		return auth.NewHTTPClient(spec.Auth)
	}
	return http.DefaultClient
}

// TODO:
//
//nolint:nolintlint,gocognit
func buildHTTPRequest(
	ctx context.Context,
	spec *types.Spec,
	collection *types.Collection,
	ep *types.Endpoint,
	params map[string]any,
	requestBody map[string]any,
) (*http.Request, error) {
	baseURL := spec.BaseURL
	if len(collection.BaseURL) > 0 {
		baseURL = collection.BaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	reqURL := baseURL + "/" + strings.TrimLeft(ep.Path, "/")

	paramByIn := func(in string) map[string]string {
		m := make(map[string]string, len(params))
		for _, p := range ep.Operation.Parameters {
			if p.In != in {
				continue
			}
			val, ok := params[p.Name]
			if !ok {
				continue
			}
			m[p.Name] = fmt.Sprintf("%v", val)
		}
		return m
	}

	// substitute path parameters
	// TODO:
	pathParams := paramByIn("path")
	for name, val := range pathParams {
		reqURL = strings.ReplaceAll(reqURL, "{"+name+"}", url.PathEscape(val))
	}

	parsedURL, parseErr := url.Parse(reqURL)
	if parseErr != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", reqURL, parseErr)
	}

	// add query parameters
	q := parsedURL.Query()
	for name, val := range paramByIn("query") {
		q.Set(name, val)
	}
	parsedURL.RawQuery = q.Encode()

	// build body
	var bodyReader io.Reader
	if requestBody != nil {
		bodyBytes, marshalErr := json.Marshal(requestBody)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", marshalErr)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq, reqErr := http.NewRequestWithContext(ctx, ep.Name, parsedURL.String(), bodyReader)
	if reqErr != nil {
		return nil, fmt.Errorf("failed to create request: %w", reqErr)
	}

	// apply spec-level headers
	for k, v := range spec.Headers {
		httpReq.Header.Set(k, v)
	}

	// apply collection-level headers
	for k, v := range collection.Headers {
		httpReq.Header.Set(k, v)
	}

	// apply header parameters from operation
	for name, val := range paramByIn("header") {
		httpReq.Header.Set(name, val)
	}

	// apply content-type for json body
	if requestBody != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	return httpReq, nil
}
