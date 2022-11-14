package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func NewRequiredRequestNonPathParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
			for paramLocation, paramItems := range operationItem.ParametersDiff.Added {
				if paramLocation == "path" {
					// it is processed in the separate check NewRequestPathParameterCheck
					continue
				}

				for _, paramName := range paramItems {
					for _, param := range operationItem.Revision.Parameters {
						if param.Value.Name == paramName {
							if param.Value.Required {
								source := (*operationsSources)[operationItem.Revision]
								result = append(result, BackwardCompatibilityError{
									Id:        "new-required-request-parameter",
									Level:     ERR,
									Text:      fmt.Sprintf("added the new required %s request parameter %s", ColorizedValue(paramLocation), ColorizedValue(paramName)),
									Operation: operation,
									Path:      path,
									Source:    source,
									ToDo:      "Add to exceptions-list.md",
								})
							}
							break
						}
					}
				}
			}
		}
	}
	return result
}
