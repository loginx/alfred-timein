package usecases

import (
	"github.com/loginx/alfred-timein/internal/domain"
)

// Geocoder defines the interface for geocoding services
type Geocoder interface {
	Geocode(query string) (*domain.Location, error)
}

// TimezoneFinder defines the interface for timezone lookup services
type TimezoneFinder interface {
	GetTimezoneName(longitude, latitude float64) (string, error)
}

// Cache defines the interface for caching services
type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Clear()
}

// OutputFormatter defines the interface for output formatting
type OutputFormatter interface {
	FormatTimezoneInfo(timezone *domain.Timezone, city string, cached bool) ([]byte, error)
	FormatTimeInfo(timezone *domain.Timezone) ([]byte, error)
	FormatError(message string) ([]byte, error)
}