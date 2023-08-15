package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMinLengthUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					minLengthDiff := paramDiff.SchemaDiff.MinLengthDiff
					if minLengthDiff == nil {
						continue
					}
					if minLengthDiff.From == nil ||
						minLengthDiff.To == nil {
						continue
					}

					id := "request-parameter-min-length-increased"
					level := ERR
					if IsDecreasedValue(minLengthDiff) {
						id = "request-parameter-min-length-decreased"
						level = INFO
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Text:        config.Localize(id, ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
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
