package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	RequestBodyBecameOptionalId = "request-body-became-optional"
	RequestBodyBecameRequiredId = "request-body-became-required"
)

func RequestBodyRequiredUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

			if operationItem.RequestBodyDiff.RequiredDiff == nil {
				continue
			}

			source := (*operationsSources)[operationItem.Revision]

			id := RequestBodyBecameOptionalId
			logLevel := INFO
			if operationItem.RequestBodyDiff.RequiredDiff.To == true {
				id = RequestBodyBecameRequiredId
				logLevel = ERR
			}

			result = append(result, ApiChange{
				Id:          id,
				Level:       logLevel,
				Operation:   operation,
				OperationId: operationItem.Revision.OperationID,
				Path:        path,
				Source:      load.NewSource(source),
			})
		}
	}
	return result
}
