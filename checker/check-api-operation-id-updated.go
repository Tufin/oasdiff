package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	apiOperationRemovedCheckId = "api-operation-id-removed"
	apiOperationAddCheckId     = "api-operation-id-added"
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

			level := INFO
			id := apiOperationRemovedCheckId
			text := fmt.Sprintf(config.i18n(id), ColorizedValue(operationItem.Base.OperationID), ColorizedValue(operationItem.Revision.OperationID))
			if operationItem.OperationIDDiff.From == nil || operationItem.OperationIDDiff.From == "" {
				id = apiOperationAddCheckId
				op = pathItem.Revision.Operations()[operation]
				text = fmt.Sprintf(config.i18n(id), ColorizedValue(operationItem.Revision.OperationID))
			}

			result = append(result, ApiChange{
				Id:          id,
				Level:       config.getLogLevel(id, level),
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
