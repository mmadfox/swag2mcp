package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

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

func TestNewRateLimitError(t *testing.T) {
	t.Parallel()

	err := NewRateLimitError(newTestError("rate limited"))
	require.Equal(t, "rate_limit", err.Code)
	require.Equal(t, "rate limited", err.Message)
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
	require.Contains(t, err.Error(), "not_found")
}

func TestMapReaderError_pathNotAllowed(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrPathNotAllowed)
	require.Contains(t, err.Error(), "validation_failed")
}

func TestMapReaderError_invalidJSONPath(t *testing.T) {
	t.Parallel()

	err := mapReaderError(reader.ErrInvalidJSONPath)
	require.Contains(t, err.Error(), "validation_failed")
}

func TestMapReaderError_default(t *testing.T) {
	t.Parallel()

	err := mapReaderError(newTestError("unknown"))
	require.Contains(t, err.Error(), "invoke_error")
}
