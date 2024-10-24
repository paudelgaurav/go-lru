# LRU Cache Implementation in Go

This is a simple implementation of an LRU (Least Recently Used) cache in Go. It is designed to store a limited number of items and automatically evicts the least recently used items when the cache reaches its maximum capacity.

## Features

- **Concurrent Safe**: The cache is thread-safe, allowing concurrent access by multiple goroutines.
- **Custom Capacity**: You can define the maximum capacity of the cache when creating a new instance.
- **Simple API**: Easy to use operations to interact with the cache.

## Usage


### Example

Here is an example of how to use the LRU cache:

```go
package main

import (
    "fmt"
    lrucache "github.com/paudelgaurav/go-lru" 
)

func main() {
    cache := lrucache.NewCache(3)

    cache.Add("key1", "value1")
    cache.Add("key2", "value2")
    cache.Add("key3", "value3")

    val, ok := cache.Get("key1")
    if ok {
        fmt.Println("Got:", val) // Outputs: Got: value1
    }

    // Adding a new item will evict the least recently used item
    cache.Add("key4", "value4")

    _, ok = cache.Get("key2")
    if !ok {
        fmt.Println("key2 was evicted") // Outputs: key2 was evicted
    }

    // Current cache length
    fmt.Println("Cache length:", cache.Len()) // Outputs: Cache length: 3

    // Clear the cache
    cache.Clear()
    fmt.Println("Cache length after clear:", cache.Len()) // Outputs: Cache length after clear: 0
}
```

### Testing

To run the tests for this implementation, you can use the following command:

```bash
go test -race 
```

This will execute the tests and check for race conditions.


## Contributing

Feel free to submit issues or pull requests for any improvements or additional features.
