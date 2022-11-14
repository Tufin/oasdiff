package checker

import (
	"fmt"
	"strconv"

	"github.com/tufin/oasdiff/diff"
)

func ResponseSuccessStatusRemoved(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
			for _, responseStatus := range operationItem.ResponsesDiff.Deleted {
				status, err := strconv.Atoi(responseStatus)
				if err != nil {
					continue
				}
				if status < 200 || status > 299 {
					continue
				}

				result = append(result, BackwardCompatibilityError{
					Id:        "response-success-status-removed",
					Level:     ERR,
					Text:      fmt.Sprintf("removed the success response with the status %s", ColorizedValue(responseStatus)),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
			}
		}
	}
	return result
}
