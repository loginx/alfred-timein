package domain

import (
	"fmt"
	"strings"
	"time"
)

// Timezone represents a validated IANA timezone
type Timezone struct {
	Name string
}

// NewTimezone creates a new Timezone after validation
func NewTimezone(name string) (*Timezone, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("timezone name cannot be empty")
	}

	// Validate that the timezone is loadable
	_, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %s", name)
	}

	return &Timezone{Name: name}, nil
}

// Location returns the Go time.Location for this timezone
func (tz *Timezone) Location() (*time.Location, error) {
	return time.LoadLocation(tz.Name)
}

// String returns the timezone name
func (tz *Timezone) String() string {
	return tz.Name
}

// CityFromTimezone extracts the city/region from an IANA timezone string
func (tz *Timezone) City() string {
	parts := strings.Split(tz.Name, "/")
	if len(parts) > 1 {
		return strings.ReplaceAll(parts[1], "_", " ")
	}
	return tz.Name
}