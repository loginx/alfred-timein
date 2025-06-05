package usecases

import (
	"testing"

	"github.com/loginx/alfred-timein/internal/domain"
)

// MockFormatter for testing
type MockFormatter struct {
	formatTimeInfoCalled bool
	formatErrorCalled    bool
	lastError           string
}

func (m *MockFormatter) FormatTimezoneInfo(timezone *domain.Timezone, city string, cached bool) ([]byte, error) {
	return []byte("mock timezone info"), nil
}

func (m *MockFormatter) FormatTimeInfo(timezone *domain.Timezone) ([]byte, error) {
	m.formatTimeInfoCalled = true
	return []byte("mock time info"), nil
}

func (m *MockFormatter) FormatError(message string) ([]byte, error) {
	m.formatErrorCalled = true
	m.lastError = message
	return []byte("mock error"), nil
}

func TestTimeinUseCase_ShouldFormatCurrentTimeForValidTimezone(t *testing.T) {
	// Given a timein use case and a valid timezone
	formatter := &MockFormatter{}
	uc := NewTimeinUseCase(formatter)
	
	// When getting timezone info for a known timezone
	output, err := uc.GetTimezoneInfo("America/New_York")
	
	// Then it should successfully format time information
	if err != nil {
		t.Fatalf("Expected successful time formatting for valid timezone, got error: %v", err)
	}
	
	// And the formatter should be called to format time info
	if !formatter.formatTimeInfoCalled {
		t.Error("Expected FormatTimeInfo to be called for valid timezone")
	}
	
	// And produce the expected output
	if string(output) != "mock time info" {
		t.Errorf("Expected formatted time output, got '%s'", string(output))
	}
}

func TestTimeinUseCase_ShouldRejectInvalidTimezoneGracefully(t *testing.T) {
	// Given a timein use case and an invalid timezone
	formatter := &MockFormatter{}
	uc := NewTimeinUseCase(formatter)
	
	// When trying to get info for an invalid timezone
	output, err := uc.GetTimezoneInfo("Invalid/Timezone")
	
	// Then it should return an error
	if err == nil {
		t.Fatal("Expected error for invalid timezone, but got none")
	}
	
	// And format an error message for the user
	if !formatter.formatErrorCalled {
		t.Error("Expected error to be formatted for user display")
	}
	
	// And provide error output
	if string(output) != "mock error" {
		t.Errorf("Expected error output, got '%s'", string(output))
	}
}

func TestTimeinUseCase_ShouldProvideCompleteTimezoneInformation(t *testing.T) {
	// Given a timein use case
	formatter := &MockFormatter{}
	uc := NewTimeinUseCase(formatter)
	
	// When requesting detailed timezone information
	info, err := uc.GetTimezoneInfoForFormatting("America/New_York")
	
	// Then it should successfully provide complete information
	if err != nil {
		t.Fatalf("Expected successful timezone info extraction, got error: %v", err)
	}
	
	// And include the correct timezone
	if info.Timezone.String() != "America/New_York" {
		t.Errorf("Expected timezone 'America/New_York', got '%s'", info.Timezone.String())
	}
	
	// And extract the human-readable city name
	if info.City != "New York" {
		t.Errorf("Expected city 'New York', got '%s'", info.City)
	}
	
	// And provide current time in that timezone
	if info.CurrentTime.IsZero() {
		t.Error("Expected current time to be set")
	}
	
	// And include timezone abbreviation for display
	if info.Abbreviation == "" {
		t.Error("Expected timezone abbreviation to be provided")
	}
}