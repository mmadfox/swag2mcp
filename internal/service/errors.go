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

type LLMError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Original string `json:"hint,omitempty"`
}

func NewValidationError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     validationFailedErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

func NewNotFoundError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     notFoundErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

func NewRateLimitError(err error) *LLMError {
	return &LLMError{
		Code:     rateLimitErrCode,
		Message:  err.Error(),
		Original: "",
	}
}

func NewInvokeError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     invokeErrorErrCode,
		Message:  msg,
		Original: formatError(err),
	}
}

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
