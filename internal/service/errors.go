package service

import "fmt"

type LLMError struct {
	Code     string `json:"code"`    // "validation_failed", "not_found", "internal_error"
	Message  string `json:"message"` // human-readable for LLM
	Original string `json:"original,omitempty"`
}

func NewValidationError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     "validation_failed",
		Message:  msg,
		Original: formatError(err),
	}
}

func NewNotFoundError(msg string, err error) *LLMError {
	return &LLMError{
		Code:     "not_found",
		Message:  msg,
		Original: formatError(err),
	}
}

func (e *LLMError) Error() string {
	return e.Message
}

func formatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v", err)
}
