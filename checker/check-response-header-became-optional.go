package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func ResponseHeaderBecameOptional(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil {
				continue
			}
			if operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.HeadersDiff == nil {
					continue
				}

				for headerName, headerDiff := range responseDiff.HeadersDiff.Modified {
					requiredDiff := headerDiff.RequiredDiff
					if requiredDiff == nil {
						continue
					}
					if requiredDiff.From != true {
						continue
					}

					result = append(result, ApiChange{
						Id:          "response-header-became-optional",
						Level:       ERR,
						Text:        config.Localize("response-header-became-optional", ColorizedValue(headerName), ColorizedValue(responseStatus)),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}
			}
		}
	}
	return result
}
