package reader

import "errors"

var (
	// ErrFileNotFound indicates the requested response file does not exist.
	ErrFileNotFound = errors.New("response file not found")
	// ErrPathNotAllowed indicates the path is outside the responses directory.
	ErrPathNotAllowed = errors.New("path is not inside the responses directory")
	// ErrInvalidJSONPath indicates the provided jsonPath syntax is not supported.
	ErrInvalidJSONPath = errors.New("invalid jsonPath")
	// ErrPathNotFound indicates the jsonPath did not resolve to any value in the file.
	ErrPathNotFound = errors.New("jsonPath did not match any value")
	// ErrInvalidLineRange indicates the requested line range is malformed or out of bounds.
	ErrInvalidLineRange = errors.New("invalid line range")
	// ErrNotJSON indicates the file is not valid JSON.
	ErrNotJSON = errors.New("file is not valid JSON")
)
