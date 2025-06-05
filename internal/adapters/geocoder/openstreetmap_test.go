package geocoder

import (
	"testing"
)

func TestOpenStreetMapGeocoder_ShouldReturnLocationForKnownCity(t *testing.T) {
	// Given a geocoder and a well-known city
	geocoder := NewOpenStreetMapGeocoder()
	
	// When geocoding the city
	location, err := geocoder.Geocode("Paris")
	
	// Then it should return a valid location without error
	if err != nil {
		t.Fatalf("Expected successful geocoding for Paris, got error: %v", err)
	}
	
	// And the location should have reasonable coordinates for Paris
	if location.Latitude < 48 || location.Latitude > 49 {
		t.Errorf("Expected Paris latitude around 48.8, got %f", location.Latitude)
	}
	if location.Longitude < 2 || location.Longitude > 3 {
		t.Errorf("Expected Paris longitude around 2.3, got %f", location.Longitude)
	}
	if location.Name != "Paris" {
		t.Errorf("Expected location name 'Paris', got '%s'", location.Name)
	}
}

func TestOpenStreetMapGeocoder_ShouldFailGracefullyForNonexistentPlace(t *testing.T) {
	// Given a geocoder and a nonsense query
	geocoder := NewOpenStreetMapGeocoder()
	
	// When geocoding a place that doesn't exist
	location, err := geocoder.Geocode("XYZ123NotARealPlace456")
	
	// Then it should return an error
	if err == nil {
		t.Fatal("Expected error for nonexistent place, got none")
	}
	
	// And no location should be returned
	if location != nil {
		t.Errorf("Expected nil location for failed geocoding, got %v", location)
	}
}

func TestOpenStreetMapGeocoder_ShouldHandleEmptyQuery(t *testing.T) {
	// Given a geocoder and an empty query
	geocoder := NewOpenStreetMapGeocoder()
	
	// When geocoding with empty string
	location, err := geocoder.Geocode("")
	
	// Then it should return an error
	if err == nil {
		t.Fatal("Expected error for empty query, got none")
	}
	
	// And no location should be returned
	if location != nil {
		t.Errorf("Expected nil location for empty query, got %v", location)
	}
}

func TestOpenStreetMapGeocoder_ShouldHandleLandmarks(t *testing.T) {
	// Given a geocoder and a famous landmark
	geocoder := NewOpenStreetMapGeocoder()
	
	// When geocoding a landmark
	location, err := geocoder.Geocode("Eiffel Tower")
	
	// Then it should return a location in Paris
	if err != nil {
		t.Fatalf("Expected successful geocoding for Eiffel Tower, got error: %v", err)
	}
	
	// The coordinates should be in Paris
	if location.Latitude < 48 || location.Latitude > 49 {
		t.Errorf("Expected Eiffel Tower latitude in Paris range, got %f", location.Latitude)
	}
}