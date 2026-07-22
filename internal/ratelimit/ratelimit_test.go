package ratelimit

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	rl := New()
	require.NotNil(t, rl)
	assert.Equal(t, 10*time.Second, rl.interval)
	assert.NotNil(t, rl.last)
}

func TestNewWithInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
		want     time.Duration
	}{
		{name: "positive interval", interval: 5 * time.Second, want: 5 * time.Second},
		{name: "zero interval falls back to default", interval: 0, want: defaultInterval},
		{name: "negative interval falls back to default", interval: -1, want: defaultInterval},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rl := NewWithInterval(tt.interval)
			require.NotNil(t, rl)
			assert.Equal(t, tt.want, rl.interval)
		})
	}
}

func TestAllow_FirstCall(t *testing.T) {
	t.Parallel()

	rl := New()
	err := rl.Allow("endpoint-1")
	assert.NoError(t, err)
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	t.Parallel()

	rl := NewWithInterval(1 * time.Hour)
	require.NoError(t, rl.Allow("endpoint-1"))

	err := rl.Allow("endpoint-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Contains(t, err.Error(), "endpoint-1")
}

func TestAllow_DifferentEndpoints(t *testing.T) {
	t.Parallel()

	rl := NewWithInterval(1 * time.Hour)
	require.NoError(t, rl.Allow("endpoint-a"))
	require.NoError(t, rl.Allow("endpoint-b"))
	require.NoError(t, rl.Allow("endpoint-c"))
}

func TestAllow_AfterIntervalExpires(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := &RateLimiter{
		last:     make(map[string]time.Time),
		interval: 10 * time.Second,
		now:      func() time.Time { return now },
	}
	require.NoError(t, rl.Allow("endpoint-1"))

	now = now.Add(15 * time.Second)

	err := rl.Allow("endpoint-1")
	assert.NoError(t, err)
}

func TestAllow_ErrorMessageFormat(t *testing.T) {
	t.Parallel()

	rl := NewWithInterval(1 * time.Hour)
	require.NoError(t, rl.Allow("my-endpoint"))

	err := rl.Allow("my-endpoint")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `"my-endpoint"`)
	assert.Contains(t, err.Error(), "seconds")
}

func TestLimiterInterface(t *testing.T) {
	t.Parallel()

	var _ Limiter = (*RateLimiter)(nil)
}
