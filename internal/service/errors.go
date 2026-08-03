/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mmadfox/swag2mcp/internal/reader"
)

const (
	validationFailedErrCode            = "validation_failed"
	invalidSpecIDErrCode               = "invalid_spec_id"
	invalidCollectionIDErrCode         = "invalid_collection_id"
	invalidTagIDErrCode                = "invalid_tag_id"
	invalidEndpointIDErrCode           = "invalid_endpoint_id"
	parameterValidationFailedErrCode   = "parameter_validation_failed"
	requestBodyValidationFailedErrCode = "request_body_validation_failed"
	searchQueryErrCode                 = "search_query_error"
	importSourceErrCode                = "import_source_error"
	notFoundErrCode                    = "not_found"
	endpointNotFoundErrCode            = "endpoint_not_found"
	specNotFoundErrCode                = "spec_not_found"
	collectionNotFoundErrCode          = "collection_not_found"
	tagNotFoundErrCode                 = "tag_not_found"
	searchNoResultsErrCode             = "search_no_results"
	exportNoCollectionsErrCode         = "export_no_collections"
	importNoMatchErrCode               = "import_no_match"
	configNotFoundErrCode              = "config_not_found"
	responseFileNotFoundErrCode        = "response_file_not_found"
	responsePathNotFoundErrCode        = "response_path_not_found"
	rateLimitErrCode                   = "rate_limit"
	globalRateLimitErrCode             = "global_rate_limit"
	invokeErrorErrCode                 = "invoke_error"
	httpRequestErrCode                 = "http_request_error"
	responseReadErrCode                = "response_read_error"
	fileWriteErrCode                   = "file_write_error"
	fileCreateErrCode                  = "file_create_error"
	dirCreateErrCode                   = "dir_create_error"
	streamErrCode                      = "stream_error"
	buildRequestErrCode                = "build_request_error"
	configErrorErrCode                 = "config_error"
	workspaceErrorErrCode              = "workspace_error"
	parseErrorErrCode                  = "parse_error"
	authErrorErrCode                   = "auth_error"
	responseRequestErrCode             = "response_request_error"
	responseJSONPathErrCode            = "response_jsonpath_error"
	responseLineRangeErrCode           = "response_line_range_error"
	responseNotJSONErrCode             = "response_not_json"
	responseReadFailedErrCode          = "response_read_failed"
	exportErrorErrCode                 = "export_error"
	importErrorErrCode                 = "import_error"
	toolLoadErrCode                    = "tool_load_error"
)

// LLMError is an error type returned to the LLM with a machine-readable code and human-readable message.
type LLMError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

// NewValidationError creates an LLMError with code "validation_failed".
func NewValidationError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    validationFailedErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewParameterValidationError creates an LLMError with code "parameter_validation_failed".
func NewParameterValidationError(err error) *LLMError {
	return &LLMError{
		Code:    parameterValidationFailedErrCode,
		Message: "Parameter validation failed. Check that all required parameters are provided and match the expected names.",
		Hint:    formatError(err),
	}
}

// NewRequestBodyValidationError creates an LLMError with code "request_body_validation_failed".
func NewRequestBodyValidationError(err error) *LLMError {
	return &LLMError{
		Code:    requestBodyValidationFailedErrCode,
		Message: "Request body validation failed. Check that all required fields are present and no unknown fields are included.",
		Hint:    formatError(err),
	}
}

// NewNotFoundError creates an LLMError with code "not_found".
func NewNotFoundError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    notFoundErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewEndpointNotFoundError creates an LLMError with code "endpoint_not_found".
func NewEndpointNotFoundError(endpointID string, _ error) *LLMError {
	return &LLMError{
		Code:    endpointNotFoundErrCode,
		Message: fmt.Sprintf("Endpoint %q was not found.", endpointID),
		Hint:    "Use the search tool to find the correct endpoint ID.",
	}
}

// NewSpecNotFoundError creates an LLMError with code "spec_not_found".
func NewSpecNotFoundError(specID string, _ error) *LLMError {
	return &LLMError{
		Code:    specNotFoundErrCode,
		Message: fmt.Sprintf("Spec %q was not found.", specID),
		Hint:    "Use spec_list to see all available specs.",
	}
}

// NewCollectionNotFoundError creates an LLMError with code "collection_not_found".
func NewCollectionNotFoundError(collectionID string, _ error) *LLMError {
	return &LLMError{
		Code:    collectionNotFoundErrCode,
		Message: fmt.Sprintf("Collection %q was not found.", collectionID),
		Hint:    "Use collection_by_spec to list collections.",
	}
}

// NewTagNotFoundError creates an LLMError with code "tag_not_found".
func NewTagNotFoundError(tagID string, _ error) *LLMError {
	return &LLMError{
		Code:    tagNotFoundErrCode,
		Message: fmt.Sprintf("Tag %q was not found.", tagID),
		Hint:    "Use tag_by_collection or tag_by_spec to list tags.",
	}
}

// NewInvalidSpecIDError creates an LLMError with code "invalid_spec_id".
func NewInvalidSpecIDError(err error) *LLMError {
	return &LLMError{
		Code:    invalidSpecIDErrCode,
		Message: "The spec ID is invalid. It must be a 32-character hex string. Use spec_list to find the correct spec ID.",
		Hint:    formatError(err),
	}
}

// NewInvalidCollectionIDError creates an LLMError with code "invalid_collection_id".
func NewInvalidCollectionIDError(err error) *LLMError {
	return &LLMError{
		Code:    invalidCollectionIDErrCode,
		Message: "The collection ID is invalid. It must be a 32-character hex string.",
		Hint:    formatError(err),
	}
}

// NewInvalidTagIDError creates an LLMError with code "invalid_tag_id".
func NewInvalidTagIDError(err error) *LLMError {
	return &LLMError{
		Code:    invalidTagIDErrCode,
		Message: "The tag ID is invalid. It must be a 32-character hex string. Use tag_by_collection or tag_by_spec to find the correct tag ID.",
		Hint:    formatError(err),
	}
}

// NewInvalidEndpointIDError creates an LLMError with code "invalid_endpoint_id".
func NewInvalidEndpointIDError(err error) *LLMError {
	return &LLMError{
		Code:    invalidEndpointIDErrCode,
		Message: "The endpoint ID is invalid. It must be a 32-character hex string. Use the search tool to find the correct endpoint ID.",
		Hint:    formatError(err),
	}
}

// NewSearchQueryError creates an LLMError with code "search_query_error".
func NewSearchQueryError(err error) *LLMError {
	return &LLMError{
		Code:    searchQueryErrCode,
		Message: "The search query format is invalid. Use simple text search (e.g. 'get pet') or field search (e.g. method:GET, tag:pets).",
		Hint:    formatError(err),
	}
}

// NewSearchNoResultsError creates an LLMError with code "search_no_results".
func NewSearchNoResultsError() *LLMError {
	return &LLMError{
		Code:    searchNoResultsErrCode,
		Message: "The search query did not match any endpoints. Try a different query.",
	}
}

// NewImportSourceError creates an LLMError with code "import_source_error".
func NewImportSourceError(err error) *LLMError {
	return &LLMError{
		Code:    importSourceErrCode,
		Message: "Import requires a source URL with filename, a spec filter, or a zip backup path.",
		Hint:    formatError(err),
	}
}

// NewExportNoCollectionsError creates an LLMError with code "export_no_collections".
func NewExportNoCollectionsError() *LLMError {
	return &LLMError{
		Code:    exportNoCollectionsErrCode,
		Message: "No collections found to export. Ensure the workspace has specs with valid collections.",
	}
}

// NewImportNoMatchError creates an LLMError with code "import_no_match".
func NewImportNoMatchError(filter string) *LLMError {
	return &LLMError{
		Code:    importNoMatchErrCode,
		Message: fmt.Sprintf("No matching specs found for filter %q.", filter),
	}
}

// NewImportSpecNotFoundError creates an LLMError with code "import_no_match"
// for a spec domain that does not exist in the config.
func NewImportSpecNotFoundError(domain string) *LLMError {
	return &LLMError{
		Code:    importNoMatchErrCode,
		Message: fmt.Sprintf("Spec %q not found in config.", domain),
	}
}

// NewConfigNotFoundError creates an LLMError with code "config_not_found".
func NewConfigNotFoundError() *LLMError {
	return &LLMError{
		Code:    configNotFoundErrCode,
		Message: "No configuration found in the workspace. Run 'swag2mcp init' first.",
	}
}

// NewResponseFileNotFoundError creates an LLMError with code "response_file_not_found".
func NewResponseFileNotFoundError(err error) *LLMError {
	return &LLMError{
		Code:    responseFileNotFoundErrCode,
		Message: "Response file not found - invoke an endpoint first.",
		Hint:    formatError(err),
	}
}

// NewResponsePathNotFoundError creates an LLMError with code "response_path_not_found".
func NewResponsePathNotFoundError(err error) *LLMError {
	return &LLMError{
		Code:    responsePathNotFoundErrCode,
		Message: "jsonPath did not match any value in the file.",
		Hint:    formatError(err),
	}
}

// NewResponseRequestError creates an LLMError with code "response_request_error".
func NewResponseRequestError(err error) *LLMError {
	return &LLMError{
		Code:    responseRequestErrCode,
		Message: "Request is invalid - path must point to a saved response file.",
		Hint:    formatError(err),
	}
}

// NewResponseJSONPathError creates an LLMError with code "response_jsonpath_error".
func NewResponseJSONPathError(err error) *LLMError {
	return &LLMError{
		Code:    responseJSONPathErrCode,
		Message: "jsonPath is invalid - use dotted path like data.0.name.",
		Hint:    formatError(err),
	}
}

// NewResponseLineRangeError creates an LLMError with code "response_line_range_error".
func NewResponseLineRangeError(err error) *LLMError {
	return &LLMError{
		Code:    responseLineRangeErrCode,
		Message: "Line or range is invalid - use 1-based line number or start-end.",
		Hint:    formatError(err),
	}
}

// NewResponseNotJSONError creates an LLMError with code "response_not_json".
func NewResponseNotJSONError(err error) *LLMError {
	return &LLMError{
		Code:    responseNotJSONErrCode,
		Message: "File is not valid JSON.",
		Hint:    formatError(err),
	}
}

// NewResponseReadFailedError creates an LLMError with code "response_read_failed".
func NewResponseReadFailedError(err error) *LLMError {
	return &LLMError{
		Code:    responseReadFailedErrCode,
		Message: "Failed to read the response file.",
		Hint:    formatError(err),
	}
}

// NewHTTPRequestError creates an LLMError with code "http_request_error".
func NewHTTPRequestError(err error) *LLMError {
	return &LLMError{
		Code:    httpRequestErrCode,
		Message: "The API request failed.",
		Hint:    formatError(err),
	}
}

// NewResponseReadError creates an LLMError with code "response_read_error".
func NewResponseReadError(err error) *LLMError {
	return &LLMError{
		Code:    responseReadErrCode,
		Message: "Failed to read the API response.",
		Hint:    formatError(err),
	}
}

// NewFileWriteError creates an LLMError with code "file_write_error".
func NewFileWriteError(err error) *LLMError {
	return &LLMError{
		Code:    fileWriteErrCode,
		Message: "Failed to write the API response to disk.",
		Hint:    formatError(err),
	}
}

// NewFileCreateError creates an LLMError with code "file_create_error".
func NewFileCreateError(err error) *LLMError {
	return &LLMError{
		Code:    fileCreateErrCode,
		Message: "Failed to create response file.",
		Hint:    formatError(err),
	}
}

// NewDirCreateError creates an LLMError with code "dir_create_error".
func NewDirCreateError(err error) *LLMError {
	return &LLMError{
		Code:    dirCreateErrCode,
		Message: "Failed to create responses directory.",
		Hint:    formatError(err),
	}
}

// NewStreamError creates an LLMError with code "stream_error".
func NewStreamError(err error) *LLMError {
	return &LLMError{
		Code:    streamErrCode,
		Message: "Failed to stream the API response to disk.",
		Hint:    formatError(err),
	}
}

// NewBuildRequestError creates an LLMError with code "build_request_error".
func NewBuildRequestError(err error) *LLMError {
	return &LLMError{
		Code:    buildRequestErrCode,
		Message: "Failed to build the HTTP request. Check the endpoint parameters and try again.",
		Hint:    formatError(err),
	}
}

// NewExportError creates an LLMError with code "export_error".
func NewExportError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    exportErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewImportError creates an LLMError with code "import_error".
func NewImportError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    importErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewRateLimitError creates an LLMError with code "rate_limit".
// This is returned when a specific endpoint is called too frequently.
// The LLM should wait for the cooldown or try a different endpoint.
func NewRateLimitError(err error) *LLMError {
	return &LLMError{
		Code:    rateLimitErrCode,
		Message: err.Error(),
		Hint:    "This endpoint was called too recently. Wait for the cooldown period to expire, then try again. You can call a different endpoint in the meantime.",
	}
}

// NewGlobalRateLimitError creates an LLMError with code "global_rate_limit".
// This is returned when the total number of invoke requests per second exceeds the limit.
// The LLM should wait before making any further requests.
func NewGlobalRateLimitError(err error) *LLMError {
	return &LLMError{
		Code:    globalRateLimitErrCode,
		Message: err.Error(),
		Hint:    "Too many API requests were made in a short time. Wait a few seconds before making any new invoke calls.",
	}
}

// NewInvokeError creates an LLMError with code "invoke_error".
func NewInvokeError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    invokeErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewConfigError creates an LLMError with code "config_error".
func NewConfigError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    configErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewWorkspaceError creates an LLMError with code "workspace_error".
func NewWorkspaceError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    workspaceErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewParseError creates an LLMError with code "parse_error".
func NewParseError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    parseErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewAuthError creates an LLMError with code "auth_error".
func NewAuthError(msg string, err error) *LLMError {
	return &LLMError{
		Code:    authErrorErrCode,
		Message: msg,
		Hint:    formatError(err),
	}
}

// NewToolLoadError creates an LLMError with code "tool_load_error".
func NewToolLoadError(err error) *LLMError {
	return &LLMError{
		Code:    toolLoadErrCode,
		Message: "Failed to load MCP tool definitions.",
		Hint:    formatError(err),
	}
}

// Error returns the JSON-encoded string representation of the LLMError.
func (e *LLMError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func formatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v", err)
}

// mapReaderError converts reader errors into LLM-facing errors.
func mapReaderError(err error) error {
	switch {
	case errors.Is(err, reader.ErrFileNotFound):
		return NewResponseFileNotFoundError(err)
	case errors.Is(err, reader.ErrPathNotAllowed):
		return NewResponseRequestError(err)
	case errors.Is(err, reader.ErrInvalidJSONPath):
		return NewResponseJSONPathError(err)
	case errors.Is(err, reader.ErrPathNotFound):
		return NewResponsePathNotFoundError(err)
	case errors.Is(err, reader.ErrInvalidLineRange):
		return NewResponseLineRangeError(err)
	case errors.Is(err, reader.ErrNotJSON):
		return NewResponseNotJSONError(err)
	default:
		return NewResponseReadFailedError(err)
	}
}
