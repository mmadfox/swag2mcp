// Package ratelimit provides a per-endpoint rate limiter for API invoke operations.
package ratelimit

// SPDX-License-Identifier: AGPL-3.0-only
//
// Use of this software is governed by the AGPL v3 license
// included in the /LICENSE file.

import (
	"fmt"
	"sync"
	"time"
)

const invokeRateInterval = 10 * time.Second

// Limiter checks whether an endpoint ID is allowed to proceed.
type Limiter interface {
	Allow(endpointID string) error
}

// RateLimiter enforces a per-endpoint cooldown period.
type RateLimiter struct {
	mu   sync.Mutex
	last map[string]time.Time
}

// New creates a new RateLimiter.
func New() *RateLimiter {
	return &RateLimiter{
		last: make(map[string]time.Time),
	}
}

// Allow returns nil if the endpointID is allowed, or an error if it is rate-limited.
func (rl *RateLimiter) Allow(endpointID string) error {
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
