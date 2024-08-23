# xmap

[![Go Reference](https://pkg.go.dev/badge/github.com/mdawar/xmap.svg)](https://pkg.go.dev/github.com/mdawar/xmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdawar/xmap)](https://goreportcard.com/report/github.com/mdawar/xmap)
[![Go Tests](https://github.com/mdawar/xmap/actions/workflows/go.yml/badge.svg?branch=main&event=push)](https://github.com/mdawar/xmap/actions)

A generic and thread-safe Go map with automatic key expiration.

## Installation

```sh
go get -u github.com/mdawar/xmap
```

## Usage

#### New Map

```go
// Create a map with the default configuration.
m := xmap.New[string, int]()
// Stop the cleanup goroutine and clear the map.
defer m.Stop()
```

#### Create

```go
// Create new entries in the map.
m.Set("a", 1, time.Minute) // Key that expires in 1 minute.
m.Set("b", 2, 0)           // Key that never expires (0 TTL).

// Replace a key.
m.Set("a", 3, time.Hour) // Replace key (New value and expiration time).
```

#### Update

```go
// Update the value without changing the expiration time.
// Reports whether the key was updated (Key exists).
ok := m.Update("b", 4)
```

#### Get

```go
// Get the value if the key exists and has not expired.
// The second return value reports whether the key exists.
value, ok := m.Get("a")

// Get the value with the expiration time.
// The third return value reports whether the key exists.
value, expiration, ok := m.GetWithExpiration("a")
// If the key never expires, it will have a zero expiration time value.
neverExpires := expiration.IsZero()
```

#### Delete

```go
// Delete a key from the map.
m.Delete("a")

// Delete all the keys from the map.
m.Clear()
```

#### Length

```go
total := m.Len()
```

#### Iteration

```go
for key, value := range m.All() {
	fmt.Println("Key:", key, "-", "Value:", value)
}
```

#### Remove Expired Keys

```go
// Expired keys are automatically removed at regular intervals.
// Additionally, the removal of expired keys can be manually triggered.
removed := m.RemoveExpired() // Returns the number of removed keys.
```

## Configuration

| Name              | Type            | Description                                                      |
| ----------------- | --------------- | ---------------------------------------------------------------- |
| `CleanupInterval` | `time.Duration` | Interval at which expired keys are removed (Default: 5 minutes). |
| `InitialCapacity` | `int`           | Initial map capacity hint (Passed to `make()`).                  |
| `TimeSource`      | `xmap.Time`     | Custom time source (Useful for testing).                         |

Example:

```go
package xmap

import (
	"time"

	"github.com/mdawar/xmap"
)

func main() {
	m := xmap.NewWithConfig[string, int](xmap.Config{
		CleanupInterval: 10 * time.Minute,
		InitialCapacity: 10_000_000,
		TimeSource:      mockTime,
	})
	defer m.Stop()
}
```

## Tests

```sh
make test
```

Or:

```sh
go test -cover -race
```

## Benchmarks

```sh
make benchmark
```

Or:

```sh
go test -bench .
```
