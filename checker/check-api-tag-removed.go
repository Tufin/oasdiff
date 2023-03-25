package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	apiTagRemovedCheckId = "api-tag-removed"
)

func APITagRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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

			if operationItem.TagsDiff == nil || len(operationItem.TagsDiff.Deleted) == 0 {
				continue
			}

			result = append(result, BackwardCompatibilityError{
				Id:        apiTagRemovedCheckId,
				Level:     ERR,
				Text:      fmt.Sprintf(config.i18n(apiTagRemovedCheckId), ColorizedValue(operationItem.TagsDiff.Deleted)),
				Operation: operation,
				Path:      path,
				Source:    source,
			})
		}
	}
	return result
}
