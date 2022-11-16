package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMinIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
					minDiff := paramDiff.SchemaDiff.MinDiff
					if minDiff == nil {
						continue
					}
					if minDiff.From == nil ||
						minDiff.To == nil {
						continue
					}

					if !IsIncreasedValue(minDiff) {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, BackwardCompatibilityError{
						Id:        "request-parameter-min-increased",
						Level:     ERR,
						Text:      fmt.Sprintf("for the %s request parameter %s, the min was increased from %s to %s", ColorizedValue(paramLocation), ColorizedValue(paramName), minDiff.From, minDiff.To),
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
