package alfred

import (
	"os"
	"path/filepath"
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
