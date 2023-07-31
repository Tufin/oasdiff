package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMaxUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if maxDiff.From == nil ||
						maxDiff.To == nil {
						continue
					}

					id := "request-parameter-max-decreased"
					level := ERR

					if !IsDecreasedValue(maxDiff) {
						id = "request-parameter-max-increased"
						level = INFO
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Text:        fmt.Sprintf(config.i18n(id), ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
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
