package api

import (
	"testing"
	"time"
)

func TestRateLimiter_canIncrement(t *testing.T) {
	const limit = 3
	rl := newRateLimiter(limit, time.Minute)
	const key = "127.0.0.1"

	for i := 0; i < limit; i++ {
		ok, err := rl.canIncrement(key)
		if err != nil {
			t.Fatalf("canIncrement() error = %v", err)
		}
		if !ok {
			t.Fatalf("canIncrement() = false after %d increments, want true", i)
		}
		if err := rl.Increment(key); err != nil {
			t.Fatalf("Increment() error = %v", err)
		}
	}

	ok, err := rl.canIncrement(key)
	if err != nil {
		t.Fatalf("canIncrement() error = %v", err)
	}
	if ok {
		t.Errorf("canIncrement() = true after %d increments, want false", limit)
	}
}
