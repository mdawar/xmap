package xmap_test

import (
	"fmt"
	"time"

	"github.com/mdawar/xmap"
)

func ExampleMap() {
	// Create a map with the default configuration.
	m := xmap.New[string, int]()
	defer m.Stop() // Stop the cleanup goroutine and clear the map.

	// Create new entries in the map.
	m.Set("a", 1, time.Minute) // Key that expires in 1 minute.
	m.Set("b", 2, 0)           // Key that never expires (0 TTL).

	// Replace a key.
	m.Set("a", 3, time.Hour) // Replace key (New value and expiration time).

	// Update the value without changing the expiration time.
	if ok := m.Update("b", 4); ok {
		fmt.Println("Key updated successfully")
	} else {
		fmt.Println("Key does not exist")
	}

	// Get the value if the key exists and has not expired.
	if value, ok := m.Get("a"); ok {
		fmt.Println("Value:", value)
	} else {
		fmt.Println("Key does not exist")
	}

	// Get the value with the expiration time.
	if value, expiration, ok := m.GetWithExpiration("a"); ok {
		// If the key never expires, it will have a zero expiration time value.
		fmt.Println("Key expires:", !expiration.IsZero())
		fmt.Println("Value:", value, "-", "Expiration:", expiration)
	} else {
		fmt.Println("Key does not exist")
	}

	total := m.Len() // Length of the map.

	fmt.Println("Total entries in the map:", total)

	m.Delete("a") // Delete a key from the map.
	m.Clear()     // Delete all the keys from the map.

	// Expired keys are automatically removed at regular intervals.
	// Additionally, the removal of expired keys can be manually triggered.
	removed := m.RemoveExpired()

	fmt.Println("Total keys removed:", removed)
}

func ExampleNewWithConfig() {
	m := xmap.NewWithConfig[string, int](xmap.Config{
		CleanupInterval: 10 * time.Minute, // Change the default cleanup interval.
		InitialCapacity: 1_000_000,        // Initial capacity hint (Passed to make).
	})
	defer m.Stop()
}

func ExampleMap_All() {
	m := xmap.New[string, int]()
	defer m.Stop()

	for k, v := range m.All() {
		fmt.Println("Key:", k, "-", "Value:", v)
	}
}
