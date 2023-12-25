package xmap_test

import (
	"testing"
	"time"

	"github.com/mdawar/xmap"
)

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
		t.Errorf("want map length 0, got %d", m.Len())
	}
}
