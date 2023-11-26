package checker

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
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
			op := pathItem.Revision.Operations()[operation]
			source := (*operationsSources)[op]

			if !op.Deprecated {
				continue
			}

			if operationDiff.ExtensionsDiff != nil && !operationDiff.ExtensionsDiff.Deleted.Empty() {
				result = append(result, ApiChange{
					Id:          APISunsetDeletedId,
					Level:       ERR,
					Text:        config.Localize(APISunsetDeletedId),
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

			rawDate, date, err := getSunsetDate(op.Extensions)
			if err != nil {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetParseId,
					Level:       ERR,
					Text:        config.Localize(APIDeprecatedSunsetParseId, rawDate, err),
					Args:        []any{rawDate, err},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}

			rawDate, baseDate, err := getSunsetDate(opBase.Extensions)
			if err != nil {
				result = append(result, ApiChange{
					Id:          APIDeprecatedSunsetParseId,
					Level:       ERR,
					Text:        config.Localize(APIDeprecatedSunsetParseId, rawDate, err),
					Args:        []any{rawDate, err},
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
					Id:          ParseErrorId,
					Level:       ERR,
					Text:        fmt.Sprintf("parsing error %s", err.Error()),
					Args:        []any{err.Error()},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}

			deprecationDays := getDeprecationDays(config, stability)

			if baseDate.After(date) && days < deprecationDays {
				result = append(result, ApiChange{
					Id:          APISunsetDateChangedTooSmallId,
					Level:       ERR,
					Text:        config.Localize(APISunsetDateChangedTooSmallId, baseDate, date, baseDate, deprecationDays),
					Args:        []any{baseDate, date, baseDate, deprecationDays},
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
