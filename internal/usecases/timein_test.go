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

func TestTimeinUseCase_GetTimezoneInfo_Valid(t *testing.T) {
	formatter := &MockFormatter{}
	uc := NewTimeinUseCase(formatter)
	
	output, err := uc.GetTimezoneInfo("America/New_York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !formatter.formatTimeInfoCalled {
		t.Errorf("expected FormatTimeInfo to be called")
	}
	
	if string(output) != "mock time info" {
		t.Errorf("expected 'mock time info', got %s", string(output))
	}
}

func TestTimeinUseCase_GetTimezoneInfo_Invalid(t *testing.T) {
	formatter := &MockFormatter{}
	uc := NewTimeinUseCase(formatter)
	
	output, err := uc.GetTimezoneInfo("Invalid/Timezone")
	if err == nil {
		t.Fatalf("expected error for invalid timezone")
	}
	
	if !formatter.formatErrorCalled {
		t.Errorf("expected FormatError to be called")
	}
	
	if string(output) != "mock error" {
		t.Errorf("expected 'mock error', got %s", string(output))
	}
}

func TestTimeinUseCase_GetTimezoneInfoForFormatting(t *testing.T) {
	formatter := &MockFormatter{}
	uc := NewTimeinUseCase(formatter)
	
	info, err := uc.GetTimezoneInfoForFormatting("America/New_York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if info.Timezone.String() != "America/New_York" {
		t.Errorf("expected 'America/New_York', got %s", info.Timezone.String())
	}
	
	if info.City != "New York" {
		t.Errorf("expected 'New York', got %s", info.City)
	}
	
	if info.CurrentTime.IsZero() {
		t.Errorf("expected non-zero current time")
	}
	
	if info.Abbreviation == "" {
		t.Errorf("expected non-empty abbreviation")
	}
}