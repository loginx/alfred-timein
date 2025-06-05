package timezonefinder

import (
	"testing"
)

func TestTzfTimezoneFinder_GetTimezoneName(t *testing.T) {
	finder, err := NewTzfTimezoneFinder()
	if err != nil {
		t.Fatalf("Failed to create timezone finder: %v", err)
	}

	tests := []struct {
		name      string
		latitude  float64
		longitude float64
		expected  string
	}{
		{
			name:      "Tokyo",
			latitude:  35.6762,
			longitude: 139.6503,
			expected:  "Asia/Tokyo",
		},
		{
			name:      "New York", 
			latitude:  40.7128,
			longitude: -74.0060,
			expected:  "America/New_York",
		},
		{
			name:      "London",
			latitude:  51.5074,
			longitude: -0.1278,
			expected:  "Europe/London",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := finder.GetTimezoneName(tt.longitude, tt.latitude)
			if err != nil {
				t.Errorf("GetTimezoneName() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("GetTimezoneName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}