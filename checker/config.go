package checker

import "log"

type Config struct {
	Checks              BackwardCompatibilityChecks
	MinSunsetBetaDays   uint
	MinSunsetStableDays uint
	LogLevels           map[string]Level
}

const (
	DefaultBetaDeprecationDays   = uint(0)
	DefaultStableDeprecationDays = uint(0)
)

// NewConfig creates a new configuration with default values.
func NewConfig() *Config {

	rules := GetAllRules()

	return &Config{
		Checks:              rulesToChecks(rules),
		LogLevels:           rulesToLevels(rules),
		MinSunsetBetaDays:   DefaultBetaDeprecationDays,
		MinSunsetStableDays: DefaultStableDeprecationDays,
	}
}

// WithOptionalCheck adds a check to the list of optional checks.
func (config *Config) WithOptionalCheck(id string) *Config {
	return config.WithOptionalChecks([]string{id})
}

// WithOptionalChecks overrides the log level of the given checks to ERR so they will appear in `oasdiff breaking`
func (config *Config) WithOptionalChecks(ids []string) *Config {
	for _, id := range ids {
		config.LogLevels[id] = ERR
	}
	return config
}

// WithDeprecation sets the number of days before sunset for deprecation warnings.
func (config *Config) WithDeprecation(deprecationDaysBeta uint, deprecationDaysStable uint) *Config {
	config.MinSunsetBetaDays = deprecationDaysBeta
	config.MinSunsetStableDays = deprecationDaysStable
	return config
}

// WithSingleCheck sets a single check to be used.
func (config *Config) WithSingleCheck(check BackwardCompatibilityCheck) *Config {
	return config.WithChecks(BackwardCompatibilityChecks{check})
}

// WithChecks sets a list of checks to be used.
func (config *Config) WithChecks(checks BackwardCompatibilityChecks) *Config {
	config.Checks = checks
	return config
}

func (config *Config) getLogLevel(checkId string) Level {
	level, ok := config.LogLevels[checkId]

	if !ok {
		log.Fatal("check id not found: ", checkId)
	}

	return level

}
