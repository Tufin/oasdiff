package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMaxLengthDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
					maxLengthDiff := paramDiff.SchemaDiff.MaxLengthDiff
					if maxLengthDiff == nil {
						continue
					}
					if maxLengthDiff.From == nil ||
						maxLengthDiff.To == nil {
						continue
					}

					if !IsDecreasedValue(maxLengthDiff) {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, BackwardCompatibilityError{
						Id:        "request-parameter-max-length-decreased",
						Level:     ERR,
						Text:      fmt.Sprintf("for the %s request parameter %s, the maxLength was decreased from %s to %s", ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(maxLengthDiff.From), ColorizedValue(maxLengthDiff.To)),
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
