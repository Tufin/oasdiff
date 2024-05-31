package checker

import "github.com/tufin/oasdiff/utils"

type Config struct {
	Checks              []BackwardCompatibilityCheck
	MinSunsetBetaDays   int
	MinSunsetStableDays int
	LogLevelOverrides   map[string]Level
}

const (
	BetaDeprecationDays   = 31
	StableDeprecationDays = 180
)

func NewConfig() *Config {
	return &Config{
		Checks:              allChecks(),
		LogLevelOverrides:   map[string]Level{},
		MinSunsetBetaDays:   BetaDeprecationDays,
		MinSunsetStableDays: StableDeprecationDays,
	}
}

func (config *Config) WithOptionalChecks(includeChecks utils.StringList) *Config {
	config.LogLevelOverrides = levelOverrides(includeChecks)
	return config
}

func (config *Config) WithDeprecation(deprecationDaysBeta int, deprecationDaysStable int) *Config {
	config.MinSunsetBetaDays = deprecationDaysBeta
	config.MinSunsetStableDays = deprecationDaysStable
	return config
}

func (c *Config) getLogLevel(checkerId string, defaultLevel Level) Level {
	if level, ok := c.LogLevelOverrides[checkerId]; ok {
		return level
	}
	return defaultLevel
}
