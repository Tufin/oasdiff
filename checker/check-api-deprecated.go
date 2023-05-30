package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func APIDeprecatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationDiff := range pathItem.OperationsDiff.Modified {
			op := pathItem.Revision.Operations()[operation]
			source := (*operationsSources)[op]

			if operationDiff.DeprecatedDiff == nil {
				continue
			}
			id := "api-path-deprecated"
			if operationDiff.DeprecatedDiff.To == false || operationDiff.DeprecatedDiff.To == nil {
				id = "api-path-reactivated"
			}

			result = append(result, BackwardCompatibilityError{
				Id:          id,
				Level:       INFO,
				Text:        config.i18n(id),
				Operation:   operation,
				OperationId: op.OperationID,
				Path:        path,
				Source:      source,
			})

		}
	}

	return result
}
