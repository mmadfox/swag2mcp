package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mmadfox/swag2mcp/internal/reader"
)

const (
	validationFailedErrCode = "validation_failed"
	notFoundErrCode         = "not_found"
	rateLimitErrCode        = "rate_limit"
	invokeErrorErrCode      = "invoke_error"
	configErrorErrCode      = "config_error"
	workspaceErrorErrCode   = "workspace_error"
	parseErrorErrCode       = "parse_error"
	authErrorErrCode        = "auth_error"
)

// LLMError is an error type returned to the LLM with a machine-readable code and human-readable message.
type LLMError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Original string `json:"hint,omitempty"`
}

// NewValidationError creates an LLMError with code "validation_failed".
func NewValidationError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     validationFailedErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

// NewNotFoundError creates an LLMError with code "not_found".
func NewNotFoundError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     notFoundErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

// NewRateLimitError creates an LLMError with code "rate_limit".
func NewRateLimitError(err error) *LLMError {
	return &LLMError{
		Code:     rateLimitErrCode,
		Message:  err.Error(),
		Original: "",
	}
}

// NewInvokeError creates an LLMError with code "invoke_error".
func NewInvokeError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     invokeErrorErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

// NewConfigError creates an LLMError with code "config_error".
func NewConfigError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     configErrorErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

// NewWorkspaceError creates an LLMError with code "workspace_error".
func NewWorkspaceError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     workspaceErrorErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

// NewParseError creates an LLMError with code "parse_error".
func NewParseError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     parseErrorErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

// NewAuthError creates an LLMError with code "auth_error".
func NewAuthError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     authErrorErrCode,
		Message:  msg,
		Original: formatError(err),
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
		return NewNotFoundError(
			"Response file not found - invoke an endpoint first.",
			err,
		)
	case errors.Is(err, reader.ErrPathNotAllowed):
		return NewValidationError(
			"Path is not inside the responses directory.",
			err,
		)
	case errors.Is(err, reader.ErrInvalidJSONPath):
		return NewValidationError(
			"jsonPath is invalid - use dotted path like data.0.name.",
			err,
		)
	case errors.Is(err, reader.ErrPathNotFound):
		return NewNotFoundError(
			"jsonPath did not match any value in the file.",
			err,
		)
	case errors.Is(err, reader.ErrInvalidLineRange):
		return NewValidationError(
			"Line or range is invalid - use 1-based line number or start-end.",
			err,
		)
	case errors.Is(err, reader.ErrNotJSON):
		return NewValidationError(
			"File is not valid JSON.",
			err,
		)
	default:
		return NewInvokeError(
			"Failed to read the response file.",
			err,
		)
	}
}
