package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMaxDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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

					if !IsDecreasedValue(maxDiff) {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, BackwardCompatibilityError{
						Id:        "request-parameter-max-decreased",
						Level:     ERR,
						Text:      fmt.Sprintf("for the %s request parameter %s, the max was decreased from %s to %s", ColorizedValue(paramLocation), ColorizedValue(paramName), maxDiff.From, maxDiff.To),
						Operation: operation,
						Path:      path,
						Source:    source,
						ToDo:      "Add to exceptions-list.md",
					})
				}
			}
		}
	}
	return result
}
