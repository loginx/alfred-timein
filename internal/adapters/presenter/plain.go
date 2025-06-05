package presenter

import (
	"fmt"
	"time"

	"github.com/loginx/alfred-timein/internal/domain"
)

// PlainFormatter formats output as plain text
type PlainFormatter struct{}

// NewPlainFormatter creates a new PlainFormatter
func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

// FormatTimezoneInfo formats timezone information as plain text
func (f *PlainFormatter) FormatTimezoneInfo(timezone *domain.Timezone, city string, cached bool) ([]byte, error) {
	return []byte(timezone.String() + "\n"), nil
}

// FormatTimeInfo formats current time information as plain text
func (f *PlainFormatter) FormatTimeInfo(tz *domain.Timezone) ([]byte, error) {
	loc, err := tz.Location()
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	// Human-friendly, locale-aware output
	humanTime := now.Format("Monday, 02 January 2006, 3:04:05 PM")
	return []byte(humanTime + "\n"), nil
}

// FormatError formats error messages as plain text
func (f *PlainFormatter) FormatError(message string) ([]byte, error) {
	return []byte(fmt.Sprintf("Error: %s\n", message)), nil
}