package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	APIOperationIdRemovedId = "api-operation-id-removed"
	APIOperationIdAddId     = "api-operation-id-added"
)

func APIOperationIdUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.OperationIDDiff == nil {
				continue
			}

			op := pathItem.Base.Operations()[operation]
			source := (*operationsSources)[op]

			id := APIOperationIdRemovedId
			text := config.Localize(id, ColorizedValue(operationItem.Base.OperationID), ColorizedValue(operationItem.Revision.OperationID))
			if operationItem.OperationIDDiff.From == nil || operationItem.OperationIDDiff.From == "" {
				id = APIOperationIdAddId
				op = pathItem.Revision.Operations()[operation]
				text = config.Localize(id, ColorizedValue(operationItem.Revision.OperationID))
			}

			result = append(result, ApiChange{
				Id:          id,
				Level:       config.getLogLevel(id, INFO),
				Text:        text,
				Operation:   operation,
				OperationId: op.OperationID,
				Path:        path,
				Source:      source,
			})
		}
	}
	return result
}
