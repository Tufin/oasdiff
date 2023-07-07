package checker

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

func APIDeprecationCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

			if operationDiff.DeprecatedDiff == nil {
				continue
			}

			if operationDiff.DeprecatedDiff.To == nil || operationDiff.DeprecatedDiff.To == false {
				// not breaking changes
				result = append(result, ApiChange{
					Id:          "endpoint-reactivated",
					Level:       INFO,
					Text:        config.i18n("endpoint-reactivated"),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}

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

			if days < deprecationDays {
				result = append(result, ApiChange{
					Id:          "api-sunset-date-too-small",
					Level:       ERR,
					Text:        fmt.Sprintf(config.i18n("api-sunset-date-too-small"), date, deprecationDays),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}

			// not breaking changes
			result = append(result, ApiChange{
				Id:          "endpoint-deprecated",
				Level:       INFO,
				Text:        config.i18n("endpoint-deprecated"),
				Operation:   operation,
				OperationId: op.OperationID,
				Path:        path,
				Source:      source,
			})
		}
	}

	return result
}
