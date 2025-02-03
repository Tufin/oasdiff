package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterReactivatedId             = "request-parameter-reactivated"
	RequestParameterDeprecatedSunsetParseId   = "request-parameter-deprecated-sunset-parse"
	RequestParameterDeprecatedSunsetMissingId = "request-parameter-deprecated-sunset-missing"
	RequestParameterInvalidStabilityLevelId   = "request-parameter-invalid-stability-level"
	RequestParameterSunsetDateTooSmallId      = "request-parameter-sunset-date-too-small"
	RequestParameterDeprecatedId              = "request-parameter-deprecated"
)

func RequestParameterDeprecationCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
							RequestParameterReactivatedId,
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
							RequestParameterDeprecatedSunsetParseId,
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
							RequestParameterSunsetDateTooSmallId,
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
						RequestParameterDeprecatedId,
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
