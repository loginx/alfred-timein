package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/loginx/alfred-timein/internal/adapters/cache"
	"github.com/loginx/alfred-timein/internal/adapters/timezonefinder"
)

type Capital struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: preseed <cache-directory>")
	}

	cacheDir := os.Args[1]
	
	// Load capitals data
	capitalsData, err := os.ReadFile("data/capitals.json")
	if err != nil {
		log.Fatalf("Failed to read capitals data: %v", err)
	}

	var capitals []Capital
	if err := json.Unmarshal(capitalsData, &capitals); err != nil {
		log.Fatalf("Failed to parse capitals data: %v", err)
	}

	// Initialize timezone finder
	tzf, err := timezonefinder.NewTzfTimezoneFinder()
	if err != nil {
		log.Fatalf("Failed to initialize timezone finder: %v", err)
	}

	// Create cache
	c := cache.NewLRUCache(1000, 0, cacheDir) // Large cache size, TTL not used for pre-seeding

	// Prepare pre-seed entries
	entries := make(map[string]string)
	
	fmt.Printf("Pre-seeding cache with %d capitals...\n", len(capitals))
	
	for _, capital := range capitals {
		timezone, err := tzf.GetTimezoneName(capital.Lng, capital.Lat)
		if err != nil {
			fmt.Printf("Warning: Failed to get timezone for %s: %v\n", capital.Name, err)
			continue
		}
		
		// Create cache key for city name only
		cityKey := strings.ToLower(capital.Name)
		entries[cityKey] = timezone
		
		fmt.Printf("  %s, %s -> %s\n", capital.Name, capital.Country, timezone)
	}

	// Pre-seed the cache
	c.PreSeed(entries)

	fmt.Printf("Successfully pre-seeded cache with %d cities in %s\n", len(entries), filepath.Join(cacheDir, "geotz_cache.json"))
}