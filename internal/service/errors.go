package service

import (
	"encoding/json"
	"fmt"
)

const (
	validationFailedErrCode = "validation_failed"
	notFoundErrCode         = "not_found"
)

type LLMError struct {
	Code     string `json:"code"`    // "validation_failed", "not_found", "internal_error"
	Message  string `json:"message"` // human-readable for LLM
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
