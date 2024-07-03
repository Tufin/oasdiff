package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	APIOperationIdRemovedId = "api-operation-id-removed"
	APIOperationIdAddId     = "api-operation-id-added"
)

func APIOperationIdUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

			op := pathItem.Base.GetOperation(operation)
			id := APIOperationIdRemovedId
			args := []any{operationItem.Base.OperationID, operationItem.Revision.OperationID}
			if operationItem.OperationIDDiff.From == nil || operationItem.OperationIDDiff.From == "" {
				id = APIOperationIdAddId
				op = pathItem.Revision.GetOperation(operation)
				args = []any{operationItem.Revision.OperationID}
			}

			result = append(result, NewApiChange(
				id,
				config.getLogLevel(id, INFO),
				args,
				"",
				operationsSources,
				op,
				operation,
				path,
			))
		}
	}
	return result
}
