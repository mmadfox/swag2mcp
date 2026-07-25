package cache

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeta(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	metaPath := filepath.Join(dir, "test.meta")

	m := fileMeta{
		Source:     "https://example.com/spec",
		SourceType: "url",
		CachedAt:   time.Now(),
		TTLSec:     3600,
	}
	require.NoError(t, writeMeta(metaPath, m))

	got, err := readMeta(metaPath)
	require.NoError(t, err)

	assert.Equal(t, m.Source, got.Source)
	assert.Equal(t, m.SourceType, got.SourceType)
	assert.Equal(t, m.TTLSec, got.TTLSec)
}

func TestMeta_expired(t *testing.T) {
	t.Parallel()
	m := fileMeta{
		CachedAt: time.Now().Add(-2 * time.Hour),
		TTLSec:   3600,
	}
	assert.True(t, m.IsExpired(), "expected expired meta")

	m2 := fileMeta{
		CachedAt: time.Now().Add(-30 * time.Minute),
		TTLSec:   3600,
	}
	assert.False(t, m2.IsExpired(), "expected non-expired meta")
}

func TestMeta_readNotFound(t *testing.T) {
	t.Parallel()
	_, err := readMeta("/nonexistent/path")
	require.Error(t, err, "expected error for non-existent meta")
}
