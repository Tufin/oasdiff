package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestParameterRequiredValueUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			for paramLocation, paramItems := range operationItem.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {
					requiredDiff := paramItem.RequiredDiff
					if requiredDiff == nil {
						continue
					}

					id := "request-parameter-became-required"
					level := ERR

					if requiredDiff.To != true {
						id = "request-parameter-became-optional"
						level = INFO
					}

					source := (*operationsSources)[operationItem.Revision]
					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Text:        config.Localize(id, ColorizedValue(paramLocation), ColorizedValue(paramName)),
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
