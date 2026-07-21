package auth

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoAuthClient_Apply(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/api", nil)
	req.Header.Set("X-Custom", "should-stay")

	var info Info
	require.NoError(t, client.Apply(req, &info), "Apply()")

	assert.Empty(t, req.Header.Get(headerAuthorization))
	assert.Equal(t, "should-stay", req.Header.Get("X-Custom"))
	assert.Nil(t, info.Headers)
	assert.Nil(t, info.QueryParams)
}
