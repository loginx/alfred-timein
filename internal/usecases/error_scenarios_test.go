package usecases

import (
	"fmt"
	"testing"

	"github.com/loginx/alfred-timein/internal/domain"
)

// MockFailingFormatter simulates formatter errors
type MockFailingFormatter struct {
	shouldFailOnTimeInfo     bool
	shouldFailOnTimezoneInfo bool
	shouldFailOnError        bool
}

func (m *MockFailingFormatter) FormatTimezoneInfo(timezone *domain.Timezone, city string, cached bool) ([]byte, error) {
	if m.shouldFailOnTimezoneInfo {
		return nil, fmt.Errorf("formatter failed on timezone info")
	}
	return []byte("formatted timezone info"), nil
}

func (m *MockFailingFormatter) FormatTimeInfo(timezone *domain.Timezone) ([]byte, error) {
	if m.shouldFailOnTimeInfo {
		return nil, fmt.Errorf("formatter failed on time info")
	}
	return []byte("formatted time info"), nil
}

func (m *MockFailingFormatter) FormatError(message string) ([]byte, error) {
	if m.shouldFailOnError {
		return nil, fmt.Errorf("formatter failed on error")
	}
	return []byte("formatted error"), nil
}

func TestTimeinUseCase_ShouldHandleFormatterFailures(t *testing.T) {
	// Given a use case with a failing formatter
	formatter := &MockFailingFormatter{shouldFailOnTimeInfo: true}
	uc := NewTimeinUseCase(formatter)
	
	// When trying to format time info
	output, err := uc.GetTimezoneInfo("America/New_York")
	
	// Then it should handle the formatter error
	if err == nil {
		t.Fatal("Expected error when formatter fails, got none")
	}
	
	// Output might be nil when formatter fails (this is acceptable)
	t.Logf("Got output: %v, error: %v", output, err)
}

func TestGeotzUseCase_ShouldHandleMultipleFailures(t *testing.T) {
	// Given a use case with all dependencies that can fail
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{shouldFail: true}
	tzFinder := &MockTimezoneFinder{shouldFail: true}
	cache := NewMockCache()
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	// When geocoding fails
	output, err := uc.GetTimezoneFromCity("Unknown City")
	
	// Then it should handle gracefully
	if err == nil {
		t.Fatal("Expected error when geocoding fails")
	}
	
	// And provide user-friendly error output
	if string(output) != "mock error" {
		t.Errorf("Expected error output, got '%s'", string(output))
	}
}

func TestGeotzUseCase_ShouldHandleCacheCorruption(t *testing.T) {
	// Given a cache with invalid timezone data
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{}
	tzFinder := &MockTimezoneFinder{}
	cache := NewMockCache()
	
	// Pre-populate cache with invalid timezone
	cache.Set("test city", "Invalid/Corrupt/Timezone")
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	// When trying to use cached data
	_, err := uc.GetTimezoneFromCity("Test City")
	
	// Then it should handle the invalid cached data
	if err == nil {
		t.Fatal("Expected error for corrupted cache data")
	}
	
	// And provide error output
	if !formatter.formatErrorCalled {
		t.Error("Expected error to be formatted for user")
	}
}

func TestGeotzUseCase_ShouldHandleWhitespaceOnlyInput(t *testing.T) {
	// Given a use case and whitespace-only input
	formatter := &MockFormatter{}
	geocoder := &MockGeocoder{}
	tzFinder := &MockTimezoneFinder{}
	cache := NewMockCache()
	
	uc := NewGeotzUseCase(geocoder, tzFinder, cache, formatter)
	
	// When providing only whitespace
	_, err := uc.GetTimezoneFromCity("   \t\n   ")
	
	// Then it should treat as empty and error
	if err == nil {
		t.Fatal("Expected error for whitespace-only input")
	}
	
	// And provide helpful error message
	if !formatter.formatErrorCalled {
		t.Error("Expected error to be formatted for user")
	}
}