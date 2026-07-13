package service

import (
	"encoding/json"
	"fmt"
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
