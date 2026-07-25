package service

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

// ResponseOutline returns a high-level structural summary of a saved response.
func (s *Service) ResponseOutline(_ context.Context, req ResponseOutlineRequest) (ResponseOutlineResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return ResponseOutlineResponse{}, NewValidationError(
			"The request is invalid — ensure path is provided and points to a saved response file.",
			err,
		)
	}

	r := reader.New(s.ws.ResponsesDir())
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
func (s *Service) ResponseCompress(_ context.Context, req ResponseCompressRequest) (ResponseCompressResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return ResponseCompressResponse{}, NewValidationError(
			"The request is invalid — ensure path and mode are provided.",
			err,
		)
	}

	r := reader.New(s.ws.ResponsesDir())
	result, err := r.Compress(req.Path, reader.CompressOptions{
		JSONPath:   req.JSONPath,
		Mode:       req.Mode,
		ArrayHead:  req.ArrayHead,
		ArrayTail:  req.ArrayTail,
		StringLen:  req.StringLen,
		SelectKeys: req.SelectKeys,
		Limit:      s.maxResponseSize,
	})
	if err != nil {
		return ResponseCompressResponse{}, mapReaderError(err)
	}

	if result.TooLarge {
		ref, saveErr := s.saveReaderResult(req.Path, result.Body)
		if saveErr != nil {
			return ResponseCompressResponse{}, NewInvokeError(
				"The compressed result is still too large and could not be saved to disk.",
				saveErr,
			)
		}
		return ResponseCompressResponse{FileRef: &ref, Hint: result.Hint}, nil
	}

	return ResponseCompressResponse{Body: result.Body, Hint: result.Hint}, nil
}

// ResponseSlice extracts a fragment of a saved response file.
func (s *Service) ResponseSlice(_ context.Context, req ResponseSliceRequest) (ResponseSliceResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return ResponseSliceResponse{}, NewValidationError(
			"The request is invalid — ensure path is provided and at least one of jsonPath, line, or range is set.",
			err,
		)
	}

	r := reader.New(s.ws.ResponsesDir())
	slice, err := r.Slice(req.Path, reader.SliceOptions{
		JSONPath: req.JSONPath,
		Line:     req.Line,
		Range:    req.Range,
		Around:   req.Around,
		Limit:    s.maxResponseSize,
	})
	if err != nil {
		return ResponseSliceResponse{}, mapReaderError(err)
	}

	fragmentBytes, _ := json.Marshal(slice.Value)
	if s.maxResponseSize > 0 && len(fragmentBytes) > s.maxResponseSize {
		ref, saveErr := s.saveReaderResult(req.Path, fragmentBytes)
		if saveErr != nil {
			return ResponseSliceResponse{}, NewInvokeError(
				"The extracted fragment is too large and could not be saved to disk.",
				saveErr,
			)
		}
		return ResponseSliceResponse{Slice: slice, FileRef: &ref}, nil
	}

	return ResponseSliceResponse{Slice: slice}, nil
}

// saveReaderResult saves an extracted or compressed JSON fragment to disk.
func (s *Service) saveReaderResult(_ string, data any) (FileReference, error) {
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

	if err := os.MkdirAll(s.ws.ResponsesDir(), 0o750); err != nil {
		return FileReference{}, fmt.Errorf("create responses dir: %w", err)
	}

	suf := randomSuffix(randSuffixLen)
	fname := fmt.Sprintf("response-fragment-%s.json", suf)
	fp := filepath.Join(s.ws.ResponsesDir(), fname)

	if err := os.WriteFile(fp, body, 0o600); err != nil {
		return FileReference{}, fmt.Errorf("write response fragment file: %w", err)
	}

	size := formatSize(len(body))
	maxSizeStr := formatSize(s.maxResponseSize)
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
