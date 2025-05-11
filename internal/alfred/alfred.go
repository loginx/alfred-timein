package alfred

import (
	"encoding/json"
)

// Icon represents the icon object in Alfred JSON.
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path"`
}

// Mod represents a modifier key object in Alfred JSON.
type Mod struct {
	Valid     *bool                  `json:"valid,omitempty"`
	Arg       interface{}            `json:"arg,omitempty"`
	Subtitle  string                 `json:"subtitle,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// Text represents the text object in Alfred JSON.
type Text struct {
	Copy      string `json:"copy,omitempty"`
	LargeType string `json:"largetype,omitempty"`
}

// Item represents a single result item in Alfred JSON.
type Item struct {
	UID          string                 `json:"uid,omitempty"`
	Title        string                 `json:"title"`
	Subtitle     string                 `json:"subtitle,omitempty"`
	Arg          interface{}            `json:"arg,omitempty"`
	Icon         *Icon                  `json:"icon,omitempty"`
	Valid        *bool                  `json:"valid,omitempty"`
	Match        string                 `json:"match,omitempty"`
	Autocomplete string                 `json:"autocomplete,omitempty"`
	Type         string                 `json:"type,omitempty"`
	Mods         map[string]Mod         `json:"mods,omitempty"`
	Action       interface{}            `json:"action,omitempty"`
	Text         *Text                  `json:"text,omitempty"`
	QuickLookURL string                 `json:"quicklookurl,omitempty"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
}

// CacheConfig represents the cache configuration in Alfred JSON.
type CacheConfig struct {
	Seconds     int  `json:"seconds"`
	LooseReload bool `json:"loosereload,omitempty"`
}

// ScriptFilterOutput is the root object for Alfred Script Filter JSON output.
type ScriptFilterOutput struct {
	Variables     map[string]interface{} `json:"variables,omitempty"`
	Items         []Item                 `json:"items"`
	Rerun         float64                `json:"rerun,omitempty"`
	Cache         *CacheConfig           `json:"cache,omitempty"`
	SkipKnowledge bool                   `json:"skipknowledge,omitempty"`
}

// NewScriptFilterOutput creates a new ScriptFilterOutput with sane defaults.
func NewScriptFilterOutput() *ScriptFilterOutput {
	return &ScriptFilterOutput{
		Items: make([]Item, 0),
	}
}

// AddItem appends an item to the ScriptFilterOutput.
func (s *ScriptFilterOutput) AddItem(item Item) {
	s.Items = append(s.Items, item)
}

// ToJSON marshals the ScriptFilterOutput to Alfred-compatible JSON.
func (s *ScriptFilterOutput) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// MustToJSON marshals the ScriptFilterOutput to JSON and panics on error (for CLI use).
func (s *ScriptFilterOutput) MustToJSON() []byte {
	data, err := s.ToJSON()
	if err != nil {
		panic(err)
	}
	return data
} 