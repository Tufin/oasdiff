package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterBecameEnumId = "request-parameter-became-enum"
)

func RequestParameterBecameEnumCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			if operationItem.ParametersDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for paramLocation, paramItems := range operationItem.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {
					if paramItem.SchemaDiff == nil {
						continue
					}

					if enumDiff := paramItem.SchemaDiff.EnumDiff; enumDiff == nil || !enumDiff.EnumAdded {
						continue
					}

					result = append(result, ApiChange{
						Id:          RequestParameterBecameEnumId,
						Level:       ERR,
						Text:        config.Localize(RequestParameterBecameEnumId, ColorizedValue(paramLocation), ColorizedValue(paramName)),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}
			}
		}
	}
	return result
}
