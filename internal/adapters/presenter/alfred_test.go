package presenter

import (
	"encoding/json"
	"testing"

	"github.com/loginx/alfred-timein/internal/domain"
)

func TestAlfredFormatter_ShouldFormatValidTimezoneInfoWithCache(t *testing.T) {
	formatter := NewAlfredFormatter()
	timezone, _ := domain.NewTimezone("Europe/Paris")

	output, err := formatter.FormatTimezoneInfo(timezone, "Paris", true)
	if err != nil {
		t.Fatalf("Expected successful formatting, got error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	items := result["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	item := items[0].(map[string]interface{})

	if item["title"] != "Europe/Paris" {
		t.Errorf("Expected title 'Europe/Paris', got '%v'", item["title"])
	}

	subtitle := item["subtitle"].(string)
	if !contains(subtitle, "Paris") || !contains(subtitle, "cached") {
		t.Errorf("Expected subtitle to contain 'Paris' and 'cached', got '%s'", subtitle)
	}

	cache := result["cache"].(map[string]interface{})
	if cache["seconds"].(float64) != 604800 {
		t.Errorf("Expected cache seconds to be 604800, got %v", cache["seconds"])
	}

	// Should have workflow icon
	icon := item["icon"].(map[string]interface{})
	if icon["path"] != "icon.png" {
		t.Errorf("Expected icon path 'icon.png', got '%v'", icon["path"])
	}
}

func TestAlfredFormatter_ShouldFormatTimeInfoWithAbbreviation(t *testing.T) {
	formatter := NewAlfredFormatter()
	timezone, _ := domain.NewTimezone("America/New_York")

	output, err := formatter.FormatTimeInfo(timezone)
	if err != nil {
		t.Fatalf("Expected successful formatting, got error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	items := result["items"].([]interface{})
	item := items[0].(map[string]interface{})

	title := item["title"].(string)
	if !contains(title, "America/New_York") {
		t.Errorf("Expected title to contain timezone, got '%s'", title)
	}

	subtitle := item["subtitle"].(string)
	if !contains(subtitle, "New York") {
		t.Errorf("Expected subtitle to contain 'New York', got '%s'", subtitle)
	}

	cache := result["cache"].(map[string]interface{})
	if cache["seconds"].(float64) != 60 {
		t.Errorf("Expected cache seconds to be 60, got %v", cache["seconds"])
	}

	// Should preserve result ordering
	if result["skipknowledge"] != true {
		t.Errorf("Expected skipknowledge to be true, got %v", result["skipknowledge"])
	}

	// Should have workflow icon
	icon := item["icon"].(map[string]interface{})
	if icon["path"] != "icon.png" {
		t.Errorf("Expected icon path 'icon.png', got '%v'", icon["path"])
	}

	// Should have text for CMD+C and CMD+L
	text := item["text"].(map[string]interface{})
	if text["copy"] != title {
		t.Errorf("Expected text.copy to match title")
	}
	if text["largetype"] != title {
		t.Errorf("Expected text.largetype to match title")
	}

	// Should have modifier keys
	mods := item["mods"].(map[string]interface{})
	cmdMod := mods["cmd"].(map[string]interface{})
	if cmdMod["arg"] != "America/New_York" {
		t.Errorf("Expected cmd mod arg to be timezone, got '%v'", cmdMod["arg"])
	}
	altMod := mods["alt"].(map[string]interface{})
	altArg := altMod["arg"].(string)
	if !contains(altArg, "T") || !contains(altArg, "-") {
		t.Errorf("Expected alt mod arg to be ISO 8601, got '%s'", altArg)
	}

	// Should have Universal Action
	if item["action"] == nil {
		t.Error("Expected action field to be set")
	}
}

func TestAlfredFormatter_ShouldFormatErrorsAsInvalidItems(t *testing.T) {
	formatter := NewAlfredFormatter()

	output, err := formatter.FormatError("Something went wrong")
	if err != nil {
		t.Fatalf("Expected successful error formatting, got error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	items := result["items"].([]interface{})
	item := items[0].(map[string]interface{})

	if item["title"] != "Error" {
		t.Errorf("Expected title 'Error', got '%v'", item["title"])
	}
	if item["subtitle"] != "Something went wrong" {
		t.Errorf("Expected subtitle 'Something went wrong', got '%v'", item["subtitle"])
	}
	if item["valid"] != false {
		t.Errorf("Expected valid to be false for errors, got %v", item["valid"])
	}

	// Should have system error icon
	icon := item["icon"].(map[string]interface{})
	if !contains(icon["path"].(string), "AlertStopIcon") {
		t.Errorf("Expected error icon, got '%v'", icon["path"])
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