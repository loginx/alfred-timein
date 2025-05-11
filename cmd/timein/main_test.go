package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

func TestTimein_Argv_Alfred(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "--format=alfred", "America/New_York")
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
		t.Fatalf("expected valid time, got error: %v", item["subtitle"])
	}
}

func TestTimein_Argv_Plain(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "--format=plain", "America/New_York")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(string(out))
	if result == "" || strings.Contains(result, "Error") || strings.Contains(result, "{") {
		t.Errorf("expected plain human time, got: %v", result)
	}
}

func TestTimein_Stdin_Alfred(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "--format=alfred")
	cmd.Stdin = strings.NewReader("Europe/London\n")
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
		t.Fatalf("expected valid time, got error: %v", item["subtitle"])
	}
}

func TestTimein_Stdin_Plain(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "--format=plain")
	cmd.Stdin = strings.NewReader("Europe/London\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(string(out))
	if result == "" || strings.Contains(result, "Error") || strings.Contains(result, "{") {
		t.Errorf("expected plain human time, got: %v", result)
	}
}

func TestTimein_Error_Alfred(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "--format=alfred")
	cmd.Stdin = strings.NewReader("\n")
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
}

func TestTimein_Error_Plain(t *testing.T) {
	cmd := exec.Command("go", "run", "./main.go", "--format=plain")
	cmd.Stdin = strings.NewReader("\n")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected error exit code for invalid input")
	}
	result := string(out)
	if !strings.Contains(result, "IANA timezone argument required") {
		t.Errorf("expected error message in stderr, got: %v", result)
	}
}
