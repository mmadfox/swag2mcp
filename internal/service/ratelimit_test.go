package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvokeRateLimiter_FirstCallAllowed(t *testing.T) {
	t.Parallel()

	rl := newInvokeRateLimiter()
	require.NoError(t, rl.allow("ep-1"))
}

func TestInvokeRateLimiter_SecondCallBlocked(t *testing.T) {
	t.Parallel()

	rl := newInvokeRateLimiter()
	require.NoError(t, rl.allow("ep-1"))
	require.Error(t, rl.allow("ep-1"))
}

func TestInvokeRateLimiter_DifferentEndpoints(t *testing.T) {
	t.Parallel()

	rl := newInvokeRateLimiter()
	require.NoError(t, rl.allow("ep-1"))
	require.NoError(t, rl.allow("ep-2"))
}
