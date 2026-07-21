package reader

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// defaultOutlineMaxDepth is the default recursion depth for outlines.
	defaultOutlineMaxDepth = 3
	// defaultOutlineMaxArrayItems is the default number of array items inspected for type info.
	defaultOutlineMaxArrayItems = 5
	// defaultValuePreviewLen is the default maximum length for scalar previews in outlines.
	defaultValuePreviewLen = 80
	// defaultSliceAround is the default number of lines included around a slice target.
	defaultSliceAround = 20
	// defaultCompressStringLen is the default maximum string length in truncate mode.
	defaultCompressStringLen = 80
	// defaultCompressArrayHead is the default number of leading array items kept in sample mode.
	defaultCompressArrayHead = 3
	// defaultCompressArrayTail is the default number of trailing array items kept in sample mode.
	defaultCompressArrayTail = 2
)

// Reader reads and navigates large JSON response files in a workspace.
type Reader interface {
	// Outline returns a high-level structural summary of the file at path.
	Outline(path string, opts OutlineOptions) (Outline, error)
	// Compress reduces a JSON value selected by jsonPath according to mode.
	Compress(path string, opts CompressOptions) (CompressResult, error)
	// Slice returns a fragment of JSON by jsonPath or line range.
	Slice(path string, opts SliceOptions) (Slice, error)
}

// New creates a Reader rooted at responsesDir.
func New(responsesDir string) Reader {
	return &reader{responsesDir: responsesDir}
}

type reader struct {
	responsesDir string
}

// validatePath ensures path exists and is inside the responses directory.
func (r *reader) validatePath(path string) error {
	absDir, err := filepath.Abs(r.responsesDir)
	if err != nil {
		return fmt.Errorf("resolve responses dir: %w", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve file path: %w", err)
	}

	rel, err := filepath.Rel(absDir, absPath)
	if err != nil {
		return fmt.Errorf("rel responses dir: %w", err)
	}

	// pathEscapesResponsesDir checks whether the relative path leaves the
	// responses directory. rel is the path from responsesDir to the target file.
	// On Unix it is like "get-users-abc.json" when inside, or ".." / "../x"
	// when outside. filepath.IsAbs catches edge cases where Rel cannot produce
	// a relative result.
	pathEscapesResponsesDir := rel == ".." ||
		strings.HasPrefix(rel, ".."+string(filepath.Separator)) ||
		filepath.IsAbs(rel)
	if pathEscapesResponsesDir {
		return ErrPathNotAllowed
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrFileNotFound
		}
		return fmt.Errorf("stat file: %w", err)
	}
	if info.IsDir() {
		return ErrFileNotFound
	}

	return nil
}

// normalizeOutlineOptions fills default values.
func normalizeOutlineOptions(opts OutlineOptions) OutlineOptions {
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = defaultOutlineMaxDepth
	}
	if opts.MaxArrayItems <= 0 {
		opts.MaxArrayItems = defaultOutlineMaxArrayItems
	}
	return opts
}

// normalizeCompressOptions fills default values.
func normalizeCompressOptions(opts CompressOptions) CompressOptions {
	if opts.ArrayHead <= 0 {
		opts.ArrayHead = defaultCompressArrayHead
	}
	if opts.ArrayTail <= 0 {
		opts.ArrayTail = defaultCompressArrayTail
	}
	if opts.StringLen <= 0 {
		opts.StringLen = defaultCompressStringLen
	}
	return opts
}

// normalizeSliceOptions fills default values.
func normalizeSliceOptions(opts SliceOptions) SliceOptions {
	if opts.Around <= 0 {
		opts.Around = defaultSliceAround
	}
	return opts
}
