package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMinUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					minDiff := paramDiff.SchemaDiff.MinDiff
					if minDiff == nil {
						continue
					}
					if minDiff.From == nil ||
						minDiff.To == nil {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					id := "request-parameter-min-increased"
					level := ERR
					if !IsIncreasedValue(minDiff) {
						id = "request-parameter-min-decreased"
						level = INFO
					}

					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Text:        config.Localize(id, ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(minDiff.From), ColorizedValue(minDiff.To)),
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
