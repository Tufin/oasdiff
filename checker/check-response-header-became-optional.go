package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponseHeaderBecameOptional(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil {
				continue
			}
			if operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.HeadersDiff == nil {
					continue
				}

				for headerName, headerDiff := range responseDiff.HeadersDiff.Modified {
					requiredDiff := headerDiff.RequiredDiff
					if requiredDiff == nil {
						continue
					}
					if requiredDiff.From != true {
						continue
					}

					result = append(result, BackwardCompatibilityError{
						Id:        "response-header-became-optional",
						Level:     ERR,
						Text:      fmt.Sprintf("the response header %s became optional for the status %s", ColorizedValue(headerName), ColorizedValue(responseStatus)),
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
