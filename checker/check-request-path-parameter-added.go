package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func NewRequestPathParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
				if paramLocation != "path" {
					continue
				}

				for _, paramName := range paramItems {
					source := (*operationsSources)[operationItem.Revision]
					result = append(result, BackwardCompatibilityError{
						Id:          "new-request-path-parameter",
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n("new-request-path-parameter"), ColorizedValue(paramName)),
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
