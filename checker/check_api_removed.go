package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
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
			if change := checkAPIRemoval(APIPathRemovedWithoutDeprecationId, APIPathRemovedBeforeSunsetId, op, operationsSources, operation, path); change != nil {
				result = append(result, change)
			}
		}
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for _, operation := range pathItem.OperationsDiff.Deleted {
			op := pathItem.Base.Operations()[operation]
			if change := checkAPIRemoval(APIRemovedWithoutDeprecationId, APIRemovedBeforeSunsetId, op, operationsSources, operation, path); change != nil {
				result = append(result, change)
			}
		}
	}

	return result
}

func checkAPIRemoval(deprecationId, sunsetId string, op *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method, path string) Change {
	if !op.Deprecated {
		return ApiChange{
			Id:          deprecationId,
			Level:       ERR,
			Operation:   method,
			OperationId: op.OperationID,
			Path:        path,
			Source:      load.NewSource((*operationsSources)[op]),
		}
	}
	sunset, ok := getSunset(op.Extensions)
	if !ok {
		// No sunset date, allow removal
		return nil
	}

	date, err := getSunsetDate(sunset)
	if err != nil {
		return getAPIPathSunsetParseId(op, operationsSources, method, path, err)
	}

	if !civil.DateOf(time.Now()).After(date) {
		return ApiChange{
			Id:          sunsetId,
			Level:       ERR,
			Args:        []any{date},
			Operation:   method,
			OperationId: op.OperationID,
			Path:        path,
			Source:      load.NewSource((*operationsSources)[op]),
		}
	}
	return nil
}

func getAPIPathSunsetParseId(operation *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method string, path string, err error) Change {
	return ApiChange{
		Id:          APIDeprecatedSunsetParseId,
		Level:       ERR,
		Args:        []any{err},
		Operation:   method,
		OperationId: operation.OperationID,
		Path:        path,
		Source:      load.NewSource((*operationsSources)[operation]),
	}
}
