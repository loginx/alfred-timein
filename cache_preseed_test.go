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
)

func TestCachePreSeedBehavior(t *testing.T) {
	testDir := fmt.Sprintf("/tmp/alfred-timein-cache-test-%d", time.Now().UnixNano())
	defer os.RemoveAll(testDir)

	// Create test cache
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

	// Pre-seed the cache
	c.PreSeed(entries)

	// Test 1: Cache should contain pre-seeded entries
	t.Run("Pre-seeded entries exist", func(t *testing.T) {
		for key, expectedValue := range entries {
			if value, ok := c.Get(key); !ok {
				t.Errorf("Expected key %s to exist in cache", key)
			} else if value != expectedValue {
				t.Errorf("Expected value %s for key %s, got %s", expectedValue, key, value)
			}
		}
	})

	// Test 2: Cache lookups should be fast for pre-seeded entries
	t.Run("Pre-seeded cache lookups are fast", func(t *testing.T) {
		testCases := []string{"london", "paris", "tokyo"}
		
		for _, city := range testCases {
			start := time.Now()
			if value, ok := c.Get(city); !ok {
				t.Errorf("Expected city %s to be in cache", city)
			} else {
				duration := time.Since(start)
				t.Logf("Cache lookup for %s took %v, returned %s", city, duration, value)
				
				// Cache hits should be very fast (under 1ms)
				if duration > 1*time.Millisecond {
					t.Errorf("Cache lookup for %s took too long: %v", city, duration)
				}
			}
		}
	})

	// Test 3: User entries should override pre-seeded entries
	t.Run("User entries override pre-seeded", func(t *testing.T) {
		// Add a user entry for London
		c.Set("london", "User/Custom_Timezone")
		
		// Should get the user entry, not the pre-seeded one
		if value, ok := c.Get("london"); !ok {
			t.Error("Expected london to be in cache")
		} else if value != "User/Custom_Timezone" {
			t.Errorf("Expected user entry to override pre-seeded entry, got %s", value)
		}
	})

	// Test 4: Pre-seeded entries should have long TTL
	t.Run("Pre-seeded entries have long TTL", func(t *testing.T) {
		// This is harder to test directly, but we can verify the entries persist
		// after some time has passed
		time.Sleep(10 * time.Millisecond)
		
		if value, ok := c.Get("paris"); !ok {
			t.Error("Expected paris to still be in cache after time delay")
		} else if !strings.Contains(value, "Europe/Paris") {
			t.Errorf("Expected Europe/Paris, got %s", value)
		}
	})
}