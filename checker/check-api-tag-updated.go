package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	APITagRemovedCheckId = "api-tag-removed"
	APITagAddedCheckId   = "api-tag-added"
)

func APITagUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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

			if operationItem.TagsDiff == nil {
				continue
			}

			for _, tag := range operationItem.TagsDiff.Deleted {
				result = append(result, BackwardCompatibilityError{
					Id:          APITagRemovedCheckId,
					Level:       config.getLogLevel(APITagRemovedCheckId, INFO),
					Text:        fmt.Sprintf(config.i18n(APITagRemovedCheckId), ColorizedValue(tag)),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})

			}

			for _, tag := range operationItem.TagsDiff.Added {
				result = append(result, BackwardCompatibilityError{
					Id:          APITagAddedCheckId,
					Level:       config.getLogLevel(APITagAddedCheckId, INFO),
					Text:        fmt.Sprintf(config.i18n(APITagAddedCheckId), ColorizedValue(tag)),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})

			}

		}
	}
	return result
}
