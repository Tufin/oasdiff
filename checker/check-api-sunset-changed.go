package checker

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

func APISunsetChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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

			if op.Deprecated != true {
				continue
			}

			if operationDiff.ExtensionsDiff != nil && operationDiff.ExtensionsDiff.Deleted != nil {
				result = append(result, BackwardCompatibilityError{
					Id:        "sunset-deleted",
					Level:     ERR,
					Text:      config.i18n("sunset-deleted"),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
			}

			if operationDiff.ExtensionsDiff == nil || operationDiff.ExtensionsDiff.Modified == nil {
				continue
			}

			opBase := pathItem.Base.Operations()[operation]

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

			baseDate, err := diff.GetSunsetDate(opBase.ExtensionProps)
			if err != nil {
				result = append(result, BackwardCompatibilityError{
					Id:        "api-deprecated-sunset-parse",
					Level:     ERR,
					Text:      "api sunset date can't be parsed for deprecated API",
					Operation: operation,
					Path:      path,
					Source:    (*operationsSources)[opBase],
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

			if baseDate.After(date) && days < deprecationDays {
				result = append(result, BackwardCompatibilityError{
					Id:        "api-sunset-date-changed-too-small",
					Level:     ERR,
					Text:      fmt.Sprintf(config.i18n("api-sunset-date-changed-too-small"), baseDate, date, baseDate, deprecationDays),
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
