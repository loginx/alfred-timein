package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultCacheDir  = ".cache"
	defaultCacheFile = "geotz_cache.json"
	defaultMaxSize   = 100
	defaultTTL       = 7 * 24 * time.Hour  // 7 days
	preseedTTL       = 90 * 24 * time.Hour // 90 days for pre-seeded entries
)

type cacheEntry struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	TTL       time.Duration `json:"ttl,omitempty"`
}

// LRUCache implements the Cache interface with LRU eviction and persistence
type LRUCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	order   []string
	max     int
	ttl     time.Duration
	path    string
}

// NewLRUCache creates a new LRUCache
func NewLRUCache(max int, ttl time.Duration, dir string) *LRUCache {
	c := &LRUCache{
		entries: make(map[string]cacheEntry),
		order:   make([]string, 0, max),
		max:     max,
		ttl:     ttl,
		path:    filepath.Join(dir, defaultCacheFile),
	}
	c.load()
	return c
}

// NewDefaultCache creates a cache with default settings
func NewDefaultCache() *LRUCache {
	return NewLRUCache(defaultMaxSize, defaultTTL, defaultCacheDir)
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return "", false
	}

	// Use entry-specific TTL if set, otherwise use cache default
	ttl := c.ttl
	if entry.TTL > 0 {
		ttl = entry.TTL
	}

	if time.Since(entry.CreatedAt) > ttl {
		c.deleteUnsafe(key)
		return "", false
	}

	// Move to front
	c.moveToFrontUnsafe(key)
	return entry.Value, true
}

// Set stores a value in the cache
func (c *LRUCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.entries[key]; !ok && len(c.entries) >= c.max {
		// Evict oldest
		if len(c.order) > 0 {
			oldest := c.order[len(c.order)-1]
			c.deleteUnsafe(oldest)
		}
	}

	c.entries[key] = cacheEntry{Value: value, CreatedAt: time.Now(), TTL: 0}
	c.moveToFrontUnsafe(key)
	c.persistUnsafe()
}

// SetWithTTL stores a value in the cache with a custom TTL
func (c *LRUCache) SetWithTTL(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.entries[key]; !ok && len(c.entries) >= c.max {
		// Evict oldest
		if len(c.order) > 0 {
			oldest := c.order[len(c.order)-1]
			c.deleteUnsafe(oldest)
		}
	}

	c.entries[key] = cacheEntry{Value: value, CreatedAt: time.Now(), TTL: ttl}
	c.moveToFrontUnsafe(key)
	c.persistUnsafe()
}

// PreSeed adds entries to the cache with long TTL, used for build-time seeding
func (c *LRUCache) PreSeed(entries map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, value := range entries {
		// Only add if not already present to avoid overwriting user cache
		if _, exists := c.entries[key]; !exists {
			c.entries[key] = cacheEntry{
				Value:     value,
				CreatedAt: time.Now(),
				TTL:       preseedTTL,
			}
			c.order = append(c.order, key)
		}
	}
	c.persistUnsafe()
}

// Clear removes all entries from the cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]cacheEntry)
	c.order = make([]string, 0, c.max)
	os.Remove(c.path)
}

// moveToFrontUnsafe moves key to front of order slice (caller must hold lock)
func (c *LRUCache) moveToFrontUnsafe(key string) {
	// Remove if exists
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
	c.order = append([]string{key}, c.order...)
}

// deleteUnsafe removes key from cache (caller must hold lock)
func (c *LRUCache) deleteUnsafe(key string) {
	delete(c.entries, key)
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}

// persistUnsafe saves cache to disk (caller must hold lock)
func (c *LRUCache) persistUnsafe() {
	if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil {
		return // Fail silently for now
	}

	data := struct {
		Max   int              `json:"max"`
		Cache [][2]interface{} `json:"cache"`
	}{
		Max:   c.max,
		Cache: make([][2]interface{}, 0, len(c.entries)),
	}

	for _, k := range c.order {
		if entry, ok := c.entries[k]; ok {
			data.Cache = append(data.Cache, [2]interface{}{k, entry})
		}
	}

	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		os.WriteFile(c.path, jsonData, 0644)
	}
}

// load restores cache from disk
func (c *LRUCache) load() {
	f, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer f.Close()

	var data struct {
		Max   int                  `json:"max"`
		Cache [][2]json.RawMessage `json:"cache"`
	}

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return
	}

	for _, pair := range data.Cache {
		var k string
		var entry cacheEntry

		if err := json.Unmarshal(pair[0], &k); err != nil {
			continue
		}
		if err := json.Unmarshal(pair[1], &entry); err != nil {
			continue
		}
		// Use entry-specific TTL if set, otherwise use cache default
		ttl := c.ttl
		if entry.TTL > 0 {
			ttl = entry.TTL
		}
		if time.Since(entry.CreatedAt) > ttl {
			continue
		}

		c.entries[k] = entry
		c.order = append(c.order, k)
	}
}