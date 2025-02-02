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

			helper := newHelper(
				config,
				operationItem.Revision,
				operationsSources,
				operation,
				path,
			)

			for paramLocation, paramItems := range operationItem.ParametersDiff.Deleted {
				for _, paramName := range paramItems {
					if change := helper.checkParameterRemoval(operationItem.Base.Parameters.GetByInAndName(paramLocation, paramName)); change != nil {
						result = append(result, change)
					}
				}
			}
		}
	}
	return result
}

func (o helper) checkParameterRemoval(param *openapi3.Parameter) Change {

	if !param.Deprecated {
		return NewApiChange(
			RequestParameterRemovedId,
			o.config,
			[]any{param.In, param.Name},
			commentId(RequestParameterRemovedId),
			o.operationsSources,
			o.operation,
			o.method,
			o.path,
		)
	}

	sunset, ok := getSunset(param.Extensions)
	if !ok {
		return NewApiChange(
			RequestParameterRemovedWithDeprecationId,
			o.config,
			[]any{param.In, param.Name},
			"",
			o.operationsSources,
			o.operation,
			o.method,
			o.path,
		)
	}

	date, err := getSunsetDate(sunset)
	if err != nil {
		return NewApiChange(
			RequestParameterSunsetParseId,
			o.config,
			[]any{param.In, param.Name, err},
			"",
			o.operationsSources,
			o.operation,
			o.method,
			o.path,
		)
	}

	if civil.DateOf(time.Now()).Before(date) {
		return NewApiChange(
			ParameterRemovedBeforeSunsetId,
			o.config,
			[]any{param.In, param.Name, date},
			"",
			o.operationsSources,
			o.operation,
			o.method,
			o.path,
		)
	}
	return nil
}
