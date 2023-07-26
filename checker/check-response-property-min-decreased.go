package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMinDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff == nil ||
					responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}
				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinDiff != nil {
						minDiff := mediaTypeDiff.SchemaDiff.MinDiff
						if minDiff.From != nil &&
							minDiff.To != nil {
							if isDecreasedValue(minDiff) {
								result = append(result, ApiChange{
									Id:          "response-body-min-decreased",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-body-min-decreased"), colorizedValue(minDiff.From), colorizedValue(minDiff.To)),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}
						}
					}

					checkModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							minDiff := propertyDiff.MinDiff
							if minDiff == nil {
								return
							}
							if minDiff.To == nil ||
								minDiff.From == nil {
								return
							}
							if !isDecreasedValue(minDiff) {
								return
							}

							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          "response-property-min-decreased",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-property-min-decreased"), colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(minDiff.From), colorizedValue(minDiff.To), colorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
				}
			}
		}
	}
	return result
}
