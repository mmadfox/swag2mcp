package service

import (
	"testing"
)

func TestInvokeRateLimiter_FirstCallAllowed(t *testing.T) {
	t.Parallel()

	rl := newInvokeRateLimiter()
	if err := rl.allow("ep-1"); err != nil {
		t.Fatalf("allow() = %v, want nil", err)
	}
}

func TestInvokeRateLimiter_SecondCallBlocked(t *testing.T) {
	t.Parallel()

	rl := newInvokeRateLimiter()
	if err := rl.allow("ep-1"); err != nil {
		t.Fatalf("first allow() = %v", err)
	}
	if err := rl.allow("ep-1"); err == nil {
		t.Fatal("second allow() = nil, want error")
	}
}

func TestInvokeRateLimiter_DifferentEndpoints(t *testing.T) {
	t.Parallel()

	rl := newInvokeRateLimiter()
	if err := rl.allow("ep-1"); err != nil {
		t.Fatalf("allow(ep-1) = %v", err)
	}
	if err := rl.allow("ep-2"); err != nil {
		t.Fatalf("allow(ep-2) = %v", err)
	}
}
