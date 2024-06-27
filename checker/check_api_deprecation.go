package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	EndpointReactivatedId        = "endpoint-reactivated"
	APIDeprecatedSunsetParseId   = "api-deprecated-sunset-parse"
	APIDeprecatedSunsetMissingId = "api-deprecated-sunset-missing"
	APIInvalidStabilityLevelId   = "api-invalid-stability-level"
	APISunsetDateTooSmallId      = "api-sunset-date-too-small"
	EndpointDeprecatedId         = "endpoint-deprecated"
)

func APIDeprecationCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationDiff := range pathItem.OperationsDiff.Modified {
			op := pathItem.Revision.GetOperation(operation)
			source := (*operationsSources)[op]

			if operationDiff.DeprecatedDiff == nil {
				continue
			}

			if operationDiff.DeprecatedDiff.To == nil || operationDiff.DeprecatedDiff.To == false {
				// not breaking changes
				result = append(result, ApiChange{
					Id:          EndpointReactivatedId,
					Level:       INFO,
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			stability, err := getStabilityLevel(op.Extensions)
			if err != nil {
				// handled in CheckBackwardCompatibility
				continue
			}

			deprecationDays := getDeprecationDays(config, stability)

			sunset, ok := getSunset(op.Extensions)
			if !ok {
				// if deprecation policy is defined and sunset is missing, it's a breaking change
				if deprecationDays > 0 {
					result = append(result, getAPIDeprecatedSunsetMissing(op, operationsSources, path, operation))
				}
				continue
			}

			date, err := getSunsetDate(sunset)
			if err != nil {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetParseId,
					Level:       ERR,
					Args:        []any{err},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			days := date.DaysSince(civil.DateOf(time.Now()))

			if days < int(deprecationDays) {
				result = append(result, ApiChange{
					Id:          APISunsetDateTooSmallId,
					Level:       ERR,
					Args:        []any{date, deprecationDays},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			// not breaking changes
			result = append(result, ApiChange{
				Id:          EndpointDeprecatedId,
				Level:       INFO,
				Operation:   operation,
				OperationId: op.OperationID,
				Path:        path,
				Source:      load.NewSource(source),
			})
		}
	}

	return result
}
