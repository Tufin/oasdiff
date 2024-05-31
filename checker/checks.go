package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

// BackwardCompatibilityCheck, or a check, is a function that receives a diff report and returns a list of changes
type BackwardCompatibilityCheck func(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes

func defaultChecks() []BackwardCompatibilityCheck {

	result := []BackwardCompatibilityCheck{}
	m := utils.StringSet{}
	for _, rule := range GetRequiredRules() {
		// functions are not comparable, so we convert them to strings
		pStr := fmt.Sprintf("%v", rule.Handler)
		if !m.Contains(pStr) {
			m.Add(pStr)
			result = append(result, rule.Handler)
		}
	}
	return result
}

func allChecks() []BackwardCompatibilityCheck {
	optionalRules := GetOptionalRules()
	optionalChecks := make([]BackwardCompatibilityCheck, len(optionalRules))
	for i, rule := range optionalRules {
		optionalChecks[i] = rule.Handler
	}

	return append(defaultChecks(), optionalChecks...)
}

func levelOverrides(includeChecks []string) map[string]Level {
	result := map[string]Level{}
	for _, s := range includeChecks {
		// if the checker was explicitly included with the `--include-checks`,
		// it means that it's output is considered a breaking change,
		// so the returned level should be overwritten to ERR.
		result[s] = ERR
	}
	return result
}
