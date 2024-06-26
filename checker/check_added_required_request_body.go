package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	AddedRequiredRequestBodyId = "request-body-added-required"
	AddedOptionalRequestBodyId = "request-body-added-optional"
)

func AddedRequestBodyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
			if operationItem.RequestBodyDiff.Added {
				id := AddedOptionalRequestBodyId
				logLevel := INFO

				if operationItem.Revision.RequestBody.Value.Required {
					id = AddedRequiredRequestBodyId
					logLevel = ERR
				}

				result = append(result, NewApiChange(
					id,
					logLevel,
					nil,
					"",
					operationsSources,
					operationItem.Revision,
					operation,
					path,
				))
			}
		}
	}
	return result
}
