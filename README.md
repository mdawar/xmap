# xmap

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
	defer m.Stop() // Stop the cleanup goroutine.

	// Create a new entry in the map.
	m.Set("a", 1, time.Minute) // Key that expires in 1 minute.
	m.Set("b", 2, 0)           // Key that never expires (0 TTL).

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
}
```
