package domain

import (
	"testing"
)

func TestNewLocation_Valid(t *testing.T) {
	loc, err := NewLocation("New York", 40.7128, -74.0060)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != "New York" {
		t.Errorf("expected 'New York', got %s", loc.Name)
	}
	if loc.Latitude != 40.7128 {
		t.Errorf("expected 40.7128, got %f", loc.Latitude)
	}
	if loc.Longitude != -74.0060 {
		t.Errorf("expected -74.0060, got %f", loc.Longitude)
	}
}

func TestNewLocation_InvalidLatitude(t *testing.T) {
	_, err := NewLocation("Test", 91.0, 0)
	if err == nil {
		t.Errorf("expected error for latitude > 90")
	}
	_, err = NewLocation("Test", -91.0, 0)
	if err == nil {
		t.Errorf("expected error for latitude < -90")
	}
}

func TestNewLocation_InvalidLongitude(t *testing.T) {
	_, err := NewLocation("Test", 0, 181.0)
	if err == nil {
		t.Errorf("expected error for longitude > 180")
	}
	_, err = NewLocation("Test", 0, -181.0)
	if err == nil {
		t.Errorf("expected error for longitude < -180")
	}
}

func TestNewLocation_EmptyName(t *testing.T) {
	_, err := NewLocation("", 0, 0)
	if err == nil {
		t.Errorf("expected error for empty name")
	}
}