package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestUserCanGetCurrentTimeInKnownCity verifies the core user story:
// "As a user, I want to get the current time in any city"
func TestUserCanGetCurrentTimeInKnownCity(t *testing.T) {
	// Given a user wants to know the time in Bangkok
	// When they run the complete workflow
	
	// First get the timezone for Bangkok
	geotzCmd := exec.Command("go", "run", "./cmd/geotz", "Bangkok")
	timezoneOutput, err := geotzCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get timezone for Bangkok: %v", err)
	}
	
	timezone := strings.TrimSpace(string(timezoneOutput))
	
	// Then get the current time for that timezone
	timeinCmd := exec.Command("go", "run", "./cmd/timein", timezone)
	timeOutput, err := timeinCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get time for %s: %v", timezone, err)
	}
	
	timeString := strings.TrimSpace(string(timeOutput))
	
	// The user should get a meaningful result
	if timezone != "Asia/Bangkok" {
		t.Errorf("Expected timezone 'Asia/Bangkok' for Bangkok, got '%s'", timezone)
	}
	
	// The time should be in human-readable format and current
	if !strings.Contains(timeString, "2025") {
		t.Errorf("Expected time to contain current year, got '%s'", timeString)
	}
	
	// Should contain day of week
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	containsWeekday := false
	for _, day := range weekdays {
		if strings.Contains(timeString, day) {
			containsWeekday = true
			break
		}
	}
	if !containsWeekday {
		t.Errorf("Expected time to contain weekday, got '%s'", timeString)
	}
}

// TestUserGetsConsistentResultsFromCache verifies caching behavior:
// "As a user, I expect repeat queries to be fast and consistent"
func TestUserGetsConsistentResultsFromCache(t *testing.T) {
	// Clean cache to start fresh
	os.Remove("geotz_cache.json")
	
	city := "London"
	
	// Ensure binary exists
	if _, err := os.Stat("bin/geotz"); os.IsNotExist(err) {
		t.Skip("Skipping cache test - bin/geotz not found. Run 'make build' first.")
	}
	
	// First query (cache miss)
	start1 := time.Now()
	cmd1 := exec.Command("./bin/geotz", city)
	result1, err := cmd1.Output()
	duration1 := time.Since(start1)
	if err != nil {
		t.Fatalf("First query failed: %v", err)
	}
	
	// Second query (cache hit)
	start2 := time.Now()
	cmd2 := exec.Command("./bin/geotz", city)
	result2, err := cmd2.Output()
	duration2 := time.Since(start2)
	if err != nil {
		t.Fatalf("Second query failed: %v", err)
	}
	
	// Results should be identical
	if string(result1) != string(result2) {
		t.Errorf("Cache should return identical results. First: '%s', Second: '%s'", 
			string(result1), string(result2))
	}
	
	// Second query should be faster or comparable (allowing for go run overhead)
	// In CI environment, go run overhead can mask cache benefits
	if duration2 > duration1*2 {
		t.Logf("First query: %v, Second query: %v", duration1, duration2)
		t.Error("Cache hit should not be significantly slower than cache miss")
	} else {
		t.Logf("Cache performance: First query: %v, Second query: %v", duration1, duration2)
	}
	
	// Result should be valid timezone
	timezone := strings.TrimSpace(string(result1))
	if timezone != "Europe/London" {
		t.Errorf("Expected 'Europe/London' for London, got '%s'", timezone)
	}
}

// TestUserGetsHelpfulErrorForInvalidCity verifies error handling:
// "As a user, I should get clear error messages for invalid input"
func TestUserGetsHelpfulErrorForInvalidCity(t *testing.T) {
	// Given a user tries to look up a nonsense city
	cmd := exec.Command("go", "run", "./cmd/geotz", "NotARealCity12345XYZ")
	output, err := cmd.CombinedOutput()
	
	// The command should exit with error code
	if err == nil {
		t.Fatal("Expected error exit code for invalid city")
	}
	
	// And provide helpful error message
	errorMsg := string(output)
	if !strings.Contains(strings.ToLower(errorMsg), "could not geocode") {
		t.Errorf("Expected helpful error message, got: '%s'", errorMsg)
	}
}

// TestAlfredWorkflowProvidesExpectedJSONStructure verifies Alfred integration:
// "As an Alfred user, I should get properly formatted JSON results"
func TestAlfredWorkflowProvidesExpectedJSONStructure(t *testing.T) {
	// Given a user queries through Alfred format
	cmd := exec.Command("go", "run", "./cmd/geotz", "--format=alfred", "Tokyo")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Alfred format query failed: %v", err)
	}
	
	// Parse the JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Alfred output should be valid JSON: %v", err)
	}
	
	// Verify required Alfred structure
	items, exists := result["items"]
	if !exists {
		t.Fatal("Alfred response must contain 'items' array")
	}
	
	itemArray := items.([]interface{})
	if len(itemArray) != 1 {
		t.Fatalf("Expected 1 result item, got %d", len(itemArray))
	}
	
	item := itemArray[0].(map[string]interface{})
	
	// Verify required item fields
	if _, exists := item["title"]; !exists {
		t.Error("Alfred item must have 'title' field")
	}
	if _, exists := item["subtitle"]; !exists {
		t.Error("Alfred item must have 'subtitle' field")
	}
	if _, exists := item["arg"]; !exists {
		t.Error("Alfred item must have 'arg' field")
	}
	
	// Title should be the timezone
	title := item["title"].(string)
	if title != "Asia/Tokyo" {
		t.Errorf("Expected title 'Asia/Tokyo' for Tokyo, got '%s'", title)
	}
	
	// Subtitle should mention the city
	subtitle := item["subtitle"].(string)
	if !strings.Contains(strings.ToLower(subtitle), "tokyo") {
		t.Errorf("Subtitle should mention Tokyo, got '%s'", subtitle)
	}
}

// TestPipelineWorkflowWorksEndToEnd verifies the common piped usage:
// "As a user, I should be able to pipe geotz output to timein"
func TestPipelineWorkflowWorksEndToEnd(t *testing.T) {
	// Clear cache for predictable timing
	os.Remove("geotz_cache.json")
	
	// Given a user wants to pipe geotz to timein, use shell to handle pipe
	cmd := exec.Command("sh", "-c", "go run ./cmd/geotz Berlin | go run ./cmd/timein")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Pipeline failed: %v", err)
	}
	
	result := strings.TrimSpace(string(output))
	
	// Should get current time in Berlin timezone
	if !strings.Contains(result, "2025") {
		t.Errorf("Expected current year in result, got '%s'", result)
	}
	
	// Should be human readable
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	containsWeekday := false
	for _, day := range weekdays {
		if strings.Contains(result, day) {
			containsWeekday = true
			break
		}
	}
	if !containsWeekday {
		t.Errorf("Expected weekday in time output, got '%s'", result)
	}
}