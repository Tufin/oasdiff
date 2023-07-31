package checker

import "github.com/tufin/oasdiff/checker/localizations"

type Config struct {
	Checks              []BackwardCompatibilityCheck
	MinSunsetBetaDays   int
	MinSunsetStableDays int
	Localizer           localizations.Localizer
	LogLevelOverrides   map[string]Level
}

func (c *Config) i18n(messageID string) string {
	return c.Localizer.Get("messages." + messageID)
}

func (c *Config) getLogLevel(checkerId string, defaultLevel Level) Level {
	if level, ok := c.LogLevelOverrides[checkerId]; ok {
		return level
	}
	return defaultLevel
}

func (c *Config) conditionalError(isConditionSatisfied bool) Level {
	if isConditionSatisfied {
		return ERR
	}

	return INFO
}
