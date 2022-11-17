package checker

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

func APIDeprecationCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationDiff := range pathItem.OperationsDiff.Modified {
			op := pathItem.Revision.Operations()[operation]
			source := (*operationsSources)[op]

			if operationDiff.DeprecatedDiff == nil {
				continue
			}

			if operationDiff.DeprecatedDiff.To == nil || operationDiff.DeprecatedDiff.To == false {
				continue
			}

			date, err := diff.GetSunsetDate(op.ExtensionProps)
			if err != nil {
				result = append(result, BackwardCompatibilityError{
					Id:        "api-deprecated-sunset-parse",
					Level:     ERR,
					Text:      "api sunset date can't be parsed for deprecated API",
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}

			days := date.DaysSince(civil.DateOf(time.Now()))
			deprecationDays := config.MinSunsetStableDays

			stability, err := getStabilityLevel(op.ExtensionProps)
			if err != nil {
				result = append(result, BackwardCompatibilityError{
					Id:        "parsing-error",
					Level:     ERR,
					Text:      fmt.Sprintf("parsing error %s", err.Error()),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			if stability == "beta" {
				deprecationDays = config.MinSunsetBetaDays
			}

			if days < deprecationDays {
				result = append(result, BackwardCompatibilityError{
					Id:        "api-sunset-date-too-small",
					Level:     ERR,
					Text:      fmt.Sprintf("api sunset date %s is too small, must be at least %d days from now", date, deprecationDays),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
			}
		}
	}

	return result
}
