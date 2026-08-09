/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestResponseService_ResponseOutline_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{}, slog.Default())
	_, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseCompress_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{}, slog.Default())
	_, err := svc.ResponseCompress(context.Background(), ResponseCompressRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseSlice_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{}, slog.Default())
	_, err := svc.ResponseSlice(context.Background(), ResponseSliceRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseOutline_success(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{"name": "test", "items": []int{1, 2, 3}}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{Path: fp})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Outline.Type)
}

func TestResponseService_ResponseOutline_fileNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(t.TempDir()).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	_, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{Path: "/nonexistent/file.json"})
	require.Error(t, err)
}

func TestResponseService_ResponseCompress_success(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{"name": "test", "items": []int{1, 2, 3}}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseCompress(context.Background(), ResponseCompressRequest{
		Path: fp,
		Mode: reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Body)
}

func TestResponseService_ResponseCompress_tooLarge(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	large := make([]int, 1000)
	for i := range large {
		large[i] = i
	}
	data := map[string]any{"items": large}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "large.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	ctx := newServiceContext()
	ctx.maxResponseSize.Store(10)

	svc := newResponseService(ctx, ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseCompress(context.Background(), ResponseCompressRequest{
		Path: fp,
		Mode: reader.CompressKeysOnly,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.FileRef)
	require.FileExists(t, resp.FileRef.Path)
}

func TestResponseService_ResponseSlice_success(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{"name": "test", "value": 42}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseSlice(context.Background(), ResponseSliceRequest{
		Path:     fp,
		JSONPath: "name",
	})
	require.NoError(t, err)
	require.Equal(t, "test", resp.Slice.Value)
}

func TestResponseService_ResponseSlice_tooLarge(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{"items": make([]int, 500)}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "large.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	ctx := newServiceContext()
	ctx.maxResponseSize.Store(10)

	svc := newResponseService(ctx, ws, fakeValidator{}, slog.Default())
	_, err = svc.ResponseSlice(context.Background(), ResponseSliceRequest{
		Path:     fp,
		JSONPath: "items",
	})
	require.Error(t, err)
}

func TestResponseService_saveReaderResult(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	ref, err := svc.saveReaderResult("", map[string]any{"key": "val"})
	require.NoError(t, err)
	require.FileExists(t, ref.Path)
	require.Contains(t, ref.Message, "saved to disk")
}

func TestResponseService_saveReaderResult_string(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	ref, err := svc.saveReaderResult("", "raw string data")
	require.NoError(t, err)
	require.FileExists(t, ref.Path)
}

func TestResponseService_saveReaderResult_bytes(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	ref, err := svc.saveReaderResult("", []byte(`{"raw": true}`))
	require.NoError(t, err)
	require.FileExists(t, ref.Path)
}

func TestResponseService_saveReaderResult_maxSizeHint(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	ctx := newServiceContext()
	ctx.maxResponseSize.Store(2048)

	svc := newResponseService(ctx, ws, fakeValidator{}, slog.Default())
	ref, err := svc.saveReaderResult("", map[string]any{"key": "val"})
	require.NoError(t, err)
	require.Contains(t, ref.MaxSizeHint, "KB")
}

func TestResponseService_ResponseFilter_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{}, slog.Default())
	_, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseFilter_fileNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(t.TempDir()).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	_, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     "/nonexistent/file.json",
		JSONPath: "items",
	})
	require.Error(t, err)
}

func TestResponseService_ResponseFilter_search(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{
		"items": []any{
			map[string]any{"name": "bitcoin", "price": 50000},
			map[string]any{"name": "ethereum", "price": 3000},
			map[string]any{"name": "litecoin", "price": 100},
		},
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		Search:   "bitcoin",
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Total)
	require.Len(t, resp.Items, 1)
	require.Equal(t, "memory", resp.Strategy)
}

func TestResponseService_ResponseFilter_filter(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{
		"items": []any{
			map[string]any{"name": "bitcoin", "price": 50000},
			map[string]any{"name": "ethereum", "price": 3000},
			map[string]any{"name": "litecoin", "price": 100},
		},
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		Filter:   "price > 1000",
	})
	require.NoError(t, err)
	require.Equal(t, 2, resp.Total)
	require.Len(t, resp.Items, 2)
}

func TestResponseService_ResponseFilter_pagination(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := make([]any, 25)
	for i := range items {
		items[i] = map[string]any{"id": i + 1, "name": fmt.Sprintf("item-%d", i+1)}
	}
	data := map[string]any{"results": items}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())

	// Page 1
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "results",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, 25, resp.Total)
	require.Equal(t, 3, resp.TotalPages)
	require.Len(t, resp.Items, 10)

	// Page 3 (last)
	resp, err = svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "results",
		Page:     3,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Len(t, resp.Items, 5)
}

func TestResponseService_ResponseFilter_emptyArray(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{"items": []any{}}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "empty.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
	})
	require.NoError(t, err)
	require.Equal(t, 0, resp.Total)
	require.Empty(t, resp.Items)
}

func TestResponseService_ResponseFilter_invalidFilter(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{"items": []any{map[string]any{"name": "test"}}}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	_, err = svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		Filter:   "invalid",
	})
	require.Error(t, err)
}

func TestResponseService_ResponseFilter_nestedJSONPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{
		"data": map[string]any{
			"items": []any{
				map[string]any{"name": "alpha"},
				map[string]any{"name": "beta"},
			},
		},
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "nested.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "data.items",
		Search:   "alpha",
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Total)
}

func TestResponseService_ResponseFilter_containsFilter(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{
		"items": []any{
			map[string]any{"name": "hello world"},
			map[string]any{"name": "goodbye world"},
		},
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		Filter:   "name contains hello",
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Total)
}

func TestResponseService_ResponseFilter_notEqualFilter(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{
		"items": []any{
			map[string]any{"status": "active"},
			map[string]any{"status": "inactive"},
		},
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		Filter:   "status != inactive",
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Total)
}

func TestResponseService_ResponseFilter_pageOutOfRange(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	data := map[string]any{
		"items": []any{
			map[string]any{"name": "only one"},
		},
	}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		Page:     100,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Total)
	require.Empty(t, resp.Items)
}

func TestResponseService_ResponseFilter_defaultPageSize(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := make([]any, 20)
	for i := range items {
		items[i] = map[string]any{"id": i + 1}
	}
	data := map[string]any{"items": items}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
	})
	require.NoError(t, err)
	require.Equal(t, 10, resp.PageSize)
	require.Len(t, resp.Items, 10)
}

func TestResponseService_ResponseFilter_maxPageSize(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := make([]any, 100)
	for i := range items {
		items[i] = map[string]any{"id": i + 1}
	}
	data := map[string]any{"items": items}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "test.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "items",
		PageSize: 200,
	})
	require.NoError(t, err)
	require.Equal(t, 50, resp.PageSize)
	require.Len(t, resp.Items, 50)
}

func TestResponseService_ResponseFilter_rootArray(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{
		map[string]any{"name": "bitcoin", "price": 50000},
		map[string]any{"name": "ethereum", "price": 3000},
		map[string]any{"name": "litecoin", "price": 100},
	}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())

	// Search with empty jsonPath (root array)
	resp, err := svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:   fp,
		Search: "bitcoin",
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Total)
	require.Len(t, resp.Items, 1)
	require.Equal(t, "memory", resp.Strategy)

	// Filter with empty jsonPath
	resp, err = svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:   fp,
		Filter: "price > 1000",
	})
	require.NoError(t, err)
	require.Equal(t, 2, resp.Total)
	require.Len(t, resp.Items, 2)

	// Pagination with empty jsonPath
	resp, err = svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		Page:     1,
		PageSize: 2,
	})
	require.NoError(t, err)
	require.Equal(t, 3, resp.Total)
	require.Len(t, resp.Items, 2)
}

func TestResponseService_ResponseFilter_rootArrayStreaming(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := make([]any, 100)
	for i := range items {
		items[i] = map[string]any{"id": i + 1, "name": fmt.Sprintf("item-%d", i+1)}
	}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root-stream.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), fakeValidator{}, slog.Default())

	// Call filterArrayStreaming directly to test streaming strategy
	matched, total, err := svc.filterArrayStreaming(fp, "", "item-50", "", 1, 10)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, matched, 1)
}

func TestResponseService_ResponseFilter_rootArrayIndex(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{
		map[string]any{"id": 1, "name": "alpha"},
		map[string]any{"id": 2, "name": "beta"},
		map[string]any{"id": 3, "name": "gamma"},
	}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())

	// jsonPath "0" on a root array resolves to an element (object), which is
	// not an array — filter must fail with an informative path-not-found error.
	_, err = svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "0",
		Search:   "alpha",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "response_path_not_found")
}

func TestResponseService_ResponseFilter_rootArrayLast(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{
		map[string]any{"id": 1, "name": "alpha"},
		map[string]any{"id": 2, "name": "beta"},
		map[string]any{"id": 3, "name": "gamma"},
	}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())

	// jsonPath "-" on a root array resolves to the last element (object), which
	// is not an array — filter must fail with an informative path-not-found error.
	_, err = svc.ResponseFilter(context.Background(), ResponseFilterRequest{
		Path:     fp,
		JSONPath: "-",
		Search:   "gamma",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "response_path_not_found")
}

func TestResponseService_ResponseSlice_rootArrayLast(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{
		map[string]any{"id": 1, "name": "alpha"},
		map[string]any{"id": 2, "name": "beta"},
		map[string]any{"id": 3, "name": "gamma"},
	}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseSlice(context.Background(), ResponseSliceRequest{
		Path:     fp,
		JSONPath: "-",
	})
	require.NoError(t, err)
	val, ok := resp.Slice.Value.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(3), val["id"])
}

func TestResponseService_ResponseOutline_rootArrayHint(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{map[string]any{"id": 1}, map[string]any{"id": 2}}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{Path: fp})
	require.NoError(t, err)
	require.Equal(t, "array", resp.Outline.Type)
	var found bool
	for _, h := range resp.Outline.CompressionHints {
		if strings.Contains(h, "root-level array") {
			found = true
			break
		}
	}
	require.True(t, found, "expected a root-level array hint, got %v", resp.Outline.CompressionHints)
}

func TestResponseService_ResponseSlice_emptyJSONPathRootArray(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{map[string]any{"id": 1}, map[string]any{"id": 2}}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseSlice(context.Background(), ResponseSliceRequest{Path: fp})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Slice.Context)
	require.True(t, resp.Slice.IsComplete)
	val, ok := resp.Slice.Value.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), val["id"])
}

func TestResponseService_ResponseSlice_atThisJSONPathRootArray(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	respDir := filepath.Join(tmpDir, "responses")
	require.NoError(t, os.MkdirAll(respDir, 0o750))

	items := []any{map[string]any{"id": 1}, map[string]any{"id": 2}}
	raw, err := json.Marshal(items)
	require.NoError(t, err)
	fp := filepath.Join(respDir, "root.json")
	require.NoError(t, os.WriteFile(fp, raw, 0o600))

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(respDir).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{}, slog.Default())
	resp, err := svc.ResponseSlice(context.Background(), ResponseSliceRequest{Path: fp, JSONPath: "@this"})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Slice.Context)
	require.True(t, resp.Slice.IsComplete)
	val, ok := resp.Slice.Value.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(1), val["id"])
}
