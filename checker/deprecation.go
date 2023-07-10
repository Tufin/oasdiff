package checker

import (
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

const (
	BetaDeprecationDays   = 31
	StableDeprecationDays = 180
)

func getSunsetDate(Extensions map[string]interface{}) (string, civil.Date, error) {
	sunset, ok := Extensions[diff.SunsetExtension].(string)
	if !ok {
		sunsetJson, ok := Extensions[diff.SunsetExtension].(json.RawMessage)
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
