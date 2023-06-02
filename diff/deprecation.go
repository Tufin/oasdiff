package diff

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
)

const SunsetExtension = "x-sunset"

func GetSunsetDate(Extensions map[string]interface{}) (string, civil.Date, error) {
	sunset, ok := Extensions[SunsetExtension].(string)
	if !ok {
		sunsetJson, ok := Extensions[SunsetExtension].(json.RawMessage)
		if !ok {
			return "", civil.Date{}, errors.New("sunset header not found")
		}
		if err := json.Unmarshal(sunsetJson, &sunset); err != nil {
			return "", civil.Date{}, errors.New("unmarshal failed")
		}
	}

	if date, err := civil.ParseDate(sunset); err == nil {
		return sunset, date, nil
	} else if date, err := time.Parse(time.RFC3339, sunset); err == nil {
		return sunset, civil.DateOf(date), nil
	}

	return sunset, civil.Date{}, errors.New("failed to parse sunset date")
}

// SunsetAllowed checks if an element can be deleted after deprecation period
func SunsetAllowed(deprecated bool, Extensions map[string]interface{}) bool {

	if !deprecated {
		return false
	}

	_, date, err := GetSunsetDate(Extensions)
	if err != nil {
		return false
	}

	return civil.DateOf(time.Now()).After(date)
}

func DeprecationPeriodSufficient(deprecationDays int, Extensions map[string]interface{}) bool {
	if deprecationDays == 0 {
		return true
	}

	_, date, err := GetSunsetDate(Extensions)
	if err != nil {
		return false
	}

	days := date.DaysSince(civil.DateOf(time.Now()))

	return days >= deprecationDays
}
