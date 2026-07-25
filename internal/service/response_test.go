package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/mmadfox/swag2mcp/internal/workspace"
	"github.com/stretchr/testify/require"
)

func TestResponseOutline(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithWorkspace(t)
	path := writeResponseFile(t, s, `{"status":"ok","data":[{"id":1},{"id":2}]}`)

	resp, err := s.ResponseOutline(context.Background(), ResponseOutlineRequest{
		Path:          path,
		MaxDepth:      3,
		MaxArrayItems: 2,
	})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Outline.Type)
	require.Equal(t, []string{"status", "data"}, resp.Outline.Keys)
	require.NotEmpty(t, resp.Outline.CompressionHints)
}

func TestResponseOutline_PathOutsideResponsesDir(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithWorkspace(t)
	outside := filepath.Join(t.TempDir(), "outside.json")
	require.NoError(t, os.WriteFile(outside, []byte(`{"a":1}`), 0o600))

	_, err := s.ResponseOutline(context.Background(), ResponseOutlineRequest{Path: outside})
	require.Error(t, err)
	var llmErr *LLMError
	require.True(t, errors.As(err, &llmErr))
	require.Equal(t, validationFailedErrCode, llmErr.Code)
}

func TestResponseCompress_FirstOfArray(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithWorkspace(t)
	path := writeResponseFile(t, s, `{"items":[{"id":1,"name":"a"},{"id":2,"name":"b"}]}`)

	resp, err := s.ResponseCompress(context.Background(), ResponseCompressRequest{
		Path:     path,
		JSONPath: "items",
		Mode:     reader.CompressFirstOfArray,
	})
	require.NoError(t, err)
	require.Nil(t, resp.FileRef)
	body, ok := resp.Body.(map[string]any)
	require.True(t, ok)
	compressed, ok := body["compressed"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), compressed["length"])
}

func TestResponseCompress_TooLarge(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithSmallLimit(t)
	path := writeResponseFile(t, s, `{"items":[{"id":1,"name":"very long string value here"}]}`)

	resp, err := s.ResponseCompress(context.Background(), ResponseCompressRequest{
		Path:     path,
		JSONPath: "items",
		Mode:     reader.CompressFirstOfArray,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.FileRef)
	require.Nil(t, resp.Body)
}

func TestResponseSlice_ByJSONPath(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithWorkspace(t)
	path := writeResponseFile(t, s, `{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`)

	resp, err := s.ResponseSlice(context.Background(), ResponseSliceRequest{
		Path:     path,
		JSONPath: "users.1",
	})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Slice.Context)
	require.Equal(t, "users.2", resp.Slice.NextPath)
	require.Equal(t, "users.0", resp.Slice.PrevPath)
	val, ok := resp.Slice.Value.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(2), val["id"])
}

func TestResponseSlice_ByLine(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithWorkspace(t)
	path := writeResponseFile(t, s, "{\n\t\"status\": \"ok\"\n}\n")

	resp, err := s.ResponseSlice(context.Background(), ResponseSliceRequest{
		Path:   path,
		Line:   2,
		Around: 1,
	})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Slice.Context)
	require.True(t, resp.Slice.IsComplete)
}

func TestResponseSlice_PathNotFound(t *testing.T) {
	t.Parallel()

	s := newTestServiceWithWorkspace(t)
	path := writeResponseFile(t, s, `{"a":1}`)

	_, err := s.ResponseSlice(context.Background(), ResponseSliceRequest{
		Path:     path,
		JSONPath: "missing",
	})
	require.Error(t, err)
	var llmErr *LLMError
	require.True(t, errors.As(err, &llmErr))
	require.Equal(t, notFoundErrCode, llmErr.Code)
}

func newTestServiceWithWorkspace(t *testing.T) *Service {
	t.Helper()
	ws, err := workspace.NewFromBase(t.TempDir())
	require.NoError(t, err)
	return newTestService(t, WithWorkspace(ws))
}

func newTestServiceWithSmallLimit(t *testing.T) *Service {
	t.Helper()
	ws, err := workspace.NewFromBase(t.TempDir())
	require.NoError(t, err)
	s := newTestService(t, WithWorkspace(ws))
	s.maxResponseSize = 10
	return s
}

func writeResponseFile(t *testing.T, s *Service, content string) string {
	t.Helper()

	require.NoError(t, os.MkdirAll(s.ws.ResponsesDir(), 0o750))
	path := filepath.Join(s.ws.ResponsesDir(), "test-response.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}
