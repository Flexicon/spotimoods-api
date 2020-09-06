package internal

import "time"

// Cache storage interface
type Cache interface {
	// Set an item to cache
	Set(item *CacheItem) error
	// Get a cached value by key
	Get(key string, value interface{}) error
	// Delete removes the given key from cache
	Delete(key string) error
	// Exists checks for the existance of the given key
	Exists(key string) bool
	// Once gets the Value for the given Key from the cache or executes, caches, and returns the results of the given Do func
	Once(item *CacheItem) error
}

// CacheItem configuration
type CacheItem struct {
	Key   string
	Value interface{}

	// TTL is the cache expiration time.
	// Default TTL is 1 hour.
	TTL time.Duration

	// Do returns value to be cached.
	Do func(*CacheItem) (interface{}, error)

	// IfExists only sets the key if it already exist.
	IfExists bool

	// IfNotExists only sets the key if it does not already exist.
	IfNotExists bool

	// SkipLocalCache skips local cache as if it is not set.
	SkipLocalCache bool
}
