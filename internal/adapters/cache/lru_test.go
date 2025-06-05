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