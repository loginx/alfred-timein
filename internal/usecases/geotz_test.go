package usecases

import (
	"fmt"
	"testing"

	"github.com/loginx/alfred-timein/internal/domain"
)

// MockGeocoder for testing
type MockGeocoder struct {
	shouldFail bool
}

func (m *MockGeocoder) Geocode(query string) (*domain.Location, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("geocoding failed")
	}
	return domain.NewLocation(query, 40.7128, -74.0060)
}

// MockTimezoneFinder for testing
type MockTimezoneFinder struct {
	shouldFail bool
}

func (m *MockTimezoneFinder) GetTimezoneName(longitude, latitude float64) (string, error) {
	if m.shouldFail {
		return "", fmt.Errorf("timezone lookup failed")
	}
	return "America/New_York", nil
}

// MockCache for testing
type MockCache struct {
	data map[string]string
}

func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]string),
	}
}

func (m *MockCache) Get(key string) (string, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *MockCache) Set(key, value string) {
	m.data[key] = value
}

func (m *MockCache) Clear() {
	m.data = make(map[string]string)
}

func TestGeotzUseCase_GetTimezoneFromCity_Valid(t *testing.T) {
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{}
	tzFinder := &MockTimezoneFinder{}
	cache := NewMockCache()
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	output, err := uc.GetTimezoneFromCity("New York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if string(output) != "mock timezone info" {
		t.Errorf("expected 'mock timezone info', got %s", string(output))
	}
	
	// Check if result was cached
	if cached, ok := cache.Get("new york"); !ok || cached != "America/New_York" {
		t.Errorf("expected result to be cached")
	}
}

func TestGeotzUseCase_GetTimezoneFromCity_Cached(t *testing.T) {
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{}
	tzFinder := &MockTimezoneFinder{}
	cache := NewMockCache()
	
	// Pre-populate cache
	cache.Set("new york", "America/New_York")
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	output, err := uc.GetTimezoneFromCity("New York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if string(output) != "mock timezone info" {
		t.Errorf("expected 'mock timezone info', got %s", string(output))
	}
}

func TestGeotzUseCase_GetTimezoneFromCity_GeocodingFails(t *testing.T) {
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{shouldFail: true}
	tzFinder := &MockTimezoneFinder{}
	cache := NewMockCache()
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	output, err := uc.GetTimezoneFromCity("Unknown City")
	if err == nil {
		t.Fatalf("expected error for failed geocoding")
	}
	
	if !formatter.formatErrorCalled {
		t.Errorf("expected FormatError to be called")
	}
	
	if string(output) != "mock error" {
		t.Errorf("expected 'mock error', got %s", string(output))
	}
}

func TestGeotzUseCase_GetTimezoneFromCity_TimezoneLookupFails(t *testing.T) {
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{}
	tzFinder := &MockTimezoneFinder{shouldFail: true}
	cache := NewMockCache()
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	_, err := uc.GetTimezoneFromCity("Unknown City")
	if err == nil {
		t.Fatalf("expected error for failed timezone lookup")
	}
	
	if !formatter.formatErrorCalled {
		t.Errorf("expected FormatError to be called")
	}
}

func TestGeotzUseCase_GetTimezoneFromCity_EmptyCity(t *testing.T) {
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{}
	tzFinder := &MockTimezoneFinder{}
	cache := NewMockCache()
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	_, err := uc.GetTimezoneFromCity("")
	if err == nil {
		t.Fatalf("expected error for empty city")
	}
	
	if !formatter.formatErrorCalled {
		t.Errorf("expected FormatError to be called")
	}
}