package diff

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

func getSunsetDate(extensionProps openapi3.ExtensionProps) (time.Time, error) {
	sunsetJson, ok := extensionProps.Extensions["x-sunset"].(json.RawMessage)
	if !ok {
		return time.Time{}, errors.New("not found")
	}

	var sunset string
	if err := json.Unmarshal(sunsetJson, &sunset); err != nil {
		return time.Time{}, errors.New("unmarshal failed")
	}

	date, err := time.Parse("2006-01-02", sunset)
	if err != nil {
		return time.Time{}, errors.New("failed to parse time")
	}

	return date, nil
}

// sunsetAllowed checks if an element can be deleted after deprecation period
func sunsetAllowed(deprecated bool, extensionProps openapi3.ExtensionProps) bool {

	if !deprecated {
		return false
	}

	date, err := getSunsetDate(extensionProps)
	if err != nil {
		return false
	}

	return time.Now().After(date)
}

func deprecationPeriodSufficient(deprecationDays int, extensionProps openapi3.ExtensionProps) bool {
	if deprecationDays == 0 {
		return true
	}

	date, err := getSunsetDate(extensionProps)
	if err != nil {
		return false
	}

	days := int(date.Sub(time.Now()).Hours() / 24)

	return days >= deprecationDays
}
