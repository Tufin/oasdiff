package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

const (
	APIPathRemovedWithoutDeprecationId = "api-path-removed-without-deprecation"
	APIPathRemovedWithDeprecationId    = "api-path-removed-with-deprecation"
	APIPathSunsetParseId               = "api-path-sunset-parse"
	APIPathRemovedBeforeSunsetId       = "api-path-removed-before-sunset"
	APIRemovedWithoutDeprecationId     = "api-removed-without-deprecation"
	APIRemovedWithDeprecationId        = "api-removed-with-deprecation"
	APIRemovedBeforeSunsetId           = "api-removed-before-sunset"
)

func APIRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	return append(
		checkRemovedPaths(diffReport.PathsDiff, operationsSources, config),
		checkRemovedOperations(diffReport.PathsDiff, operationsSources, config)...,
	)
}

func checkRemovedPaths(pathsDiff *diff.PathsDiff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {

	if pathsDiff == nil {
		return nil
	}

	result := make(Changes, 0)
	for _, path := range pathsDiff.Deleted {
		if pathsDiff.Base.Value(path) == nil {
			continue
		}

		for operation := range pathsDiff.Base.Value(path).Operations() {
			op := pathsDiff.Base.Value(path).GetOperation(operation)
			if change := checkAPIRemoval(config, true, op, operationsSources, operation, path); change != nil {
				result = append(result, change)
			}
		}
	}
	return result
}

func checkRemovedOperations(pathsDiff *diff.PathsDiff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	if pathsDiff == nil {
		return nil
	}

	result := make(Changes, 0)

	for path, pathItem := range pathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for _, operation := range pathItem.OperationsDiff.Deleted {
			op := pathItem.Base.GetOperation(operation)
			if change := checkAPIRemoval(config, false, op, operationsSources, operation, path); change != nil {
				result = append(result, change)
			}
		}
	}

	return result
}

func checkAPIRemoval(config *Config, isPath bool, op *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method, path string) Change {
	if !op.Deprecated {
		return NewApiChange(
			getWithoutDeprecationId(isPath),
			config,
			nil,
			"",
			operationsSources,
			op,
			method,
			path,
		)
	}
	sunset, ok := getSunset(op.Extensions)
	if !ok {
		return NewApiChange(
			getWithDeprecationId(isPath),
			config,
			nil,
			"",
			operationsSources,
			op,
			method,
			path,
		)
	}

	date, err := getSunsetDate(sunset)
	if err != nil {
		return getAPIPathSunsetParse(config, op, operationsSources, method, path, err)
	}

	if civil.DateOf(time.Now()).Before(date) {
		return NewApiChange(
			getBeforeSunsetId(isPath),
			config,
			[]any{date},
			"",
			operationsSources,
			op,
			method,
			path,
		)
	}
	return nil
}

func getAPIPathSunsetParse(config *Config, operation *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method string, path string, err error) Change {
	return NewApiChange(
		APIPathSunsetParseId,
		config,
		[]any{err},
		"",
		operationsSources,
		operation,
		method,
		path,
	)
}

func getWithDeprecationId(isPath bool) string {
	if isPath {
		return APIPathRemovedWithDeprecationId
	}
	return APIRemovedWithDeprecationId
}

func getWithoutDeprecationId(isPath bool) string {
	if isPath {
		return APIPathRemovedWithoutDeprecationId
	}
	return APIRemovedWithoutDeprecationId
}

func getBeforeSunsetId(isPath bool) string {
	if isPath {
		return APIPathRemovedBeforeSunsetId
	}
	return APIRemovedBeforeSunsetId
}
