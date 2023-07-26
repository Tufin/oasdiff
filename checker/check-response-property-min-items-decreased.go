package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMinItemsDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinItemsDiff != nil {
						minItemsDiff := mediaTypeDiff.SchemaDiff.MinItemsDiff
						if minItemsDiff.From != nil &&
							minItemsDiff.To != nil {
							if isDecreasedValue(minItemsDiff) {
								result = append(result, ApiChange{
									Id:          "response-body-min-items-decreased",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-body-min-items-decreased"), colorizedValue(minItemsDiff.From), colorizedValue(minItemsDiff.To)),
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
							minItemsDiff := propertyDiff.MinItemsDiff
							if minItemsDiff == nil {
								return
							}
							if minItemsDiff.To == nil ||
								minItemsDiff.From == nil {
								return
							}
							if !isDecreasedValue(minItemsDiff) {
								return
							}

							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          "response-property-min-items-decreased",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-property-min-items-decreased"), colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(minItemsDiff.From), colorizedValue(minItemsDiff.To), colorizedValue(responseStatus)),
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
