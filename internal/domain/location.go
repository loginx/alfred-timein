package domain

import "fmt"

// Location represents a geographic location
type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

// NewLocation creates a new Location
func NewLocation(name string, lat, lng float64) (*Location, error) {
	if name == "" {
		return nil, fmt.Errorf("location name cannot be empty")
	}
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("invalid latitude: %f", lat)
	}
	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("invalid longitude: %f", lng)
	}

	return &Location{
		Name:      name,
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

// String returns the location name
func (l *Location) String() string {
	return l.Name
}