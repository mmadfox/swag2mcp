package service

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmadfox/swag2mcp/internal/reader"
)

// ResponseOutlineRequest requests a structural outline of a saved response file.
type ResponseOutlineRequest struct {
	Path          string `json:"path"          validate:"required"            jsonschema:"required,Absolute path returned in fileRef.path"`
	MaxDepth      int    `json:"maxDepth,omitempty"      jsonschema:"optional,Max recursion depth (default 3)"`
	MaxArrayItems int    `json:"maxArrayItems,omitempty" jsonschema:"optional,How many array items to inspect (default 5)"`
}

// ResponseOutlineResponse returns the structural summary of a response file.
type ResponseOutlineResponse struct {
	Outline reader.Outline `json:"outline"`
}

// ResponseCompressRequest requests compression of a JSON value in a saved response file.
type ResponseCompressRequest struct {
	Path       string              `json:"path"       validate:"required"            jsonschema:"required,Absolute path returned in fileRef.path"`
	JSONPath   string              `json:"jsonPath,omitempty"   jsonschema:"optional,Path to the value to compress (e.g. data or data.0)"`
	Mode       reader.CompressMode `json:"mode"       validate:"required"            jsonschema:"required,Compression mode: first_of_array, sample_array, truncate_strings, keys_only, select_keys"`
	ArrayHead  int                 `json:"arrayHead,omitempty"  jsonschema:"optional,Number of leading array items for sample_array mode"`
	ArrayTail  int                 `json:"arrayTail,omitempty"  jsonschema:"optional,Number of trailing array items for sample_array mode"`
	StringLen  int                 `json:"stringLen,omitempty"  jsonschema:"optional,Maximum string length for truncate_strings mode"`
	SelectKeys []string            `json:"selectKeys,omitempty" jsonschema:"optional,Keys to keep for select_keys mode"`
}

// ResponseCompressResponse returns a compressed JSON body or a file reference.
type ResponseCompressResponse struct {
	Body    any            `json:"body,omitempty"`
	FileRef *FileReference `json:"fileRef,omitempty"`
	Hint    string         `json:"hint,omitempty"`
}

// ResponseSliceRequest requests a fragment of a saved response file.
type ResponseSliceRequest struct {
	Path     string `json:"path"     validate:"required"            jsonschema:"required,Absolute path returned in fileRef.path"`
	JSONPath string `json:"jsonPath,omitempty" jsonschema:"optional,Logical path to the value (e.g. data.3.name)"`
	Line     int    `json:"line,omitempty"   jsonschema:"optional,1-based line number to center the fragment on"`
	Range    string `json:"range,omitempty"  jsonschema:"optional,Line range as start-end (e.g. 120-240)"`
	Around   int    `json:"around,omitempty" jsonschema:"optional,Lines to include around line (default 20)"`
}

// ResponseSliceResponse returns a JSON fragment and its context.
type ResponseSliceResponse struct {
	Slice   reader.Slice   `json:"slice"`
	FileRef *FileReference `json:"fileRef,omitempty"`
}

type responseService struct {
	ctx *serviceContext
	ws  WorkspaceOps
	v   RequestValidator
}

func newResponseService(
	ctx *serviceContext,
	ws WorkspaceOps,
	v RequestValidator,
) *responseService {
	return &responseService{ctx: ctx, ws: ws, v: v}
}

// ResponseOutline returns a high-level structural summary of a saved response.
func (rs *responseService) ResponseOutline(_ context.Context, req ResponseOutlineRequest) (ResponseOutlineResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseOutlineResponse{}, NewValidationError(
			"Request is invalid - path must point to a saved response file.",
			err,
		)
	}

	r := reader.New(rs.ws.ResponsesDir())
	outline, err := r.Outline(req.Path, reader.OutlineOptions{
		MaxDepth:      req.MaxDepth,
		MaxArrayItems: req.MaxArrayItems,
	})
	if err != nil {
		return ResponseOutlineResponse{}, mapReaderError(err)
	}

	return ResponseOutlineResponse{Outline: outline}, nil
}

// ResponseCompress reduces a JSON value in a saved response file.
func (rs *responseService) ResponseCompress(_ context.Context, req ResponseCompressRequest) (ResponseCompressResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseCompressResponse{}, NewValidationError(
			"Request is invalid - path and mode are required.",
			err,
		)
	}

	r := reader.New(rs.ws.ResponsesDir())
	result, err := r.Compress(req.Path, reader.CompressOptions{
		JSONPath:   req.JSONPath,
		Mode:       req.Mode,
		ArrayHead:  req.ArrayHead,
		ArrayTail:  req.ArrayTail,
		StringLen:  req.StringLen,
		SelectKeys: req.SelectKeys,
		Limit:      int(rs.ctx.maxResponseSize.Load()),
	})
	if err != nil {
		return ResponseCompressResponse{}, mapReaderError(err)
	}

	if result.TooLarge {
		ref, saveErr := rs.saveReaderResult(req.Path, result.Body)
		if saveErr != nil {
			return ResponseCompressResponse{}, NewInvokeError(
				"Compressed result is too large and could not be saved.",
				saveErr,
			)
		}
		return ResponseCompressResponse{FileRef: &ref, Hint: result.Hint}, nil
	}

	return ResponseCompressResponse{Body: result.Body, Hint: result.Hint}, nil
}

// ResponseSlice extracts a fragment of a saved response file.
func (rs *responseService) ResponseSlice(_ context.Context, req ResponseSliceRequest) (ResponseSliceResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseSliceResponse{}, NewValidationError(
			"Request is invalid - provide path and jsonPath, line, or range.",
			err,
		)
	}

	r := reader.New(rs.ws.ResponsesDir())
	slice, err := r.Slice(req.Path, reader.SliceOptions{
		JSONPath: req.JSONPath,
		Line:     req.Line,
		Range:    req.Range,
		Around:   req.Around,
		Limit:    int(rs.ctx.maxResponseSize.Load()),
	})
	if err != nil {
		return ResponseSliceResponse{}, mapReaderError(err)
	}

	maxSize := int(rs.ctx.maxResponseSize.Load())
	fragmentBytes, _ := json.Marshal(slice.Value)
	if maxSize > 0 && len(fragmentBytes) > maxSize {
		ref, saveErr := rs.saveReaderResult(req.Path, fragmentBytes)
		if saveErr != nil {
			return ResponseSliceResponse{}, NewInvokeError(
				"Extracted fragment is too large and could not be saved.",
				saveErr,
			)
		}
		return ResponseSliceResponse{Slice: slice, FileRef: &ref}, nil
	}

	return ResponseSliceResponse{Slice: slice}, nil
}

// saveReaderResult saves an extracted or compressed JSON fragment to disk.
func (rs *responseService) saveReaderResult(_ string, data any) (FileReference, error) {
	var body []byte
	switch v := data.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		var err error
		body, err = json.Marshal(v)
		if err != nil {
			return FileReference{}, fmt.Errorf("marshal result: %w", err)
		}
	}

	if err := os.MkdirAll(rs.ws.ResponsesDir(), 0o750); err != nil {
		return FileReference{}, fmt.Errorf("create responses dir: %w", err)
	}

	suf := randomSuffix(randSuffixLen)
	fname := fmt.Sprintf("response-fragment-%s.json", suf)
	fp := filepath.Join(rs.ws.ResponsesDir(), fname)

	if err := os.WriteFile(fp, body, 0o600); err != nil {
		return FileReference{}, fmt.Errorf("write response fragment file: %w", err)
	}

	size := formatSize(len(body))
	maxSizeStr := formatSize(int(rs.ctx.maxResponseSize.Load()))
	msg := fmt.Sprintf(
		"Response fragment (%s) saved to disk. Use the path with response_slice or response_outline.",
		size,
	)

	return FileReference{
		Path:        fp,
		Size:        len(body),
		SizeHint:    size,
		MaxSizeHint: maxSizeStr,
		Message:     msg,
		OpenCmd:     openCommand(fp),
	}, nil
}
