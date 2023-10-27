package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/utils"
)

func GetDefaultChecks() Config {
	return GetChecks(utils.StringList{})
}

func GetChecks(includeChecks utils.StringList) Config {
	return getBackwardCompatibilityCheckConfig(allChecks(), LevelOverrides(includeChecks), BetaDeprecationDays, StableDeprecationDays)
}

func LevelOverrides(includeChecks utils.StringList) map[string]Level {
	result := map[string]Level{}
	for _, s := range includeChecks {
		// if the checker was explicitly included with the `--include-checks`,
		// it means that it's output is considered a breaking change,
		// so the returned level should overwritten to ERR.
		result[s] = ERR
	}
	return result
}

func GetAllChecks(includeChecks utils.StringList, deprecationDaysBeta int, deprecationDaysStable int) Config {
	return getBackwardCompatibilityCheckConfig(allChecks(), LevelOverrides(includeChecks), deprecationDaysBeta, deprecationDaysStable)
}

func getBackwardCompatibilityCheckConfig(checks []BackwardCompatibilityCheck, levelOverrides map[string]Level, minSunsetBetaDays int, minSunsetStableDays int) Config {
	return Config{
		Checks:              checks,
		LogLevelOverrides:   levelOverrides,
		MinSunsetBetaDays:   minSunsetBetaDays,
		MinSunsetStableDays: minSunsetStableDays,
		Localize:            NewLocalizer("en", "en"),
	}
}

func optionalChecks() map[string]BackwardCompatibilityCheck {
	optionalRules := GetOptionalRules()

	result := make(map[string]BackwardCompatibilityCheck, len(optionalRules))
	for _, rule := range optionalRules {
		result[rule.Id] = rule.Handler
	}
	return result
}

func GetOptionalChecks() []string {
	optionalChecks := optionalChecks()

	result := make([]string, len(optionalChecks))
	i := 0
	for key := range optionalChecks {
		result[i] = key
		i++
	}

	return result
}

func defaultChecks() []BackwardCompatibilityCheck {

	result := []BackwardCompatibilityCheck{}
	m := utils.StringSet{}
	for _, rule := range GetRequiredRules() {
		pStr := fmt.Sprintf("%v", rule.Handler)
		if !m.Contains(pStr) {
			m.Add(pStr)
			result = append(result, rule.Handler)
		}
	}
	return result
}

func allChecks() []BackwardCompatibilityCheck {
	defaultChecks := defaultChecks()
	optionalChecks := optionalChecks()

	for _, v := range optionalChecks {
		defaultChecks = append(defaultChecks, v)
	}
	return defaultChecks
}
