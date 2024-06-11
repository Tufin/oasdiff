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

// NewConfig creates a new configuration with default values.
func NewConfig() *Config {
	return &Config{
		Checks:              allChecks(),
		LogLevelOverrides:   map[string]Level{},
		MinSunsetBetaDays:   BetaDeprecationDays,
		MinSunsetStableDays: StableDeprecationDays,
	}
}

// WithOptionalCheck adds a check to the list of optional checks.
func (config *Config) WithOptionalCheck(id string) *Config {
	return config.WithOptionalChecks([]string{id})
}

// WithOptionalChecks adds a list of checks to the list of optional checks.
func (config *Config) WithOptionalChecks(ids []string) *Config {
	config.LogLevelOverrides = levelOverrides(ids)
	return config
}

// WithDeprecation sets the number of days before sunset for deprecation warnings.
func (config *Config) WithDeprecation(deprecationDaysBeta int, deprecationDaysStable int) *Config {
	config.MinSunsetBetaDays = deprecationDaysBeta
	config.MinSunsetStableDays = deprecationDaysStable
	return config
}

// WithSingleCheck sets a single check to be used.
func (config *Config) WithSingleCheck(check BackwardCompatibilityCheck) *Config {
	return config.WithChecks([]BackwardCompatibilityCheck{check})
}

// WithChecks sets a list of checks to be used.
func (config *Config) WithChecks(checks []BackwardCompatibilityCheck) *Config {
	config.Checks = checks
	return config
}

func (c *Config) getLogLevel(checkerId string, defaultLevel Level) Level {
	if level, ok := c.LogLevelOverrides[checkerId]; ok {
		return level
	}
	return defaultLevel
}
