package alfred

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	cacheDir     = ".cache"
	cacheFile    = "geotz_cache.json"
	maxCacheSize = 100
	defaultTTL   = 7 * 24 * time.Hour // 7 days
)

type cacheEntry struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type lruCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	order   []string
	max     int
	ttl     time.Duration
	path    string
}

func NewLRUCache(max int, ttl time.Duration, dir string) *lruCache {
	c := &lruCache{
		entries: make(map[string]cacheEntry),
		order:   make([]string, 0, max),
		max:     max,
		ttl:     ttl,
		path:    filepath.Join(dir, cacheFile),
	}
	c.load()
	return c
}

func (c *lruCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok || time.Since(entry.CreatedAt) > c.ttl {
		if ok {
			c.delete(key)
		}
		return "", false
	}
	// Move to front
	c.moveToFront(key)
	return entry.Value, true
}

func (c *lruCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.entries[key]; !ok && len(c.entries) >= c.max {
		// Evict oldest
		oldest := c.order[len(c.order)-1]
		c.delete(oldest)
	}
	c.entries[key] = cacheEntry{Value: value, CreatedAt: time.Now()}
	c.moveToFront(key)
	c.persist()
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]cacheEntry)
	c.order = make([]string, 0, c.max)
	os.Remove(c.path)
}

func (c *lruCache) moveToFront(key string) {
	// Remove if exists
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
	c.order = append([]string{key}, c.order...)
}

func (c *lruCache) delete(key string) {
	delete(c.entries, key)
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}

func (c *lruCache) persist() {
	os.MkdirAll(filepath.Dir(c.path), 0755)
	data := struct {
		Max   int              `json:"max"`
		Cache [][2]interface{} `json:"cache"`
	}{
		Max:   c.max,
		Cache: make([][2]interface{}, 0, len(c.entries)),
	}
	for _, k := range c.order {
		entry := c.entries[k]
		data.Cache = append(data.Cache, [2]interface{}{k, entry})
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(c.path, b, 0644)
}

func (c *lruCache) load() {
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
		if time.Since(entry.CreatedAt) > c.ttl {
			continue
		}
		c.entries[k] = entry
		c.order = append(c.order, k)
	}
}

func DefaultGeotzCache() *lruCache {
	return NewLRUCache(maxCacheSize, defaultTTL, cacheDir)
}
