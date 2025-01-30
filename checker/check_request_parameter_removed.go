package checker

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterRemovedId                = "request-parameter-removed"
	RequestParameterRemovedWithDeprecationId = "request-parameter-removed-with-deprecation"
	ParameterRemovedBeforeSunsetId           = "request-parameter-removed-before-sunset"
)

func RequestParameterRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ParametersDiff == nil {
				continue
			}

			for paramLocation, paramItems := range operationItem.ParametersDiff.Deleted {
				for _, paramName := range paramItems {
					param := operationItem.Base.Parameters.GetByInAndName(paramLocation, paramName)
					if change := checkParameterRemoval(config, operationItem.Revision, operationsSources, param, operation, path, paramLocation, paramName); change != nil {
						result = append(result, change)
					}
				}
			}
		}
	}
	return result
}

func checkParameterRemoval(config *Config, op *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, param *openapi3.Parameter, method, path, paramLocation, paramName string) Change {
	if !param.Deprecated {
		return NewApiChange(
			RequestParameterRemovedId,
			config,
			[]any{paramLocation, paramName},
			"",
			operationsSources,
			op,
			method,
			path,
		)
	}

	sunset, ok := getSunset(param.Extensions)
	if !ok {
		return NewApiChange(
			RequestParameterRemovedWithDeprecationId,
			config,
			[]any{paramLocation, paramName},
			"",
			operationsSources,
			op,
			method,
			path,
		)
	}

	date, err := getSunsetDate(sunset)
	if err != nil {
		return NewApiChange(
			APIPathSunsetParseId,
			config,
			[]any{err},
			"",
			operationsSources,
			op,
			method,
			path,
		)
	}

	if civil.DateOf(time.Now()).Before(date) {
		return NewApiChange(
			ParameterRemovedBeforeSunsetId,
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
