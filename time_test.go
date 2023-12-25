package xmap_test

import (
	"sync"
	"time"

	"github.com/mdawar/xmap"
)

var _ xmap.Time = (*mockTime)(nil)

// mockTime is a mock time source.
type mockTime struct {
	// Mutex to synchronize access to the tickers slice.
	sync.RWMutex
	// Mocked current time.
	mock time.Time
	// Tickers created by NewTicker.
	tickers []*mockTicker
}

// newMockTime creates a new mock time source for testing.
func newMockTime(now time.Time) *mockTime {
	return &mockTime{mock: now}
}

// Now returns the mocked current time.
func (mt *mockTime) Now() time.Time {
	return mt.mock
}

// Advance advances the mocked time by the specified duration.
func (mt *mockTime) Advance(d time.Duration) {
	mt.mock = mt.mock.Add(d)
}

// Set changes the mocked time returned by the Now method to the specified time.
func (mt *mockTime) Set(t time.Time) {
	mt.mock = t
}

// NewTicker returns an [xmap.Ticker] that sends ticks with the mocked time when Tick is called.
func (mt *mockTime) NewTicker(d time.Duration) xmap.Ticker {
	ticker := newMockTicker()
	mt.Lock()
	mt.tickers = append(mt.tickers, ticker)
	mt.Unlock()
	return ticker
}

// Tick sends a tick on all the tickers created with NewTicker.
//
// The ticks are sent manually, no need to overcomplicate things.
//
// The mock time should be set first using Advance() or Set() then the tick should be sent
// using this method which sends a tick on all the created tickers with the current time.
func (mt *mockTime) Tick() {
	mt.RLock()
	defer mt.RUnlock()

	for _, ticker := range mt.tickers {
		ticker.Tick(mt.mock)
	}
}

// mockTicker is a mock [xmap.Ticker] for testing.
type mockTicker struct {
	// Channel on which the ticks are delivered.
	c chan time.Time
}

// newMockTicker creates a new mock ticker for testing.
func newMockTicker() *mockTicker {
	// Use a buffered channel like [time.Ticker].
	return &mockTicker{make(chan time.Time, 1)}
}

// C returns the channel on which the ticks are delivered.
func (mt *mockTicker) C() <-chan time.Time {
	return mt.c
}

func (mt *mockTicker) Stop() {}

// Tick sends a tick on the ticker channel with the specified time.
func (mt *mockTicker) Tick(now time.Time) {
	mt.c <- now
}
