package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	APITagRemovedCheckId = "api-tag-removed"
	APITagAddedCheckId   = "api-tag-added"
)

func APITagUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
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
				result = append(result, ApiChange{
					Id:          APITagRemovedCheckId,
					Level:       config.getLogLevel(APITagRemovedCheckId, INFO),
					Text:        config.Localize(APITagRemovedCheckId, ColorizedValue(tag)),
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      source,
				})

			}

			for _, tag := range operationItem.TagsDiff.Added {
				result = append(result, ApiChange{
					Id:          APITagAddedCheckId,
					Level:       config.getLogLevel(APITagAddedCheckId, INFO),
					Text:        config.Localize(APITagAddedCheckId, ColorizedValue(tag)),
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
