package checker

import (
	"fmt"
	"strconv"

	"github.com/tufin/oasdiff/diff"
)

func ResponseSuccessStatusRemoved(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	success := func(status int) bool {
		return status >= 200 && status <= 299
	}

	return ResponseSuccessRemoved(diffReport, operationsSources, config, success, "response-success-status-removed")
}

func ResponseNonSuccessStatusRemoved(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	notSuccess := func(status int) bool {
		return status < 200 || status > 299
	}

	return ResponseSuccessRemoved(diffReport, operationsSources, config, notSuccess, "response-non-success-status-removed")
}

func ResponseSuccessRemoved(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig, filter func(int) bool, id string) []BackwardCompatibilityError {
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

				if filter(status) {
					result = append(result, BackwardCompatibilityError{
						Id:          id,
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n(id), ColorizedValue(responseStatus)),
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
