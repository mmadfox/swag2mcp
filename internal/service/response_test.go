package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestResponseService_ResponseOutline_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{})
	_, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseCompress_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{})
	_, err := svc.ResponseCompress(context.Background(), ResponseCompressRequest{})
	require.Error(t, err)
}

func TestResponseService_ResponseSlice_validationError(t *testing.T) {
	t.Parallel()

	svc := newResponseService(newServiceContext(), NewMockWorkspaceOps(gomock.NewController(t)), strictValidator{})
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

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
	resp, err := svc.ResponseOutline(context.Background(), ResponseOutlineRequest{Path: fp})
	require.NoError(t, err)
	require.Equal(t, "object", resp.Outline.Type)
}

func TestResponseService_ResponseOutline_fileNotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ws := NewMockWorkspaceOps(ctrl)
	ws.EXPECT().ResponsesDir().Return(t.TempDir()).AnyTimes()

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
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

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
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

	svc := newResponseService(ctx, ws, fakeValidator{})
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

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
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

	svc := newResponseService(ctx, ws, fakeValidator{})
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

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
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

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
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

	svc := newResponseService(newServiceContext(), ws, fakeValidator{})
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

	svc := newResponseService(ctx, ws, fakeValidator{})
	ref, err := svc.saveReaderResult("", map[string]any{"key": "val"})
	require.NoError(t, err)
	require.Contains(t, ref.MaxSizeHint, "KB")
}
