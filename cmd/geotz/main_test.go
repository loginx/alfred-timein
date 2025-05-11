package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGeotz_ValidCity_Alfred(t *testing.T) {
	os.RemoveAll(".cache")
	cmd := exec.Command("go", "run", "./main.go", "--format=alfred", "Berlin")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	items := parsed["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	item := items[0].(map[string]interface{})
	if item["title"] == "Error" {
		t.Fatalf("expected valid timezone, got error: %v", item["subtitle"])
	}
	if item["arg"] == "" {
		t.Errorf("arg should not be empty")
	}
	if !strings.Contains(strings.ToLower(item["subtitle"].(string)), "berlin") {
		t.Errorf("subtitle should contain city name")
	}
	if strings.Contains(item["subtitle"].(string), "cached") {
		t.Errorf("first lookup should not be cached")
	}
}

func TestGeotz_ValidCity_Plain(t *testing.T) {
	os.RemoveAll(".cache")
	cmd := exec.Command("go", "run", "./main.go", "--format=plain", "Paris")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(string(out))
	if result == "" || strings.Contains(result, "Error") || strings.Contains(result, "{") {
		t.Errorf("expected plain IANA timezone, got: %v", result)
	}
}

func TestGeotz_CacheHit_Alfred(t *testing.T) {
	os.RemoveAll(".cache")
	city := "Paris"
	cmd := exec.Command("go", "run", "./main.go", "--format=alfred", city)
	_, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmd2 := exec.Command("go", "run", "./main.go", "--format=alfred", city)
	out2, err := cmd2.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(out2, &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	item := parsed["items"].([]interface{})[0].(map[string]interface{})
	if !strings.Contains(item["subtitle"].(string), "cached") {
		t.Errorf("expected cache hit, subtitle: %v", item["subtitle"])
	}
}

func TestGeotz_InvalidCity_Alfred(t *testing.T) {
	os.RemoveAll(".cache")
	cmd := exec.Command("go", "run", "./main.go", "--format=alfred", "NotARealCity123456")
	out, _ := cmd.Output()
	var parsed map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(out), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	items := parsed["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	item := items[0].(map[string]interface{})
	if item["title"] != "Error" {
		t.Errorf("expected error title, got %v", item["title"])
	}
	if item["valid"] != false {
		t.Errorf("expected valid=false for error item")
	}
}

func TestGeotz_InvalidCity_Plain(t *testing.T) {
	os.RemoveAll(".cache")
	cmd := exec.Command("go", "run", "./main.go", "--format=plain", "NotARealCity123456")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error exit code for invalid city")
	}
	result := string(out)
	if !strings.Contains(result, "Could not geocode") {
		t.Errorf("expected error message in stderr, got: %v", result)
	}
}
