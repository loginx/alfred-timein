package domain

import (
	"testing"
)

func TestNewTimezone_Valid(t *testing.T) {
	tz, err := NewTimezone("America/New_York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz.Name != "America/New_York" {
		t.Errorf("expected 'America/New_York', got %s", tz.Name)
	}
}

func TestNewTimezone_Invalid(t *testing.T) {
	_, err := NewTimezone("Invalid/Timezone")
	if err == nil {
		t.Errorf("expected error for invalid timezone")
	}
}

func TestNewTimezone_Empty(t *testing.T) {
	_, err := NewTimezone("")
	if err == nil {
		t.Errorf("expected error for empty timezone")
	}
	_, err = NewTimezone("   ")
	if err == nil {
		t.Errorf("expected error for whitespace-only timezone")
	}
}

func TestTimezone_Location(t *testing.T) {
	tz, _ := NewTimezone("Europe/London")
	loc, err := tz.Location()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Europe/London" {
		t.Errorf("expected 'Europe/London', got %s", loc.String())
	}
}

func TestTimezone_City(t *testing.T) {
	tests := []struct {
		timezone string
		expected string
	}{
		{"America/New_York", "New York"},
		{"Europe/London", "London"},
		{"Asia/Ho_Chi_Minh", "Ho Chi Minh"},
		{"UTC", "UTC"},
	}

	for _, test := range tests {
		tz, _ := NewTimezone(test.timezone)
		city := tz.City()
		if city != test.expected {
			t.Errorf("for timezone %s, expected city %s, got %s", test.timezone, test.expected, city)
		}
	}
}