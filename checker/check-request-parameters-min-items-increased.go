package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMinItemsIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
					minItemsDiff := paramDiff.SchemaDiff.MinItemsDiff
					if minItemsDiff == nil {
						continue
					}
					if minItemsDiff.From == nil ||
						minItemsDiff.To == nil {
						continue
					}

					if !IsIncreasedValue(minItemsDiff) {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, BackwardCompatibilityError{
						Id:        "request-parameter-min-items-increased",
						Level:     ERR,
						Text:      fmt.Sprintf("for the %s request parameter %s, the minItems was increased from %s to %s", ColorizedValue(paramLocation), ColorizedValue(paramName), minItemsDiff.From, minItemsDiff.To),
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
