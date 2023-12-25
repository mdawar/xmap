package xmap_test

import "time"

// retryUntil keeps calling the function f until it returns true or the deadline d has been reached.
func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)

	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
