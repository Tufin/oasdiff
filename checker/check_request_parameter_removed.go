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
	RequestParameterSunsetParseId            = "request-parameter-sunset-parse"
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

			opInfo := newOpInfo(
				config,
				operationItem.Revision,
				operationsSources,
				operation,
				path,
			)

			for paramLocation, paramItems := range operationItem.ParametersDiff.Deleted {
				for _, paramName := range paramItems {
					if change := checkParameterRemoval(opInfo, operationItem.Base.Parameters.GetByInAndName(paramLocation, paramName)); change != nil {
						result = append(result, change)
					}
				}
			}
		}
	}
	return result
}

func checkParameterRemoval(oi opInfo, param *openapi3.Parameter) Change {

	if !param.Deprecated {
		return NewApiChange(
			RequestParameterRemovedId,
			oi.config,
			[]any{param.In, param.Name},
			commentId(RequestParameterRemovedId),
			oi.operationsSources,
			oi.operation,
			oi.method,
			oi.path,
		)
	}

	sunset, ok := getSunset(param.Extensions)
	if !ok {
		return NewApiChange(
			RequestParameterRemovedWithDeprecationId,
			oi.config,
			[]any{param.In, param.Name},
			"",
			oi.operationsSources,
			oi.operation,
			oi.method,
			oi.path,
		)
	}

	date, err := getSunsetDate(sunset)
	if err != nil {
		return NewApiChange(
			RequestParameterSunsetParseId,
			oi.config,
			[]any{param.In, param.Name, err},
			"",
			oi.operationsSources,
			oi.operation,
			oi.method,
			oi.path,
		)
	}

	if civil.DateOf(time.Now()).Before(date) {
		return NewApiChange(
			ParameterRemovedBeforeSunsetId,
			oi.config,
			[]any{param.In, param.Name, date},
			"",
			oi.operationsSources,
			oi.operation,
			oi.method,
			oi.path,
		)
	}
	return nil
}
