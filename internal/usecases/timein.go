package usecases

import (
	"time"

	"github.com/loginx/alfred-timein/internal/domain"
	"github.com/tkuchiki/go-timezone"
)

// TimeinUseCase handles timezone information retrieval
type TimeinUseCase struct {
	formatter OutputFormatter
}

// NewTimeinUseCase creates a new TimeinUseCase
func NewTimeinUseCase(formatter OutputFormatter) *TimeinUseCase {
	return &TimeinUseCase{
		formatter: formatter,
	}
}

// GetTimezoneInfo gets current time information for a timezone
func (uc *TimeinUseCase) GetTimezoneInfo(timezoneStr string) ([]byte, error) {
	tz, err := domain.NewTimezone(timezoneStr)
	if err != nil {
		output, _ := uc.formatter.FormatError(err.Error())
		return output, err
	}

	return uc.formatter.FormatTimeInfo(tz)
}

// TimezoneInfo represents timezone information for formatting
type TimezoneInfo struct {
	Timezone     *domain.Timezone
	CurrentTime  time.Time
	City         string
	Abbreviation string
}

// GetTimezoneInfoForFormatting gets timezone info structured for formatting
func (uc *TimeinUseCase) GetTimezoneInfoForFormatting(timezoneStr string) (*TimezoneInfo, error) {
	tz, err := domain.NewTimezone(timezoneStr)
	if err != nil {
		return nil, err
	}

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

	return &TimezoneInfo{
		Timezone:     tz,
		CurrentTime:  now,
		City:         city,
		Abbreviation: abbr,
	}, nil
}