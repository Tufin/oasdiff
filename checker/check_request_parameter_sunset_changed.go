package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterSunsetDeletedId             = "request-parameter-sunset-deleted"
	RequestParameterSunsetDateChangedTooSmallId = "request-parameter-sunset-date-changed-too-small"
)

func RequestParameterSunsetChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationDiff := range pathItem.OperationsDiff.Modified {
			opRevision := pathItem.Revision.GetOperation(operation)
			opBase := pathItem.Base.GetOperation(operation)

			// op := pathItem.Revision.GetOperation(operation)
			// opInfo := newOpInfo(config, op, operationsSources, operation, path)

			if operationDiff.ParametersDiff == nil {
				continue
			}

			for paramLocation, paramItems := range operationDiff.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {

					paramBase := opBase.Parameters.GetByInAndName(paramLocation, paramName)
					paramRevision := opRevision.Parameters.GetByInAndName(paramLocation, paramName)

					if !paramRevision.Deprecated {
						continue
					}

					if paramItem.ExtensionsDiff == nil {
						continue
					}

					if paramItem.ExtensionsDiff.Deleted.Contains(diff.SunsetExtension) {
						result = append(result, NewApiChange(
							RequestParameterSunsetDeletedId,
							config,
							[]any{paramLocation, paramName},
							"",
							operationsSources,
							opRevision,
							operation,
							path,
						))
						continue
					}

					if _, ok := paramItem.ExtensionsDiff.Modified[diff.SunsetExtension]; !ok {
						continue
					}

					date, err := getSunsetDate(paramRevision.Extensions[diff.SunsetExtension])
					if err != nil {
						opInfo := newOpInfo(config, opRevision, operationsSources, operation, path)
						result = append(result, getRequestParameterSunsetParse(opInfo, paramRevision, err))
						continue
					}

					baseDate, err := getSunsetDate(paramBase.Extensions[diff.SunsetExtension])
					if err != nil {
						opInfo := newOpInfo(config, opBase, operationsSources, operation, path)
						result = append(result, getRequestParameterSunsetParse(opInfo, paramBase, err))
						continue
					}

					days := date.DaysSince(civil.DateOf(time.Now()))

					stability, err := getStabilityLevel(opRevision.Extensions)
					if err != nil {
						// handled in CheckBackwardCompatibility
						continue
					}

					deprecationDays := getDeprecationDays(config, stability)

					if baseDate.After(date) && days < int(deprecationDays) {
						result = append(result, NewApiChange(
							RequestParameterSunsetDateChangedTooSmallId,
							config,
							[]any{paramRevision.In, paramRevision.Name, baseDate, date, baseDate, deprecationDays},
							"",
							operationsSources,
							opRevision,
							operation,
							path,
						))
					}
				}
			}
		}
	}

	return result
}
