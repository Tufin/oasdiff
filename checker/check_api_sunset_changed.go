package checker

import (
	"time"

	"cloud.google.com/go/civil"
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
			source := (*operationsSources)[opRevision]

			if !opRevision.Deprecated {
				continue
			}

			if operationDiff.ExtensionsDiff != nil && !operationDiff.ExtensionsDiff.Deleted.Empty() {
				result = append(result, ApiChange{
					Id:          APISunsetDeletedId,
					Level:       ERR,
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}

			if operationDiff.ExtensionsDiff == nil || operationDiff.ExtensionsDiff.Modified.Empty() {
				continue
			}

			opBase := pathItem.Base.Operations()[operation]

			sunsetRevision, ok := getSunset(opRevision.Extensions)
			if !ok {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetMissingId,
					Level:       ERR,
					Args:        []any{},
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			date, err := getSunsetDate(sunsetRevision)
			if err != nil {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetParseId,
					Level:       ERR,
					Args:        []any{err},
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			sunsetBase, ok := getSunset(opBase.Extensions)
			if !ok {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetMissingId,
					Level:       ERR,
					Args:        []any{},
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource((*operationsSources)[opBase]),
				})
				continue
			}

			baseDate, err := getSunsetDate(sunsetBase)
			if err != nil {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetParseId,
					Level:       ERR,
					Args:        []any{err},
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource((*operationsSources)[opBase]),
				})
				continue
			}

			days := date.DaysSince(civil.DateOf(time.Now()))

			stability, err := getStabilityLevel(opRevision.Extensions)
			if err != nil {
				// handled in CheckBackwardCompatibility
				continue
			}

			deprecationDays := getDeprecationDays(config, stability)

			if baseDate.After(date) && days < deprecationDays {
				result = append(result, ApiChange{
					Id:          APISunsetDateChangedTooSmallId,
					Level:       ERR,
					Args:        []any{baseDate, date, baseDate, deprecationDays},
					Operation:   operation,
					OperationId: opRevision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
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

func getDeprecationDays(config *Config, stability string) int {
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
