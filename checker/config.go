package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/checker/localizations"
)

type Localizer func(key string, args ...interface{}) string

func NewDefaultLocalizer() Localizer {
	return NewLocalizer(localizations.LangDefault)
}

func NewLocalizer(locale string) Localizer {
	locales := localizations.New(locale, localizations.LangDefault)

	return func(originalKey string, args ...interface{}) string {
		key := "messages." + originalKey
		pattern := locales.Get(key)

		// if key not found, return original key
		// TODO: improve localizations to return error when key not found
		if pattern == key {
			return originalKey
		}

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
