package cache

import "fmt"

// LocationError describes an error accessing a location.
type LocationError struct {
	Location string
	Type     string // "file" or "url"
	Err      error
}

func (e *LocationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Err)
}

func (e *LocationError) Unwrap() error {
	return e.Err
}
