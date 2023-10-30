package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterMaxSetId = "request-parameter-max-set"
)

func RequestParameterMaxSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}
					maxDiff := paramDiff.SchemaDiff.MaxDiff
					if maxDiff == nil {
						continue
					}
					if maxDiff.From != nil ||
						maxDiff.To == nil {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, ApiChange{
						Id:          RequestParameterMaxSetId,
						Level:       WARN,
						Text:        config.Localize(RequestParameterMaxSetId, ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(maxDiff.To)),
						Comment:     config.Localize(RequestParameterMaxSetId + "-comment"),
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
