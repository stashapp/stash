package config

import (
	"testing"
	"time"
)

func TestSetGuestUserEnabledDoesNotDeadlock(t *testing.T) {
	i := InitializeEmpty()

	done := make(chan struct{})
	go func() {
		i.SetGuestUserEnabled(true)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("SetGuestUserEnabled deadlocked")
	}

	if !i.GetGuestUserEnabled() {
		t.Error("expected guest user to be enabled")
	}
}
