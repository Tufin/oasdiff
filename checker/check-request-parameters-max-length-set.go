package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterMaxLengthSetId = "request-parameter-max-length-set"
)

func RequestParameterMaxLengthSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					maxLengthDiff := paramDiff.SchemaDiff.MaxLengthDiff
					if maxLengthDiff == nil {
						continue
					}
					if maxLengthDiff.From != nil ||
						maxLengthDiff.To == nil {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, ApiChange{
						Id:          RequestParameterMaxLengthSetId,
						Level:       WARN,
						Text:        config.Localize(RequestParameterMaxLengthSetId, ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(maxLengthDiff.To)),
						Args:        []any{},
						Comment:     config.Localize(comment(RequestParameterMaxLengthSetId)),
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
