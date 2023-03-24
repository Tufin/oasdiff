package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	id = "api-operation-id-removed"
)

func APIOperationIdRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			op := pathItem.Base.Operations()[operation]
			source := (*operationsSources)[op]

			if operationItem.OperationIDDiff == nil {
				continue
			}

			result = append(result, BackwardCompatibilityError{
				Id:        id,
				Level:     ERR,
				Text:      fmt.Sprintf(config.i18n(id), ColorizedValue(operationItem.Base.OperationID), ColorizedValue(operationItem.Revision.OperationID)),
				Operation: operation,
				Path:      path,
				Source:    source,
			})
		}
	}
	return result
}
