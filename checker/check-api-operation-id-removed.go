package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	apiOperationRemovedCheckId = "api-operation-id-removed"
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
				Id:          apiOperationRemovedCheckId,
				Level:       ERR,
				Text:        fmt.Sprintf(config.i18n(apiOperationRemovedCheckId), ColorizedValue(operationItem.Base.OperationID), ColorizedValue(operationItem.Revision.OperationID)),
				Operation:   operation,
				OperationId: op.OperationID,
				Path:        path,
				Source:      source,
			})
		}
	}
	return result
}
