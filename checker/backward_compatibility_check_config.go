package checker

import "github.com/tufin/oasdiff/checker/localizations"

type BackwardCompatibilityCheckConfig struct {
	Checks              []BackwardCompatibilityCheck
	MinSunsetBetaDays   int
	MinSunsetStableDays int
	Localizer           localizations.Localizer
	LogLevelOverrides   map[string]Level
}

func (c *BackwardCompatibilityCheckConfig) i18n(messageID string) string {
	return c.Localizer.Get("messages." + messageID)
}

func (c *BackwardCompatibilityCheckConfig) getLogLevel(checkerId string, defaultLevel Level) Level {
	if level, ok := c.LogLevelOverrides[checkerId]; ok {
		return level
	}
	return defaultLevel
}
