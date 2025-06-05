package geocoder

import (
	"fmt"

	"github.com/codingsince1985/geo-golang"
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"github.com/loginx/alfred-timein/internal/domain"
)

// OpenStreetMapGeocoder implements the Geocoder interface using OpenStreetMap
type OpenStreetMapGeocoder struct {
	geocoder geo.Geocoder
}

// NewOpenStreetMapGeocoder creates a new OpenStreetMapGeocoder
func NewOpenStreetMapGeocoder() *OpenStreetMapGeocoder {
	return &OpenStreetMapGeocoder{
		geocoder: openstreetmap.Geocoder(),
	}
}

// Geocode converts a location query to coordinates
func (g *OpenStreetMapGeocoder) Geocode(query string) (*domain.Location, error) {
	result, err := g.geocoder.Geocode(query)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("no results found for: %s", query)
	}

	return domain.NewLocation(query, float64(result.Lat), float64(result.Lng))
}