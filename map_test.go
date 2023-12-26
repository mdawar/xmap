package xmap_test

import (
	"context"
	"maps"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/goleak"

	"github.com/mdawar/xmap"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMapSetThenGet(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	keyName := "abc"
	wantValue := 10

	m.Set(keyName, wantValue, 0)

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	gotValue, ok := m.Get(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %d, got %d", wantValue, gotValue)
	}
}

func TestMapSetThenGetWithExpiration(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, string]()
	defer m.Stop()

	keyName := "keyName"
	wantValue := "Testing"

	m.Set(keyName, wantValue, time.Hour)

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	gotValue, gotExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %q, got %q", wantValue, gotValue)
	}

	if !gotExpiration.After(time.Now()) {
		t.Errorf("want expiration time in the future, got %v", gotExpiration)
	}
}

func TestMapGetNonExistingKeyReturnsFalse(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	value, ok := m.Get("doesNotExist")
	if ok {
		t.Fatal("want false getting a non existing key, got true")
	}

	if value != 0 {
		t.Errorf("want zero value for non existing key %d, got %d", 0, value)
	}
}

func TestMapGetWithExpirationNonExistingKeyReturnsFalse(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	value, expiration, ok := m.GetWithExpiration("doesNotExist")
	if ok {
		t.Fatal("want false getting a non existing key, got true")
	}

	if value != 0 {
		t.Errorf("want zero value for non existing key %d, got %d", 0, value)
	}

	if !expiration.IsZero() {
		t.Errorf("want zero time value expiration for non existing key, got %v", expiration)
	}
}

func TestMapSetWithZeroTTLThenGetWithExpiration(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	keyName := "keyName"
	wantValue := 123456

	m.Set(keyName, wantValue, 0) // Non expiring key (0 TTL).

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	gotValue, gotExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %d, got %d", wantValue, gotValue)
	}

	if !gotExpiration.IsZero() {
		t.Errorf("want key with 0 TTL to have zero expiration time, got %v", gotExpiration)
	}
}

func TestMapSetReplacesExistingValue(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	keyName := "abc:123"
	wantValue := 123

	m.Set(keyName, wantValue, 0) // Non expiring key (0 TTL).

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	gotValue, gotExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %d, got %d", wantValue, gotValue)
	}

	// Replace the key with a new value and expiration time.
	wantNewValue := 456
	m.Set(keyName, wantNewValue, time.Hour)

	gotNewValue, gotNewExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantNewValue != gotNewValue {
		t.Errorf("want new value %d, got %d", wantNewValue, gotNewValue)
	}

	if gotNewExpiration.Equal(gotExpiration) {
		t.Errorf("want different expiration time, got same expiration time %v", gotExpiration)
	}

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}
}

func TestMapUpdateNonExistingKeyReturnsFalse(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	if ok := m.Update("doesNotExist", 100); ok {
		t.Error("want false updating a non existin key, got true")
	}

	if m.Len() != 0 {
		t.Errorf("want map length %d, got %d", 0, m.Len())
	}
}

func TestMapUpdateReplacesTheValueOnly(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	keyName := "abc"
	wantValue := 111

	m.Set(keyName, wantValue, time.Minute)

	gotValue, gotExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %d, got %d", wantValue, gotValue)
	}

	// Update the value and keep the expiration time.
	wantNewValue := 999
	if ok := m.Update(keyName, wantNewValue); !ok {
		t.Fatal("key was not updated, does not exist")
	}

	gotNewValue, gotNewExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantNewValue != gotNewValue {
		t.Errorf("want new value %d, got %d", wantNewValue, gotNewValue)
	}

	if !gotNewExpiration.Equal(gotExpiration) {
		t.Errorf("want same expiration time after value update %v, got %v", gotExpiration, gotNewExpiration)
	}

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}
}

func TestMapDeleteRemovesTheKey(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	keyName := "abc123"

	m.Set(keyName, 5, time.Hour)

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	m.Delete(keyName)

	if _, ok := m.Get(keyName); ok {
		t.Errorf("key %q was not removed from the map", keyName)
	}

	if m.Len() != 0 {
		t.Fatalf("want map length %d, got %d", 0, m.Len())
	}
}

func TestMapClearRemovesAllTheKeys(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	entries := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	wantLen := len(entries)

	for k, v := range entries {
		m.Set(k, v, time.Hour)
	}

	if wantLen != m.Len() {
		t.Fatalf("want map length %d, got %d", wantLen, m.Len())
	}

	m.Clear()

	for k := range entries {
		if _, ok := m.Get(k); ok {
			t.Errorf("key %q was not removed from the map", k)
		}
	}

	if m.Len() != 0 {
		t.Fatalf("want map length %d, got %d", 0, m.Len())
	}
}

func TestMapKeyExpirationAndCleanup(t *testing.T) {
	t.Parallel()

	m := xmap.NewWithConfig[string, int](xmap.Config{
		CleanupInterval: 50 * time.Millisecond,
	})
	defer m.Stop()

	m.Set("a", 1, 10*time.Millisecond)
	m.Set("b", 2, 0) // Never expires.
	m.Set("c", 3, 30*time.Millisecond)
	m.Set("d", 4, 50*time.Millisecond)

	if m.Len() != 4 {
		t.Fatalf("want map length %d, got %d", 4, m.Len())
	}

	checkIfExpired := func(key string) {
		t.Helper()
		if _, exp, ok := m.GetWithExpiration(key); ok {
			t.Errorf("key %q with expiration time %v did not expire at %v", key, exp, time.Now())
		}
	}

	time.Sleep(15 * time.Millisecond)
	checkIfExpired("a")

	time.Sleep(20 * time.Millisecond) // 35 Milliseconds passed.
	checkIfExpired("c")

	time.Sleep(20 * time.Millisecond) // 55 Milliseconds passed.
	checkIfExpired("d")

	if value, ok := m.Get("b"); !ok {
		t.Errorf("key %q with 0 TTL must not be removed from the map", "b")
	} else if value != 2 {
		t.Errorf("want key %q value %d, got %d", "b", 2, value)
	}

	// Only the key with 0 TTL should be left in the map.
	if m.Len() != 1 {
		t.Errorf("want map length %d, got %d", 1, m.Len())
	}
}

func TestMapCallingStopMultipleTimesDoesNotPanic(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()

	if m.Stopped() {
		t.Fatal("map is stopped before calling Stop")
	}

	m.Stop()

	if !m.Stopped() {
		t.Fatal("map was not stopped after calling Stop")
	}

	// Stop should be safe to be called multiple times.
	m.Stop()
}

func TestMapStopClearsTheMap(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	entries := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	wantLen := len(entries)

	for k, v := range entries {
		m.Set(k, v, time.Hour)
	}

	if wantLen != m.Len() {
		t.Fatalf("want map length %d, got %d", wantLen, m.Len())
	}

	// Should stop the cleanup goroutine and clear the map.
	m.Stop()

	for k := range entries {
		if _, ok := m.Get(k); ok {
			t.Errorf("key %q was not removed from the map", k)
		}
	}

	if m.Len() != 0 {
		t.Fatalf("want map length %d, got %d", 0, m.Len())
	}
}

func TestMapGetAndUpdateExpiredKey(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testTime := newMockTime(now)

	m := xmap.NewWithConfig[string, int](xmap.Config{
		TimeSource: testTime,
	})
	defer m.Stop()

	keyName := "abc:123"
	wantValue := 1122
	wantExpiration := now.Add(time.Hour)

	m.Set(keyName, wantValue, time.Hour)

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	gotValue, gotExpiration, ok := m.GetWithExpiration(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %d, got %d", wantValue, gotValue)
	}

	// Verify the exact expiration time.
	if !wantExpiration.Equal(gotExpiration) {
		t.Fatalf("want expiration time %v, got %v", wantExpiration, gotExpiration)
	}

	// Set the current time to the exact expiration time.
	testTime.Set(wantExpiration)

	if _, ok := m.Get(keyName); !ok {
		t.Errorf("key %q should not expire at the exact expiration time", keyName)
	}

	// Advance the time by 1 nanosecond to make the key expire.
	testTime.Advance(time.Nanosecond)

	// The key should have expired.
	if gotValue, ok := m.Get(keyName); ok {
		t.Errorf("key %q did not expire on time", keyName)
	} else if gotValue != 0 {
		t.Errorf("expired key %q should return zero value %d, got %d", keyName, 0, gotValue)
	}

	// We should not be able to update the key after it expires.
	if ok := m.Update(keyName, 100); ok {
		t.Errorf("key %q should not be updated on expiration", keyName)
	}
}

func TestMapKeyExpirationAndRemoval(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testTime := newMockTime(now)

	m := xmap.NewWithConfig[string, int](xmap.Config{
		TimeSource: testTime,
	})
	defer m.Stop()

	// Wait until the cleanup goroutine is active.
	if isActive := retryUntil(20*time.Millisecond, func() bool {
		return m.CleanupActive()
	}); !isActive {
		t.Fatal("cleanup goroutine did not start in time")
	}

	keyName := "abc"
	wantValue := 11
	wantExpiration := now.Add(time.Hour)

	m.Set(keyName, wantValue, time.Hour)

	if m.Len() != 1 {
		t.Fatalf("want map length %d, got %d", 1, m.Len())
	}

	if _, gotExpiration, ok := m.GetWithExpiration(keyName); !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	} else if !wantExpiration.Equal(gotExpiration) {
		t.Fatalf("want expiration time %v, got %v", wantExpiration, gotExpiration)
	}

	// Set the current time to the exact expiration time.
	testTime.Set(wantExpiration)

	if _, ok := m.Get(keyName); !ok {
		t.Errorf("key %q should not expire at the exact expiration time", keyName)
	}

	// Advance the time by 1 nanosecond to make the key expire.
	testTime.Advance(time.Nanosecond)

	// The key should have expired.
	if _, ok := m.Get(keyName); ok {
		t.Errorf("key %q did not expire on time", keyName)
	}

	// Send a tick on the created tickers.
	// The cleanup goroutine must be ready to receive before sending the tick.
	testTime.Tick()

	// Wait until the key is removed.
	if keyRemoved := retryUntil(time.Second, func() bool {
		return m.Len() == 0
	}); !keyRemoved {
		t.Errorf("want map length %d, got %d", 0, m.Len())
	}
}

func TestMapKeyWithZeroTTLNeverExpires(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testTime := newMockTime(now)

	m := xmap.NewWithConfig[string, int](xmap.Config{
		TimeSource: testTime,
	})
	defer m.Stop()

	// Wait until the cleanup goroutine is active.
	if isActive := retryUntil(20*time.Millisecond, func() bool {
		return m.CleanupActive()
	}); !isActive {
		t.Fatal("cleanup goroutine did not start in time")
	}

	keyName := "abc"
	wantValue := 11

	m.Set(keyName, wantValue, 0) // Never expires (0 TTL).
	// Add another key that expires so we can test the length of the map at the end.
	// Needed to be able to wait for the cleanup to finish, otherwise we would have a flaky test.
	m.Set("expiringKey", 1, time.Hour)

	if m.Len() != 2 {
		t.Fatalf("want map length %d, got %d", 2, m.Len())
	}

	if _, gotExpiration, ok := m.GetWithExpiration(keyName); !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	} else if !gotExpiration.IsZero() {
		t.Errorf("want zero time value expiration for key with 0 TTL, got %v", gotExpiration)
	}

	// Advance the time 1 year.
	testTime.Advance(24 * 365 * time.Hour)

	// Send a tick on the created tickers.
	// The cleanup goroutine must be ready to receive before sending the tick.
	testTime.Tick()

	// Wait until 1 key is removed.
	if cleanupDone := retryUntil(time.Second, func() bool {
		return m.Len() == 1
	}); !cleanupDone {
		t.Errorf("want map length %d after cleanup, got %d", 1, m.Len())
	}

	if _, ok := m.Get(keyName); !ok {
		t.Errorf("key %q with 0 TTL should not expire", keyName)
	}
}

func TestMapManualExpiredKeysRemoval(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testTime := newMockTime(now)

	m := xmap.NewWithConfig[string, int](xmap.Config{
		TimeSource: testTime,
	})
	// Since we're using a mock time source, the cleanup goroutine
	// won't receive ticks unless we send them.
	defer m.Stop()

	if removed := m.RemoveExpired(); removed != 0 {
		t.Fatalf("want %d key removals for empty map, got %d", 0, removed)
	}

	entries := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	wantLen := len(entries)

	for k, v := range entries {
		m.Set(k, v, time.Hour)
	}

	if wantLen != m.Len() {
		t.Fatalf("want map length %d, got %d", wantLen, m.Len())
	}

	// Advance to the exact expiration time.
	testTime.Advance(time.Hour)

	if removed := m.RemoveExpired(); removed != 0 {
		t.Fatalf("want %d key removals at the exact expiration time, got %d", 0, removed)
	}

	// Advance 1 more nanosecond to make the keys expire.
	testTime.Advance(time.Nanosecond)

	if removed := m.RemoveExpired(); removed != wantLen {
		t.Fatalf("want %d key removals on expiration, got %d", wantLen, removed)
	}

	if m.Len() != 0 {
		t.Fatalf("want map length %d, got %d", 0, m.Len())
	}
}

func TestMapIterateOverMapEntries(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testTime := newMockTime(now)

	m := xmap.NewWithConfig[string, int](xmap.Config{
		TimeSource: testTime,
	})
	defer m.Stop()

	entries := []struct {
		key   string
		value int
		ttl   time.Duration
	}{
		{"a", 1, 10 * time.Second},
		{"b", 2, time.Minute},
		{"c", 3, time.Hour},
		{"d", 4, 0}, // Never expires.
	}

	for _, entry := range entries {
		m.Set(entry.key, entry.value, entry.ttl)
	}

	// Loop over the entries and verify the returned elements.
	checkEntriesRange := func(wantEntries map[string]int) {
		t.Helper()

		gotEntries := make(map[string]int)
		for entry := range m.Entries(context.Background()) {
			gotEntries[entry.Key] = entry.Value
		}

		if !maps.Equal(wantEntries, gotEntries) {
			t.Errorf("want %v, got %v", wantEntries, gotEntries)
		}
	}

	// No keys have expired yet.
	checkEntriesRange(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})

	// Advance the time to make "a" expire.
	testTime.Advance(15 * time.Second)
	checkEntriesRange(map[string]int{"b": 2, "c": 3, "d": 4})

	// Advance the time to make "b" expire.
	testTime.Advance(time.Minute) // 1 minute and 15 seconds have passed.
	checkEntriesRange(map[string]int{"c": 3, "d": 4})

	// Advance the time to make "c" expire.
	testTime.Advance(time.Hour) // 1 hour, 1 minute and 15 seconds have passed.
	checkEntriesRange(map[string]int{"d": 4})

	// Advance the time 1 year.
	testTime.Advance(24 * 365 * time.Hour)
	checkEntriesRange(map[string]int{"d": 4})
}

func TestMapPartialIterationOverEntries(t *testing.T) {
	t.Parallel()

	now := time.Now()
	testTime := newMockTime(now)

	m := xmap.NewWithConfig[string, int](xmap.Config{
		TimeSource: testTime,
	})
	defer m.Stop()

	entries := []struct {
		key   string
		value int
		ttl   time.Duration
	}{
		{"a", 1, 10 * time.Second},
		{"b", 2, time.Minute},
		{"c", 3, time.Hour},
		{"d", 4, 0}, // Never expires.
	}

	for _, entry := range entries {
		m.Set(entry.key, entry.value, entry.ttl)
	}

	// Number of entries consumed.
	var gotEntries atomic.Int32

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for range m.Entries(ctx) {
		gotEntries.Add(1)
		cancel() // Must cancel the context to release the lock.
		break    // Stop after consuming 1 entry.
	}

	if gotEntries.Load() != 1 {
		t.Errorf("want to consume 1 entry, got %d", gotEntries.Load())
	}
}
