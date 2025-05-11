package alfred

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestScriptFilterOutput_JSON(t *testing.T) {
	output := NewScriptFilterOutput()
	output.Variables = map[string]interface{}{"foo": "bar"}
	output.Rerun = 1.5
	output.Cache = &CacheConfig{Seconds: 3600, LooseReload: true}
	output.SkipKnowledge = true

	item := Item{
		UID:      "test-uid",
		Title:    "Test Title",
		Subtitle: "Test Subtitle",
		Arg:      "test-arg",
		Icon:     &Icon{Type: "fileicon", Path: "~/Desktop"},
		Valid:    boolPtr(true),
		Match:    "test match",
		Autocomplete: "Test Auto",
		Type:     "file",
		Mods: map[string]Mod{
			"alt": {
				Valid:    boolPtr(false),
				Arg:      "alt-arg",
				Subtitle: "alt subtitle",
				Variables: map[string]interface{}{"altvar": 1},
			},
		},
		Action: map[string]interface{}{"text": "action text"},
		Text:   &Text{Copy: "copy text", LargeType: "large type"},
		QuickLookURL: "https://example.com/ql",
		Variables: map[string]interface{}{"itemvar": 2},
	}
	output.AddItem(item)

	jsonBytes, err := output.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Spot check a few fields
	if parsed["skipknowledge"] != true {
		t.Errorf("skipknowledge not set correctly")
	}
	if parsed["rerun"] != 1.5 {
		t.Errorf("rerun not set correctly")
	}
	if parsed["variables"].(map[string]interface{})["foo"] != "bar" {
		t.Errorf("variables not set correctly")
	}
	if parsed["cache"].(map[string]interface{})["seconds"] != float64(3600) {
		t.Errorf("cache.seconds not set correctly")
	}
	if parsed["cache"].(map[string]interface{})["loosereload"] != true {
		t.Errorf("cache.loosereload not set correctly")
	}

	items := parsed["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	itemMap := items[0].(map[string]interface{})
	if itemMap["title"] != "Test Title" {
		t.Errorf("item.title not set correctly")
	}
	if itemMap["uid"] != "test-uid" {
		t.Errorf("item.uid not set correctly")
	}
	if itemMap["valid"] != true {
		t.Errorf("item.valid not set correctly")
	}
	if itemMap["type"] != "file" {
		t.Errorf("item.type not set correctly")
	}
	if itemMap["autocomplete"] != "Test Auto" {
		t.Errorf("item.autocomplete not set correctly")
	}
	if itemMap["quicklookurl"] != "https://example.com/ql" {
		t.Errorf("item.quicklookurl not set correctly")
	}
	if !reflect.DeepEqual(itemMap["variables"].(map[string]interface{}), map[string]interface{}{"itemvar": float64(2)}) {
		t.Errorf("item.variables not set correctly")
	}
}

func boolPtr(b bool) *bool {
	return &b
} 