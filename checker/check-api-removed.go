package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
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
		if diffReport.PathsDiff.Base[path] == nil || diffReport.PathsDiff.Base[path].Operations() == nil {
			continue
		}
		for operation := range diffReport.PathsDiff.Base[path].Operations() {
			op := diffReport.PathsDiff.Base[path].Operations()[operation]
			if !op.Deprecated {
				source := "original_source=" + (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathRemovedWithoutDeprecationId,
					Level:       ERR,
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}
			rawDate, date, err := getSunsetDate(op.Extensions)
			if err != nil {
				source := "original_source=" + (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathSunsetParseId,
					Level:       ERR,
					Args:        []any{rawDate, err},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
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
					Source:      source,
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
					Source:      source,
				})
				continue
			}
			rawDate, date, err := getSunsetDate(op.Extensions)
			if err != nil {
				source := (*operationsSources)[op]
				result = append(result, ApiChange{
					Id:          APIPathSunsetParseId,
					Level:       ERR,
					Args:        []any{rawDate, err},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
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
					Source:      source,
				})
			}
		}
	}

	return result
}
