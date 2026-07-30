/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

package cache

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	cli := defaultHTTPClient()
	require.NotNil(t, cli, "defaultHTTPClient() returned nil")
	require.NotNil(t, cli.client, "http.Client is nil")
}

func TestHTTPClientGet_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	data, err := cli.Get(context.Background(), srv.URL)
	require.NoError(t, err, "Get()")
	assert.Equal(t, "hello world", string(data))
}

func TestHTTPClientGet_Non200(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	_, err := cli.Get(context.Background(), srv.URL)
	require.Error(t, err, "expected error for 404")
}

func TestHTTPClientGet_EmptyBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cli := defaultHTTPClient()
	_, err := cli.Get(context.Background(), srv.URL)
	require.Error(t, err, "expected error for empty body")
}

func TestHTTPClientSetClient(t *testing.T) {
	t.Parallel()

	cli := defaultHTTPClient()
	newCli := &http.Client{Timeout: 10}
	cli.SetClient(newCli)
	assert.Same(t, newCli, cli.client)
}
