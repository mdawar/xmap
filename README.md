# xmap

A thread-safe, generic Go map with automatic key expiration.

## Installation

```sh
go get -u github.com/mdawar/xmap
```

## Usage

```go
m := xmap.New[string, int]()
```
