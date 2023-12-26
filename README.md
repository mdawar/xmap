# xmap
[![Go Reference](https://pkg.go.dev/badge/github.com/mdawar/xmap.svg)](https://pkg.go.dev/github.com/mdawar/xmap)

A generic and thread-safe Go map with automatic key expiration.

## Installation

```sh
go get -u github.com/mdawar/xmap
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/mdawar/xmap"
)

func main() {
	// Create a map with the default configuration.
	m := xmap.New[string, int]()
	defer m.Stop() // Stop the cleanup goroutine and clear the map.

	// Create new entries in the map.
	m.Set("a", 1, time.Minute) // Key that expires in 1 minute.
	m.Set("b", 2, 0)           // Key that never expires (0 TTL).

	// Replace a key.
	m.Set("a", 3, time.Hour) // Replace key (New value and expiration time).

	// Update the value without changing the expiration time.
	// Reports whether the key was updated (Key exists).
	ok := m.Update("b", 4)

	// Get the value if the key exists and has not expired.
	// The second return value reports whether the key exists.
	value, ok := m.Get("a")

	// Get the value with the expiration time.
	// The third return value reports whether the key exists.
	value, expiration, ok := m.GetWithExpiration("a")
	// If the key never expires, it will have a zero expiration time value.
	neverExpires := expiration.IsZero()

	// Length of the map.
	total := m.Len()

	// Delete a key from the map.
	m.Delete("a")

	// Delete all the keys from the map.
	m.Clear()

	// Expired keys are automatically removed at regular intervals.
	// Additionally, the removal of expired keys can be manually triggered.
	removed := m.RemoveExpired() // Returns the number of removed keys.
}
```

## Tests and Benchmarks

Run the tests:

```sh
make test
# Or
go test ./... -cover -race
```

Run the benchmarks:

```sh
make benchmark
# Or
go test ./... -bench .
```
