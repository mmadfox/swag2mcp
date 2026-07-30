/*
 * SPDX-FileCopyrightText: © 2025-2026 Sergey "mmadfox" Liskonog
 * SPDX-License-Identifier: Apache-2.0
 */

// Package ratelimit provides a per-endpoint rate limiter for API invoke operations.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// defaultInterval is the default per-endpoint cooldown period (10s).
const defaultInterval = 10 * time.Second

// Limiter checks whether an endpoint ID is allowed to proceed.
type Limiter interface {
	Allow(endpointID string) error
}

// RateLimiter enforces a per-endpoint cooldown period.
type RateLimiter struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
	now      func() time.Time
}

// New creates a new RateLimiter with the default interval (10s).
func New() *RateLimiter {
	return &RateLimiter{
		last:     make(map[string]time.Time),
		interval: defaultInterval,
		now:      time.Now,
	}
}

// NewWithInterval creates a new RateLimiter with the given interval.
func NewWithInterval(interval time.Duration) *RateLimiter {
	if interval <= 0 {
		return New()
	}
	return &RateLimiter{
		last:     make(map[string]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Allow returns nil if the endpointID is allowed, or an error if it is rate-limited.
func (rl *RateLimiter) Allow(endpointID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()
	if last, ok := rl.last[endpointID]; ok && now.Sub(last) < rl.interval {
		remaining := rl.interval - now.Sub(last)
		return fmt.Errorf(
			"rate limit exceeded for endpoint %q: try again in %.0f seconds",
			endpointID, remaining.Seconds(),
		)
	}

	rl.last[endpointID] = now
	return nil
}
