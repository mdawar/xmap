package xmap

import "time"

// Time represents a time source.
type Time interface {
	// Now returns the current time.
	Now() time.Time
	// NewTicker returns a new [Ticker] containing a channel that
	// will send the current time on the channel after each tick.
	NewTicker(time.Duration) Ticker
}

// A Ticker holds a channel that delivers "ticks" of a clock at intervals.
type Ticker interface {
	// C returns the channel on which the ticks are delivered.
	C() <-chan time.Time
	// Stop turns off the ticker.
	Stop()
}

var _ Time = (*systemTime)(nil)

// systemTime is a time source using the system time.
type systemTime struct{}

// Now returns the current system local time.
func (t *systemTime) Now() time.Time {
	return time.Now()
}

// NewTicker returns a new system time [Ticker].
func (t *systemTime) NewTicker(d time.Duration) Ticker {
	return &systemTicker{time.NewTicker(d)}
}

var _ Ticker = (*systemTicker)(nil)

// systemTicker is a wrapper for [time.Ticker] that implements the [Ticker] interface.
type systemTicker struct {
	// The wrapped system time ticker.
	ticker *time.Ticker
}

// C returns the channel on which the "ticks" are delivered.
func (t *systemTicker) C() <-chan time.Time {
	return t.ticker.C
}

// Stop turns off the ticker.
func (t *systemTicker) Stop() {
	t.ticker.Stop()
}
