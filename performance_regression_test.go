//go:build bdd
// +build bdd

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Performance thresholds - these are our SLA requirements
const (
	cacheHitMaxDuration  = 100 * time.Millisecond // Pre-seeded cache hits must be under 100ms
	cacheMissMinDuration = 200 * time.Millisecond // Cache misses should take at least 200ms (proves geocoding happened)
)

// TestPerformanceRegression ensures cache pre-seeding provides expected performance improvements
func TestPerformanceRegression(t *testing.T) {
	// Ensure we have binaries built
	if _, err := os.Stat("bin/geotz"); os.IsNotExist(err) {
		t.Skip("Skipping performance test - bin/geotz not found. Run 'make build' first.")
	}

	// Ensure we have pre-seeded cache
	if _, err := os.Stat("geotz_cache.json"); os.IsNotExist(err) {
		t.Skip("Skipping performance test - geotz_cache.json not found. Run 'make preseed' first.")
	}

	t.Run("Pre-seeded cities meet cache hit performance SLA", func(t *testing.T) {
		preSeededCities := []string{
			"london", "paris", "tokyo", "berlin", "madrid", 
			"rome", "amsterdam", "stockholm", "oslo",
		}

		var failures []string
		
		for _, city := range preSeededCities {
			duration := measureCityLookup(t, city)
			
			if duration > cacheHitMaxDuration {
				failure := fmt.Sprintf("%s took %v (exceeds %v SLA)", 
					city, duration, cacheHitMaxDuration)
				failures = append(failures, failure)
				t.Errorf("❌ %s", failure)
			} else {
				t.Logf("✅ %s: %v (under %v SLA)", city, duration, cacheHitMaxDuration)
			}
		}

		if len(failures) > 0 {
			t.Errorf("Performance SLA failures:\n%s", strings.Join(failures, "\n"))
		}
	})

	t.Run("Cache misses still work but are appropriately slow", func(t *testing.T) {
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

	t.Run("Cache effectiveness ratio", func(t *testing.T) {
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
}

// TestCacheWarmupPerformance tests the entire cache pre-seeding workflow
func TestCacheWarmupPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cache warmup test in short mode")
	}

	testDir := "/tmp/alfred-timein-warmup-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	defer os.RemoveAll(testDir)

	// Setup clean environment
	if err := os.MkdirAll(testDir+"/workflow", 0755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	if err := os.MkdirAll(testDir+"/bin", 0755); err != nil {
		t.Fatalf("Failed to create bin dir: %v", err)
	}

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(testDir)

	// Copy required files
	copyFile(t, originalDir+"/bin/geotz", "bin/geotz")
	copyFile(t, originalDir+"/data/capitals.json", "data/capitals.json")

	// Build preseed tool
	buildCmd := exec.Command("go", "build", "-o", "preseed", originalDir+"/cmd/preseed")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build preseed: %v", err)
	}

	t.Run("Cache pre-seeding completes in reasonable time", func(t *testing.T) {
		start := time.Now()
		
		cmd := exec.Command("./preseed")
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)
		
		if err != nil {
			t.Fatalf("Preseed failed: %v, output: %s", err, output)
		}

		maxPreseedTime := 30 * time.Second // Should complete within 30 seconds
		if duration > maxPreseedTime {
			t.Errorf("Cache pre-seeding took %v, expected under %v", duration, maxPreseedTime)
		} else {
			t.Logf("✅ Cache pre-seeding completed in %v", duration)
		}
	})

	t.Run("Pre-seeded cache provides immediate performance benefit", func(t *testing.T) {
		// Test a few pre-seeded cities to ensure they're fast
		testCities := []string{"london", "paris", "tokyo"}
		
		for _, city := range testCities {
			duration := measureCityLookup(t, city)
			
			if duration > cacheHitMaxDuration {
				t.Errorf("Pre-seeded city %s took %v, expected under %v", 
					city, duration, cacheHitMaxDuration)
			} else {
				t.Logf("✅ %s: %v", city, duration)
			}
		}
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