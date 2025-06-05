package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_IntegrationWithCacheKeys(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(100, time.Hour, dir)
	
	// Pre-seed cache with coordinate-based keys (matching geotz workflow)
	preseedEntries := map[string]string{
		"51.507400,-0.127800": "Europe/London",   // London coordinates
		"35.676200,139.650300": "Asia/Tokyo",     // Tokyo coordinates  
	}
	cache.PreSeed(preseedEntries)
	
	// Test 1: Coordinate-based cache lookup (simulating geotz coordinate caching)
	start := time.Now()
	value, found := cache.Get("51.507400,-0.127800")
	duration1 := time.Since(start)
	
	if !found {
		t.Errorf("expected pre-seeded coordinate key to be found in cache")
	}
	if value != "Europe/London" {
		t.Errorf("expected Europe/London, got %s", value)
	}
	if duration1 > 5*time.Millisecond {
		t.Errorf("expected cache hit to be very fast, took %v", duration1)
	}
	
	// Test 2: City-based cache lookup (simulating city name caching)
	cache.Set("london", "Europe/London")
	
	start = time.Now()
	value, found = cache.Get("london")
	duration2 := time.Since(start)
	
	if !found {
		t.Errorf("expected city name key to be found in cache")
	}
	if value != "Europe/London" {
		t.Errorf("expected Europe/London, got %s", value)
	}
	if duration2 > 5*time.Millisecond {
		t.Errorf("expected cache hit to be very fast, took %v", duration2)
	}
	
	// Test 3: Cache miss simulation
	start = time.Now()
	_, found = cache.Get("non-existent-key")
	duration3 := time.Since(start)
	
	if found {
		t.Errorf("expected cache miss for non-existent key")
	}
	if duration3 > 5*time.Millisecond {
		t.Errorf("expected cache miss to be very fast, took %v", duration3)
	}
}

func TestCache_TTLBehaviorInIntegration(t *testing.T) {
	dir := t.TempDir()
	
	// Create cache with very short TTL for testing
	cache := NewLRUCache(100, 50*time.Millisecond, dir)
	
	// Simulate regular cache entry (should expire quickly)
	cache.Set("40.7128,-74.0060", "America/New_York")
	
	// Simulate pre-seeded entry (should last longer)
	cache.SetWithTTL("48.8566,2.3522", "Europe/Paris", time.Hour)
	
	// Both should be available immediately
	if value, ok := cache.Get("40.7128,-74.0060"); !ok || value != "America/New_York" {
		t.Errorf("expected regular cache entry to be available immediately")
	}
	if value, ok := cache.Get("48.8566,2.3522"); !ok || value != "Europe/Paris" {
		t.Errorf("expected pre-seeded entry to be available immediately")
	}
	
	// Wait for regular entry to expire
	time.Sleep(70 * time.Millisecond)
	
	// Regular entry should be expired
	if _, ok := cache.Get("40.7128,-74.0060"); ok {
		t.Errorf("expected regular cache entry to expire after short TTL")
	}
	
	// Pre-seeded entry should still be available
	if value, ok := cache.Get("48.8566,2.3522"); !ok || value != "Europe/Paris" {
		t.Errorf("expected pre-seeded entry to still be available with longer TTL")
	}
}

func TestCache_PerformanceWithPreseeding(t *testing.T) {
	dir := t.TempDir()
	cache := NewLRUCache(1000, time.Hour, dir)
	
	// Pre-seed with many entries (simulating real pre-seeding)
	manyEntries := make(map[string]string)
	for i := 0; i < 100; i++ {
		lat := 50.0 + float64(i)*0.1
		lng := 0.0 + float64(i)*0.1
		key := formatCoordinate(lat, lng)
		manyEntries[key] = "Europe/London" // Simplified for test
	}
	cache.PreSeed(manyEntries)
	
	// Test that cache performance remains good with many entries
	start := time.Now()
	for i := 0; i < 100; i++ {
		lat := 50.0 + float64(i)*0.1
		lng := 0.0 + float64(i)*0.1
		key := formatCoordinate(lat, lng)
		if _, ok := cache.Get(key); !ok {
			t.Errorf("expected pre-seeded entry to be available: %s", key)
		}
	}
	totalDuration := time.Since(start)
	
	avgDuration := totalDuration / 100
	if avgDuration > 1*time.Millisecond {
		t.Errorf("expected cache lookups to be very fast, average was %v", avgDuration)
	}
}

// Helper function to format coordinates consistently
func formatCoordinate(lat, lng float64) string {
	return fmt.Sprintf("%.6f,%.6f", lat, lng)
}