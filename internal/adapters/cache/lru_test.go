package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLRUCache_SetGet(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(3, time.Hour, dir)
	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("c", "3")
	if v, ok := cache.Get("a"); !ok || v != "1" {
		t.Errorf("expected to get '1', got '%v'", v)
	}
	if v, ok := cache.Get("b"); !ok || v != "2" {
		t.Errorf("expected to get '2', got '%v'", v)
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(2, time.Hour, dir)
	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("c", "3") // should evict "a"
	if _, ok := cache.Get("a"); ok {
		t.Errorf("expected 'a' to be evicted")
	}
	if v, ok := cache.Get("b"); !ok || v != "2" {
		t.Errorf("expected to get '2', got '%v'", v)
	}
	if v, ok := cache.Get("c"); !ok || v != "3" {
		t.Errorf("expected to get '3', got '%v'", v)
	}
}

func TestLRUCache_TTL(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(2, 10*time.Millisecond, dir)
	cache.Set("a", "1")
	time.Sleep(20 * time.Millisecond)
	if _, ok := cache.Get("a"); ok {
		t.Errorf("expected 'a' to expire by TTL")
	}
}

func TestLRUCache_PersistAndLoad(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(2, time.Hour, dir)
	cache.Set("a", "1")
	cache.Set("b", "2")
	cache2 := NewLRUCache(2, time.Hour, dir)
	if v, ok := cache2.Get("a"); !ok || v != "1" {
		t.Errorf("expected to get '1' after reload, got '%v'", v)
	}
	if v, ok := cache2.Get("b"); !ok || v != "2" {
		t.Errorf("expected to get '2' after reload, got '%v'", v)
	}
}

func TestLRUCache_Clear(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(2, time.Hour, dir)
	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Clear()
	if _, ok := cache.Get("a"); ok {
		t.Errorf("expected 'a' to be cleared")
	}
	if _, ok := cache.Get("b"); ok {
		t.Errorf("expected 'b' to be cleared")
	}
	if _, err := os.Stat(filepath.Join(dir, "geotz_cache.json")); !os.IsNotExist(err) {
		t.Errorf("expected cache file to be deleted")
	}
}

func TestLRUCache_ConcurrentAccess(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(100, time.Hour, dir)
	
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100
	
	// Test concurrent reads and writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := fmt.Sprintf("value-%d-%d", id, j)
				cache.Set(key, value)
				if v, ok := cache.Get(key); ok && v != value {
					t.Errorf("concurrent access issue: expected %s, got %s", value, v)
				}
			}
		}(i)
	}
	
	wg.Wait()
}

func TestLRUCache_PreSeed(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, time.Hour, dir)
	
	// Pre-seed with some entries
	entries := map[string]string{
		"london":  "Europe/London",
		"tokyo":   "Asia/Tokyo",
		"newyork": "America/New_York",
	}
	cache.PreSeed(entries)
	
	// Check that pre-seeded entries are available
	for key, expectedValue := range entries {
		if value, ok := cache.Get(key); !ok || value != expectedValue {
			t.Errorf("expected pre-seeded key %s to have value %s, got %s", key, expectedValue, value)
		}
	}
	
	// Add a user entry with same key - should not be overwritten by pre-seed
	cache.Set("london", "user-value")
	cache.PreSeed(map[string]string{"london": "preseed-value"})
	
	if value, ok := cache.Get("london"); !ok || value != "user-value" {
		t.Errorf("expected user value to take precedence over pre-seed, got %s", value)
	}
}

func TestLRUCache_SetWithTTL(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, time.Hour, dir)
	
	// Set entry with custom short TTL
	cache.SetWithTTL("short", "value", 10*time.Millisecond)
	
	// Should be available immediately
	if value, ok := cache.Get("short"); !ok || value != "value" {
		t.Errorf("expected entry with custom TTL to be available immediately")
	}
	
	// Should expire after custom TTL
	time.Sleep(20 * time.Millisecond)
	if _, ok := cache.Get("short"); ok {
		t.Errorf("expected entry with custom TTL to expire")
	}
}

func TestLRUCache_PreSeedPersistence(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, time.Hour, dir)
	
	// Pre-seed with entries
	entries := map[string]string{
		"paris": "Europe/Paris",
		"rome":  "Europe/Rome",
	}
	cache.PreSeed(entries)
	
	// Create new cache instance (simulates restart)
	cache2 := NewLRUCache(10, time.Hour, dir)
	
	// Check that pre-seeded entries are still available
	for key, expectedValue := range entries {
		if value, ok := cache2.Get(key); !ok || value != expectedValue {
			t.Errorf("expected pre-seeded key %s to persist across restart, got %s", key, value)
		}
	}
}

func TestLRUCache_DefaultTTLBehavior(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, 50*time.Millisecond, dir)
	
	// Set entry with default TTL (TTL=0 should use cache default)
	cache.Set("default-ttl", "value")
	
	// Should be available immediately
	if value, ok := cache.Get("default-ttl"); !ok || value != "value" {
		t.Errorf("expected entry to be available immediately")
	}
	
	// Should still be available before TTL expires
	time.Sleep(30 * time.Millisecond)
	if value, ok := cache.Get("default-ttl"); !ok || value != "value" {
		t.Errorf("expected entry to be available before TTL expires")
	}
	
	// Should expire after default TTL
	time.Sleep(30 * time.Millisecond) // Total: 60ms > 50ms TTL
	if _, ok := cache.Get("default-ttl"); ok {
		t.Errorf("expected entry to expire after default TTL")
	}
}

func TestLRUCache_ZeroTTLUsesDefault(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, time.Hour, dir)
	
	// Manually create entry with TTL=0 (should use cache default)
	cache.entries["zero-ttl"] = cacheEntry{
		Value:     "test-value",
		CreatedAt: time.Now(),
		TTL:       0, // This should fall back to cache default
	}
	cache.order = append(cache.order, "zero-ttl")
	
	// Should be available since TTL=0 uses cache default (1 hour)
	if value, ok := cache.Get("zero-ttl"); !ok || value != "test-value" {
		t.Errorf("expected zero TTL entry to use cache default TTL")
	}
}

func TestLRUCache_PreSeedVsRegularEntryBehavior(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, time.Hour, dir)
	
	// Pre-seed with some entries 
	preseedEntries := map[string]string{
		"london": "Europe/London",
		"tokyo":  "Asia/Tokyo",
	}
	cache.PreSeed(preseedEntries)
	
	// Add regular cache entry
	cache.Set("berlin", "Europe/Berlin")
	
	// Both types should be available
	if value, ok := cache.Get("london"); !ok || value != "Europe/London" {
		t.Errorf("expected pre-seeded entry to be available")
	}
	if value, ok := cache.Get("berlin"); !ok || value != "Europe/Berlin" {
		t.Errorf("expected regular entry to be available")
	}
	
	// Check that pre-seeded entries have longer TTL by examining internal state
	londonEntry := cache.entries["london"]
	berlinEntry := cache.entries["berlin"]
	
	if londonEntry.TTL == 0 {
		t.Errorf("expected pre-seeded entry to have custom TTL set")
	}
	if berlinEntry.TTL != 0 {
		t.Errorf("expected regular entry to have TTL=0 (uses default)")
	}
	if londonEntry.TTL <= time.Hour {
		t.Errorf("expected pre-seeded entry TTL to be longer than regular cache default")
	}
}

func TestLRUCache_CacheHitMissScenarios(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(10, time.Hour, dir)
	
	// Test cache miss for non-existent entry
	if _, ok := cache.Get("non-existent"); ok {
		t.Errorf("expected cache miss for non-existent entry")
	}
	
	// Add entry and test cache hit
	cache.Set("existing", "value")
	if value, ok := cache.Get("existing"); !ok || value != "value" {
		t.Errorf("expected cache hit for existing entry")
	}
	
	// Test cache miss after expiry
	shortCache := NewLRUCache(10, 10*time.Millisecond, dir)
	shortCache.Set("expires", "value")
	time.Sleep(20 * time.Millisecond)
	if _, ok := shortCache.Get("expires"); ok {
		t.Errorf("expected cache miss after entry expiry")
	}
}

func TestLRUCache_ComplexCacheWorkflow(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(6, time.Hour, dir) // Larger cache to avoid immediate eviction
	
	// Pre-seed with capitals
	capitals := map[string]string{
		"london": "Europe/London",
		"paris":  "Europe/Paris",
		"tokyo":  "Asia/Tokyo",
	}
	cache.PreSeed(capitals)
	
	// Verify pre-seeded entries are available
	for key, expectedValue := range capitals {
		if value, ok := cache.Get(key); !ok || value != expectedValue {
			t.Errorf("expected pre-seeded entry %s to be available", key)
		}
	}
	
	// Add regular entries
	cache.Set("berlin", "Europe/Berlin")
	cache.Set("madrid", "Europe/Madrid")
	cache.Set("rome", "Europe/Rome")
	
	// Cache should be at capacity (6 entries)
	if len(cache.entries) != 6 {
		t.Errorf("expected cache to have 6 entries, got %d", len(cache.entries))
	}
	
	// Add one more entry - should trigger eviction of least recently used
	cache.Set("amsterdam", "Europe/Amsterdam")
	
	// Should still have 6 entries (oldest evicted)
	if len(cache.entries) != 6 {
		t.Errorf("expected cache to maintain max size of 6 entries, got %d", len(cache.entries))
	}
	
	// Most recently added entry should be available
	if value, ok := cache.Get("amsterdam"); !ok || value != "Europe/Amsterdam" {
		t.Errorf("expected newest entry to be available")
	}
	
	// Recently accessed entries should still be available
	if value, ok := cache.Get("rome"); !ok || value != "Europe/Rome" {
		t.Errorf("expected recently added entry to be available")
	}
}

func TestLRUCache_EvictionRespectsTTL(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(3, time.Hour, dir) // Very small cache for clear eviction testing
	
	// Add entries that will be evicted by LRU
	cache.Set("first", "value1")
	cache.Set("second", "value2") 
	cache.Set("third", "value3") // Cache is now full
	
	// Add entry with custom TTL
	cache.SetWithTTL("long-lived", "special", 24*time.Hour)
	
	// "first" should be evicted (oldest)
	if _, ok := cache.Get("first"); ok {
		t.Errorf("expected oldest entry to be evicted")
	}
	
	// Long-lived entry should be available
	if value, ok := cache.Get("long-lived"); !ok || value != "special" {
		t.Errorf("expected long-lived entry to be available")
	}
}