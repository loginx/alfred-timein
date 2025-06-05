package domain

import (
	"testing"
)

func TestTimezone_ShouldHandleUTCTimezone(t *testing.T) {
	// Given UTC timezone
	tz, err := NewTimezone("UTC")
	
	// Then it should be valid
	if err != nil {
		t.Fatalf("Expected UTC to be valid timezone, got error: %v", err)
	}
	
	// And city should be UTC itself
	if tz.City() != "UTC" {
		t.Errorf("Expected UTC city to be 'UTC', got '%s'", tz.City())
	}
}

func TestTimezone_ShouldHandleComplexCityNames(t *testing.T) {
	// Given timezone with complex city name
	tz, err := NewTimezone("America/Argentina/Buenos_Aires")
	
	// Then it should be valid
	if err != nil {
		t.Fatalf("Expected complex timezone to be valid, got error: %v", err)
	}
	
	// And extract the middle part (current implementation behavior)
	if tz.City() != "Argentina" {
		t.Errorf("Expected city 'Argentina', got '%s'", tz.City())
	}
}

func TestTimezone_ShouldHandleWhitespaceInInput(t *testing.T) {
	// Given timezone with leading/trailing whitespace
	tz, err := NewTimezone("  Europe/Paris  ")
	
	// Then it should be valid (trimmed)
	if err != nil {
		t.Fatalf("Expected trimmed timezone to be valid, got error: %v", err)
	}
	
	// And name should be clean
	if tz.String() != "Europe/Paris" {
		t.Errorf("Expected clean timezone name, got '%s'", tz.String())
	}
}

func TestLocation_ShouldHandleBoundaryCoordinates(t *testing.T) {
	// Given locations at coordinate boundaries
	testCases := []struct {
		name string
		lat  float64
		lng  float64
		valid bool
	}{
		{"North Pole", 90.0, 0.0, true},
		{"South Pole", -90.0, 0.0, true},
		{"International Date Line", 0.0, 180.0, true},
		{"Prime Meridian", 0.0, -180.0, true},
		{"Just over north", 90.1, 0.0, false},
		{"Just under south", -90.1, 0.0, false},
		{"Just over east", 0.0, 180.1, false},
		{"Just under west", 0.0, -180.1, false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loc, err := NewLocation(tc.name, tc.lat, tc.lng)
			
			if tc.valid {
				if err != nil {
					t.Errorf("Expected valid location for %s, got error: %v", tc.name, err)
				}
				if loc == nil {
					t.Errorf("Expected location object for valid coordinates")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for invalid coordinates %s, got none", tc.name)
				}
			}
		})
	}
}