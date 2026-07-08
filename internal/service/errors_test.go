package service

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestLLMError_ValidationError(t *testing.T) {
	t.Parallel()

	llmErr := NewValidationError("invalid input", nil)
	if llmErr.Code != validationFailedErrCode {
		t.Errorf("Code = %q, want %q", llmErr.Code, validationFailedErrCode)
	}
	if llmErr.Message != "invalid input" {
		t.Errorf("Message = %q, want %q", llmErr.Message, "invalid input")
	}
}

func TestLLMError_NotFoundError(t *testing.T) {
	t.Parallel()

	llmErr := NewNotFoundError("not found", nil)
	if llmErr.Code != notFoundErrCode {
		t.Errorf("Code = %q, want %q", llmErr.Code, notFoundErrCode)
	}
}

func TestLLMError_RateLimitError(t *testing.T) {
	t.Parallel()

	llmErr := NewRateLimitError(errors.New("rate limit exceeded"))
	if llmErr.Code != rateLimitErrCode {
		t.Errorf("Code = %q, want %q", llmErr.Code, rateLimitErrCode)
	}
	if llmErr.Message != "rate limit exceeded" {
		t.Errorf("Message = %q, want %q", llmErr.Message, "rate limit exceeded")
	}
}

func TestLLMError_JSONSerialization(t *testing.T) {
	t.Parallel()

	llmErr := NewValidationError("test message", nil)
	data, marshalErr := json.Marshal(llmErr)
	if marshalErr != nil {
		t.Fatalf("json.Marshal() = %v", marshalErr)
	}

	var decoded LLMError
	if unmarshalErr := json.Unmarshal(data, &decoded); unmarshalErr != nil {
		t.Fatalf("json.Unmarshal() = %v", unmarshalErr)
	}
	if decoded.Code != validationFailedErrCode {
		t.Errorf("Code = %q, want %q", decoded.Code, validationFailedErrCode)
	}
	if decoded.Message != "test message" {
		t.Errorf("Message = %q, want %q", decoded.Message, "test message")
	}
}

func TestLLMError_ErrorString(t *testing.T) {
	t.Parallel()

	llmErr := NewNotFoundError("spec not found", nil)
	errStr := llmErr.Error()
	if errStr == "" {
		t.Fatal("Error() returned empty string")
	}
}
