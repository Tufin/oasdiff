package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterRemovedId = "request-parameter-removed"
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
					source := (*operationsSources)[operationItem.Revision]
					result = append(result, ApiChange{
						Id:          RequestParameterRemovedId,
						Level:       WARN,
						Text:        config.Localize(RequestParameterRemovedId, ColorizedValue(paramLocation), ColorizedValue(paramName)),
						Args:        []any{paramLocation, paramName},
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
