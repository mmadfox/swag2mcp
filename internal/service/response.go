/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/mmadfox/swag2mcp/internal/config"
	"github.com/mmadfox/swag2mcp/internal/reader"
	"github.com/tidwall/gjson"
)

const (
	streamThreshold = 10 * 1024 * 1024 // 10 MB — switch to streaming strategy
	maxPageSize     = 50
	maxPage         = 100
)

// ResponseFilterRequest requests filtering, searching, and pagination of a saved response file.
type ResponseFilterRequest struct {
	Path     string `json:"path"     validate:"required"            jsonschema:"required,Absolute path returned in fileRef.path"`
	JSONPath string `json:"jsonPath,omitempty"   jsonschema:"optional,Path to the array to filter (e.g. pets, data.items). Leave empty for root arrays."`
	Search   string `json:"search,omitempty"   jsonschema:"optional,Full-text search across all fields of each item"`
	Filter   string `json:"filter,omitempty"   jsonschema:"optional,Structured filter condition (e.g. status = active, price > 100)"`
	Page     int    `json:"page,omitempty"     jsonschema:"optional,Page number starting from 1 (default 1)"`
	PageSize int    `json:"pageSize,omitempty" jsonschema:"optional,Items per page (max 50, default 10)"`
}

// ResponseFilterResponse returns filtered and paginated results.
type ResponseFilterResponse struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Total      int    `json:"total"`
	TotalPages int    `json:"totalPages"`
	Items      []any  `json:"items"`
	Strategy   string `json:"strategy"`
}

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
	log *slog.Logger
}

func newResponseService(
	ctx *serviceContext,
	ws WorkspaceOps,
	v RequestValidator,
	log *slog.Logger,
) *responseService {
	return &responseService{ctx: ctx, ws: ws, v: v, log: log}
}

// ResponseOutline returns a high-level structural summary of a saved response.
func (rs *responseService) ResponseOutline(ctx context.Context, req ResponseOutlineRequest) (ResponseOutlineResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseOutlineResponse{}, NewResponseRequestError(err)
	}

	r := reader.New(rs.ws.ResponsesDir())
	outline, err := r.Outline(req.Path, reader.OutlineOptions{
		MaxDepth:      req.MaxDepth,
		MaxArrayItems: req.MaxArrayItems,
	})
	if err != nil {
		rs.log.ErrorContext(ctx, "response_outline failed", "path", req.Path, "error", err)
		return ResponseOutlineResponse{}, mapReaderError(err)
	}

	return ResponseOutlineResponse{Outline: outline}, nil
}

// ResponseCompress reduces a JSON value in a saved response file.
func (rs *responseService) ResponseCompress(ctx context.Context, req ResponseCompressRequest) (ResponseCompressResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseCompressResponse{}, NewResponseRequestError(err)
	}

	r := reader.New(rs.ws.ResponsesDir())
	result, err := r.Compress(req.Path, reader.CompressOptions{
		JSONPath:   req.JSONPath,
		Mode:       req.Mode,
		ArrayHead:  req.ArrayHead,
		ArrayTail:  req.ArrayTail,
		StringLen:  req.StringLen,
		SelectKeys: req.SelectKeys,
		Limit:      rs.ctx.MaxResponseSize(),
	})
	if err != nil {
		rs.log.ErrorContext(ctx, "response_compress failed", "path", req.Path, "json_path", req.JSONPath, "error", err)
		return ResponseCompressResponse{}, mapReaderError(err)
	}

	if result.TooLarge {
		ref, saveErr := rs.saveReaderResult(req.Path, result.Body)
		if saveErr != nil {
			rs.log.ErrorContext(ctx, "response_compress failed: save result", "path", req.Path, "error", saveErr)
			return ResponseCompressResponse{}, NewResponseReadFailedError(saveErr)
		}
		return ResponseCompressResponse{FileRef: &ref, Hint: result.Hint}, nil
	}

	return ResponseCompressResponse{Body: result.Body, Hint: result.Hint}, nil
}

// ResponseSlice extracts a fragment of a saved response file.
func (rs *responseService) ResponseSlice(ctx context.Context, req ResponseSliceRequest) (ResponseSliceResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseSliceResponse{}, NewResponseRequestError(err)
	}

	r := reader.New(rs.ws.ResponsesDir())
	slice, err := r.Slice(req.Path, reader.SliceOptions{
		JSONPath: req.JSONPath,
		Line:     req.Line,
		Range:    req.Range,
		Around:   req.Around,
		Limit:    rs.ctx.MaxResponseSize(),
	})
	if err != nil {
		rs.log.ErrorContext(ctx, "response_slice failed", "path", req.Path, "json_path", req.JSONPath, "error", err)
		return ResponseSliceResponse{}, mapReaderError(err)
	}

	maxSize := rs.ctx.MaxResponseSize()
	fragmentBytes, _ := json.Marshal(slice.Value)
	if maxSize > 0 && len(fragmentBytes) > maxSize {
		ref, saveErr := rs.saveReaderResult(req.Path, fragmentBytes)
		if saveErr != nil {
			rs.log.ErrorContext(ctx, "response_slice failed: save result", "path", req.Path, "error", saveErr)
			return ResponseSliceResponse{}, NewResponseReadFailedError(saveErr)
		}
		return ResponseSliceResponse{Slice: slice, FileRef: &ref}, nil
	}

	return ResponseSliceResponse{Slice: slice}, nil
}

// ResponseFilter filters, searches, and paginates through arrays in saved response files.
func (rs *responseService) ResponseFilter(ctx context.Context, req ResponseFilterRequest) (ResponseFilterResponse, error) {
	if err := rs.v.Struct(req); err != nil {
		return ResponseFilterResponse{}, NewResponseRequestError(err)
	}

	page := max(1, req.Page)
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	pageSize = min(pageSize, maxPageSize)
	page = min(page, maxPage)

	info, err := os.Stat(req.Path)
	if err != nil {
		rs.log.ErrorContext(ctx, "response_filter failed: file not found", "path", req.Path, "error", err)
		return ResponseFilterResponse{}, NewResponseFileNotFoundError(err)
	}

	var (
		items    []any
		total    int
		strategy string
	)

	if info.Size() >= streamThreshold {
		strategy = "streaming"
		items, total, err = rs.filterArrayStreaming(req.Path, req.JSONPath, req.Search, req.Filter, page, pageSize)
	} else {
		strategy = "memory"
		items, total, err = rs.filterArrayInMemory(req.Path, req.JSONPath, req.Search, req.Filter, page, pageSize)
	}
	if err != nil {
		rs.log.ErrorContext(ctx, "response_filter failed", "path", req.Path, "json_path", req.JSONPath, "error", err)
		return ResponseFilterResponse{}, mapReaderError(err)
	}

	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}

	return ResponseFilterResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Items:      items,
		Strategy:   strategy,
	}, nil
}

// filterArrayInMemory loads the entire file into memory and filters the array using gjson.
func (rs *responseService) filterArrayInMemory(
	path, jsonPath, search, filter string,
	page, pageSize int,
) ([]any, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, err
	}

	gp := jsonPath
	if gp == "" {
		gp = "@this"
	}

	result := gjson.GetBytes(data, gp)
	if !result.IsArray() {
		return nil, 0, reader.ErrPathNotFound
	}

	parsedFilter, filterErr := parseFilter(filter)
	if filterErr != nil && filter != "" {
		return nil, 0, filterErr
	}

	var matched []any
	for _, r := range result.Array() {
		item := r.Value()
		if search != "" && !searchItem(item, search) {
			continue
		}
		if parsedFilter != nil && !applyFilter(item, parsedFilter) {
			continue
		}
		matched = append(matched, item)
	}

	total := len(matched)
	offset := (page - 1) * pageSize
	if offset >= total {
		return []any{}, total, nil
	}

	end := min(offset+pageSize, total)

	return matched[offset:end], total, nil
}

// filterArrayStreaming uses [json.Decoder] to stream through the file without loading it entirely.
func (rs *responseService) filterArrayStreaming(
	path, jsonPath, search, filter string,
	page, pageSize int,
) ([]any, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	if err := navigateDecoderToArray(dec, jsonPath); err != nil {
		return nil, 0, err
	}

	parsedFilter, filterErr := parseFilter(filter)
	if filterErr != nil && filter != "" {
		return nil, 0, filterErr
	}

	offset := (page - 1) * pageSize
	end := offset + pageSize

	return streamFilterItems(dec, search, parsedFilter, offset, end)
}

// navigateDecoderToArray advances the decoder through a dotted jsonPath to reach an array.
// An empty jsonPath means the root value is expected to be an array.
func navigateDecoderToArray(dec *json.Decoder, jsonPath string) error {
	if jsonPath == "" {
		tok, err := dec.Token()
		if err != nil {
			return reader.ErrPathNotFound
		}
		if tok != json.Delim('[') {
			return reader.ErrPathNotFound
		}
		return nil
	}

	for part := range strings.SplitSeq(jsonPath, ".") {
		tok, err := dec.Token()
		if err != nil {
			return reader.ErrPathNotFound
		}
		if tok != json.Delim('{') {
			return reader.ErrPathNotFound
		}

		if err := seekObjectKey(dec, part); err != nil {
			return err
		}
	}

	tok, err := dec.Token()
	if err != nil {
		return reader.ErrPathNotFound
	}
	if tok != json.Delim('[') {
		return reader.ErrPathNotFound
	}
	return nil
}

// seekObjectKey advances the decoder through an object to find a specific key.
func seekObjectKey(dec *json.Decoder, key string) error {
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return reader.ErrPathNotFound
		}
		k, ok := keyTok.(string)
		if !ok {
			return reader.ErrPathNotFound
		}
		if k == key {
			return nil
		}
		if err := skipValue(dec); err != nil {
			return err
		}
	}
	return reader.ErrPathNotFound
}

// streamFilterItems reads items from a decoder, applies search/filter, and returns a page.
func streamFilterItems(
	dec *json.Decoder,
	search string,
	parsedFilter *filterCondition,
	offset, end int,
) ([]any, int, error) {
	var matched []any
	index := 0

	for dec.More() {
		var item any
		if err := dec.Decode(&item); err != nil {
			return nil, 0, err
		}

		match := true
		if search != "" {
			match = searchItem(item, search)
		}
		if match && parsedFilter != nil {
			match = applyFilter(item, parsedFilter)
		}

		if match {
			if index >= offset && index < end {
				matched = append(matched, item)
			}
			index++
		}
	}

	return matched, index, nil
}

// searchItem recursively searches all fields of an item for a substring.
func searchItem(item any, query string) bool {
	query = strings.ToLower(query)
	var searchRecursive func(any) bool
	searchRecursive = func(v any) bool {
		switch val := v.(type) {
		case string:
			return strings.Contains(strings.ToLower(val), query)
		case map[string]any:
			for _, fv := range val {
				if searchRecursive(fv) {
					return true
				}
			}
		case []any:
			return slices.ContainsFunc(val, func(ev any) bool {
				return searchRecursive(ev)
			})
		case float64:
			return strings.Contains(strconv.FormatFloat(val, 'f', -1, 64), query)
		case bool:
			return strings.Contains(strings.ToLower(strconv.FormatBool(val)), query)
		}
		return false
	}
	return searchRecursive(item)
}

type filterCondition struct {
	field    string
	operator string
	value    string
}

// parseFilter parses a filter string like "field = value" or "price > 100".
// Returns nil when filter is empty (no filtering).
//
//nolint:nilnil // nil, nil means "no filter" which is semantically correct.
func parseFilter(filter string) (*filterCondition, error) {
	if filter == "" {
		return nil, nil
	}

	operators := []string{">=", "<=", "!=", "=", ">", "<", " contains "}
	for _, op := range operators {
		idx := strings.Index(filter, op)
		if idx < 0 {
			continue
		}
		field := strings.TrimSpace(filter[:idx])
		value := strings.TrimSpace(filter[idx+len(op):])
		op = strings.TrimSpace(op)
		if field == "" {
			return nil, errors.New("invalid filter: empty field name")
		}
		return &filterCondition{field: field, operator: op, value: value}, nil
	}

	return nil, fmt.Errorf("invalid filter: no operator found in %q", filter)
}

// applyFilter checks if an item matches a filter condition.
func applyFilter(item any, cond *filterCondition) bool {
	m, ok := item.(map[string]any)
	if !ok {
		return false
	}

	val, ok := m[cond.field]
	if !ok {
		return false
	}

	switch cond.operator {
	case "=":
		return fmt.Sprintf("%v", val) == cond.value
	case "!=":
		return fmt.Sprintf("%v", val) != cond.value
	case "contains":
		s, ok := val.(string)
		if !ok {
			return false
		}
		return strings.Contains(strings.ToLower(s), strings.ToLower(cond.value))
	case ">", "<", ">=", "<=":
		fv, ok := toFloat64(val)
		if !ok {
			return false
		}
		fc, err := strconv.ParseFloat(cond.value, 64)
		if err != nil {
			return false
		}
		switch cond.operator {
		case ">":
			return fv > fc
		case "<":
			return fv < fc
		case ">=":
			return fv >= fc
		case "<=":
			return fv <= fc
		}
	}

	return false
}

// toFloat64 converts a JSON number to float64 for comparison.
func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

// skipValue skips a single JSON value in a streaming decoder.
func skipValue(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}

	switch tok {
	case json.Delim('{'), json.Delim('['):
		depth := 1
		for depth > 0 {
			tok, err := dec.Token()
			if err != nil {
				return err
			}
			switch tok {
			case json.Delim('{'), json.Delim('['):
				depth++
			case json.Delim('}'), json.Delim(']'):
				depth--
			}
		}
	}

	return nil
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

	suf := randomSuffix(config.RandSuffixLen)
	fname := fmt.Sprintf("response-fragment-%s.json", suf)
	fp := filepath.Join(rs.ws.ResponsesDir(), fname)

	if err := os.WriteFile(fp, body, 0o600); err != nil {
		return FileReference{}, fmt.Errorf("write response fragment file: %w", err)
	}

	size := formatSize(len(body))
	maxSizeStr := formatSize(rs.ctx.MaxResponseSize())
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
