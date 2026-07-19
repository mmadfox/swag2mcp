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
			"The response file was not found — invoke an endpoint first and use the fileRef.path returned.",
			err,
		)
	case errors.Is(err, reader.ErrPathNotAllowed):
		return NewValidationError(
			"The path is not inside the responses directory — only saved response files may be read.",
			err,
		)
	case errors.Is(err, reader.ErrInvalidJSONPath):
		return NewValidationError(
			"The jsonPath is invalid — use a dotted path such as data.0.name.",
			err,
		)
	case errors.Is(err, reader.ErrPathNotFound):
		return NewNotFoundError(
			"The jsonPath did not match any value in the file — check the outline and try a different path.",
			err,
		)
	case errors.Is(err, reader.ErrInvalidLineRange):
		return NewValidationError(
			"The line or range is invalid — use a 1-based line number or a start-end range.",
			err,
		)
	case errors.Is(err, reader.ErrNotJSON):
		return NewValidationError(
			"The file is not valid JSON — only JSON response files can be outlined or sliced.",
			err,
		)
	default:
		return NewInvokeError(
			"Failed to read the response file — verify the path and try again.",
			err,
		)
	}
}
