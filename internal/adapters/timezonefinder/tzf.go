package timezonefinder

import (
	"fmt"

	"github.com/ringsaturn/tzf"
)

// TzfTimezoneFinder implements the TimezoneFinder interface using tzf
type TzfTimezoneFinder struct {
	finder tzf.F
}

// NewTzfTimezoneFinder creates a new TzfTimezoneFinder
func NewTzfTimezoneFinder() (*TzfTimezoneFinder, error) {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize timezone finder: %w", err)
	}

	return &TzfTimezoneFinder{
		finder: finder,
	}, nil
}

// GetTimezoneName returns the timezone name for given coordinates
func (tf *TzfTimezoneFinder) GetTimezoneName(longitude, latitude float64) (string, error) {
	tz := tf.finder.GetTimezoneName(longitude, latitude)
	if tz == "" {
		return "", fmt.Errorf("no timezone found for coordinates: %f, %f", latitude, longitude)
	}
	return tz, nil
}