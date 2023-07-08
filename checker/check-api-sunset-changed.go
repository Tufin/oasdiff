package checker

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

func APISunsetChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
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

			if !op.Deprecated {
				continue
			}

			if operationDiff.ExtensionsDiff != nil && !operationDiff.ExtensionsDiff.Deleted.Empty() {
				result = append(result, ApiChange{
					Id:          "sunset-deleted",
					Level:       ERR,
					Text:        config.i18n("sunset-deleted"),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
			}

			if operationDiff.ExtensionsDiff == nil || operationDiff.ExtensionsDiff.Modified.Empty() {
				continue
			}

			opBase := pathItem.Base.Operations()[operation]

			rawDate, date, err := diff.GetSunsetDate(op.Extensions)
			if err != nil {
				result = append(result, ApiChange{
					Id:          "api-deprecated-sunset-parse",
					Level:       ERR,
					Text:        fmt.Sprintf(config.i18n("api-deprecated-sunset-parse"), rawDate, err),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}

			rawDate, baseDate, err := diff.GetSunsetDate(opBase.Extensions)
			if err != nil {
				result = append(result, ApiChange{
					Id:          "api-deprecated-sunset-parse",
					Level:       ERR,
					Text:        fmt.Sprintf(config.i18n("api-deprecated-sunset-parse"), rawDate, err),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      (*operationsSources)[opBase],
				})
				continue
			}

			days := date.DaysSince(civil.DateOf(time.Now()))

			stability, err := getStabilityLevel(op.Extensions)
			if err != nil {
				result = append(result, ApiChange{
					Id:          "parsing-error",
					Level:       ERR,
					Text:        fmt.Sprintf("parsing error %s", err.Error()),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}

			deprecationDays := getDeperacationDays(config, stability)

			if baseDate.After(date) && days < deprecationDays {
				result = append(result, ApiChange{
					Id:          "api-sunset-date-changed-too-small",
					Level:       ERR,
					Text:        fmt.Sprintf(config.i18n("api-sunset-date-changed-too-small"), baseDate, date, baseDate, deprecationDays),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
			}
		}
	}

	return result
}

func getDeperacationDays(config Config, stability string) int {
	if stability == "beta" {
		return config.MinSunsetBetaDays
	}

	return config.MinSunsetStableDays
}
