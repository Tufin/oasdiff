package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/checker/localizations"
)

type Localizer func(key string, args ...interface{}) string

func NewLocalizer(locale string, fallbackLocale string) Localizer {
	locales := localizations.New(locale, fallbackLocale)

	return func(key string, args ...interface{}) string {
		pattern := locales.Get("messages." + key)

		return fmt.Sprintf(pattern, args...)
	}
}

type Config struct {
	Checks              []BackwardCompatibilityCheck
	MinSunsetBetaDays   int
	MinSunsetStableDays int
	Localize            Localizer
	LogLevelOverrides   map[string]Level
}

func (c *Config) getLogLevel(checkerId string, defaultLevel Level) Level {
	if level, ok := c.LogLevelOverrides[checkerId]; ok {
		return level
	}
	return defaultLevel
}

func ConditionalError(isConditionSatisfied bool, defaultLevel Level) Level {
	if isConditionSatisfied {
		return ERR
	}

	return defaultLevel
}
