package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func DefaultChecks() []BackwardCompatibilityCheck {
	checks := make([]BackwardCompatibilityCheck, 0)
	checks = append(checks, parameterRemovedCheck)
	return checks
}

func parameterRemovedCheck(diff *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diff.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diff.PathsDiff.Modified {
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
					result = append(result, BackwardCompatibilityError{
						Id:        "parameter-removed",
						Level:     WARN,
						Text:      fmt.Sprintf("deleted %s request parameter %s", paramLocation, paramName),
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