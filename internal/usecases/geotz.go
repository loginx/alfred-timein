package usecases

import (
	"fmt"
	"strings"

	"github.com/loginx/alfred-timein/internal/domain"
)

// GeotzUseCase handles geocoding to timezone conversion
type GeotzUseCase struct {
	geocoder        Geocoder
	timezoneFinder  TimezoneFinder
	cache          Cache
	formatter      OutputFormatter
}

// NewGeotzUseCase creates a new GeotzUseCase
func NewGeotzUseCase(geocoder Geocoder, timezoneFinder TimezoneFinder, cache Cache, formatter OutputFormatter) *GeotzUseCase {
	return &GeotzUseCase{
		geocoder:       geocoder,
		timezoneFinder: timezoneFinder,
		cache:         cache,
		formatter:     formatter,
	}
}

// GetTimezoneFromCity converts a city name to timezone
func (uc *GeotzUseCase) GetTimezoneFromCity(city string) ([]byte, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		output, _ := uc.formatter.FormatError("City or landmark argument required.")
		return output, fmt.Errorf("city or landmark argument required")
	}

	// Check cache first
	cacheKey := strings.ToLower(city)
	if tz, ok := uc.cache.Get(cacheKey); ok {
		timezone, err := domain.NewTimezone(tz)
		if err != nil {
			output, _ := uc.formatter.FormatError(err.Error())
			return output, err
		}
		return uc.formatter.FormatTimezoneInfo(timezone, city, true)
	}

	// Geocode the city
	location, err := uc.geocoder.Geocode(city)
	if err != nil {
		output, _ := uc.formatter.FormatError("Could not geocode: " + city)
		return output, fmt.Errorf("could not geocode: %s", city)
	}

	// Find timezone for the location
	tz, err := uc.timezoneFinder.GetTimezoneName(location.Longitude, location.Latitude)
	if err != nil || tz == "" {
		output, _ := uc.formatter.FormatError("Could not resolve timezone for: " + city)
		return output, fmt.Errorf("could not resolve timezone for: %s", city)
	}

	timezone, err := domain.NewTimezone(tz)
	if err != nil {
		output, _ := uc.formatter.FormatError(err.Error())
		return output, err
	}

	// Cache the result
	uc.cache.Set(cacheKey, tz)

	return uc.formatter.FormatTimezoneInfo(timezone, city, false)
}