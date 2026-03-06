package mcp

import (
	"os"
	"sync"
	"time"
)

// CacheEntry stores a cached value along with the file mod time used to validate it.
type CacheEntry struct {
	Value   any
	ModTime time.Time
}

// Cache provides in-memory caching with file mod-time invalidation.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]CacheEntry),
	}
}

// Get retrieves a cached value if the file hasn't been modified since caching.
// Returns (value, true) on hit, (nil, false) on miss or stale.
func (c *Cache) Get(key string, filePaths ...string) (any, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	// Check if any file has been modified since the cache was set
	for _, fp := range filePaths {
		info, err := os.Stat(fp)
		if err != nil {
			return nil, false
		}
		if info.ModTime().After(entry.ModTime) {
			return nil, false
		}
	}

	return entry.Value, true
}

// Set stores a value in the cache with the current timestamp.
func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		Value:   value,
		ModTime: time.Now(),
	}
}

// Invalidate removes a specific key from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// InvalidateAll clears the entire cache.
func (c *Cache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]CacheEntry)
}
