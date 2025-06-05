package timezonefinder

import (
	"testing"
)

func TestTzfTimezoneFinder_ShouldFindTimezoneForNewYorkCoordinates(t *testing.T) {
	// Given a timezone finder and coordinates for New York City
	finder, err := NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to create timezone finder: %v", err)
	}
	
	// When looking up timezone for NYC coordinates
	timezone, err := finder.GetTimezoneName(-74.0060, 40.7128) // NYC coordinates
	
	// Then it should return the correct timezone
	if err != nil {
		t.Fatalf("Expected successful timezone lookup, got error: %v", err)
	}
	
	if timezone != "America/New_York" {
		t.Errorf("Expected 'America/New_York' for NYC coordinates, got '%s'", timezone)
	}
}

func TestTzfTimezoneFinder_ShouldFindTimezoneForTokyoCoordinates(t *testing.T) {
	// Given a timezone finder and coordinates for Tokyo
	finder, err := NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to create timezone finder: %v", err)
	}
	
	// When looking up timezone for Tokyo coordinates
	timezone, err := finder.GetTimezoneName(139.6917, 35.6895) // Tokyo coordinates
	
	// Then it should return the correct timezone
	if err != nil {
		t.Fatalf("Expected successful timezone lookup, got error: %v", err)
	}
	
	if timezone != "Asia/Tokyo" {
		t.Errorf("Expected 'Asia/Tokyo' for Tokyo coordinates, got '%s'", timezone)
	}
}

func TestTzfTimezoneFinder_ShouldHandleOceanCoordinates(t *testing.T) {
	// Given a timezone finder and coordinates in the middle of the ocean
	finder, err := NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to create timezone finder: %v", err)
	}
	
	// When looking up timezone for ocean coordinates
	timezone, err := finder.GetTimezoneName(-160.0, 25.0) // Middle of Pacific
	
	// Then it might return an error or a UTC-based timezone
	// (The exact behavior depends on the tzf library implementation)
	if err != nil {
		// This is acceptable - ocean coordinates might not have timezones
		return
	}
	
	// If it returns a timezone, it should be a valid IANA string
	if timezone == "" {
		t.Error("Expected either error or valid timezone for ocean coordinates")
	}
}

func TestTzfTimezoneFinder_ShouldHandleInvalidCoordinates(t *testing.T) {
	// Given a timezone finder and invalid coordinates
	finder, err := NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to create timezone finder: %v", err)
	}
	
	// When looking up timezone for coordinates outside valid range
	timezone, err := finder.GetTimezoneName(200.0, 100.0) // Invalid coordinates
	
	// Then it should handle gracefully (error or empty string)
	if err != nil {
		// Error is acceptable for invalid coordinates
		return
	}
	
	if timezone == "" {
		// Empty string is also acceptable
		return
	}
	
	// If it returns something, it shouldn't crash
	t.Logf("Got timezone '%s' for invalid coordinates (library handled gracefully)", timezone)
}

func TestTzfTimezoneFinder_ShouldFindTimezoneForLondonCoordinates(t *testing.T) {
	// Given a timezone finder and coordinates for London
	finder, err := NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to create timezone finder: %v", err)
	}
	
	// When looking up timezone for London coordinates
	timezone, err := finder.GetTimezoneName(-0.1276, 51.5074) // London coordinates
	
	// Then it should return the correct timezone
	if err != nil {
		t.Fatalf("Expected successful timezone lookup, got error: %v", err)
	}
	
	if timezone != "Europe/London" {
		t.Errorf("Expected 'Europe/London' for London coordinates, got '%s'", timezone)
	}
}