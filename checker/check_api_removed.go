package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	APIPathRemovedWithoutDeprecationId = "api-path-removed-without-deprecation"
	APIPathSunsetParseId               = "api-path-sunset-parse"
	APIPathRemovedBeforeSunsetId       = "api-path-removed-before-sunset"
	APIRemovedWithoutDeprecationId     = "api-removed-without-deprecation"
	APIRemovedBeforeSunsetId           = "api-removed-before-sunset"
)

func APIRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for _, path := range diffReport.PathsDiff.Deleted {
		if diffReport.PathsDiff.Base.Value(path) == nil || diffReport.PathsDiff.Base.Value(path).Operations() == nil {
			continue
		}
		for operation := range diffReport.PathsDiff.Base.Value(path).Operations() {
			op := diffReport.PathsDiff.Base.Value(path).Operations()[operation]
			if !op.Deprecated {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathRemovedWithoutDeprecationId,
					Level:       ERR,
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			sunset, ok := getSunset(op.Extensions)
			if !ok {
				// No sunset date, allow removal
				continue
			}

			date, err := getSunsetDate(sunset)
			if err != nil {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathSunsetParseId,
					Level:       ERR,
					Args:        []any{err},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}

			if !civil.DateOf(time.Now()).After(date) {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathRemovedBeforeSunsetId,
					Level:       ERR,
					Args:        []any{date},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}
		}
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for _, operation := range pathItem.OperationsDiff.Deleted {
			op := pathItem.Base.Operations()[operation]
			if !op.Deprecated {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIRemovedWithoutDeprecationId,
					Level:       ERR,
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}
			sunset, ok := getSunset(op.Extensions)
			if !ok {
				// No sunset date, allow removal
				continue
			}

			date, err := getSunsetDate(sunset)
			if err != nil {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathSunsetParseId,
					Level:       ERR,
					Args:        []any{err},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}
			if !civil.DateOf(time.Now()).After(date) {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIRemovedBeforeSunsetId,
					Level:       ERR,
					Args:        []any{date},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}
		}
	}

	return result
}
