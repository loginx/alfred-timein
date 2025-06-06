//go:build bdd
// +build bdd

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/loginx/alfred-timein/internal/adapters/cache"
	"github.com/loginx/alfred-timein/internal/adapters/geocoder"
	"github.com/loginx/alfred-timein/internal/adapters/presenter"
	"github.com/loginx/alfred-timein/internal/adapters/timezonefinder"
	"github.com/loginx/alfred-timein/internal/usecases"
)

func TestGeotzIntegrationWithCachePreSeeding(t *testing.T) {
	testDir := fmt.Sprintf("/tmp/alfred-timein-integration-test-%d", time.Now().UnixNano())
	defer os.RemoveAll(testDir)

	// Create and set up cache with pre-seeded data
	c := cache.NewLRUCache(200, 24*time.Hour, testDir)
	
	// Pre-seed with test capital data
	entries := map[string]string{
		"london":                 "Europe/London",
		"paris":                  "Europe/Paris",
		"tokyo":                  "Asia/Tokyo",
		"51.507400,-0.127800":    "Europe/London",  // London coordinates
		"48.856600,2.352200":     "Europe/Paris",   // Paris coordinates
		"35.676200,139.650300":   "Asia/Tokyo",     // Tokyo coordinates
	}
	c.PreSeed(entries)

	// Set up use case with pre-seeded cache
	geocoderAdapter := geocoder.NewOpenStreetMapGeocoder()
	tzf, err := timezonefinder.NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to initialize timezone finder: %v", err)
	}
	formatter := presenter.NewPlainFormatter()
	
	useCase := usecases.NewGeotzUseCase(geocoderAdapter, tzf, c, formatter)

	t.Run("Cache hit for pre-seeded city", func(t *testing.T) {
		start := time.Now()
		
		result, err := useCase.GetTimezoneFromCity("London")
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		output := string(result)
		t.Logf("Lookup for London took %v, result: %s", duration, output)
		
		if !strings.Contains(output, "Europe/London") {
			t.Errorf("Expected result to contain Europe/London, got: %s", output)
		}
		
		// Cache hits should be fast (under 100ms)
		if duration > 100*time.Millisecond {
			t.Errorf("Expected fast cache hit, but took %v", duration)
		}
	})

	t.Run("Cache hit for pre-seeded city case insensitive", func(t *testing.T) {
		start := time.Now()
		
		result, err := useCase.GetTimezoneFromCity("TOKYO")
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		output := string(result)
		t.Logf("Lookup for TOKYO took %v, result: %s", duration, output)
		
		if !strings.Contains(output, "Asia/Tokyo") {
			t.Errorf("Expected result to contain Asia/Tokyo, got: %s", output)
		}
		
		// Cache hits should be fast (under 100ms)
		if duration > 100*time.Millisecond {
			t.Errorf("Expected fast cache hit, but took %v", duration)
		}
	})

	t.Run("Cache miss requires geocoding", func(t *testing.T) {
		start := time.Now()
		
		result, err := useCase.GetTimezoneFromCity("New York")
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		
		output := string(result)
		t.Logf("Lookup for New York took %v, result: %s", duration, output)
		
		// Should contain New York timezone
		if !strings.Contains(output, "America/New_York") {
			t.Errorf("Expected result to contain America/New_York, got: %s", output)
		}
		
		// Cache misses will take longer (network geocoding required)
		if duration < 100*time.Millisecond {
			t.Logf("Unexpectedly fast for cache miss: %v", duration)
		}
	})

	t.Run("Subsequent lookup of new city should be cached", func(t *testing.T) {
		// First lookup should be slow (cache miss)
		start1 := time.Now()
		result1, err1 := useCase.GetTimezoneFromCity("Berlin")
		duration1 := time.Since(start1)
		
		if err1 != nil {
			t.Fatalf("Expected no error on first lookup, got: %v", err1)
		}
		
		// Second lookup should be fast (cache hit)
		start2 := time.Now()
		result2, err2 := useCase.GetTimezoneFromCity("Berlin")
		duration2 := time.Since(start2)
		
		if err2 != nil {
			t.Fatalf("Expected no error on second lookup, got: %v", err2)
		}
		
		// Results should be the same
		if string(result1) != string(result2) {
			t.Errorf("Expected same results, got %s vs %s", string(result1), string(result2))
		}
		
		t.Logf("First Berlin lookup: %v, Second: %v", duration1, duration2)
		
		// Second lookup should be much faster
		if duration2 > duration1/2 {
			t.Errorf("Expected second lookup to be much faster. First: %v, Second: %v", duration1, duration2)
		}
	})
}