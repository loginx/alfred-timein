package presenter

import (
	"fmt"
	"time"

	"github.com/loginx/alfred-timein/internal/alfred"
	"github.com/loginx/alfred-timein/internal/domain"
	"github.com/tkuchiki/go-timezone"
)

const alfredCacheSeconds = 604800 // 7 days

// AlfredFormatter formats output for Alfred Script Filter
type AlfredFormatter struct{}

// NewAlfredFormatter creates a new AlfredFormatter
func NewAlfredFormatter() *AlfredFormatter {
	return &AlfredFormatter{}
}

// FormatTimezoneInfo formats timezone information for Alfred
func (f *AlfredFormatter) FormatTimezoneInfo(timezone *domain.Timezone, city string, cached bool) ([]byte, error) {
	out := alfred.NewScriptFilterOutput()
	out.Cache = &alfred.CacheConfig{Seconds: alfredCacheSeconds}

	subtitle := city
	if cached {
		subtitle += " (cached)"
	}

	item := alfred.Item{
		Title:    timezone.String(),
		Subtitle: subtitle,
		Arg:      timezone.String(),
		Variables: map[string]interface{}{
			"city": city,
		},
	}

	out.AddItem(item)
	return out.ToJSON()
}

// FormatTimeInfo formats current time information for Alfred
func (f *AlfredFormatter) FormatTimeInfo(tz *domain.Timezone) ([]byte, error) {
	loc, err := tz.Location()
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	city := tz.City()

	// Get timezone abbreviation
	tzlib := timezone.New()
	isDST := tzlib.IsDST(now)
	abbr, err := tzlib.GetTimezoneAbbreviation(tz.String(), isDST)
	if err != nil || abbr == "" {
		abbr = now.Format("MST")
	}

	title := fmt.Sprintf("%s - %s", tz.String(), now.Format("Mon, Jan 2, 3:04 PM"))
	subtitle := fmt.Sprintf("Current time in %s (%s)", city, abbr)

	out := alfred.NewScriptFilterOutput()
	out.Cache = &alfred.CacheConfig{Seconds: 60}

	item := alfred.Item{
		Title:    title,
		Subtitle: subtitle,
		Arg:      title,
		Variables: map[string]interface{}{
			"timezone": tz.String(),
		},
	}

	out.AddItem(item)
	return out.ToJSON()
}

// FormatError formats error messages for Alfred
func (f *AlfredFormatter) FormatError(message string) ([]byte, error) {
	out := alfred.NewScriptFilterOutput()
	item := alfred.Item{
		Title:    "Error",
		Subtitle: message,
		Valid:    boolPtr(false),
	}
	out.AddItem(item)
	return out.ToJSON()
}

func boolPtr(b bool) *bool {
	return &b
}