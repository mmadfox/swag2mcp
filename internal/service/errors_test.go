/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
)

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func newTestError(msg string) error {
	return &testError{msg: msg}
}

func TestLLMError_Error(t *testing.T) {
	t.Parallel()

	err := NewValidationError("test message", nil)
	require.Contains(t, err.Error(), "validation_failed")
	require.Contains(t, err.Error(), "test message")
}

func TestNewValidationError(t *testing.T) {
	t.Parallel()

	err := NewValidationError("invalid input", nil)
	require.Equal(t, "validation_failed", err.Code)
	require.Equal(t, "invalid input", err.Message)
}

func TestNewNotFoundError(t *testing.T) {
	t.Parallel()

	err := NewNotFoundError("not found", nil)
	require.Equal(t, "not_found", err.Code)
}

func TestNewEndpointNotFoundError(t *testing.T) {
	t.Parallel()

	err := NewEndpointNotFoundError("ep1", newTestError("index miss"))
	require.Equal(t, "endpoint_not_found", err.Code)
	require.Contains(t, err.Message, "ep1")
	require.Contains(t, err.Hint, "search tool")
}

func TestNewSpecNotFoundError(t *testing.T) {
	t.Parallel()

	err := NewSpecNotFoundError("s1", newTestError("index miss"))
	require.Equal(t, "spec_not_found", err.Code)
	require.Contains(t, err.Message, "s1")
	require.Contains(t, err.Hint, "spec_list")
}

func TestNewCollectionNotFoundError(t *testing.T) {
	t.Parallel()

	err := NewCollectionNotFoundError("c1", newTestError("index miss"))
	require.Equal(t, "collection_not_found", err.Code)
	require.Contains(t, err.Message, "c1")
	require.Contains(t, err.Hint, "collection_by_spec")
}

func TestNewTagNotFoundError(t *testing.T) {
	t.Parallel()

	err := NewTagNotFoundError("t1", newTestError("index miss"))
	require.Equal(t, "tag_not_found", err.Code)
	require.Contains(t, err.Message, "t1")
	require.Contains(t, err.Hint, "tag_by_collection")
}

func TestNewParameterValidationError(t *testing.T) {
	t.Parallel()

	err := NewParameterValidationError(newTestError("missing required param"))
	require.Equal(t, "parameter_validation_failed", err.Code)
	require.Contains(t, err.Message, "required parameters")
	require.Contains(t, err.Hint, "missing required param")
}

func TestNewRequestBodyValidationError(t *testing.T) {
	t.Parallel()

	err := NewRequestBodyValidationError(newTestError("unknown field"))
	require.Equal(t, "request_body_validation_failed", err.Code)
	require.Contains(t, err.Message, "required fields")
	require.Contains(t, err.Hint, "unknown field")
}

func TestNewRateLimitError(t *testing.T) {
	t.Parallel()

	err := NewRateLimitError(newTestError("rate limited"))
	require.Equal(t, "rate_limit", err.Code)
	require.Equal(t, "rate limited", err.Message)
}

func TestNewRateLimitError_Hint(t *testing.T) {
	t.Parallel()

	err := NewRateLimitError(newTestError("rate limit exceeded for endpoint \"ep1\": try again in 8 seconds"))
	require.Equal(t, "rate_limit", err.Code)
	require.Contains(t, err.Message, "try again in 8 seconds")
	require.Contains(t, err.Hint, "Wait for the cooldown period")
	require.Contains(t, err.Hint, "different endpoint")
}

func TestNewGlobalRateLimitError(t *testing.T) {
	t.Parallel()

	err := NewGlobalRateLimitError(newTestError("global rate limit exceeded"))
	require.Equal(t, "global_rate_limit", err.Code)
	require.Equal(t, "global rate limit exceeded", err.Message)
	require.Contains(t, err.Hint, "Too many API requests")
	require.Contains(t, err.Hint, "Wait a few seconds")
}

func TestNewInvokeError(t *testing.T) {
	t.Parallel()

	err := NewInvokeError("api call failed", nil)
	require.Equal(t, "invoke_error", err.Code)
}

func TestNewConfigError(t *testing.T) {
	t.Parallel()

	err := NewConfigError("config error", nil)
	require.Equal(t, "config_error", err.Code)
}

func TestNewWorkspaceError(t *testing.T) {
	t.Parallel()

	err := NewWorkspaceError("workspace error", nil)
	require.Equal(t, "workspace_error", err.Code)
}

func TestNewParseError(t *testing.T) {
	t.Parallel()

	err := NewParseError("parse error", nil)
	require.Equal(t, "parse_error", err.Code)
}

func TestNewAuthError(t *testing.T) {
	t.Parallel()

	err := NewAuthError("auth error", nil)
	require.Equal(t, "auth_error", err.Code)
}

func TestMapReaderError_fileNotFound(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrFileNotFound)
	require.Contains(t, err.Error(), "response_file_not_found")
}

func TestMapReaderError_pathNotAllowed(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrPathNotAllowed)
	require.Contains(t, err.Error(), "response_request_error")
}

func TestMapReaderError_invalidJSONPath(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrInvalidJSONPath)
	require.Contains(t, err.Error(), "response_jsonpath_error")
}

func TestMapReaderError_default(t *testing.T) {
	t.Parallel()

	err := mapReaderError(newTestError("unknown"))
	require.Contains(t, err.Error(), "response_read_failed")
}

func TestMapReaderError_pathNotFound(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrPathNotFound)
	require.Contains(t, err.Error(), "response_path_not_found")
}

func TestMapReaderError_invalidLineRange(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrInvalidLineRange)
	require.Contains(t, err.Error(), "response_line_range_error")
}

func TestMapReaderError_notJSON(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrNotJSON)
	require.Contains(t, err.Error(), "response_not_json")
}
