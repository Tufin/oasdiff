package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
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
		if diffReport.PathsDiff.Base.Value(path) == nil {
			continue
		}

		for operation := range diffReport.PathsDiff.Base.Value(path).Operations() {
			op := diffReport.PathsDiff.Base.Value(path).GetOperation(operation)
			if change := checkAPIRemoval(config, APIPathRemovedWithoutDeprecationId, APIPathRemovedBeforeSunsetId, op, operationsSources, operation, path); change != nil {
				result = append(result, change)
			}
		}
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for _, operation := range pathItem.OperationsDiff.Deleted {
			op := pathItem.Base.GetOperation(operation)
			if change := checkAPIRemoval(config, APIRemovedWithoutDeprecationId, APIRemovedBeforeSunsetId, op, operationsSources, operation, path); change != nil {
				result = append(result, change)
			}
		}
	}

	return result
}

func checkAPIRemoval(config *Config, deprecationId, sunsetId string, op *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method, path string) Change {
	if !op.Deprecated {
		return NewApiChange(
			deprecationId,
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
		// No sunset date, allow removal
		return nil
	}

	date, err := getSunsetDate(sunset)
	if err != nil {
		return getAPIPathSunsetParse(config, op, operationsSources, method, path, err)
	}

	if civil.DateOf(time.Now()).Before(date) {
		return NewApiChange(
			sunsetId,
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
