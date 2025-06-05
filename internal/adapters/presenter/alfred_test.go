package presenter

import (
	"encoding/json"
	"testing"

	"github.com/loginx/alfred-timein/internal/domain"
)

func TestAlfredFormatter_ShouldFormatValidTimezoneInfoWithCache(t *testing.T) {
	// Given an Alfred formatter and a timezone
	formatter := NewAlfredFormatter()
	timezone, _ := domain.NewTimezone("Europe/Paris")
	
	// When formatting timezone info with cached flag
	output, err := formatter.FormatTimezoneInfo(timezone, "Paris", true)
	
	// Then it should produce valid Alfred JSON
	if err != nil {
		t.Fatalf("Expected successful formatting, got error: %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}
	
	// And it should contain the expected structure
	items := result["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}
	
	item := items[0].(map[string]interface{})
	
	// Title should be the timezone
	if item["title"] != "Europe/Paris" {
		t.Errorf("Expected title 'Europe/Paris', got '%v'", item["title"])
	}
	
	// Subtitle should indicate it's cached
	subtitle := item["subtitle"].(string)
	if !contains(subtitle, "Paris") || !contains(subtitle, "cached") {
		t.Errorf("Expected subtitle to contain 'Paris' and 'cached', got '%s'", subtitle)
	}
	
	// Should have cache configuration
	cache := result["cache"].(map[string]interface{})
	if cache["seconds"].(float64) != 604800 { // 7 days
		t.Errorf("Expected cache seconds to be 604800, got %v", cache["seconds"])
	}
}

func TestAlfredFormatter_ShouldFormatTimeInfoWithAbbreviation(t *testing.T) {
	// Given an Alfred formatter and a timezone
	formatter := NewAlfredFormatter()
	timezone, _ := domain.NewTimezone("America/New_York")
	
	// When formatting time info
	output, err := formatter.FormatTimeInfo(timezone)
	
	// Then it should produce valid Alfred JSON with time
	if err != nil {
		t.Fatalf("Expected successful formatting, got error: %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}
	
	items := result["items"].([]interface{})
	item := items[0].(map[string]interface{})
	
	// Title should contain timezone and formatted time
	title := item["title"].(string)
	if !contains(title, "America/New_York") {
		t.Errorf("Expected title to contain timezone, got '%s'", title)
	}
	
	// Subtitle should contain city and timezone abbreviation
	subtitle := item["subtitle"].(string)
	if !contains(subtitle, "New York") {
		t.Errorf("Expected subtitle to contain 'New York', got '%s'", subtitle)
	}
	
	// Should have shorter cache (60 seconds for time)
	cache := result["cache"].(map[string]interface{})
	if cache["seconds"].(float64) != 60 {
		t.Errorf("Expected cache seconds to be 60, got %v", cache["seconds"])
	}
}

func TestAlfredFormatter_ShouldFormatErrorsAsInvalidItems(t *testing.T) {
	// Given an Alfred formatter and an error message
	formatter := NewAlfredFormatter()
	
	// When formatting an error
	output, err := formatter.FormatError("Something went wrong")
	
	// Then it should produce valid Alfred JSON
	if err != nil {
		t.Fatalf("Expected successful error formatting, got error: %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}
	
	items := result["items"].([]interface{})
	item := items[0].(map[string]interface{})
	
	// Should be marked as an error
	if item["title"] != "Error" {
		t.Errorf("Expected title 'Error', got '%v'", item["title"])
	}
	
	// Should contain the error message
	if item["subtitle"] != "Something went wrong" {
		t.Errorf("Expected subtitle 'Something went wrong', got '%v'", item["subtitle"])
	}
	
	// Should be marked as invalid (not actionable)
	if item["valid"] != false {
		t.Errorf("Expected valid to be false for errors, got %v", item["valid"])
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || 
		   len(s) > len(substr) && s[:len(substr)] == substr ||
		   len(s) > len(substr) && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}