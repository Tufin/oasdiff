package checker

type Config struct {
	Checks              []BackwardCompatibilityCheck
	MinSunsetBetaDays   int
	MinSunsetStableDays int
	LogLevelOverrides   map[string]Level
	Tags                []string
}

func (c *Config) getLogLevel(checkerId string, defaultLevel Level) Level {
	if level, ok := c.LogLevelOverrides[checkerId]; ok {
		return level
	}
	return defaultLevel
}
