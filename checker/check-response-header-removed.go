package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequiredResponseHeaderRemovedId = "required-response-header-removed"
	OptionalResponseHeaderRemovedId = "optional-response-header-removed"
)

func ResponseHeaderRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

				for _, headerName := range responseDiff.HeadersDiff.Deleted {
					if responseDiff.Base.Headers[headerName] == nil {
						continue
					}
					required := responseDiff.Base.Headers[headerName].Value.Required
					if required {
						result = append(result, ApiChange{
							Id:          RequiredResponseHeaderRemovedId,
							Level:       ERR,
							Text:        config.Localize(RequiredResponseHeaderRemovedId, ColorizedValue(headerName), ColorizedValue(responseStatus)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					} else {
						result = append(result, ApiChange{
							Id:          OptionalResponseHeaderRemovedId,
							Level:       WARN,
							Text:        config.Localize(OptionalResponseHeaderRemovedId, ColorizedValue(headerName), ColorizedValue(responseStatus)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}
			}
		}
	}
	return result
}
