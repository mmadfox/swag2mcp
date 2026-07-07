package service

import (
	"fmt"
	"sync"
	"time"
)

const invokeRateInterval = 10 * time.Second

type invokeRateLimiter struct {
	mu   sync.Mutex
	last map[string]time.Time
}

func newInvokeRateLimiter() *invokeRateLimiter {
	return &invokeRateLimiter{
		last: make(map[string]time.Time),
	}
}

func (rl *invokeRateLimiter) allow(endpointID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if last, ok := rl.last[endpointID]; ok && now.Sub(last) < invokeRateInterval {
		remaining := invokeRateInterval - now.Sub(last)
		return fmt.Errorf(
			"rate limit exceeded for endpoint %q: try again in %.0f seconds",
			endpointID, remaining.Seconds(),
		)
	}

	rl.last[endpointID] = now
	return nil
}
