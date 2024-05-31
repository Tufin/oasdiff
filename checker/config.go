package checker

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

func (config *Config) WithOptionalCheck(id string) *Config {
	config.LogLevelOverrides = levelOverrides([]string{id})
	return config
}

func (config *Config) WithOptionalChecks(ids []string) *Config {
	config.LogLevelOverrides = levelOverrides(ids)
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
