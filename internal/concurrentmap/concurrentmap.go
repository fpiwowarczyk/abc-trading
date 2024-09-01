package concurrentmap

import "sync"

// ConcurrentMap is a simple map that should be thread safe.
type ConcurrentMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func New[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data: make(map[K]V),
		mu:   sync.RWMutex{},
	}
}

// Set given key with given value.
func (c *ConcurrentMap[K, V]) Set(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = val
}

// Get value for given key, if key does not exist, return false as second arg.
func (c *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[key]
	return val, ok
}
