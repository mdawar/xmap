package xmap_test

import (
	"testing"

	"github.com/mdawar/xmap"
)

// 1. Test create and get a value
// 2. Test key expiration time
// 4. Test no expiration with 0 TTL
// 5. Test create replaces an existing value
// 6. Test update value without changing the expiration time
// 7. Test delete removes the key
// 8. Test keys are automatically removed on expiration

func TestMapSetThenGet(t *testing.T) {
	t.Parallel()

	m := xmap.New[string, int]()
	defer m.Stop()

	keyName := "abc"
	wantValue := 10

	m.Set(keyName, wantValue, 0)

	gotValue, ok := m.Get(keyName)
	if !ok {
		t.Fatalf("key %q does not exist in the map", keyName)
	}

	if wantValue != gotValue {
		t.Errorf("want value %d, got %d", wantValue, gotValue)
	}
}
