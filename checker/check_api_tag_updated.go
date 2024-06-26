package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	APITagRemovedId = "api-tag-removed"
	APITagAddedId   = "api-tag-added"
)

func APITagUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			op := pathItem.Base.GetOperation(operation)
			source := (*operationsSources)[op]

			if operationItem.TagsDiff == nil {
				continue
			}

			for _, tag := range operationItem.TagsDiff.Deleted {
				result = append(result, ApiChange{
					Id:          APITagRemovedId,
					Level:       config.getLogLevel(APITagRemovedId, INFO),
					Args:        []any{tag},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})

			}

			for _, tag := range operationItem.TagsDiff.Added {
				result = append(result, ApiChange{
					Id:          APITagAddedId,
					Level:       config.getLogLevel(APITagAddedId, INFO),
					Args:        []any{tag},
					Operation:   operation,
					OperationId: op.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}
		}
	}
	return result
}
