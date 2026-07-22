package cache

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationError_Error(t *testing.T) {
	t.Parallel()

	err := &LocationError{
		Location: "/path/to/spec.yaml",
		Type:     "file",
		Err:      errors.New("file not found"),
	}
	assert.NotEmpty(t, err.Error())
}

func TestLocationError_Unwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner error")
	err := &LocationError{
		Location: "/path",
		Type:     "file",
		Err:      inner,
	}
	assert.True(t, errors.Is(err, inner), "errors.Is() should match inner error")
}
