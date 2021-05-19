package utils

import "time"

// Timeout executes the provided todo function, and waits for it to return. If
// the function does not return before the waitTime duration is elapsed, then
// onTimeout is executed, passing a channel that will be closed when the
// function returns.
func Timeout(todo func(), waitTime time.Duration, onTimeout func(done chan struct{})) {
	done := make(chan struct{})

	go func() {
		todo()
		close(done)
	}()

	select {
	case <-done: // on time, just exit
	case <-time.After(waitTime):
		onTimeout(done)
	}
}
