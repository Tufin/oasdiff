package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

const (
	ParameterReactivatedId             = "parameter-reactivated"
	ParameterDeprecatedSunsetParseId   = "parameter-deprecated-sunset-parse"
	ParameterDeprecatedSunsetMissingId = "parameter-deprecated-sunset-missing"
	ParameterInvalidStabilityLevelId   = "parameter-invalid-stability-level"
	ParameterSunsetDateTooSmallId      = "parameter-sunset-date-too-small"
	ParameterDeprecatedId              = "parameter-deprecated"
)

func ParameterDeprecationCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationDiff := range pathItem.OperationsDiff.Modified {
			if operationDiff.ParametersDiff == nil {
				continue
			}

			op := pathItem.Revision.GetOperation(operation)

			for paramLocation, paramItems := range operationDiff.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {

					if paramItem.DeprecatedDiff == nil {
						continue
					}

					param := paramItem.Revision

					if paramItem.DeprecatedDiff.To == nil || paramItem.DeprecatedDiff.To == false {
						// not breaking changes
						result = append(result, NewApiChange(
							ParameterReactivatedId,
							config,
							[]any{paramLocation, paramName},
							"",
							operationsSources,
							op,
							operation,
							path,
						))
						continue
					}

					stability, err := getStabilityLevel(op.Extensions) // TODO: how to handle stability?
					if err != nil {
						// handled in CheckBackwardCompatibility
						continue
					}

					deprecationDays := getDeprecationDays(config, stability)

					sunset, ok := getSunset(param.Extensions)
					if !ok {
						// if deprecation policy is defined and sunset is missing, it's a breaking change
						if deprecationDays > 0 {
							result = append(result, getAPIDeprecatedSunsetMissing(newOpInfo(config, op, operationsSources, operation, path)))
						}
						continue
					}

					date, err := getSunsetDate(sunset)
					if err != nil {
						result = append(result, NewApiChange(
							ParameterDeprecatedSunsetParseId,
							config,
							[]any{err},
							"",
							operationsSources,
							op,
							operation,
							path,
						))
						continue
					}

					days := date.DaysSince(civil.DateOf(time.Now()))

					if days < int(deprecationDays) {
						result = append(result, NewApiChange(
							ParameterSunsetDateTooSmallId,
							config,
							[]any{date, deprecationDays},
							"",
							operationsSources,
							op,
							operation,
							path,
						))
						continue
					}

					// not breaking changes
					result = append(result, NewApiChange(
						ParameterDeprecatedId,
						config,
						[]any{paramLocation, paramName},
						"",
						operationsSources,
						op,
						operation,
						path,
					))
				}
			}
		}
	}

	return result
}
