package presenter

import (
	"strings"
	"testing"

	"github.com/loginx/alfred-timein/internal/domain"
)

func TestPlainFormatter_ShouldFormatTimezoneInfoWithNewline(t *testing.T) {
	// Given a plain formatter and a timezone
	formatter := NewPlainFormatter()
	timezone, _ := domain.NewTimezone("Asia/Tokyo")
	
	// When formatting timezone info
	output, err := formatter.FormatTimezoneInfo(timezone, "Tokyo", false)
	
	// Then it should return the timezone string with newline
	if err != nil {
		t.Fatalf("Expected successful formatting, got error: %v", err)
	}
	
	expected := "Asia/Tokyo\n"
	if string(output) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(output))
	}
}

func TestPlainFormatter_ShouldFormatTimeInfoAsHumanReadable(t *testing.T) {
	// Given a plain formatter and a timezone
	formatter := NewPlainFormatter()
	timezone, _ := domain.NewTimezone("Europe/London")
	
	// When formatting time info
	output, err := formatter.FormatTimeInfo(timezone)
	
	// Then it should return human-readable time with newline
	if err != nil {
		t.Fatalf("Expected successful formatting, got error: %v", err)
	}
	
	result := string(output)
	
	// Should end with newline
	if !strings.HasSuffix(result, "\n") {
		t.Error("Expected output to end with newline")
	}
	
	// Should contain readable date format
	if !strings.Contains(result, "2025") {
		t.Error("Expected output to contain current year")
	}
	
	// Should contain day of week
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	containsDay := false
	for _, day := range days {
		if strings.Contains(result, day) {
			containsDay = true
			break
		}
	}
	if !containsDay {
		t.Errorf("Expected output to contain day of week, got '%s'", result)
	}
}

func TestPlainFormatter_ShouldFormatErrorsWithPrefix(t *testing.T) {
	// Given a plain formatter and an error message
	formatter := NewPlainFormatter()
	
	// When formatting an error
	output, err := formatter.FormatError("Connection failed")
	
	// Then it should return error with prefix and newline
	if err != nil {
		t.Fatalf("Expected successful error formatting, got error: %v", err)
	}
	
	expected := "Error: Connection failed\n"
	if string(output) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(output))
	}
}

func TestPlainFormatter_ShouldIgnoreCacheFlag(t *testing.T) {
	// Given a plain formatter and a timezone
	formatter := NewPlainFormatter()
	timezone, _ := domain.NewTimezone("Australia/Sydney")
	
	// When formatting with and without cache flag
	outputCached, _ := formatter.FormatTimezoneInfo(timezone, "Sydney", true)
	outputNotCached, _ := formatter.FormatTimezoneInfo(timezone, "Sydney", false)
	
	// Then both outputs should be identical
	if string(outputCached) != string(outputNotCached) {
		t.Error("Plain formatter should ignore cache flag, but outputs differ")
	}
	
	// And both should be the timezone string
	expected := "Australia/Sydney\n"
	if string(outputCached) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(outputCached))
	}
}