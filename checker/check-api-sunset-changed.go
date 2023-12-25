package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

const (
	APISunsetDeletedId             = "sunset-deleted"
	APISunsetDateChangedTooSmallId = "api-sunset-date-changed-too-small"
)

func APISunsetChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	changeGetter := newApiChangeGetter(config, operationsSources)
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

			if !op.Deprecated {
				continue
			}

			if operationDiff.ExtensionsDiff != nil && !operationDiff.ExtensionsDiff.Deleted.Empty() {
				result = append(result, changeGetter(
					APISunsetDeletedId,
					ERR,
					nil,
					"",
					operation,
					op,
					path,
					op,
				))
			}

			if operationDiff.ExtensionsDiff == nil || operationDiff.ExtensionsDiff.Modified.Empty() {
				continue
			}

			rawDate, date, err := getSunsetDate(op.Extensions)
			if err != nil {
				result = append(result, changeGetter(
					APIDeprecatedSunsetParseId,
					ERR,
					[]any{rawDate, err},
					"",
					operation,
					op,
					path,
					op,
				))
				continue
			}

			opBase := pathItem.Base.Operations()[operation]
			rawDate, baseDate, err := getSunsetDate(opBase.Extensions)
			if err != nil {
				result = append(result, changeGetter(
					APIDeprecatedSunsetParseId,
					ERR,
					[]any{rawDate, err},
					"",
					operation,
					op,
					path,
					opBase,
				))
				continue
			}

			days := date.DaysSince(civil.DateOf(time.Now()))

			stability, err := getStabilityLevel(op.Extensions)
			if err != nil {
				result = append(result, changeGetter(
					ParseErrorId,
					ERR,
					[]any{err.Error()},
					"",
					operation,
					op,
					path,
					op,
				))
				continue
			}

			deprecationDays := getDeprecationDays(config, stability)

			if baseDate.After(date) && days < deprecationDays {
				result = append(result, changeGetter(
					APISunsetDateChangedTooSmallId,
					ERR,
					[]any{baseDate, date, baseDate, deprecationDays},
					"",
					operation,
					op,
					path,
					op,
				))
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
