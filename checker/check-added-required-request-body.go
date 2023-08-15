package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func AddedRequiredRequestBodyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.RequestBodyDiff == nil {
				continue
			}
			if operationItem.RequestBodyDiff.Added &&
				operationItem.Revision.RequestBody.Value.Required {
				source := (*operationsSources)[operationItem.Revision]
				result = append(result, ApiChange{
					Id:          "added-required-request-body",
					Level:       ERR,
					Text:        config.Localize("added-required-request-body"),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}
		}
	}
	return result
}
