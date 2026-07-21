package cache

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import "fmt"

// LocationError describes an error accessing a location.
type LocationError struct {
	Location string
	Type     string // "file" or "url"
	Err      error
}

// Error returns a human-readable description of the location error.
func (e *LocationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Err)
}

// Unwrap returns the underlying error for use with [errors.Is] and [errors.As].
func (e *LocationError) Unwrap() error {
	return e.Err
}
