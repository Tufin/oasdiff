package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponseMediaTypeRemoved(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
			for responseStatus, responsesDiff := range operationItem.ResponsesDiff.Modified {
				if responsesDiff.ContentDiff == nil {
					continue
				}
				if responsesDiff.ContentDiff.MediaTypeDeleted == nil {
					continue
				}
				for _, mediaType := range responsesDiff.ContentDiff.MediaTypeDeleted {
					result = append(result, BackwardCompatibilityError{
						Id:          "response-media-type-removed",
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n("response-media-type-removed"), ColorizedValue(mediaType), ColorizedValue(responseStatus)),
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
