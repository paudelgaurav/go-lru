package lrucache

import (
	"sync"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	cache := NewCache(3)

	// Test adding and retrieving elements
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected key1 to have value 'value1', got %v", val)
	}

	val, ok = cache.Get("key2")
	if !ok || val != "value2" {
		t.Errorf("expected key2 to have value 'value2', got %v", val)
	}

	// Test if LRU order is updated
	cache.Get("key1") // key1 should now be the most recent

	cache.Add("key4", "value4") // should evict key3
	_, ok = cache.Get("key3")
	if ok {
		t.Error("expected key3 to be evicted, but it still exists")
	}

	val, ok = cache.Get("key4")
	if !ok || val != "value4" {
		t.Errorf("expected key4 to have value 'value4', got %v", val)
	}
}

func TestEvictionPolicy(t *testing.T) {
	cache := NewCache(2)

	// Add two elements
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")

	// Add another element, evict the least recently used (key1)
	cache.Add("key3", "value3")

	_, ok := cache.Get("key1")
	if ok {
		t.Error("expected key1 to be evicted, but it still exists")
	}

	val, ok := cache.Get("key2")
	if !ok || val != "value2" {
		t.Errorf("expected key2 to have value 'value2', got %v", val)
	}

	val, ok = cache.Get("key3")
	if !ok || val != "value3" {
		t.Errorf("expected key3 to have value 'value3', got %v", val)
	}
}

func TestRemove(t *testing.T) {
	cache := NewCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")

	// Test removal of key
	removed := cache.Remove("key1")
	if !removed {
		t.Error("expected key1 to be removed")
	}

	_, ok := cache.Get("key1")
	if ok {
		t.Error("expected key1 to be removed, but it still exists")
	}
}

func TestLen(t *testing.T) {
	cache := NewCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")

	if cache.Len() != 2 {
		t.Errorf("expected length to be 2, got %d", cache.Len())
	}

	// Test after eviction
	cache.Add("key3", "value3")
	if cache.Len() != 2 {
		t.Errorf("expected length to be 2 after eviction, got %d", cache.Len())
	}
}

func TestClear(t *testing.T) {
	cache := NewCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")

	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("expected length to be 0 after clearing, got %d", cache.Len())
	}

	_, ok := cache.Get("key1")
	if ok {
		t.Error("expected key1 to be cleared, but it still exists")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(100)

	wg := sync.WaitGroup{}

	// Spawn 50 goroutines to add keys concurrently
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Add(i, i*10)
		}(i)
	}

	// Spawn 50 goroutines to get keys concurrently
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Get(i)
		}(i)
	}

	wg.Wait()

	// Expect that no more than 50 items have been added
	if cache.Len() > 50 {
		t.Errorf("expected cache length to be less than or equal to 50, got %d", cache.Len())
	}
}

func TestDuplicateEntry(t *testing.T) {
	cache := NewCache(3)

	cache.Add("key1", "value1")

	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected value 'value1', got %v", val)
	}

	cache.Add("key1", "value2")

	val, ok = cache.Get("key1")
	if !ok || val != "value2" {
		t.Errorf("expected updated value 'value2', got %v", val)
	}

	if cache.Len() > 1 {
		t.Errorf("expected cache length to be less than or equal to 3, got %d", cache.Len())
	}
}
