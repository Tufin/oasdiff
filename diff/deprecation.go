package diff

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
)

const SunsetExtension = "x-sunset"

func getSunsetDate(Extensions map[string]interface{}) (civil.Date, error) {
	sunsetJson, ok := Extensions[SunsetExtension].(json.RawMessage)
	if !ok {
		return civil.Date{}, errors.New("not found")
	}

	var sunset string
	if err := json.Unmarshal(sunsetJson, &sunset); err != nil {
		return civil.Date{}, errors.New("unmarshal failed")
	}

	date, err := civil.ParseDate(sunset)
	if err != nil {
		return civil.Date{}, errors.New("failed to parse time")
	}

	return date, nil
}

// sunsetAllowed checks if an element can be deleted after deprecation period
func sunsetAllowed(deprecated bool, Extensions map[string]interface{}) bool {

	if !deprecated {
		return false
	}

	date, err := getSunsetDate(Extensions)
	if err != nil {
		return false
	}

	return civil.DateOf(time.Now()).After(date)
}

func deprecationPeriodSufficient(deprecationDays int, Extensions map[string]interface{}) bool {
	if deprecationDays == 0 {
		return true
	}

	date, err := getSunsetDate(Extensions)
	if err != nil {
		return false
	}

	days := date.DaysSince(civil.DateOf(time.Now()))

	return days >= deprecationDays
}
