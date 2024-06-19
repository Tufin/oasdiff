package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	APISunsetDeletedId             = "sunset-deleted"
	APISunsetDateChangedTooSmallId = "api-sunset-date-changed-too-small"
)

func APISunsetChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationDiff := range pathItem.OperationsDiff.Modified {
			opRevision := pathItem.Revision.Operations()[operation]
			opBase := pathItem.Base.Operations()[operation]

			if !opRevision.Deprecated {
				continue
			}

			if operationDiff.ExtensionsDiff == nil {
				continue
			}

			if operationDiff.ExtensionsDiff.Deleted.Contains(diff.SunsetExtension) {
				result = append(result, ApiChange{
					Id:          APISunsetDeletedId,
					Level:       ERR,
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource((*operationsSources)[opRevision]),
				})
				continue
			}

			if _, ok := operationDiff.ExtensionsDiff.Modified[diff.SunsetExtension]; !ok {
				continue
			}

			date, err := getSunsetDate(opRevision.Extensions[diff.SunsetExtension])
			if err != nil {
				result = append(result, getAPIPathSunsetParse(opRevision, operationsSources, path, operation, err))
				continue
			}

			baseDate, err := getSunsetDate(opBase.Extensions[diff.SunsetExtension])
			if err != nil {
				result = append(result, getAPIPathSunsetParse(opBase, operationsSources, path, operation, err))
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
				result = append(result, ApiChange{
					Id:          APISunsetDateChangedTooSmallId,
					Level:       ERR,
					Args:        []any{baseDate, date, baseDate, deprecationDays},
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource((*operationsSources)[opRevision]),
				})
			}
		}
	}

	return result
}

const (
	STABILITY_DRAFT  = "draft"
	STABILITY_ALPHA  = "alpha"
	STABILITY_BETA   = "beta"
	STABILITY_STABLE = "stable"
)

func getDeprecationDays(config *Config, stability string) uint {
	switch stability {
	case STABILITY_DRAFT:
		return 0
	case STABILITY_ALPHA:
		return 0
	case STABILITY_BETA:
		return config.MinSunsetBetaDays
	case STABILITY_STABLE:
		return config.MinSunsetStableDays
	default:
		return config.MinSunsetStableDays
	}
}

func getAPIDeprecatedSunsetMissing(operation *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method string, path string) Change {
	return ApiChange{
		Id:          APIDeprecatedSunsetMissingId,
		Level:       ERR,
		Args:        []any{},
		Operation:   method,
		OperationId: operation.OperationID,
		Path:        path,
		Source:      load.NewSource((*operationsSources)[operation]),
	}
}
