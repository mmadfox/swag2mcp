/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package auth

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

func TestNoAuthClient_New(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	require.NoError(t, client.New())
}

func TestNoAuthClient_Type(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	assert.Equal(t, NoAuth, client.Type())
}

func TestNoAuthClient_Validate(t *testing.T) {
	t.Parallel()

	client := NewNoAuthClient()
	require.NoError(t, client.Validate())
}
