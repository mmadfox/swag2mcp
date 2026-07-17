package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLLMError_ValidationError(t *testing.T) {
	t.Parallel()

	llmErr := NewValidationError("invalid input", nil)
	require.Equal(t, validationFailedErrCode, llmErr.Code)
	require.Equal(t, "invalid input", llmErr.Message)
}

func TestLLMError_NotFoundError(t *testing.T) {
	t.Parallel()

	llmErr := NewNotFoundError("not found", nil)
	require.Equal(t, notFoundErrCode, llmErr.Code)
}

func TestLLMError_RateLimitError(t *testing.T) {
	t.Parallel()

	llmErr := NewRateLimitError(errors.New("rate limit exceeded"))
	require.Equal(t, rateLimitErrCode, llmErr.Code)
	require.Equal(t, "rate limit exceeded", llmErr.Message)
}

func TestLLMError_InvokeError(t *testing.T) {
	t.Parallel()

	llmErr := NewInvokeError("request failed", errors.New("connection refused"))
	require.Equal(t, invokeErrorErrCode, llmErr.Code)
	require.Equal(t, "request failed", llmErr.Message)
	require.Equal(t, "connection refused", llmErr.Original)
}

func TestLLMError_InvokeError_NilError(t *testing.T) {
	t.Parallel()

	llmErr := NewInvokeError("request failed", nil)
	require.Empty(t, llmErr.Original)
}

func TestLLMError_JSONSerialization(t *testing.T) {
	t.Parallel()

	llmErr := NewValidationError("test message", nil)
	data, err := json.Marshal(llmErr)
	require.NoError(t, err)

	var decoded LLMError
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	require.Equal(t, validationFailedErrCode, decoded.Code)
	require.Equal(t, "test message", decoded.Message)
}

func TestLLMError_ErrorString(t *testing.T) {
	t.Parallel()

	llmErr := NewNotFoundError("spec not found", nil)
	require.NotEmpty(t, llmErr.Error())
}
