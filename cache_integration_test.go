//go:build bdd
// +build bdd

package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Performance thresholds
const (
	cacheHitMaxDuration  = 100 * time.Millisecond // Pre-seeded cache hits must be under 100ms
	cacheMissMinDuration = 200 * time.Millisecond // Cache misses should take at least 200ms (proves geocoding happened)
)

// TestCacheIntegration tests cache behavior with real CLI and cache files
func TestCacheIntegration(t *testing.T) {
	// Ensure we have binaries built
	if _, err := os.Stat("bin/geotz"); os.IsNotExist(err) {
		t.Skip("Skipping cache integration test - bin/geotz not found. Run 'make build' first.")
	}

	// Regenerate pre-seeded cache to ensure clean state
	t.Log("Regenerating pre-seeded cache for cache integration test...")
	cmd := exec.Command("make", "preseed")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to regenerate cache: %v, output: %s", err, output)
	}

	t.Run("Pre-seeded cities are generally fast", func(t *testing.T) {
		preSeededCities := []string{
			"paris", "tokyo", "berlin", "madrid", 
			"rome", "amsterdam", "stockholm",
		}

		var fastLookups int
		
		for _, city := range preSeededCities {
			duration := measureCityLookup(t, city)
			
			if duration <= cacheHitMaxDuration {
				fastLookups++
				t.Logf("✅ %s: %v (fast)", city, duration)
			} else {
				t.Logf("⚠ %s: %v (slower than expected)", city, duration)
			}
		}

		// At least 80% of pre-seeded cities should be fast
		expectedFast := int(float64(len(preSeededCities)) * 0.8)
		if fastLookups < expectedFast {
			t.Errorf("Only %d/%d pre-seeded cities were fast, expected at least %d", 
				fastLookups, len(preSeededCities), expectedFast)
		} else {
			t.Logf("✅ %d/%d pre-seeded cities were fast", fastLookups, len(preSeededCities))
		}
	})

	t.Run("Cache misses work but are appropriately slow", func(t *testing.T) {
		// Use a city that's definitely not pre-seeded
		uniqueCity := "TestCity" + strconv.FormatInt(time.Now().UnixNano(), 10)
		
		start := time.Now()
		cmd := exec.Command("./bin/geotz", "--format=alfred", uniqueCity)
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)
		
		// Should fail (invalid city) but should take time to fail (proves geocoding attempted)
		if err == nil {
			t.Errorf("Expected error for invalid city %s, but got output: %s", uniqueCity, output)
		}
		
		if duration < cacheMissMinDuration {
			t.Errorf("Cache miss for %s was too fast (%v), expected at least %v (proves geocoding happened)", 
				uniqueCity, duration, cacheMissMinDuration)
		} else {
			t.Logf("✅ Cache miss appropriately slow: %v", duration)
		}
	})

	t.Run("Cache provides significant speedup", func(t *testing.T) {
		// Measure the performance difference between cache hit and cache miss
		cacheHitTime := measureCityLookup(t, "london")    // Pre-seeded
		
		// Use a guaranteed cache miss - invalid city that will fail but take time to geocode
		uniqueCity := "InvalidCity" + strconv.FormatInt(time.Now().UnixNano(), 10)
		cacheMissTime := measureInvalidCityLookup(t, uniqueCity) // Guaranteed cache miss
		
		if cacheMissTime <= cacheHitTime {
			t.Errorf("Cache miss (%v) should be slower than cache hit (%v)", cacheMissTime, cacheHitTime)
		}
		
		speedupRatio := float64(cacheMissTime) / float64(cacheHitTime)
		expectedMinSpeedup := 10.0 // Cache should be at least 10x faster
		
		if speedupRatio < expectedMinSpeedup {
			t.Errorf("Cache speedup ratio %.1fx is below expected minimum %.1fx", 
				speedupRatio, expectedMinSpeedup)
		} else {
			t.Logf("✅ Cache provides %.1fx speedup (%v vs %v)", 
				speedupRatio, cacheHitTime, cacheMissTime)
		}
	})

	t.Run("Non-pre-seeded cities work and get cached", func(t *testing.T) {
		// Test a real city that's not pre-seeded
		start := time.Now()
		cmd := exec.Command("./bin/geotz", "--format=alfred", "zurich")
		output, err := cmd.Output()
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("CLI command failed for non-pre-seeded city: %v", err)
		}

		result := strings.TrimSpace(string(output))
		if !strings.Contains(result, "Europe/Zurich") {
			t.Fatalf("Expected Europe/Zurich, got %s", result)
		}

		t.Logf("✓ zurich -> %s (%v) [cache miss, now cached]", result, duration)

		// Second call should be fast (cached)
		start = time.Now()
		cmd = exec.Command("./bin/geotz", "--format=alfred", "zurich")
		output, err = cmd.Output()
		duration = time.Since(start)
		
		if err != nil {
			t.Fatalf("CLI command failed for cached city: %v", err)
		}

		if duration > 500*time.Millisecond {
			t.Errorf("Cached city zurich took %v, expected under 500ms", duration)
		}

		t.Logf("✓ zurich -> cached (%v) [cache hit]", duration)
	})
}

func measureCityLookup(t *testing.T, city string) time.Duration {
	start := time.Now()
	
	cmd := exec.Command("./bin/geotz", "--format=alfred", city)
	output, err := cmd.Output()
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Failed to lookup %s: %v", city, err)
	}
	
	result := strings.TrimSpace(string(output))
	if result == "" {
		t.Fatalf("Empty output for %s", city)
	}
	
	return duration
}

func measureInvalidCityLookup(t *testing.T, city string) time.Duration {
	start := time.Now()
	
	cmd := exec.Command("./bin/geotz", "--format=alfred", city)
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)
	
	// Should fail (invalid city) but should take time to fail (proves geocoding attempted)
	if err == nil {
		t.Logf("Warning: Expected error for invalid city %s, but got output: %s", city, output)
	}
	
	return duration
}